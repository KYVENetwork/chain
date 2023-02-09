package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// UpdateMetadata allows a staker to change basic metadata like moniker, address, logo, etc.
// The update is performed immediately.
func (k msgServer) UpdateMetadata(goCtx context.Context, msg *types.MsgUpdateMetadata) (*types.MsgUpdateMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the sender is a protocol node (aka has staked into this pool).
	if !k.DoesStakerExist(ctx, msg.Creator) {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Apply new metadata to staker
	k.UpdateStakerMetadata(ctx, msg.Creator, msg.Moniker, msg.Website, msg.Logo)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventUpdateMetadata{
		Staker:  msg.Creator,
		Moniker: msg.Moniker,
		Website: msg.Website,
		Logo:    msg.Logo,
	})

	return &types.MsgUpdateMetadataResponse{}, nil
}
