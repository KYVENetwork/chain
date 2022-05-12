package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Funder returns all funder info
func (k Keeper) Funder(goCtx context.Context, req *types.QueryFunderRequest) (*types.QueryFunderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryFunderResponse{}

	// Load pool
	_, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	// Load funder
	funder, found := k.GetFunder(ctx, req.Funder, req.PoolId)

	if !found {
		return nil, sdkErrors.ErrNotFound
	}

	response.Funder = &funder

	return &response, nil
}
