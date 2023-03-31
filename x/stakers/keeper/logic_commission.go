package keeper

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// orderNewCommissionChange inserts a new change entry into the queue.
// The queue is checked in every endBlock and when the commissionChangeTime
// is over the new commission will be applied to the user.
// If another entry is currently in the queue it will be removed.
func (k Keeper) orderNewCommissionChange(ctx sdk.Context, staker string, commission sdk.Dec) {
	// Remove existing queue entry
	queueEntry, found := k.GetCommissionChangeEntryByIndex2(ctx, staker)
	if found {
		k.RemoveCommissionChangeEntry(ctx, &queueEntry)
	}

	queueIndex := k.getNextQueueSlot(ctx, types.QUEUE_IDENTIFIER_COMMISSION)

	commissionChangeEntry := types.CommissionChangeEntry{
		Index:        queueIndex,
		Staker:       staker,
		Commission:   commission,
		CreationDate: ctx.BlockTime().Unix(),
	}

	k.SetCommissionChangeEntry(ctx, commissionChangeEntry)
}

// ProcessCommissionChangeQueue checks the queue for entries which are due
// and can be executed. If this is the case, the new commission
// will be applied to the staker
func (k Keeper) ProcessCommissionChangeQueue(ctx sdk.Context) {
	k.processQueue(ctx, types.QUEUE_IDENTIFIER_COMMISSION, func(index uint64) bool {
		// Get queue entry in question
		queueEntry, found := k.GetCommissionChangeEntry(ctx, index)

		if !found {
			// continue with the next entry
			return true
		} else if queueEntry.CreationDate+int64(k.GetCommissionChangeTime(ctx)) <= ctx.BlockTime().Unix() {

			k.RemoveCommissionChangeEntry(ctx, &queueEntry)

			k.UpdateStakerCommission(ctx, queueEntry.Staker, queueEntry.Commission)

			_ = ctx.EventManager().EmitTypedEvent(&types.EventUpdateCommission{
				Staker:     queueEntry.Staker,
				Commission: queueEntry.Commission,
			})

			// Continue with next entry
			return true
		}

		// Stop queue processing
		return false
	})
}
