package keeper

import (
	"context"

	sdkErrors "cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Redelegate lets a user redelegate from one staker to another staker
// The user has N redelegation spells. When this transaction is executed
// one spell is used. When all spells are consumed the transaction fails.
// The user then needs to wait for the oldest spell to expire to call
// this transaction again.
// It's only possible to redelegate to stakers which are at least in one pool.
func (k msgServer) Redelegate(goCtx context.Context, msg *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the sender is a delegator
	if !k.DoesDelegatorExist(ctx, msg.FromStaker, msg.Creator) {
		return nil, sdkErrors.Wrapf(types.ErrNotADelegator, "%s does not delegate to %s", msg.Creator, msg.FromStaker)
	}

	// Check if destination staker exists
	if !k.stakersKeeper.DoesStakerExist(ctx, msg.ToStaker) {
		return nil, sdkErrors.WithType(types.ErrStakerDoesNotExist, msg.ToStaker)
	}

	if len(k.stakersKeeper.GetPoolAccountsFromStaker(ctx, msg.ToStaker)) == 0 {
		return nil, sdkErrors.WithType(types.ErrRedelegationToInactiveStaker, msg.ToStaker)
	}

	// Check if the sender is trying to undelegate more than he has delegated.
	if delegationAmount := k.GetDelegationAmountOfDelegator(ctx, msg.FromStaker, msg.Creator); msg.Amount > delegationAmount {
		return nil, types.ErrNotEnoughDelegation.Wrapf("%d > %d", msg.Amount, delegationAmount)
	}

	// Only errors if all spells are currently on cooldown
	if err := k.consumeRedelegationSpell(ctx, msg.Creator); err != nil {
		return nil, err
	}

	// The redelegation is translated into an undelegation from the old staker ...
	if actualAmount := k.performUndelegation(ctx, msg.FromStaker, msg.Creator, msg.Amount); actualAmount != msg.Amount {
		return nil, types.ErrNotEnoughDelegation.Wrapf("%d != %d", msg.Amount, actualAmount)
	}
	// ... and a new delegation to the new staker
	k.performDelegation(ctx, msg.ToStaker, msg.Creator, msg.Amount)

	// Emit a delegation event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventRedelegate{
		Address:    msg.Creator,
		FromStaker: msg.FromStaker,
		ToStaker:   msg.ToStaker,
		Amount:     msg.Amount,
	})

	return &types.MsgRedelegateResponse{}, nil
}
