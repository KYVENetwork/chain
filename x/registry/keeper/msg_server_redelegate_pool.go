package keeper

import (
	"context"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RedelegatePool lets a user redelegate from one staker to another staker
func (k msgServer) RedelegatePool(goCtx context.Context, msg *types.MsgRedelegatePool) (*types.MsgRedelegatePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if cooldowns are over,
	// Remove all expired entries
	for _, creationDate := range k.GetRedelegationCooldownEntries(ctx, msg.Creator) {
		if ctx.BlockTime().Unix()-int64(creationDate) > int64(k.RedelegationCooldown(ctx)) {
			k.RemoveRedelegationCooldown(ctx, msg.Creator, creationDate)
		} else {
			break
		}
	}

	// Get list of active cooldowns
	creationDates := k.GetRedelegationCooldownEntries(ctx, msg.Creator)

	// Check if there are still free blocks
	if len(creationDates) >= int(k.RedelegationMaxAmount(ctx)) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrRedelegationOnCooldown.Error())
	}
	// Check that now Redelegation occurred in this block, as it will lead to errors, as
	// the block-time is used for an index key.
	if len(creationDates) > 0 && creationDates[len(creationDates)-1] == uint64(ctx.BlockTime().Unix()) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrMultipleRedelegationInSameBlock.Error())
	}

	// All checks passed, create cooldown entry
	k.SetRedelegationCooldown(ctx, types.RedelegationCooldown{
		Address:      msg.Creator,
		CreationDate: uint64(ctx.BlockTime().Unix())})

	// Perform undelegation
	if err := k.Undelegate(ctx, msg.FromStaker, msg.FromPoolId, msg.Creator, msg.Amount); err != nil {
		return nil, err
	}

	// Perform undelegation
	if err := k.Delegate(ctx, msg.ToStaker, msg.ToPoolId, msg.Creator, msg.Amount); err != nil {
		return nil, err
	}

	// Emit a delegation event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventRedelegatePool{
		Address:  msg.Creator,
		FromPool: msg.FromPoolId,
		FromNode: msg.FromStaker,
		ToPool:   msg.ToPoolId,
		ToNode:   msg.ToStaker,
		Amount:   msg.Amount,
	}); errEmit != nil {
		return nil, errEmit
	}

	return &types.MsgRedelegatePoolResponse{}, nil
}
