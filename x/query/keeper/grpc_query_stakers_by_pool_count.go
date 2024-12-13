package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakersByPoolCount(c context.Context, req *types.QueryStakersByPoolCountRequest) (*types.QueryStakersByPoolCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// ToDo no - op

	return &types.QueryStakersByPoolCountResponse{}, nil
}
