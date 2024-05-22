package v1_5

import (
	"context"
	"fmt"
	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_pool_types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"

	"cosmossdk.io/math"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types"
	"github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
	poolkeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v1.5.0"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, cdc codec.Codec, storeKeys []storetypes.StoreKey, bundlesKeeper keeper.Keeper, poolKeeper *poolkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		if err := migrateStorageCosts(sdkCtx, bundlesKeeper, poolKeeper, storeKeys, cdc); err != nil {
			return nil, err
		}

		if err := migrateInflationShareWeight(sdkCtx, poolKeeper, storeKeys, cdc); err != nil {
			return nil, err
		}

		// TODO: migrate gov params

		// TODO: migrate fundings

		// TODO: migrate delegation outstanding rewards

		// TODO: migrate network fee and whitelist weights

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func migrateStorageCosts(sdkCtx sdk.Context, bundlesKeeper keeper.Keeper, poolKeeper *poolkeeper.Keeper, storeKeys []storetypes.StoreKey, cdc codec.Codec) error {
	var bundlesStoreKey storetypes.StoreKey
	for _, k := range storeKeys {
		if k.Name() == "bundles" {
			bundlesStoreKey = k
			break
		}
	}
	if bundlesStoreKey == nil {
		return fmt.Errorf("store key not found: bundles")
	}

	// Copy storage cost from old params to new params
	// The storage cost of all storage providers will be the same after this migration
	oldParams := v1_4_types.GetParams(sdkCtx, bundlesStoreKey, cdc)
	newParams := bundlestypes.Params{
		UploadTimeout: oldParams.UploadTimeout,
		StorageCosts: []bundlestypes.StorageCost{
			// TODO: define value for storage provider id 1 and 2
			{StorageProviderId: 1, Cost: math.LegacyMustNewDecFromStr("0.00")},
			{StorageProviderId: 2, Cost: math.LegacyMustNewDecFromStr("0.00")},
		},
		NetworkFee: oldParams.NetworkFee,
		MaxPoints:  oldParams.MaxPoints,
	}

	bundlesKeeper.SetParams(sdkCtx, newParams)
	return nil
}

func migrateInflationShareWeight(sdkCtx sdk.Context, poolKeeper *poolkeeper.Keeper, storeKeys []storetypes.StoreKey, cdc codec.Codec) error {
	var poolStoreKey storetypes.StoreKey
	for _, k := range storeKeys {
		if k.Name() == "pool" {
			poolStoreKey = k
			break
		}
	}
	if poolStoreKey == nil {
		return fmt.Errorf("store key not found: pool")
	}

	pools := v1_4_pool_types.GetAllPools(sdkCtx, poolStoreKey, cdc)
	for _, pool := range pools {
		var newPool pooltypes.Pool

		var protocol *pooltypes.Protocol
		if pool.Protocol != nil {
			protocol = &pooltypes.Protocol{
				Version:     pool.Protocol.Version,
				Binaries:    pool.Protocol.Binaries,
				LastUpgrade: pool.Protocol.LastUpgrade,
			}
		}
		var upgradePlan *pooltypes.UpgradePlan
		if pool.UpgradePlan != nil {
			upgradePlan = &pooltypes.UpgradePlan{
				Version:     pool.UpgradePlan.Version,
				Binaries:    pool.UpgradePlan.Binaries,
				ScheduledAt: pool.UpgradePlan.ScheduledAt,
				Duration:    pool.UpgradePlan.Duration,
			}
		}

		newPool = pooltypes.Pool{
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
			// Convert inflation share weight to new decimal type
			InflationShareWeight:     math.LegacyNewDec(int64(pool.InflationShareWeight)),
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
	}
	return nil
}
