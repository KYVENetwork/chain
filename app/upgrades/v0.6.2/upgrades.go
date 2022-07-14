package v0_6_2

import (
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func migrateStakers(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, staker := range registryKeeper.GetAllStaker(ctx) {
		if staker.Status == types.STAKER_STATUS_UNSPECIFIED {
			staker.Status = types.STAKER_STATUS_ACTIVE
			registryKeeper.SetStaker(ctx, staker)
		}
	}
}

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		migrateStakers(registryKeeper, ctx)

		return vm, nil
	}
}
