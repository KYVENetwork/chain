package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Undelegate handles the transaction of undelegating a given amount from the delegated tokens
// The Undelegation is not performed immediately, instead an unbonding entry is created and pushed
// to a queue. When the unbonding timeout is reached the actual undelegation is performed.
// If the delegator got slashed during the unbonding only the remaining tokens will be returned.
func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Do not allow to undelegate more than currently delegated
	if delegationAmount := k.GetDelegationAmountOfDelegator(ctx, msg.Staker, msg.Creator); msg.Amount > delegationAmount {
		return nil, types.ErrNotEnoughDelegation.Wrapf("%d > %d", msg.Amount, delegationAmount)
	}

	// Create and insert unbonding queue entry.
	k.StartUnbondingDelegator(ctx, msg.Staker, msg.Creator, msg.Amount)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventStartUndelegation{
		Address:                   msg.Creator,
		Staker:                    msg.Staker,
		Amount:                    msg.Amount,
		EstimatedUndelegationDate: uint64(ctx.BlockTime().Unix()) + k.GetUnbondingDelegationTime(ctx),
	})

	return &types.MsgUndelegateResponse{}, nil
}
