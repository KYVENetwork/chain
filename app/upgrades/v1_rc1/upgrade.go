package v1rc1

import (
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Bundles
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bundlesKeeper bundlesKeeper.Keeper,
	bundlesStoreKey storeTypes.StoreKey,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// TODO(@john): Do we need this check?
		if ctx.ChainID() != "kaon-1" {
			return vm, nil
		}

		ReinitialiseBundlesParams(ctx, bundlesKeeper, bundlesStoreKey)
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// ReinitialiseBundlesParams ...
func ReinitialiseBundlesParams(
	ctx sdk.Context,
	bundlesKeeper bundlesKeeper.Keeper,
	_ storeTypes.StoreKey,
) {
	// TODO(@john): Do we need to delete the old store?

	// store := ctx.KVStore(bundlesStoreKey)
	// store.Delete(bundlesTypes.ParamsKey)

	params := bundlesTypes.DefaultParams()
	bundlesKeeper.SetParams(ctx, params)
}
