package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakersByPoolCount(c context.Context, req *types.QueryStakersByPoolCountRequest) (*types.QueryStakersByPoolCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	result, pageRes, err := k.delegationKeeper.GetPaginatedActiveStakersByPoolCountAndDelegation(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]types.FullStaker, len(result))

	for i := 0; i < len(result); i++ {
		data[i] = *k.GetFullStaker(ctx, result[i])
	}

	return &types.QueryStakersByPoolCountResponse{Stakers: data, Pagination: pageRes}, nil
}
