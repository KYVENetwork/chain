package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakersByPoolCount(c context.Context, req *types.QueryStakersByPoolCountRequest) (*types.QueryStakersByPoolCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	data := make([]types.FullStaker, 0)

	accumulator := func(address string, accumulate bool) bool {
		fullStaker, err := k.GetFullStaker(ctx, address)
		if err != nil {
			return false
		}

		if accumulate {
			data = append(data, *fullStaker)
		}
		return true
	}

	var pageRes *query.PageResponse
	var err error

	pageRes, err = k.stakerKeeper.GetPaginatedStakersByPoolCount(ctx, req.Pagination, accumulator)

	if err != nil {
		return nil, err
	}

	return &types.QueryStakersByPoolCountResponse{Stakers: data, Pagination: pageRes}, nil
}
