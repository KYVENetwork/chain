package v0_6_0

import (
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func migrateStakers(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, staker := range registryKeeper.GetAllStaker(ctx) {
		staker.Status = types.STAKER_STATUS_ACTIVE
		registryKeeper.GetStaker(ctx, staker.Account, staker.PoolId)
	}
}

func createRedelegationParameters(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	registryKeeper.ParamStore().Set(ctx, types.KeyRedelegationCooldown, types.DefaultRedelegationCooldown)

	registryKeeper.ParamStore().Set(ctx, types.KeyRedelegationMaxAmount, types.DefaultRedelegationMaxAmount)
}

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		registryKeeper.UpgradeHelperV060MigrateSecondIndex(ctx)

		migrateStakers(registryKeeper, ctx)

		createRedelegationParameters(registryKeeper, ctx)

		return vm, nil
	}
}
