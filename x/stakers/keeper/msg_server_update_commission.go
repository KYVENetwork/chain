package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// UpdateCommission creates a queue entry to update the staker commission.
// After the `CommissionChangeTime` is over the new commission will be applied.
// If an update is currently in the queue it will get removed from the queue
// and the user needs to wait again for the full time to pass.
func (k msgServer) UpdateCommission(goCtx context.Context, msg *types.MsgUpdateCommission) (*types.MsgUpdateCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the sender is a protocol node (aka has staked into this pool).
	if !k.DoesStakerExist(ctx, msg.Creator) {
		return nil, errors.Wrap(errorsTypes.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Validate commission.
	commission, err := sdk.NewDecFromStr(msg.Commission)
	if err != nil {
		return nil, errors.Wrapf(errorsTypes.ErrLogic, types.ErrInvalidCommission.Error(), msg.Commission)
	}

	if commission.LT(sdk.NewDec(int64(0))) || commission.GT(sdk.NewDec(int64(1))) {
		return nil, errors.Wrapf(errorsTypes.ErrLogic, types.ErrInvalidCommission.Error(), msg.Commission)
	}

	// Insert commission change into queue
	k.orderNewCommissionChange(ctx, msg.Creator, msg.Commission)

	return &types.MsgUpdateCommissionResponse{}, nil
}
