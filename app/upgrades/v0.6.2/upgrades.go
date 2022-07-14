package v0_6_2

import (
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func fixPools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, pool := range registryKeeper.GetAllPool(ctx) {

		// Immediately unstake all stakers that exceed the list
		if len(pool.Stakers) > 50 {

			keepStakers := pool.Stakers[0:50]
			refundStakers := pool.Stakers[50:]

			pool.Stakers = keepStakers

			for _, stakerAddress := range refundStakers {
				staker, found := registryKeeper.GetStaker(ctx, stakerAddress, pool.Id)
				if found {
					transferError := registryKeeper.TransferToAddress(ctx, staker.Account, staker.Amount)
					if transferError != nil {
						registryKeeper.PanicHalt(ctx, "Not enough money in module: "+transferError.Error())
					}

					ctx.EventManager().EmitTypedEvent(&types.EventUnstakePool{
						PoolId:  pool.Id,
						Address: staker.Account,
						Amount:  staker.Amount,
					})

					registryKeeper.RemoveStaker(ctx, stakerAddress, 0)
				}
			}
		}

		// Fix unspecified state
		activeStakers := make([]string, 0)

		for _, stakerAddress := range pool.Stakers {
			staker, found := registryKeeper.GetStaker(ctx, stakerAddress, pool.Id)
			if found {
				staker.Status = types.STAKER_STATUS_ACTIVE
				registryKeeper.SetStaker(ctx, staker)
				activeStakers = append(activeStakers, stakerAddress)
			}
		}
		pool.Stakers = activeStakers

		// Update new total stake
		var totalStake uint64
		for _, stakerAddress := range pool.Stakers {
			staker, _ := registryKeeper.GetStaker(ctx, stakerAddress, pool.Id)
			totalStake += staker.Amount
		}
		pool.TotalStake = totalStake

		var totalInactiveStake uint64
		for _, inactiveStakerAddress := range pool.InactiveStakers {
			inactiveStaker, _ := registryKeeper.GetStaker(ctx, inactiveStakerAddress, pool.Id)
			totalInactiveStake += inactiveStaker.Amount
		}
		pool.TotalInactiveStake = totalInactiveStake

		registryKeeper.UpdateLowestStaker(ctx, &pool)
		registryKeeper.SetPool(ctx, pool)
	}
}

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		fixPools(registryKeeper, ctx)

		return vm, nil
	}
}
