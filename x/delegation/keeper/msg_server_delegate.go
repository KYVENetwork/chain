package keeper

import (
	"context"

	sdkErrors "cosmossdk.io/errors"

	"github.com/KYVENetwork/chain/util"

	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegate handles the transaction of delegating a specific amount of $KYVE to a staker
// The only requirement for the transaction to succeed is that the staker exists
// and the user has enough balance.
func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.stakersKeeper.DoesStakerExist(ctx, msg.Staker) {
		return nil, sdkErrors.WithType(types.ErrStakerDoesNotExist, msg.Staker)
	}

	// Performs logical delegation without transferring the amount
	k.performDelegation(ctx, msg.Staker, msg.Creator, msg.Amount)

	// Transfer tokens from sender to this module.
	if transferErr := util.TransferFromAddressToModule(k.bankKeeper, ctx, msg.Creator, types.ModuleName, msg.Amount); transferErr != nil {
		return nil, transferErr
	}

	// Emit a delegation event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventDelegate{
		Address: msg.Creator,
		Staker:  msg.Staker,
		Amount:  msg.Amount,
	})

	return &types.MsgDelegateResponse{}, nil
}
