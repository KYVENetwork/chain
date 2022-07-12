package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) orderNewCommissionChange(ctx sdk.Context, poolId uint64, staker string, commission string) (error error) {

	// unbondingState stores the start and the end of the queue with all unbonding entries
	// the queue is ordered by time
	commissionChangeState := k.GetCommissionChangeQueueState(ctx)

	// Increase topIndex as a new entry is about to be appended
	commissionChangeState.HighIndex += 1
	k.SetCommissionChangeQueueState(ctx, commissionChangeState)

	// Remove existing queue entry
	queueEntry, found := k.GetCommissionChangeQueueEntryByIndex2(ctx, staker, poolId)
	if found {
		k.RemoveCommissionChangeQueueEntry(ctx, &queueEntry)
	}

	// ...
	commissionChangeEntry := types.CommissionChangeQueueEntry{
		Index:        commissionChangeState.HighIndex,
		Staker:       staker,
		PoolId:       poolId,
		Commission:   commission,
		CreationDate: ctx.BlockTime().Unix(),
	}

	k.SetCommissionChangeQueueEntry(ctx, commissionChangeEntry)

	return nil
}

// ProcessCommissionChangeUnbondingQueue ...
func (k Keeper) ProcessCommissionChangeUnbondingQueue(ctx sdk.Context) {

	// Get Queue information
	queueState := k.GetCommissionChangeQueueState(ctx)

	// flag for computing every entry at the end of the queue which is due.
	// start processing the end of the queue
	for commissionChangePerformed := true; commissionChangePerformed; {
		commissionChangePerformed = false

		// Get end of queue
		queueEntry, found := k.GetCommissionChangeQueueEntry(ctx, queueState.LowIndex+1)

		if !found {
			// If there are still entries in the queue, continue with processing
			if queueState.LowIndex < queueState.HighIndex {
				queueState.LowIndex += 1
				commissionChangePerformed = true
			}
		} else if queueEntry.CreationDate+int64(k.CommissionChangeTime(ctx)) <= ctx.BlockTime().Unix() {

			queueState.LowIndex += 1
			commissionChangePerformed = true

			k.RemoveCommissionChangeQueueEntry(ctx, &queueEntry)

			staker, stakerFound := k.GetStaker(ctx, queueEntry.Staker, queueEntry.PoolId)
			if stakerFound {
				staker.Commission = queueEntry.Commission
				k.SetStaker(ctx, staker)
			}

			// Event an event.
			ctx.EventManager().EmitTypedEvent(&types.EventUpdateCommission{
				PoolId:     queueEntry.PoolId,
				Address:    queueEntry.Staker,
				Commission: queueEntry.Commission,
			})
		}

	}
	k.SetCommissionChangeQueueState(ctx, queueState)
}
