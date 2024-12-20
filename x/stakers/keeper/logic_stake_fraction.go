package keeper

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// orderNewStakeFractionChange inserts a new change entry into the queue.
// The queue is checked in every endBlock and when the stakeFractionChangeTime
// is over the new stake fraction will be applied to the user.
// If another entry is currently in the queue it will be removed.
func (k Keeper) orderNewStakeFractionChange(ctx sdk.Context, staker string, poolId uint64, stakeFraction math.LegacyDec) {
	// Remove existing queue entry
	queueEntry, found := k.GetStakeFractionChangeEntryByIndex2(ctx, staker, poolId)
	if found {
		k.RemoveStakeFractionEntry(ctx, &queueEntry)
	}

	queueIndex := k.getNextQueueSlot(ctx, types.QUEUE_IDENTIFIER_STAKE_FRACTION)

	stakeFractionChangeEntry := types.StakeFractionChangeEntry{
		Index:         queueIndex,
		Staker:        staker,
		PoolId:        poolId,
		StakeFraction: stakeFraction,
		CreationDate:  ctx.BlockTime().Unix(),
	}

	k.SetStakeFractionChangeEntry(ctx, stakeFractionChangeEntry)
}

// ProcessStakeFractionChangeQueue checks the queue for entries which are due
// and can be executed. If this is the case, the new stake fraction
// will be applied to the staker and the pool
func (k Keeper) ProcessStakeFractionChangeQueue(ctx sdk.Context) {
	k.processQueue(ctx, types.QUEUE_IDENTIFIER_STAKE_FRACTION, func(index uint64) bool {
		// Get queue entry in question
		queueEntry, found := k.GetStakeFractionChangeEntry(ctx, index)
		if !found {
			// continue with the next entry
			return true
		}

		if queueEntry.CreationDate+int64(k.GetStakeFractionChangeTime(ctx)) <= ctx.BlockTime().Unix() {
			k.RemoveStakeFractionEntry(ctx, &queueEntry)

			valaccount, valaccountFound := k.GetValaccount(ctx, queueEntry.PoolId, queueEntry.Staker)
			if !valaccountFound {
				// continue with the next entry
				return true
			}

			valaccount.StakeFraction = queueEntry.StakeFraction
			k.SetValaccount(ctx, valaccount)

			_ = ctx.EventManager().EmitTypedEvent(&types.EventUpdateStakeFraction{
				Staker:        queueEntry.Staker,
				PoolId:        queueEntry.PoolId,
				StakeFraction: queueEntry.StakeFraction,
			})

			// Continue with next entry
			return true
		}

		// Stop queue processing
		return false
	})
}
