package keeper

import (
	"context"
	"strings"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Stakers(c context.Context, req *types.QueryStakersRequest) (*types.QueryStakersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	data := make([]types.FullStaker, 0)
	req.Search = strings.ToLower(req.Search)

	accumulator := func(address string, accumulate bool) bool {
		fullStaker := k.GetFullStaker(ctx, address)

		searchAddress := strings.ToLower(fullStaker.Address)
		searchMoniker := strings.ToLower(fullStaker.Validator.GetMoniker())

		if strings.Contains(searchAddress, req.Search) || strings.Contains(searchMoniker, req.Search) {
			if accumulate {
				data = append(data, *fullStaker)
			}
			return true
		}

		return false
	}

	var pageRes *query.PageResponse
	var err error

	pageRes, err = k.stakerKeeper.GetPaginatedStakersByPoolStake(ctx, req.Pagination, req.Status, accumulator)

	if err != nil {
		return nil, err
	}

	return &types.QueryStakersResponse{Stakers: data, Pagination: pageRes}, nil
}

func (k Keeper) Staker(c context.Context, req *types.QueryStakerRequest) (*types.QueryStakerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if _, exists := k.stakerKeeper.GetValidator(ctx, req.Address); !exists {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryStakerResponse{Staker: *k.GetFullStaker(ctx, req.Address)}, nil
}
