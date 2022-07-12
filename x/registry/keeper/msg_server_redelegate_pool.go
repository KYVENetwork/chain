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

	// Check if cooldowns are over
	for _, block := range k.GetRedelegationCooldownEntries(ctx, msg.Creator) {
		creationTime := ctx.WithBlockHeight(int64(block)).BlockTime()
		if ctx.BlockTime().Unix()-creationTime.Unix() > int64(k.RedelegationCooldown(ctx)) {
			k.RemoveRedelegationCooldown(ctx, msg.Creator, block)
		} else {
			break
		}
	}

	blocks := k.GetRedelegationCooldownEntries(ctx, msg.Creator)

	if len(blocks) >= int(k.RedelegationMaxAmount(ctx)) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrRedelegationOnCooldown.Error())
	}
	if blocks[len(blocks)-1] == uint64(ctx.BlockHeight()) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrMultipleRedelegationInSameBlock.Error())
	}

	// All checks passed, create cooldown entry
	k.SetRedelegationCooldown(ctx, types.RedelegationCooldown{
		Address:      msg.Creator,
		CreatedBlock: uint64(ctx.BlockHeight())})

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
