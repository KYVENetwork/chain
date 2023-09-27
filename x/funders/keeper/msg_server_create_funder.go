package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// CreateFunder creates a new funder object and stores it in the store.
// If the funder already exists, an error is returned.
// TODO(rapha): this can be spammed right now. Someone can just created a bunch of funders to get displayed on the funders page. We should probably add a fee to this.
func (k msgServer) CreateFunder(goCtx context.Context, msg *types.MsgCreateFunder) (*types.MsgCreateFunderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if funder already exists
	if k.DoesFunderExist(ctx, msg.Creator) {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrFunderAlreadyExists.Error(), msg.Creator)
	}

	// Create new funder
	k.setFunder(ctx, &types.Funder{
		Address:     msg.Creator,
		Moniker:     msg.Moniker,
		Identity:    msg.Identity,
		Logo:        msg.Logo,
		Website:     msg.Website,
		Contact:     msg.Contact,
		Description: msg.Description,
	})

	// Emit a create funder event
	_ = ctx.EventManager().EmitTypedEvent(&types.EventCreateFunder{
		Address:     msg.Creator,
		Moniker:     msg.Moniker,
		Identity:    msg.Identity,
		Logo:        msg.Logo,
		Website:     msg.Website,
		Contact:     msg.Contact,
		Description: msg.Description,
	})

	return &types.MsgCreateFunderResponse{}, nil
}
