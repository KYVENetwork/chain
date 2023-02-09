package keeper

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) orderLeavePool(ctx sdk.Context, staker string, poolId uint64) error {
	// Remove existing queue entry
	if k.DoesLeavePoolEntryExistByIndex2(ctx, staker, poolId) {
		return sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrPoolLeaveAlreadyInProgress.Error())
	}

	queueIndex := k.getNextQueueSlot(ctx, types.QUEUE_IDENTIFIER_LEAVE)

	leavePoolEntry := types.LeavePoolEntry{
		Index:        queueIndex,
		Staker:       staker,
		PoolId:       poolId,
		CreationDate: ctx.BlockTime().Unix(),
	}

	k.SetLeavePoolEntry(ctx, leavePoolEntry)

	return nil
}

// ProcessLeavePoolQueue ...
func (k Keeper) ProcessLeavePoolQueue(ctx sdk.Context) {
	k.processQueue(ctx, types.QUEUE_IDENTIFIER_LEAVE, func(index uint64) bool {
		// Get queue entry in question
		queueEntry, found := k.GetLeavePoolEntry(ctx, index)

		if !found {
			// continue with the next entry
			return true
		} else if queueEntry.CreationDate+int64(k.GetLeavePoolTime(ctx)) <= ctx.BlockTime().Unix() {

			k.RemoveLeavePoolEntry(ctx, &queueEntry)
			k.LeavePool(ctx, queueEntry.Staker, queueEntry.PoolId)

			return true
		}
		return false
	})
}
