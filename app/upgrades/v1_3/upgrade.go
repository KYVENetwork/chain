package v1_3

import (
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// TODO UpdateBundlesVersionMap(bundlesKeeper, ctx)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func UpdateBundlesVersionMap(keeper bundlesKeeper.Keeper, ctx sdk.Context) {
	keeper.SetBundleVersionMap(ctx, bundlesTypes.BundleVersionMap{
		Versions: []*bundlesTypes.BundleVersionEntry{
			{
				Height:  uint64(ctx.BlockHeight()),
				Version: 2,
			},
		},
	})
}
