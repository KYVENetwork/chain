package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// LeavePool handles the SDK message of preparing a pool leave.
// Stakers can not leave a pool immediately. Instead, they need
// to notify the system that they want to leave a pool.
// The actual leaving happens after `LeavePoolTime` is over.
func (k msgServer) LeavePool(goCtx context.Context, msg *types.MsgLeavePool) (*types.MsgLeavePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valaccount, active := k.GetValaccount(ctx, msg.PoolId, msg.Creator)
	if !active {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrAlreadyLeftPool.Error())
	}

	valaccount.IsLeaving = true
	k.SetValaccount(ctx, valaccount)

	// Creates the queue entry to leave a pool. Does nothing further
	if err := k.orderLeavePool(ctx, msg.Creator, msg.PoolId); err != nil {
		return nil, err
	}

	return &types.MsgLeavePoolResponse{}, nil
}
