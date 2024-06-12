package v1_5

import (
	"context"
	"fmt"

	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"

	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/bundles"
	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/delegation"
	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/funders"
	v1_4_pool "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/pool"
	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/stakers"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"
	funderTypes "github.com/KYVENetwork/chain/x/funders/types"
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

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, cdc codec.Codec, storeKeys []storetypes.StoreKey, bundlesKeeper bundlesKeeper.Keeper, delegationKeeper delegationKeeper.Keeper, fundersKeeper fundersKeeper.Keeper, stakersKeeper *stakersKeeper.Keeper, poolKeeper *poolKeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// TODO: migrate gov params
		// TODO emit events if necessary

		// migrate fundings
		MigrateFundersModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, fundersTypes.StoreKey), fundersKeeper)

		// migrate delegations
		migrateDelegationModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, delegationTypes.StoreKey), delegationKeeper)

		// migrate stakers
		migrateStakersModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, stakersTypes.StoreKey), stakersKeeper)

		// migrate bundles
		migrateBundlesModule(sdkCtx, cdc, MustGetStoreKey(storeKeys, bundlesTypes.StoreKey), bundlesKeeper)

		// migrate pool
		migrateMaxVotingPowerInPool(sdkCtx, cdc, MustGetStoreKey(storeKeys, poolTypes.StoreKey), *poolKeeper)
		migrateInflationShareWeight(sdkCtx, cdc, MustGetStoreKey(storeKeys, poolTypes.StoreKey), poolKeeper)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func MustGetStoreKey(storeKeys []storetypes.StoreKey, storeName string) storetypes.StoreKey {
	for _, k := range storeKeys {
		if k.Name() == storeName {
			return k
		}
	}

	panic(fmt.Sprintf("failed to find store key: %s", fundersTypes.StoreKey))
}

func MigrateFundersModule(sdkCtx sdk.Context, cdc codec.Codec, fundersStoreKey storetypes.StoreKey, fundersKeeper fundersKeeper.Keeper) {
	// migrate params
	// TODO: define final prices and initial whitelisted coins
	oldParams := funders.GetParams(sdkCtx, cdc, fundersStoreKey)
	fundersKeeper.SetParams(sdkCtx, fundersTypes.Params{
		CoinWhitelist: []*fundersTypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				CoinDecimals:              uint32(6),
				MinFundingAmount:          math.NewIntFromUint64(oldParams.MinFundingAmount),
				MinFundingAmountPerBundle: math.NewIntFromUint64(oldParams.MinFundingAmountPerBundle),
				CoinWeight:                math.LegacyMustNewDecFromStr("0.06"),
			},
		},
		MinFundingMultiple: oldParams.MinFundingMultiple,
	})

	// migrate fundings
	oldFundings := funders.GetAllFundings(sdkCtx, cdc, fundersStoreKey)
	for _, f := range oldFundings {
		fundersKeeper.SetFunding(sdkCtx, &funderTypes.Funding{
			FunderAddress:    f.FunderAddress,
			PoolId:           f.PoolId,
			Amounts:          sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.Amount))),
			AmountsPerBundle: sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.AmountPerBundle))),
			TotalFunded:      sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.TotalFunded))),
		})
	}
}

func migrateDelegationModule(sdkCtx sdk.Context, cdc codec.Codec, delegationStoreKey storetypes.StoreKey, delegationKeeper delegationKeeper.Keeper) {
	// migrate delegation entries
	oldDelegationEntries := delegation.GetAllDelegationEntries(sdkCtx, cdc, delegationStoreKey)
	for _, d := range oldDelegationEntries {
		delegationKeeper.SetDelegationEntry(sdkCtx, delegationTypes.DelegationEntry{
			Staker: d.Staker,
			KIndex: d.KIndex,
			Value:  sdk.NewDecCoins(sdk.NewDecCoinFromDec(globalTypes.Denom, d.Value)),
		})
	}

	// migrate delegation data
	oldDelegationData := delegation.GetAllDelegationData(sdkCtx, cdc, delegationStoreKey)
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
}

func migrateStakersModule(sdkCtx sdk.Context, cdc codec.Codec, stakersStoreKey storetypes.StoreKey, stakersKeeper *stakersKeeper.Keeper) {
	// migrate stakers
	oldStakers := stakers.GetAllStakers(sdkCtx, cdc, stakersStoreKey)
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
}

func migrateBundlesModule(sdkCtx sdk.Context, cdc codec.Codec, bundlesStoreKey storetypes.StoreKey, bundlesKeeper bundlesKeeper.Keeper) {
	oldParams := bundles.GetParams(sdkCtx, cdc, bundlesStoreKey)

	// TODO: define final storage cost prices
	bundlesKeeper.SetParams(sdkCtx, bundlesTypes.Params{
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
	})
}

func migrateMaxVotingPowerInPool(sdkCtx sdk.Context, cdc codec.Codec, poolStoreKey storetypes.StoreKey, poolKeeper poolKeeper.Keeper) {
	oldParams := v1_4_pool.GetParams(sdkCtx, cdc, poolStoreKey)

	poolKeeper.SetParams(sdkCtx, poolTypes.Params{
		ProtocolInflationShare:  oldParams.ProtocolInflationShare,
		PoolInflationPayoutRate: oldParams.PoolInflationPayoutRate,
		MaxVotingPowerPerPool:   math.LegacyMustNewDecFromStr("0.5"),
	})
}

func migrateInflationShareWeight(sdkCtx sdk.Context, cdc codec.Codec, poolStoreKey storetypes.StoreKey, poolKeeper *poolKeeper.Keeper) {
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
}
