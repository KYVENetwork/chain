package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CanPropose(c context.Context, req *types.QueryCanProposeRequest) (*types.QueryCanProposeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := k.bundleKeeper.AssertCanPropose(ctx, req.PoolId, req.Staker, req.Proposer, req.FromIndex); err != nil {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   err.Error(),
		}, nil
	}

	return &types.QueryCanProposeResponse{
		Possible: true,
		Reason:   "",
	}, nil
}
