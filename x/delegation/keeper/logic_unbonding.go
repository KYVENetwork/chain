package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StartUnbondingDelegator creates a queue entry to schedule the unbonding.
// After the DelegationTime is reached the actual unbonding will be performed
// The actual unbonding is then performed by `func ProcessDelegatorUnbondingQueue(...)`
func (k Keeper) StartUnbondingDelegator(ctx sdk.Context, staker string, delegatorAddress string, amount uint64) {
	// the queue is ordered by time
	queueState := k.GetQueueState(ctx)

	// Increase topIndex as a new entry is about to be appended
	queueState.HighIndex += 1

	k.SetQueueState(ctx, queueState)

	// UnbondingEntry stores all the information which are needed to perform
	// the undelegation at the end of the unbonding time
	undelegationQueueEntry := types.UndelegationQueueEntry{
		Delegator:    delegatorAddress,
		Index:        queueState.HighIndex,
		Staker:       staker,
		Amount:       amount,
		CreationTime: uint64(ctx.BlockTime().Unix()),
	}

	k.SetUndelegationQueueEntry(ctx, undelegationQueueEntry)
}

// ProcessDelegatorUnbondingQueue is called in the end block and
// checks the queue for entries that have surpassed the unbonding time.
// If the unbonding time is reached, the actual unbonding is performed
// and the entry is removed from the queue.
func (k Keeper) ProcessDelegatorUnbondingQueue(ctx sdk.Context) {
	// Get Queue information
	queueState := k.GetQueueState(ctx)

	// flag for computing every entry at the end of the queue which is due.
	// start processing the end of the queue
	for continueProcessing := true; continueProcessing; {
		continueProcessing = false

		// Get end of queue
		undelegationEntry, found := k.GetUndelegationQueueEntry(ctx, queueState.LowIndex+1)

		if !found {
			if queueState.LowIndex < queueState.HighIndex {
				queueState.LowIndex += 1
				continueProcessing = true
			}
		} else
		// Check if unbonding time is over
		if undelegationEntry.CreationTime+k.GetUnbondingDelegationTime(ctx) <= uint64(ctx.BlockTime().Unix()) {

			// Perform undelegation and save undelegated amount to then transfer back to the user
			undelegatedAmount := k.performUndelegation(ctx, undelegationEntry.Staker, undelegationEntry.Delegator, undelegationEntry.Amount)

			// Transfer the money
			if err := util.TransferFromModuleToAddress(
				k.bankKeeper,
				ctx,
				types.ModuleName,
				undelegationEntry.Delegator,
				undelegatedAmount,
			); err != nil {
				util.PanicHalt(k.upgradeKeeper, ctx, "Not enough money in delegation module - logic_unbonding")
			}

			// Emit a delegation event.
			_ = ctx.EventManager().EmitTypedEvent(&types.EventUndelegate{
				Address: undelegationEntry.Delegator,
				Staker:  undelegationEntry.Staker,
				Amount:  undelegatedAmount,
			})

			k.RemoveUndelegationQueueEntry(ctx, &undelegationEntry)

			continueProcessing = true
			queueState.LowIndex += 1
		}
	}
	k.SetQueueState(ctx, queueState)
}
