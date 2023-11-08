package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// UpdateFunder allows a funder to change basic metadata like moniker, address, logo, etc.
func (k msgServer) UpdateFunder(goCtx context.Context, msg *types.MsgUpdateFunder) (*types.MsgUpdateFunderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Error if funder does not exist
	if !k.DoesFunderExist(ctx, msg.Creator) {
		return nil, errors.Wrap(errorsTypes.ErrNotFound, types.ErrFunderDoesNotExist.Error())
	}

	// Update funder
	k.SetFunder(ctx, &types.Funder{
		Address:     msg.Creator,
		Moniker:     msg.Moniker,
		Identity:    msg.Identity,
		Website:     msg.Website,
		Contact:     msg.Contact,
		Description: msg.Description,
	})

	// Emit a update funder event
	_ = ctx.EventManager().EmitTypedEvent(&types.EventUpdateFunder{
		Address:     msg.Creator,
		Moniker:     msg.Moniker,
		Identity:    msg.Identity,
		Website:     msg.Website,
		Contact:     msg.Contact,
		Description: msg.Description,
	})

	return &types.MsgUpdateFunderResponse{}, nil
}
