package integration

import (
	"time"

	"github.com/KYVENetwork/chain/x/bundles"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/KYVENetwork/chain/x/delegation"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/KYVENetwork/chain/x/pool"
	querytypes "github.com/KYVENetwork/chain/x/query/types"
	"github.com/KYVENetwork/chain/x/stakers"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/KYVENetwork/chain/x/team"
	"github.com/cosmos/cosmos-sdk/types/query"
	. "github.com/onsi/gomega"

	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) PerformValidityChecks() {
	// verify pool module
	suite.VerifyPoolModuleAssetsIntegrity()
	suite.VerifyPoolTotalFunds()

	if suite.options.VerifyPoolQueries {
		suite.VerifyPoolQueries()
	}

	suite.VerifyPoolGenesisImportExport()

	// verify stakers module
	suite.VerifyStakersGenesisImportExport()
	suite.VerifyStakersModuleAssetsIntegrity()
	suite.VerifyPoolTotalStake()
	suite.VerifyStakersQueries()
	suite.VerifyActiveStakers()

	// verify bundles module
	suite.VerifyBundlesQueries()
	suite.VerifyBundlesGenesisImportExport()

	// verify delegation module
	suite.VerifyDelegationQueries()
	suite.VerifyDelegationModuleIntegrity()
	suite.VerifyDelegationGenesisImportExport()

	// verify team module
	// TODO(@troy): implement team funds integrity checks
	suite.VerifyTeamGenesisImportExport()
}

// ==================
// pool module checks
// ==================

func (suite *KeeperTestSuite) VerifyPoolModuleAssetsIntegrity() {
	expectedBalance := uint64(0)
	actualBalance := uint64(0)

	for _, pool := range suite.App().PoolKeeper.GetAllPools(suite.Ctx()) {
		// pool funds should be in pool module
		for _, funder := range pool.Funders {
			expectedBalance += funder.Amount
		}
	}

	moduleAcc := suite.App().AccountKeeper.GetModuleAccount(suite.Ctx(), pooltypes.ModuleName).GetAddress()
	actualBalance = suite.App().BankKeeper.GetBalance(suite.Ctx(), moduleAcc, globalTypes.Denom).Amount.Uint64()

	Expect(actualBalance).To(Equal(expectedBalance))
}

func (suite *KeeperTestSuite) VerifyPoolTotalFunds() {
	for _, pool := range suite.App().PoolKeeper.GetAllPools(suite.Ctx()) {
		expectedBalance := uint64(0)
		actualBalance := pool.TotalFunds

		for _, funder := range pool.Funders {
			expectedBalance += funder.Amount
		}

		Expect(actualBalance).To(Equal(expectedBalance))
	}
}

func (suite *KeeperTestSuite) VerifyPoolQueries() {
	poolsState := suite.App().PoolKeeper.GetAllPools(suite.Ctx())

	poolsQuery := make([]querytypes.PoolResponse, 0)

	activePoolsQuery, activePoolsQueryErr := suite.App().QueryKeeper.Pools(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryPoolsRequest{})
	disabledPoolsQuery, disabledPoolsQueryErr := suite.App().QueryKeeper.Pools(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryPoolsRequest{
		Disabled: true,
	})

	poolsQuery = append(poolsQuery, activePoolsQuery.Pools...)
	poolsQuery = append(poolsQuery, disabledPoolsQuery.Pools...)

	Expect(activePoolsQueryErr).To(BeNil())
	Expect(disabledPoolsQueryErr).To(BeNil())

	Expect(poolsQuery).To(HaveLen(len(poolsState)))

	for i := range poolsState {
		bundleProposalState, _ := suite.App().BundlesKeeper.GetBundleProposal(suite.Ctx(), poolsState[i].Id)
		stakersState := suite.App().StakersKeeper.GetAllStakerAddressesOfPool(suite.Ctx(), poolsState[i].Id)
		totalDelegationState := suite.App().DelegationKeeper.GetDelegationOfPool(suite.Ctx(), poolsState[i].Id)

		Expect(poolsQuery[i].Id).To(Equal(poolsState[i].Id))
		Expect(*poolsQuery[i].Data).To(Equal(poolsState[i]))
		Expect(*poolsQuery[i].BundleProposal).To(Equal(bundleProposalState))
		Expect(poolsQuery[i].Stakers).To(Equal(stakersState))
		Expect(poolsQuery[i].TotalDelegation).To(Equal(totalDelegationState))

		// test pool by id
		poolByIdQuery, poolByIdQueryErr := suite.App().QueryKeeper.Pool(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryPoolRequest{
			Id: poolsState[i].Id,
		})

		Expect(poolByIdQueryErr).To(BeNil())
		Expect(poolByIdQuery.Pool.Id).To(Equal(poolsState[i].Id))
		Expect(*poolByIdQuery.Pool.Data).To(Equal(poolsState[i]))
		Expect(*poolByIdQuery.Pool.BundleProposal).To(Equal(bundleProposalState))
		Expect(poolByIdQuery.Pool.Stakers).To(Equal(stakersState))
		Expect(poolByIdQuery.Pool.TotalDelegation).To(Equal(totalDelegationState))

		// test stakers by pool
		valaccounts := suite.App().StakersKeeper.GetAllValaccountsOfPool(suite.Ctx(), poolsState[i].Id)
		stakersByPoolState := make([]querytypes.StakerPoolResponse, 0)

		for _, valaccount := range valaccounts {
			staker, stakerFound := suite.App().StakersKeeper.GetStaker(suite.Ctx(), valaccount.Staker)

			if stakerFound {
				stakersByPoolState = append(stakersByPoolState, querytypes.StakerPoolResponse{
					Staker:     suite.App().QueryKeeper.GetFullStaker(suite.Ctx(), staker.Address),
					Valaccount: valaccount,
				})
			}
		}

		stakersByPoolQuery, stakersByPoolQueryErr := suite.App().QueryKeeper.StakersByPool(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryStakersByPoolRequest{
			PoolId: poolsState[i].Id,
		})

		Expect(stakersByPoolQueryErr).To(BeNil())
		Expect(stakersByPoolQuery.Stakers).To(HaveLen(len(stakersByPoolState)))

		for s := range stakersByPoolState {
			Expect(stakersByPoolQuery.Stakers[s]).To(Equal(stakersByPoolState[s]))
		}
	}
}

func (suite *KeeperTestSuite) VerifyPoolGenesisImportExport() {
	genState := pool.ExportGenesis(suite.Ctx(), suite.App().PoolKeeper)

	// Delete all entries in Pool Store
	store := suite.Ctx().KVStore(suite.App().PoolKeeper.StoreKey())
	suite.deleteStore(store)

	err := genState.Validate()
	Expect(err).To(BeNil())
	pool.InitGenesis(suite.Ctx(), suite.App().PoolKeeper, *genState)
}

// =====================
// stakers module checks
// =====================

func (suite *KeeperTestSuite) VerifyStakersModuleAssetsIntegrity() {
	expectedBalance := uint64(0)
	actualBalance := uint64(0)

	for _, staker := range suite.App().StakersKeeper.GetAllStakers(suite.Ctx()) {
		expectedBalance += staker.CommissionRewards
	}

	moduleAcc := suite.App().AccountKeeper.GetModuleAccount(suite.Ctx(), stakertypes.ModuleName).GetAddress()
	actualBalance = suite.App().BankKeeper.GetBalance(suite.Ctx(), moduleAcc, globalTypes.Denom).Amount.Uint64()

	Expect(actualBalance).To(Equal(expectedBalance))
}

func (suite *KeeperTestSuite) VerifyPoolTotalStake() {
	for _, pool := range suite.App().PoolKeeper.GetAllPools(suite.Ctx()) {
		expectedBalance := uint64(0)
		actualBalance := suite.App().DelegationKeeper.GetDelegationOfPool(suite.Ctx(), pool.Id)

		for _, stakerAddress := range suite.App().StakersKeeper.GetAllStakerAddressesOfPool(suite.Ctx(), pool.Id) {
			expectedBalance += suite.App().DelegationKeeper.GetDelegationAmount(suite.Ctx(), stakerAddress)
		}

		Expect(actualBalance).To(Equal(expectedBalance))
	}
}

func (suite *KeeperTestSuite) VerifyActiveStakers() {
	totalDelegation := uint64(0)
	for _, delegator := range suite.App().DelegationKeeper.GetAllDelegators(suite.Ctx()) {
		if len(suite.App().StakersKeeper.GetValaccountsFromStaker(suite.Ctx(), delegator.Staker)) > 0 {
			totalDelegation += suite.App().DelegationKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), delegator.Staker, delegator.Delegator)

			validators, _ := suite.App().StakersKeeper.GetDelegations(suite.ctx, delegator.Delegator)
			Expect(validators).To(ContainElement(delegator.Staker))
		}
	}
	Expect(suite.App().StakersKeeper.TotalBondedTokens(suite.Ctx()).Uint64()).To(Equal(totalDelegation))
}

func (suite *KeeperTestSuite) VerifyStakersQueries() {
	stakersState := suite.App().StakersKeeper.GetAllStakers(suite.Ctx())
	stakersQuery, stakersQueryErr := suite.App().QueryKeeper.Stakers(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryStakersRequest{
		Pagination: &query.PageRequest{
			Limit: 1000,
		},
	})

	stakersMap := make(map[string]querytypes.FullStaker, 0)
	for _, staker := range stakersQuery.Stakers {
		stakersMap[staker.Address] = staker
	}

	Expect(stakersQueryErr).To(BeNil())
	Expect(stakersQuery.Stakers).To(HaveLen(len(stakersState)))

	for i := range stakersState {
		address := stakersState[i].Address
		suite.verifyFullStaker(stakersMap[address], address)

		stakerByAddressQuery, stakersByAddressQueryErr := suite.App().QueryKeeper.Staker(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryStakerRequest{
			Address: address,
		})

		Expect(stakersByAddressQueryErr).To(BeNil())
		suite.verifyFullStaker(stakerByAddressQuery.Staker, address)
	}
}

func (suite *KeeperTestSuite) VerifyStakersGenesisImportExport() {
	genState := stakers.ExportGenesis(suite.Ctx(), suite.App().StakersKeeper)

	// Delete all entries in Stakers Store
	store := suite.Ctx().KVStore(suite.App().StakersKeeper.StoreKey())
	iterator := store.Iterator(nil, nil)
	keys := make([][]byte, 0)
	for ; iterator.Valid(); iterator.Next() {
		key := make([]byte, len(iterator.Key()))
		copy(key, iterator.Key())
		keys = append(keys, key)
	}
	iterator.Close()
	for _, key := range keys {
		store.Delete(key)
	}

	err := genState.Validate()
	Expect(err).To(BeNil())
	stakers.InitGenesis(suite.Ctx(), suite.App().StakersKeeper, *genState)
}

// =====================
// bundles module checks
// =====================

func checkFinalizedBundle(queryBundle querytypes.FinalizedBundle, rawBundle bundlesTypes.FinalizedBundle) {
	Expect(queryBundle.Id).To(Equal(rawBundle.Id))
	Expect(queryBundle.PoolId).To(Equal(rawBundle.PoolId))
	Expect(queryBundle.StorageId).To(Equal(rawBundle.StorageId))
	Expect(queryBundle.Uploader).To(Equal(rawBundle.Uploader))
	Expect(queryBundle.FromIndex).To(Equal(rawBundle.FromIndex))
	Expect(queryBundle.ToIndex).To(Equal(rawBundle.ToIndex))
	Expect(queryBundle.ToKey).To(Equal(rawBundle.ToKey))
	Expect(queryBundle.BundleSummary).To(Equal(rawBundle.BundleSummary))
	Expect(queryBundle.DataHash).To(Equal(rawBundle.DataHash))
	Expect(queryBundle.FinalizedAt.Height.Uint64()).To(Equal(rawBundle.FinalizedAt.Height))
	date, dateErr := time.Parse(time.RFC3339, queryBundle.FinalizedAt.Timestamp)
	Expect(dateErr).To(BeNil())
	Expect(uint64(date.Unix())).To(Equal(rawBundle.FinalizedAt.Timestamp))
	Expect(queryBundle.FromKey).To(Equal(rawBundle.FromKey))
	Expect(queryBundle.StorageProviderId).To(Equal(uint64(rawBundle.StorageProviderId)))
	Expect(queryBundle.CompressionId).To(Equal(uint64(rawBundle.CompressionId)))
	Expect(queryBundle.StakeSecurity.ValidVotePower.Uint64()).To(Equal(rawBundle.StakeSecurity.ValidVotePower))
	Expect(queryBundle.StakeSecurity.TotalVotePower.Uint64()).To(Equal(rawBundle.StakeSecurity.TotalVotePower))
}

func (suite *KeeperTestSuite) VerifyBundlesQueries() {
	pools := suite.App().PoolKeeper.GetAllPools(suite.Ctx())

	for _, pool := range pools {
		finalizedBundlesState := suite.App().BundlesKeeper.GetFinalizedBundlesByPool(suite.Ctx(), pool.Id)
		finalizedBundlesQuery, finalizedBundlesQueryErr := suite.App().QueryKeeper.FinalizedBundlesQuery(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryFinalizedBundlesRequest{
			PoolId: pool.Id,
		})

		Expect(finalizedBundlesQueryErr).To(BeNil())
		Expect(finalizedBundlesQuery.FinalizedBundles).To(HaveLen(len(finalizedBundlesState)))

		for i := range finalizedBundlesState {

			finalizedBundle, finalizedBundleQueryErr := suite.App().QueryKeeper.FinalizedBundleQuery(sdk.WrapSDKContext(suite.Ctx()), &querytypes.QueryFinalizedBundleRequest{
				PoolId: pool.Id,
				Id:     finalizedBundlesState[i].Id,
			})

			Expect(finalizedBundleQueryErr).To(BeNil())

			checkFinalizedBundle(*finalizedBundle, finalizedBundlesState[i])
		}
	}
}

func (suite *KeeperTestSuite) VerifyBundlesGenesisImportExport() {
	genState := bundles.ExportGenesis(suite.Ctx(), suite.App().BundlesKeeper)
	err := genState.Validate()
	Expect(err).To(BeNil())
	bundles.InitGenesis(suite.Ctx(), suite.App().BundlesKeeper, *genState)
}

// ========================
// delegation module checks
// ========================

func (suite *KeeperTestSuite) VerifyDelegationQueries() {
	goCtx := sdk.WrapSDKContext(suite.Ctx())
	for _, delegator := range suite.App().DelegationKeeper.GetAllDelegators(suite.Ctx()) {

		// Query: delegator/{staker}/{delegator}
		resD, errD := suite.App().QueryKeeper.Delegator(goCtx, &querytypes.QueryDelegatorRequest{
			Staker:    delegator.Staker,
			Delegator: delegator.Delegator,
		})
		Expect(errD).To(BeNil())
		Expect(resD.Delegator.Delegator).To(Equal(delegator.Delegator))
		Expect(resD.Delegator.Staker).To(Equal(delegator.Staker))
		Expect(resD.Delegator.DelegationAmount).To(Equal(suite.App().DelegationKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), delegator.Staker, delegator.Delegator)))
		Expect(resD.Delegator.CurrentReward).To(Equal(suite.App().DelegationKeeper.GetOutstandingRewards(suite.Ctx(), delegator.Staker, delegator.Delegator)))

		// Query: stakers_by_delegator/{delegator}
		resSbD, errSbD := suite.App().QueryKeeper.StakersByDelegator(goCtx, &querytypes.QueryStakersByDelegatorRequest{
			Pagination: nil,
			Delegator:  delegator.Delegator,
		})
		Expect(errSbD).To(BeNil())
		Expect(resSbD.Delegator).To(Equal(delegator.Delegator))
		for _, sRes := range resSbD.Stakers {
			Expect(sRes.DelegationAmount).To(Equal(suite.App().DelegationKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), sRes.Staker.Address, delegator.Delegator)))
			Expect(sRes.CurrentReward).To(Equal(suite.App().DelegationKeeper.GetOutstandingRewards(suite.Ctx(), sRes.Staker.Address, delegator.Delegator)))
			suite.verifyFullStaker(*sRes.Staker, sRes.Staker.Address)
		}
	}

	stakersDelegators := make(map[string]map[string]delegationtypes.Delegator)
	for _, d := range suite.App().DelegationKeeper.GetAllDelegators(suite.Ctx()) {
		if stakersDelegators[d.Staker] == nil {
			stakersDelegators[d.Staker] = map[string]delegationtypes.Delegator{}
		}
		stakersDelegators[d.Staker][d.Delegator] = d
	}

	for _, staker := range suite.App().StakersKeeper.GetAllStakers(suite.Ctx()) {
		// Query: delegators_by_staker/{staker}
		resDbS, errDbS := suite.App().QueryKeeper.DelegatorsByStaker(goCtx, &querytypes.QueryDelegatorsByStakerRequest{
			Pagination: nil,
			Staker:     staker.Address,
		})
		Expect(errDbS).To(BeNil())

		delegationData, _ := suite.App().DelegationKeeper.GetDelegationData(suite.Ctx(), staker.Address)
		Expect(resDbS.TotalDelegatorCount).To(Equal(delegationData.DelegatorCount))
		Expect(resDbS.TotalDelegation).To(Equal(suite.App().DelegationKeeper.GetDelegationAmount(suite.Ctx(), staker.Address)))

		for _, delegator := range resDbS.Delegators {
			Expect(stakersDelegators[delegator.Staker][delegator.Delegator]).ToNot(BeNil())
			Expect(delegator.DelegationAmount).To(Equal(suite.App().DelegationKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), delegator.Staker, delegator.Delegator)))
			Expect(delegator.CurrentReward).To(Equal(suite.App().DelegationKeeper.GetOutstandingRewards(suite.Ctx(), delegator.Staker, delegator.Delegator)))
		}
	}
}

func (suite *KeeperTestSuite) VerifyDelegationModuleIntegrity() {
	expectedBalance := uint64(0)

	for _, delegator := range suite.App().DelegationKeeper.GetAllDelegators(suite.Ctx()) {
		expectedBalance += suite.App().DelegationKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), delegator.Staker, delegator.Delegator)
		expectedBalance += suite.App().DelegationKeeper.GetOutstandingRewards(suite.Ctx(), delegator.Staker, delegator.Delegator)
	}

	// Due to rounding errors the delegation module will get a very few nKYVE over the time.
	// As long as it is guaranteed that it's always the user who gets paid out less in case of
	// rounding, everything is fine.
	difference := suite.GetBalanceFromModule(delegationtypes.ModuleName) - expectedBalance
	//nolint:all
	Expect(difference >= 0).To(BeTrue())

	// 10 should be enough for testing
	Expect(difference <= 10).To(BeTrue())
}

func (suite *KeeperTestSuite) VerifyDelegationGenesisImportExport() {
	genState := delegation.ExportGenesis(suite.Ctx(), suite.App().DelegationKeeper)
	err := genState.Validate()
	Expect(err).To(BeNil())
	delegation.InitGenesis(suite.Ctx(), suite.App().DelegationKeeper, *genState)
}

// =========================
// team module checks
// =========================

func (suite *KeeperTestSuite) VerifyTeamGenesisImportExport() {
	genState := team.ExportGenesis(suite.Ctx(), suite.App().TeamKeeper)

	// Delete all entries in Stakers Store
	store := suite.Ctx().KVStore(suite.App().TeamKeeper.StoreKey())
	iterator := store.Iterator(nil, nil)
	keys := make([][]byte, 0)
	for ; iterator.Valid(); iterator.Next() {
		key := make([]byte, len(iterator.Key()))
		copy(key, iterator.Key())
		keys = append(keys, key)
	}
	iterator.Close()
	for _, key := range keys {
		store.Delete(key)
	}

	err := genState.Validate()
	Expect(err).To(BeNil())
	team.InitGenesis(suite.Ctx(), suite.App().TeamKeeper, *genState)
}

// ========================
// helpers
// ========================

func (suite *KeeperTestSuite) verifyFullStaker(fullStaker querytypes.FullStaker, stakerAddress string) {
	Expect(fullStaker.Address).To(Equal(stakerAddress))

	staker, found := suite.App().StakersKeeper.GetStaker(suite.Ctx(), stakerAddress)
	Expect(found).To(BeTrue())
	Expect(fullStaker.SelfDelegation).To(Equal(suite.App().DelegationKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), stakerAddress, stakerAddress)))

	selfDelegationUnbonding := uint64(0)
	for _, entry := range suite.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(suite.Ctx(), fullStaker.Address) {
		if entry.Staker == stakerAddress {
			selfDelegationUnbonding += entry.Amount
		}
	}

	Expect(fullStaker.SelfDelegationUnbonding).To(Equal(selfDelegationUnbonding))
	Expect(fullStaker.Metadata.Identity).To(Equal(staker.Identity))
	Expect(fullStaker.Metadata.SecurityContact).To(Equal(staker.SecurityContact))
	Expect(fullStaker.Metadata.Details).To(Equal(staker.Details))
	Expect(fullStaker.Metadata.Website).To(Equal(staker.Website))
	Expect(fullStaker.Metadata.Commission).To(Equal(staker.Commission))
	Expect(fullStaker.Metadata.Moniker).To(Equal(staker.Moniker))

	pendingCommissionChange, found := suite.App().StakersKeeper.GetCommissionChangeEntryByIndex2(suite.Ctx(), stakerAddress)
	if found {
		Expect(fullStaker.Metadata.PendingCommissionChange.Commission).To(Equal(pendingCommissionChange.Commission))
		Expect(fullStaker.Metadata.PendingCommissionChange.CreationDate).To(Equal(pendingCommissionChange.CreationDate))
	} else {
		Expect(fullStaker.Metadata.PendingCommissionChange).To(BeNil())
	}

	delegationData, _ := suite.App().DelegationKeeper.GetDelegationData(suite.Ctx(), stakerAddress)
	Expect(fullStaker.DelegatorCount).To(Equal(delegationData.DelegatorCount))

	Expect(fullStaker.TotalDelegation).To(Equal(suite.App().DelegationKeeper.GetDelegationAmount(suite.Ctx(), stakerAddress)))

	poolIds := make(map[uint64]bool)

	for _, poolMembership := range fullStaker.Pools {
		poolIds[poolMembership.Pool.Id] = true
		valaccount, found := suite.App().StakersKeeper.GetValaccount(suite.Ctx(), poolMembership.Pool.Id, stakerAddress)
		Expect(found).To(BeTrue())

		Expect(poolMembership.Valaddress).To(Equal(valaccount.Valaddress))
		Expect(poolMembership.IsLeaving).To(Equal(valaccount.IsLeaving))
		Expect(poolMembership.Points).To(Equal(valaccount.Points))

		pool, found := suite.App().PoolKeeper.GetPool(suite.Ctx(), valaccount.PoolId)
		Expect(found).To(BeTrue())
		Expect(poolMembership.Pool.Id).To(Equal(pool.Id))
		Expect(poolMembership.Pool.Logo).To(Equal(pool.Logo))
		Expect(poolMembership.Pool.TotalFunds).To(Equal(pool.TotalFunds))
		Expect(poolMembership.Pool.Name).To(Equal(pool.Name))
		Expect(poolMembership.Pool.Runtime).To(Equal(pool.Runtime))
		Expect(poolMembership.Pool.Status).To(Equal(suite.App().QueryKeeper.GetPoolStatus(suite.Ctx(), &pool)))
	}

	// Reverse check the pool memberships
	for _, valaccount := range suite.App().StakersKeeper.GetValaccountsFromStaker(suite.Ctx(), stakerAddress) {
		Expect(poolIds[valaccount.PoolId]).To(BeTrue())
	}
}

func (suite *KeeperTestSuite) deleteStore(store sdk.KVStore) {
	iterator := store.Iterator(nil, nil)
	keys := make([][]byte, 0)
	for ; iterator.Valid(); iterator.Next() {
		key := make([]byte, len(iterator.Key()))
		copy(key, iterator.Key())
		keys = append(keys, key)
	}
	iterator.Close()
	for _, key := range keys {
		store.Delete(key)
	}
}
