package v1_5

import (
	"context"
	"fmt"

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
		logger := sdk.UnwrapSDKContext(ctx).Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		migrateStorageCosts(ctx, bundlesKeeper, poolKeeper, storeKeys, cdc)

		// TODO: migrate gov params

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func migrateStorageCosts(ctx context.Context, bundlesKeeper keeper.Keeper, poolKeeper *poolkeeper.Keeper, storeKeys []storetypes.StoreKey, cdc codec.Codec) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var bundlesStoreKey storetypes.StoreKey
	for _, k := range storeKeys {
		if k.Name() == "bundles" {
			bundlesStoreKey = k
			break
		}
	}
	if bundlesStoreKey == nil {
		panic("bundles store key not found")
	}

	// Get all storage providers
	storageIds := map[uint32]struct{}{}
	for _, pool := range poolKeeper.GetAllPools(sdkCtx) {
		storageIds[pool.CurrentStorageProviderId] = struct{}{}
	}

	// Copy storage cost from old params to new params
	// The storage cost of all storage providers will be the same after this migration
	oldParams := v1_4_types.GetParams(sdkCtx, bundlesStoreKey, cdc)
	newParams := bundlestypes.Params{
		UploadTimeout: oldParams.UploadTimeout,
		StorageCosts:  make(map[uint32]bundlestypes.LegacyDecValue),
		NetworkFee:    oldParams.NetworkFee,
		MaxPoints:     oldParams.MaxPoints,
	}
	for storageId := range storageIds {
		newParams.StorageCosts[storageId] = bundlestypes.LegacyDecValue{
			Value: oldParams.StorageCost,
		}
	}

	bundlesKeeper.SetParams(sdkCtx, newParams)
}
