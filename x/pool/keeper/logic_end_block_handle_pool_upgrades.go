package keeper

import (
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandlePoolUpgrades handles to end-block logic for detecting and performing schedules pool-upgrades.
func (k Keeper) HandlePoolUpgrades(ctx sdk.Context) {
	for _, pool := range k.GetAllPools(ctx) {
		// PoolUpgrade is scheduled if `ScheduledAt` is not zero and smaller than the current block-time
		if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {

			// Check if pool upgrade already has been applied
			if pool.Protocol.Version != pool.UpgradePlan.Version || pool.Protocol.Binaries != pool.UpgradePlan.Binaries {
				// perform pool upgrade
				pool.Protocol.Version = pool.UpgradePlan.Version
				pool.Protocol.Binaries = pool.UpgradePlan.Binaries
				pool.Protocol.LastUpgrade = pool.UpgradePlan.ScheduledAt
			}

			// Check if upgrade duration was reached
			if uint64(ctx.BlockTime().Unix()) >= (pool.UpgradePlan.ScheduledAt + pool.UpgradePlan.Duration) {
				// reset upgrade plan to default values
				pool.UpgradePlan = &types.UpgradePlan{}
			}

			k.SetPool(ctx, pool)
		}
	}
}
