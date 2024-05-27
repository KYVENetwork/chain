package v1_5

import (
	"context"
	"fmt"

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

		// TODO: migrate gov params

		// TODO: migrate fundings

		// TODO: migrate delegation outstanding rewards

		// TODO: migrate network fee and whitelist weights

		// TODO: migrate MaxVotingPowerPerPool in pool params

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
