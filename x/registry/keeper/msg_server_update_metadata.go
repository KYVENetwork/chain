package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// UpdateMetadata handles the logic of an SDK message that allows protocol nodes to update their node's metadata.
func (k msgServer) UpdateMetadata(
	goCtx context.Context, msg *types.MsgUpdateMetadata,
) (*types.MsgUpdateMetadataResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is a protocol node (aka has staked into this pool).
	staker, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Validate commission.
	commission, err := sdk.NewDecFromStr(msg.Commission)
	if err != nil {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidCommission.Error(), msg.Commission)
	}

	if commission.LT(sdk.NewDec(int64(0))) || commission.GT(sdk.NewDec(int64(1))) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidCommission.Error(), msg.Commission)
	}

	// Update and return.
	staker.Commission = msg.Commission
	staker.Moniker = msg.Moniker
	staker.Website = msg.Website
	staker.Logo = msg.Logo

	k.SetStaker(ctx, staker)

	// Event an event.
	types.EmitUpdateMetadata(ctx, msg.Creator, msg.Id, msg.Commission, msg.Moniker, msg.Website, msg.Logo)

	return &types.MsgUpdateMetadataResponse{}, nil
}
