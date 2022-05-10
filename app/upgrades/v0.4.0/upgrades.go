package v0_4_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
)

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		for _, pool := range registryKeeper.GetAllPool(ctx) {
			pool.MaxBundleSize = 100
			registryKeeper.SetPool(ctx, pool)
		}

		// Return.
		return vm, nil
	}
}
