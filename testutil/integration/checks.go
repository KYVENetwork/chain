package integration

import (
	"fmt"
	"sort"
	"time"

	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/store"
	storeTypes "cosmossdk.io/store/types"

	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	teamtypes "github.com/KYVENetwork/chain/x/team/types"

	"github.com/KYVENetwork/chain/x/funders"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"

	"github.com/KYVENetwork/chain/x/bundles"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	poolmodule "github.com/KYVENetwork/chain/x/pool"
	querytypes "github.com/KYVENetwork/chain/x/query/types"
	"github.com/KYVENetwork/chain/x/stakers"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/KYVENetwork/chain/x/team"
	. "github.com/onsi/gomega"
)

func (suite *KeeperTestSuite) PerformValidityChecks() {
	// verify pool module
	suite.VerifyPoolModuleFundingStates()
	suite.VerifyPoolQueries()
	suite.VerifyPoolGenesisImportExport()

	// verify funders module
	suite.VerifyFundersGenesisImportExport()
	suite.VerifyFundersModuleIntegrity()
	suite.VerifyFundersModuleAssetsIntegrity()

	// verify stakers module
	suite.VerifyStakersGenesisImportExport()
	suite.VerifyStakersModuleAssetsIntegrity()
	suite.VerifyPoolTotalStake()
	suite.VerifyStakersQueries()

	// verify bundles module
	suite.VerifyBundlesQueries()
	suite.VerifyBundlesGenesisImportExport()

	// verify delegation module
	suite.VerifyDelegationModuleIntegrity()

	// verify team module
	// TODO(@troy): implement team funds integrity checks
	suite.VerifyTeamGenesisImportExport()
}

// ==================
// pool module checks
// ==================

func (suite *KeeperTestSuite) VerifyPoolModuleFundingStates() {
	// every pool must have a funding state
	for _, p := range suite.App().PoolKeeper.GetAllPools(suite.Ctx()) {
		found := suite.App().FundersKeeper.DoesFundingStateExist(suite.Ctx(), p.Id)
		Expect(found).To(BeTrue())
	}
}

func (suite *KeeperTestSuite) VerifyPoolQueries() {
	poolsState := suite.App().PoolKeeper.GetAllPools(suite.Ctx())

	poolsQuery := make([]querytypes.PoolResponse, 0)

	activePoolsQuery, activePoolsQueryErr := suite.App().QueryKeeper.Pools(suite.Ctx(), &querytypes.QueryPoolsRequest{})
	disabledPoolsQuery, disabledPoolsQueryErr := suite.App().QueryKeeper.Pools(suite.Ctx(), &querytypes.QueryPoolsRequest{
		Disabled: true,
	})

	poolsQuery = append(poolsQuery, activePoolsQuery.Pools...)
	poolsQuery = append(poolsQuery, disabledPoolsQuery.Pools...)

	// sort pools by id
	for i := range poolsQuery {
		for j := range poolsQuery {
			if poolsQuery[i].Id < poolsQuery[j].Id {
				poolsQuery[i], poolsQuery[j] = poolsQuery[j], poolsQuery[i]
			}
		}
	}

	Expect(activePoolsQueryErr).To(BeNil())
	Expect(disabledPoolsQueryErr).To(BeNil())

	Expect(poolsQuery).To(HaveLen(len(poolsState)))

	for i := range poolsState {
		bundleProposalState, _ := suite.App().BundlesKeeper.GetBundleProposal(suite.Ctx(), poolsState[i].Id)
		stakersState := suite.App().StakersKeeper.GetAllStakerAddressesOfPool(suite.Ctx(), poolsState[i].Id)
		totalDelegationState := suite.App().StakersKeeper.GetTotalStakeOfPool(suite.Ctx(), poolsState[i].Id)

		Expect(poolsQuery[i].Id).To(Equal(poolsState[i].Id))
		Expect(*poolsQuery[i].Data).To(Equal(poolsState[i]))
		Expect(*poolsQuery[i].BundleProposal).To(Equal(bundleProposalState))
		Expect(poolsQuery[i].Stakers).To(Equal(stakersState))
		Expect(poolsQuery[i].TotalDelegation).To(Equal(totalDelegationState))

		// test pool by id
		poolByIdQuery, poolByIdQueryErr := suite.App().QueryKeeper.Pool(suite.Ctx(), &querytypes.QueryPoolRequest{
			Id: poolsState[i].Id,
		})

		Expect(poolByIdQueryErr).To(BeNil())
		Expect(poolByIdQuery.Pool.Id).To(Equal(poolsState[i].Id))
		Expect(*poolByIdQuery.Pool.Data).To(Equal(poolsState[i]))
		Expect(*poolByIdQuery.Pool.BundleProposal).To(Equal(bundleProposalState))
		Expect(poolByIdQuery.Pool.Stakers).To(Equal(stakersState))
		Expect(poolByIdQuery.Pool.TotalDelegation).To(Equal(totalDelegationState))

		// test stakers by pool
		valaccounts := suite.App().StakersKeeper.GetAllPoolAccountsOfPool(suite.Ctx(), poolsState[i].Id)
		stakersByPoolState := make([]querytypes.FullStaker, 0)

		for _, valaccount := range valaccounts {
			if _, stakerFound := suite.App().StakersKeeper.GetValidator(suite.Ctx(), valaccount.Staker); stakerFound {
				stakersByPoolState = append(stakersByPoolState, *suite.App().QueryKeeper.GetFullStaker(suite.Ctx(), valaccount.Staker))
			}
		}

		sort.SliceStable(stakersByPoolState, func(a, b int) bool {
			return suite.App().StakersKeeper.GetValidatorPoolStake(suite.Ctx(), stakersByPoolState[a].Address, poolsState[i].Id) > suite.App().StakersKeeper.GetValidatorPoolStake(suite.Ctx(), stakersByPoolState[b].Address, poolsState[i].Id)
		})

		stakersByPoolQuery, stakersByPoolQueryErr := suite.App().QueryKeeper.StakersByPool(suite.Ctx(), &querytypes.QueryStakersByPoolRequest{
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
	genState := poolmodule.ExportGenesis(suite.Ctx(), suite.App().PoolKeeper)

	// Delete all entries in Pool Store
	suite.deleteStore(suite.getStoreByKeyName(pooltypes.StoreKey))
	err := genState.Validate()
	Expect(err).To(BeNil())
	poolmodule.InitGenesis(suite.Ctx(), suite.App().PoolKeeper, *genState)
}

// =====================
// stakers module checks
// =====================

func (suite *KeeperTestSuite) VerifyStakersModuleAssetsIntegrity() {
	expectedBalance := sdk.NewCoins()

	moduleAcc := suite.App().AccountKeeper.GetModuleAccount(suite.Ctx(), stakerstypes.ModuleName).GetAddress()
	actualBalance := suite.App().BankKeeper.GetAllBalances(suite.Ctx(), moduleAcc)

	Expect(actualBalance.String()).To(Equal(expectedBalance.String()))
}

func (suite *KeeperTestSuite) VerifyPoolTotalStake() {
	for _, pool := range suite.App().PoolKeeper.GetAllPools(suite.Ctx()) {
		expectedBalance := uint64(0)
		actualBalance := suite.App().StakersKeeper.GetTotalStakeOfPool(suite.Ctx(), pool.Id)

		for _, stakerAddress := range suite.App().StakersKeeper.GetAllStakerAddressesOfPool(suite.Ctx(), pool.Id) {
			expectedBalance += suite.App().StakersKeeper.GetValidatorPoolStake(suite.Ctx(), stakerAddress, pool.Id)
		}

		Expect(actualBalance).To(Equal(expectedBalance))
	}
}

func (suite *KeeperTestSuite) VerifyStakersQueries() {
	validators, _ := suite.App().StakingKeeper.GetBondedValidatorsByPower(suite.Ctx())
	stakersQuery, stakersQueryErr := suite.App().QueryKeeper.Stakers(suite.Ctx(), &querytypes.QueryStakersRequest{
		Pagination: &query.PageRequest{
			Limit: 1000,
		},
	})

	stakersMap := make(map[string]querytypes.FullStaker, 0)
	for _, staker := range stakersQuery.Stakers {
		stakersMap[staker.Address] = staker
	}

	Expect(stakersQueryErr).To(BeNil())
	Expect(stakersQuery.Stakers).To(HaveLen(len(validators)))

	for i := range validators {
		address := util.MustAccountAddressFromValAddress(validators[i].OperatorAddress)
		suite.verifyFullStaker(stakersMap[address], address)

		stakerByAddressQuery, stakersByAddressQueryErr := suite.App().QueryKeeper.Staker(suite.Ctx(), &querytypes.QueryStakerRequest{
			Address: address,
		})

		Expect(stakersByAddressQueryErr).To(BeNil())
		suite.verifyFullStaker(stakerByAddressQuery.Staker, address)
	}
}

func (suite *KeeperTestSuite) VerifyStakersGenesisImportExport() {
	genState := stakers.ExportGenesis(suite.Ctx(), suite.App().StakersKeeper)

	// Delete all entries in Stakers Store
	st := suite.getStoreByKeyName(stakerstypes.StoreKey)
	iterator := st.Iterator(nil, nil)
	keys := make([][]byte, 0)
	for ; iterator.Valid(); iterator.Next() {
		key := make([]byte, len(iterator.Key()))
		copy(key, iterator.Key())
		keys = append(keys, key)
	}
	iterator.Close()
	for _, key := range keys {
		st.Delete(key)
	}

	err := genState.Validate()
	Expect(err).To(BeNil())
	stakers.InitGenesis(suite.Ctx(), suite.App().StakersKeeper, *genState)
}

// =====================
// bundles module checks
// =====================

func checkFinalizedBundle(queryBundle querytypes.QueryFinalizedBundleResponse, rawBundle bundlesTypes.FinalizedBundle) {
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
		finalizedBundlesQuery, finalizedBundlesQueryErr := suite.App().QueryKeeper.FinalizedBundlesQuery(suite.Ctx(), &querytypes.QueryFinalizedBundlesRequest{
			PoolId: pool.Id,
		})

		Expect(finalizedBundlesQueryErr).To(BeNil())
		Expect(finalizedBundlesQuery.FinalizedBundles).To(HaveLen(len(finalizedBundlesState)))

		for i := range finalizedBundlesState {

			finalizedBundle, finalizedBundleQueryErr := suite.App().QueryKeeper.FinalizedBundleQuery(suite.Ctx(), &querytypes.QueryFinalizedBundleRequest{
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

func (suite *KeeperTestSuite) VerifyDelegationModuleIntegrity() {
	expectedBalance := sdk.NewCoins()

	for _, delegator := range suite.App().DelegationKeeper.GetAllDelegators(suite.Ctx()) {
		expectedBalance = expectedBalance.Add(
			sdk.NewInt64Coin(globalTypes.Denom,
				int64(suite.App().DelegationKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), delegator.Staker, delegator.Delegator)),
			)).Add(
			suite.App().DelegationKeeper.GetOutstandingRewards(suite.Ctx(), delegator.Staker, delegator.Delegator)...,
		)
	}

	// Due to rounding errors the delegation module will get a very few nKYVE over the time.
	// As long as it is guaranteed that it's always the user who gets paid out less in case of
	// rounding, everything is fine.
	difference := suite.GetCoinsFromModule(delegationtypes.ModuleName).Sub(expectedBalance...)
	//nolint:all
	Expect(difference.IsAnyNegative()).To(BeFalse())

	// 10 should be enough for testing, these are left-over tokens due to rounding issues
	for _, coin := range difference {
		Expect(coin.Amount.Uint64() < 10).To(BeTrue())
	}
}

// =========================
// team module checks
// =========================

func (suite *KeeperTestSuite) VerifyTeamGenesisImportExport() {
	genState := team.ExportGenesis(suite.Ctx(), suite.App().TeamKeeper)

	// Delete all entries in Stakers Store
	st := suite.getStoreByKeyName(teamtypes.StoreKey)
	iterator := st.Iterator(nil, nil)
	keys := make([][]byte, 0)
	for ; iterator.Valid(); iterator.Next() {
		key := make([]byte, len(iterator.Key()))
		copy(key, iterator.Key())
		keys = append(keys, key)
	}
	iterator.Close()
	for _, key := range keys {
		st.Delete(key)
	}

	err := genState.Validate()
	Expect(err).To(BeNil())
	team.InitGenesis(suite.Ctx(), suite.App().TeamKeeper, *genState)
}

// ========================
// funders module checks
// ========================

func (suite *KeeperTestSuite) VerifyFundersGenesisImportExport() {
	genState := funders.ExportGenesis(suite.Ctx(), suite.App().FundersKeeper)

	// Delete all entries in Funders Store
	suite.deleteStore(suite.getStoreByKeyName(funderstypes.StoreKey))

	err := genState.Validate()
	Expect(err).To(BeNil())
	funders.InitGenesis(suite.Ctx(), suite.App().FundersKeeper, *genState)
}

func (suite *KeeperTestSuite) VerifyFundersModuleIntegrity() {
	funderAddresses := make(map[string]bool)
	for _, funder := range suite.App().FundersKeeper.GetAllFunders(suite.Ctx()) {
		funderAddresses[funder.Address] = true
	}

	allActiveFundings := make(map[string]bool)
	for _, funding := range suite.App().FundersKeeper.GetAllFundings(suite.Ctx()) {
		// check if funding has a valid funder
		_, found := funderAddresses[funding.FunderAddress]
		Expect(found).To(BeTrue())

		// check if funding is active
		if !funding.Amounts.IsZero() {
			key := string(funderstypes.FundingKeyByFunder(funding.FunderAddress, funding.PoolId))
			allActiveFundings[key] = true
		}

		// check if pool exists
		_, found = suite.App().PoolKeeper.GetPool(suite.Ctx(), funding.PoolId)
		Expect(found).To(BeTrue())
	}

	for _, fundingState := range suite.App().FundersKeeper.GetAllFundingStates(suite.Ctx()) {
		fsActiveAddresses := make(map[string]bool)
		for _, funderAddress := range fundingState.ActiveFunderAddresses {
			// check if funding has a valid funder
			key := funderstypes.FundingKeyByFunder(funderAddress, fundingState.PoolId)
			_, found := allActiveFundings[string(key)]
			Expect(found).To(BeTrue())

			// check if funder is not already in the list
			Expect(fsActiveAddresses[funderAddress]).To(BeFalse())
			fsActiveAddresses[funderAddress] = true
		}

		// check if the amount of active fundings is equal to the amount of active funder addresses
		activeFundings := suite.App().FundersKeeper.GetActiveFundings(suite.Ctx(), fundingState)
		Expect(activeFundings).To(HaveLen(len(fundingState.ActiveFunderAddresses)))

		// be lower or equal to max funders
		Expect(len(fundingState.ActiveFunderAddresses)).To(BeNumerically("<=", funderstypes.MaxFunders))
	}
}

func (suite *KeeperTestSuite) VerifyFundersModuleAssetsIntegrity() {
	expectedBalance := sdk.NewCoins()
	for _, funding := range suite.App().FundersKeeper.GetAllFundings(suite.Ctx()) {
		expectedBalance = expectedBalance.Add(funding.Amounts...)
	}

	expectedFundingStateTotalAmount := sdk.NewCoins()
	for _, fundingState := range suite.App().FundersKeeper.GetAllFundingStates(suite.Ctx()) {
		activeFundings := suite.App().FundersKeeper.GetActiveFundings(suite.Ctx(), fundingState)
		totalAmount := sdk.NewCoins()
		for _, activeFunding := range activeFundings {
			totalAmount = totalAmount.Add(activeFunding.Amounts...)
		}
		totalActiveFunding := suite.App().FundersKeeper.GetTotalActiveFunding(suite.ctx, fundingState.PoolId)
		Expect(totalAmount.String()).To(Equal(totalActiveFunding.String()))
		expectedFundingStateTotalAmount = expectedFundingStateTotalAmount.Add(totalAmount...)
	}

	// total amount of fundings should be equal to the amount of the funders module account
	moduleAcc := suite.App().AccountKeeper.GetModuleAccount(suite.Ctx(), funderstypes.ModuleName).GetAddress()
	actualBalance := suite.App().BankKeeper.GetAllBalances(suite.Ctx(), moduleAcc)
	Expect(actualBalance.String()).To(Equal(expectedBalance.String()))
	Expect(actualBalance.String()).To(Equal(expectedFundingStateTotalAmount.String()))
}

// ========================
// helpers
// ========================

func (suite *KeeperTestSuite) verifyFullStaker(fullStaker querytypes.FullStaker, stakerAddress string) {
	Expect(fullStaker.Address).To(Equal(stakerAddress))

	// TODO after reworking the API Queries
	//staker, found := suite.App().StakersKeeper.GetValidator(suite.Ctx(), stakerAddress)
	//Expect(found).To(BeTrue())
	//Expect(fullStaker.SelfDelegation).To(Equal(suite.App().StakersKeeper.GetDelegationAmountOfDelegator(suite.Ctx(), stakerAddress, stakerAddress)))
	//
	//selfDelegationUnbonding := uint64(0)
	//for _, entry := range suite.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(suite.Ctx(), fullStaker.Address) {
	//	if entry.Staker == stakerAddress {
	//		selfDelegationUnbonding += entry.Amount
	//	}
	//}
	//
	//Expect(fullStaker.SelfDelegationUnbonding).To(Equal(selfDelegationUnbonding))
	//Expect(fullStaker.Metadata.Identity).To(Equal(staker.Description.Identity))
	//Expect(fullStaker.Metadata.SecurityContact).To(Equal(staker.Description.SecurityContact))
	//Expect(fullStaker.Metadata.Details).To(Equal(staker.Description.Details))
	//Expect(fullStaker.Metadata.Website).To(Equal(staker.Description.Website))
	//Expect(fullStaker.Metadata.Commission).To(Equal(staker.Commission))
	//Expect(fullStaker.Metadata.Moniker).To(Equal(staker.Description.Moniker))

	// TODO rework after commission was implemented
	//pendingCommissionChange, found := suite.App().StakersKeeper.GetCommissionChangeEntryByIndex2(suite.Ctx(), stakerAddress)
	//if found {
	//	Expect(fullStaker.Metadata.PendingCommissionChange.Commission).To(Equal(pendingCommissionChange.Commission))
	//	Expect(fullStaker.Metadata.PendingCommissionChange.CreationDate).To(Equal(pendingCommissionChange.CreationDate))
	//} else {
	//	Expect(fullStaker.Metadata.PendingCommissionChange).To(BeNil())
	//}

	poolIds := make(map[uint64]bool)

	for _, poolMembership := range fullStaker.Pools {
		poolIds[poolMembership.Pool.Id] = true
		valaccount, active := suite.App().StakersKeeper.GetPoolAccount(suite.Ctx(), stakerAddress, poolMembership.Pool.Id)
		Expect(active).To(BeTrue())

		Expect(poolMembership.PoolAddress).To(Equal(valaccount.PoolAddress))
		Expect(poolMembership.IsLeaving).To(Equal(valaccount.IsLeaving))
		Expect(poolMembership.Points).To(Equal(valaccount.Points))

		pool, found := suite.App().PoolKeeper.GetPool(suite.Ctx(), valaccount.PoolId)
		Expect(found).To(BeTrue())
		Expect(poolMembership.Pool.Id).To(Equal(pool.Id))
		Expect(poolMembership.Pool.Logo).To(Equal(pool.Logo))

		fundingState, found := suite.App().FundersKeeper.GetFundingState(suite.Ctx(), poolMembership.Pool.Id)
		Expect(found).To(BeTrue())
		Expect(poolMembership.Pool.TotalFunds).To(Equal(suite.App().FundersKeeper.GetTotalActiveFunding(suite.Ctx(), fundingState.PoolId)))
		Expect(poolMembership.Pool.Name).To(Equal(pool.Name))
		Expect(poolMembership.Pool.Runtime).To(Equal(pool.Runtime))
		Expect(poolMembership.Pool.Status).To(Equal(suite.App().QueryKeeper.GetPoolStatus(suite.Ctx(), &pool)))
		Expect(poolMembership.Pool.Status).To(Equal(suite.App().QueryKeeper.GetPoolStatus(suite.Ctx(), &pool)))
	}

	// Reverse check the pool memberships
	for _, valaccount := range suite.App().StakersKeeper.GetPoolAccountsFromStaker(suite.Ctx(), stakerAddress) {
		Expect(poolIds[valaccount.PoolId]).To(BeTrue())
	}
}

func (suite *KeeperTestSuite) deleteStore(store store.KVStore) {
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

func (suite *KeeperTestSuite) getStoreByKeyName(keyName string) storeTypes.KVStore {
	keys := suite.app.GetStoreKeys()
	for _, key := range keys {
		if key.Name() == keyName {
			return suite.Ctx().KVStore(key)
		}
	}
	panic(fmt.Errorf("store with name %s not found", keyName))
}
