package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/funders/types"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefundPool handles the logic to defund a pool.
// If the user is a funder, it will subtract the provided amount
// and send the tokens back. If the amount equals the current funding amount
// the funder is removed completely.
func (k msgServer) DefundPool(goCtx context.Context, msg *types.MsgDefundPool) (*types.MsgDefundPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Pool has to exist
	pool, err := k.pookKeeper.GetPoolWithError(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// Sender needs to be a funder in the pool
	funderAmount := pool.GetFunderAmount(msg.Creator)
	if funderAmount == 0 {
		return nil, errorsTypes.ErrNotFound
	}

	// Check if the sender is trying to defund more than they have funded.
	if msg.Amount > funderAmount {
		return nil, errors.Wrapf(errorsTypes.ErrLogic, types.ErrDefundTooHigh.Error(), msg.Creator)
	}

	// Update state variables (or completely remove if fully defunding).
	pool.SubtractAmountFromFunder(msg.Creator, msg.Amount)

	// Transfer tokens from this module to sender.
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Creator, msg.Amount); err != nil {
		return nil, err
	}

	// Emit a defund event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
		PoolId:  msg.PoolId,
		Address: msg.Creator,
		Amount:  msg.Amount,
	})

	k.SetPool(ctx, pool)

	return &types.MsgDefundPoolResponse{}, nil
}
