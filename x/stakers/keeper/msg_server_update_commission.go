package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/stakers/types"
)

// UpdateCommission creates a queue entry to update the staker commission.
// After the `CommissionChangeTime` is over the new commission will be applied.
// If an update is currently in the queue it will get removed from the queue
// and the user needs to wait again for the full time to pass.
func (k msgServer) UpdateCommission(goCtx context.Context, msg *types.MsgUpdateCommission) (*types.MsgUpdateCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if creator is active in the pool
	if _, active := k.GetValaccount(ctx, msg.PoolId, msg.Creator); !active {
		return nil, errors.Wrap(errorsTypes.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Insert commission change into queue
	k.orderNewCommissionChange(ctx, msg.Creator, msg.PoolId, msg.Commission)

	return &types.MsgUpdateCommissionResponse{}, nil
}
