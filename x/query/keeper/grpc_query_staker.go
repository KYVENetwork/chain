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
		searchMoniker := strings.ToLower(fullStaker.Metadata.Moniker)

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

	switch req.Status {
	case types.STAKER_STATUS_ACTIVE:
		pageRes, err = k.delegationKeeper.GetPaginatedActiveStakersByDelegation(ctx, req.Pagination, accumulator)
	case types.STAKER_STATUS_INACTIVE:
		pageRes, err = k.delegationKeeper.GetPaginatedInactiveStakersByDelegation(ctx, req.Pagination, accumulator)
	default:
		pageRes, err = k.delegationKeeper.GetPaginatedStakersByDelegation(ctx, req.Pagination, accumulator)
	}

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

	if !k.stakerKeeper.DoesStakerExist(ctx, req.Address) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryStakerResponse{Staker: *k.GetFullStaker(ctx, req.Address)}, nil
}
