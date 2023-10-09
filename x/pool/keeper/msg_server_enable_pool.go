package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

func (k msgServer) EnablePool(goCtx context.Context, req *types.MsgEnablePool) (*types.MsgEnablePoolResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, req.Id)

	if !found {
		return nil, errors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.Id)
	}
	if !pool.Disabled {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, "Pool is already enabled.")
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventPoolEnabled{Id: req.Id})

	pool.Disabled = false
	k.SetPool(ctx, pool)

	return &types.MsgEnablePoolResponse{}, nil
}
