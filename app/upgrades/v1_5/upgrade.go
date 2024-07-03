package v1_5

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"

	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	v1_4_bundles "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/bundles"
	v1_4_delegation "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/delegation"
	v1_4_funders "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/funders"
	v1_4_gov "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/gov"
	v1_4_pool "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/pool"
	v1_4_stakers "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/stakers"

	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"

	"cosmossdk.io/math"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v1.5.0"
)

var logger log.Logger

// MustGetStoreKey is a helper to directly access the KV-Store
func MustGetStoreKey(storeKeys []storetypes.StoreKey, storeName string) storetypes.StoreKey {
	for _, k := range storeKeys {
		if k.Name() == storeName {
			return k
		}
	}

	panic(fmt.Sprintf("failed to find store key: %s", fundersTypes.StoreKey))
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.Codec,
	storeKeys []storetypes.StoreKey,
	bundlesKeeper bundlesKeeper.Keeper,
	delegationKeeper delegationKeeper.Keeper,
	fundersKeeper fundersKeeper.Keeper,
	stakersKeeper *stakersKeeper.Keeper,
	poolKeeper *poolKeeper.Keeper,
	govKeeper *govKeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger = sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// Run KYVE migrations

		// migrate fundings
		migrateFundersModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, fundersTypes.StoreKey), fundersKeeper)

		// migrate delegations
		migrateDelegationModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, delegationTypes.StoreKey), delegationKeeper)

		// migrate stakers
		migrateStakersModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, stakersTypes.StoreKey), stakersKeeper)

		// migrate bundles
		migrateBundlesModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, bundlesTypes.StoreKey), bundlesKeeper)

		// migrate pool
		migratePoolModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, poolTypes.StoreKey), poolKeeper)

		// Run cosmos migrations
		migratedVersionMap, err := mm.RunMigrations(ctx, configurator, fromVM)

		// migrate gov params
		migrateGovParams(sdkCtx, govKeeper)

		// migrate old MsgCreatePool gov proposals
		migrateOldGovProposals(sdkCtx, cdc, MustGetStoreKey(storeKeys, govTypes.StoreKey))

		return migratedVersionMap, err
	}
}

func migrateGovParams(sdkCtx sdk.Context, govKeeper *govKeeper.Keeper) {
	params, err := govKeeper.Params.Get(sdkCtx)
	if err != nil {
		logger.Error("failed to get gov params (err=%s)", err)
	}
	params.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, 50_000_000_000))
	expeditedVotingPeriod := 30 * time.Minute
	params.ExpeditedVotingPeriod = &expeditedVotingPeriod

	err = govKeeper.Params.Set(sdkCtx, params)
	if err != nil {
		logger.Error("failed to update gov params (err=%s)", err)
	}

	logger.Info("migrated Gov params")
}

// migrateOldGovProposals migrated the MsgCreatePool in all proposals to the new schema.
// i.e. changing inflation_share_weight from uint64 to Dec
func migrateOldGovProposals(sdkCtx sdk.Context, cdc codec.Codec, govStoreKey storetypes.StoreKey) {
	proposalStore := prefix.NewStore(sdkCtx.KVStore(govStoreKey), govTypes.ProposalsKeyPrefix)

	proposalIterator := storetypes.KVStorePrefixIterator(proposalStore, []byte{})
	defer proposalIterator.Close()

	// Iterate all existing gov proposals
	migratedMessagesCounter := 0
	for ; proposalIterator.Valid(); proposalIterator.Next() {
		var proposal v1_4_gov.Proposal

		if err := cdc.Unmarshal(proposalIterator.Value(), &proposal); err != nil {
			logger.Error(fmt.Sprintf("could not unmarshal gov proposal %s", hex.EncodeToString(proposalIterator.Key())))
			continue
		}

		// Iterate the messages of each proposal
		for idx := range proposal.Messages {
			// Check if message needs to be migrated
			if proposal.Messages[idx].TypeUrl == "/kyve.pool.v1beta1.MsgCreatePool" {
				var oldMsgCreatePool v1_4_pool.MsgCreatePool
				if err := cdc.Unmarshal(proposal.Messages[idx].Value, &oldMsgCreatePool); err != nil {
					logger.Error(fmt.Sprintf("could not unmarshal gov message (proposal=%d, message_idx=%d)", proposal.Id, idx))
					continue
				}

				newMsgCreatePool := poolTypes.MsgCreatePool{
					Authority:            oldMsgCreatePool.Authority,
					Name:                 oldMsgCreatePool.Name,
					Runtime:              oldMsgCreatePool.Runtime,
					Logo:                 oldMsgCreatePool.Logo,
					Config:               oldMsgCreatePool.Config,
					StartKey:             oldMsgCreatePool.StartKey,
					UploadInterval:       oldMsgCreatePool.UploadInterval,
					InflationShareWeight: math.LegacyNewDec(int64(oldMsgCreatePool.InflationShareWeight)),
					MinDelegation:        oldMsgCreatePool.MinDelegation,
					MaxBundleSize:        oldMsgCreatePool.MaxBundleSize,
					Version:              oldMsgCreatePool.Version,
					Binaries:             oldMsgCreatePool.Binaries,
					StorageProviderId:    oldMsgCreatePool.StorageProviderId,
					CompressionId:        oldMsgCreatePool.CompressionId,
					EndKey:               "",
				}

				newMsgCreatePoolBytes, err := newMsgCreatePool.Marshal()
				if err != nil {
					logger.Error(fmt.Sprintf("could not marshal migrated gov message (proposal=%d, message_idx=%d)", proposal.Id, idx))
					continue
				}

				proposal.Messages[idx].Value = newMsgCreatePoolBytes
			}

			migratedMessagesCounter += 1
		}

		// Update Proposal Metadata
		if proposalData, ok := govProposalsData[proposal.Id]; ok {
			proposal.Title = proposalData.title
			proposal.Summary = proposalData.summary
			proposal.Proposer = proposalData.proposer
		}

		newProposalBytes, err := proposal.Marshal()
		if err != nil {
			logger.Error(fmt.Sprintf("could not marshal migrated gov proposal (proposal=%d)", proposal.Id))
			continue
		}

		proposalStore.Set(proposalIterator.Key(), newProposalBytes)
	}

	logger.Info(fmt.Sprintf("migrated MsgCreatePool messages in all gov proposals (message_count=%d)", migratedMessagesCounter))
}

func migrateFundersModule(sdkCtx sdk.Context, cdc codec.Codec, fundersStoreKey storetypes.StoreKey, fundersKeeper fundersKeeper.Keeper) {
	// migrate params
	oldParams := v1_4_funders.GetParams(sdkCtx, cdc, fundersStoreKey)

	newParams := fundersTypes.Params{
		CoinWhitelist: []*fundersTypes.WhitelistCoinEntry{
			// Prices were obtained on 03.07.2024

			// KYVE
			{
				CoinDenom:                 globalTypes.Denom,
				CoinDecimals:              uint32(6),
				MinFundingAmount:          math.NewIntFromUint64(oldParams.MinFundingAmount),
				MinFundingAmountPerBundle: math.NewIntFromUint64(oldParams.MinFundingAmountPerBundle),
				CoinWeight:                math.LegacyMustNewDecFromStr("0.0358"),
			},
			// Andromeda
			{
				CoinDenom:                 "ibc/58EDC95E791161D711F4CF012ACF30A5DA8DDEB40A484F293A52B1968903F643",
				CoinDecimals:              uint32(6),
				MinFundingAmount:          math.NewInt(1000_000_000),
				MinFundingAmountPerBundle: math.NewInt(100_000),
				CoinWeight:                math.LegacyMustNewDecFromStr("0.1007"),
			},
			// Source Protocol
			{
				CoinDenom:                 "ibc/0D2ABDF58A5DBA3D2A90398F8737D16ECAC0DDE58F9792B2918495D499400672",
				CoinDecimals:              uint32(6),
				MinFundingAmount:          math.NewInt(1000_000_000),
				MinFundingAmountPerBundle: math.NewInt(100_000),
				CoinWeight:                math.LegacyMustNewDecFromStr("0.0207"),
			},
		},
		MinFundingMultiple: oldParams.MinFundingMultiple,
	}

	fundersKeeper.SetParams(sdkCtx, newParams)

	_ = sdkCtx.EventManager().EmitTypedEvent(&fundersTypes.EventUpdateParams{
		OldParams: fundersTypes.Params{},
		NewParams: newParams,
		Payload:   "{}",
	})

	// migrate fundings
	oldFundings := v1_4_funders.GetAllFundings(sdkCtx, cdc, fundersStoreKey)
	for _, f := range oldFundings {
		fundersKeeper.SetFunding(sdkCtx, &fundersTypes.Funding{
			FunderAddress:    f.FunderAddress,
			PoolId:           f.PoolId,
			Amounts:          sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.Amount))),
			AmountsPerBundle: sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.AmountPerBundle))),
			TotalFunded:      sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.TotalFunded))),
		})
	}

	logger.Info("migrated Funders module")
}

func migrateDelegationModule(sdkCtx sdk.Context, cdc codec.Codec, delegationStoreKey storetypes.StoreKey, delegationKeeper delegationKeeper.Keeper) {
	// migrate delegation entries
	oldDelegationEntries := v1_4_delegation.GetAllDelegationEntries(sdkCtx, cdc, delegationStoreKey)
	for _, d := range oldDelegationEntries {
		delegationKeeper.SetDelegationEntry(sdkCtx, delegationTypes.DelegationEntry{
			Staker: d.Staker,
			KIndex: d.KIndex,
			Value:  sdk.NewDecCoins(sdk.NewDecCoinFromDec(globalTypes.Denom, d.Value)),
		})
	}

	// migrate delegation data
	oldDelegationData := v1_4_delegation.GetAllDelegationData(sdkCtx, cdc, delegationStoreKey)
	for _, d := range oldDelegationData {
		delegationKeeper.SetDelegationData(sdkCtx, delegationTypes.DelegationData{
			Staker:                     d.Staker,
			CurrentRewards:             sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(d.CurrentRewards))),
			TotalDelegation:            d.TotalDelegation,
			LatestIndexK:               d.LatestIndexK,
			DelegatorCount:             d.DelegatorCount,
			LatestIndexWasUndelegation: d.LatestIndexWasUndelegation,
		})
	}

	logger.Info("migrated Delegation module")
}

func migrateStakersModule(sdkCtx sdk.Context, cdc codec.Codec, stakersStoreKey storetypes.StoreKey, stakersKeeper *stakersKeeper.Keeper) {
	// migrate stakers
	oldStakers := v1_4_stakers.GetAllStakers(sdkCtx, cdc, stakersStoreKey)
	for _, s := range oldStakers {
		stakersKeeper.Migration_SetStaker(sdkCtx, stakersTypes.Staker{
			Address:           s.Address,
			Commission:        s.Commission,
			Moniker:           s.Moniker,
			Website:           s.Website,
			Identity:          s.Identity,
			SecurityContact:   s.SecurityContact,
			Details:           s.Details,
			CommissionRewards: sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(s.CommissionRewards))),
		})
	}

	logger.Info("migrated Stakers module")
}

func migrateBundlesModule(sdkCtx sdk.Context, cdc codec.Codec, bundlesStoreKey storetypes.StoreKey, bundlesKeeper bundlesKeeper.Keeper) {
	oldParams := v1_4_bundles.GetParams(sdkCtx, cdc, bundlesStoreKey)

	// TODO: define final storage cost prices
	newParams := bundlesTypes.Params{
		UploadTimeout: oldParams.UploadTimeout,
		StorageCosts: []bundlesTypes.StorageCost{
			// Arweave: https://arweave.net/price/1048576 -> 699 winston/byte * 40 USD/AR * 1.5 / 10**12
			{StorageProviderId: 1, Cost: math.LegacyMustNewDecFromStr("0.000000004194")},
			// Irys: https://node1.bundlr.network/price/1048576 -> 1048 winston/byte * 40 USD/AR * 1.5 / 10**12
			{StorageProviderId: 2, Cost: math.LegacyMustNewDecFromStr("0.000000006288")},
			// KYVE Storage Provider: zero since it is free to use for testnet participants
			{StorageProviderId: 3, Cost: math.LegacyMustNewDecFromStr("0.00")},
		},
		NetworkFee: oldParams.NetworkFee,
		MaxPoints:  oldParams.MaxPoints,
	}

	bundlesKeeper.SetParams(sdkCtx, newParams)

	_ = sdkCtx.EventManager().EmitTypedEvent(&bundlesTypes.EventUpdateParams{
		OldParams: bundlesTypes.Params{},
		NewParams: newParams,
		Payload:   "{}",
	})

	logger.Info("migrated Bundles module")
}

func migratePoolModule(sdkCtx sdk.Context, cdc codec.Codec, poolStoreKey storetypes.StoreKey, poolKeeper *poolKeeper.Keeper) {
	oldParams := v1_4_pool.GetParams(sdkCtx, cdc, poolStoreKey)

	poolKeeper.SetParams(sdkCtx, poolTypes.Params{
		ProtocolInflationShare:  oldParams.ProtocolInflationShare,
		PoolInflationPayoutRate: oldParams.PoolInflationPayoutRate,
		MaxVotingPowerPerPool:   math.LegacyMustNewDecFromStr("0.5"),
	})

	pools := v1_4_pool.GetAllPools(sdkCtx, poolStoreKey, cdc)
	for _, pool := range pools {
		var newPool poolTypes.Pool

		var protocol *poolTypes.Protocol
		if pool.Protocol != nil {
			protocol = &poolTypes.Protocol{
				Version:     pool.Protocol.Version,
				Binaries:    pool.Protocol.Binaries,
				LastUpgrade: pool.Protocol.LastUpgrade,
			}
		}
		var upgradePlan *poolTypes.UpgradePlan
		if pool.UpgradePlan != nil {
			upgradePlan = &poolTypes.UpgradePlan{
				Version:     pool.UpgradePlan.Version,
				Binaries:    pool.UpgradePlan.Binaries,
				ScheduledAt: pool.UpgradePlan.ScheduledAt,
				Duration:    pool.UpgradePlan.Duration,
			}
		}

		newPool = poolTypes.Pool{
			Id:             pool.Id,
			Name:           pool.Name,
			Runtime:        pool.Runtime,
			Logo:           pool.Logo,
			Config:         pool.Config,
			StartKey:       pool.StartKey,
			CurrentKey:     pool.CurrentKey,
			CurrentSummary: pool.CurrentSummary,
			CurrentIndex:   pool.CurrentIndex,
			TotalBundles:   pool.TotalBundles,
			UploadInterval: pool.UploadInterval,
			// Currently all pools have int64(1_000_000) as inflation_share weight.
			// Set this to 1 as we now support decimals
			InflationShareWeight:     math.LegacyMustNewDecFromStr("1"),
			MinDelegation:            pool.MinDelegation,
			MaxBundleSize:            pool.MaxBundleSize,
			Disabled:                 pool.Disabled,
			Protocol:                 protocol,
			UpgradePlan:              upgradePlan,
			CurrentStorageProviderId: pool.CurrentStorageProviderId,
			CurrentCompressionId:     pool.CurrentCompressionId,
			EndKey:                   pool.EndKey,
		}
		poolKeeper.SetPool(sdkCtx, newPool)

		_ = sdkCtx.EventManager().EmitTypedEvent(&poolTypes.EventPoolUpdated{
			Id:                   pool.Id,
			RawUpdateString:      "{\"inflation_share_weight\":\"1.0\"}",
			Name:                 pool.Name,
			Runtime:              pool.Runtime,
			Logo:                 pool.Logo,
			Config:               pool.Config,
			UploadInterval:       pool.UploadInterval,
			InflationShareWeight: math.LegacyMustNewDecFromStr("1"),
			MinDelegation:        pool.MinDelegation,
			MaxBundleSize:        pool.MaxBundleSize,
			StorageProviderId:    pool.CurrentStorageProviderId,
			CompressionId:        pool.CurrentCompressionId,
		})
	}

	logger.Info("migrated Pool module")
}
