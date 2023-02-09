package keeper

import (
	"context"

	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Pools(c context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	pools, pageRes, err := k.poolKeeper.GetPaginatedPoolsQuery(ctx, req.Pagination, req.Search, req.Runtime, req.Disabled, req.StorageProviderId)
	if err != nil {
		return nil, err
	}

	data := make([]types.PoolResponse, 0)
	for i := range pools {
		data = append(data, k.parsePoolResponse(ctx, &pools[i]))
	}

	return &types.QueryPoolsResponse{Pools: data, Pagination: pageRes}, nil
}

func (k Keeper) Pool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.poolKeeper.GetPool(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryPoolResponse{Pool: k.parsePoolResponse(ctx, &pool)}, nil
}

func (k Keeper) parsePoolResponse(ctx sdk.Context, pool *pooltypes.Pool) types.PoolResponse {
	bundleProposal, _ := k.bundleKeeper.GetBundleProposal(ctx, pool.Id)
	stakers := k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, pool.Id)

	totalSelfDelegation := uint64(0)
	for _, address := range stakers {
		totalSelfDelegation += k.delegationKeeper.GetDelegationAmountOfDelegator(ctx, address, address)
	}

	totalDelegation := k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id)

	return types.PoolResponse{
		Id:                  pool.Id,
		Data:                pool,
		BundleProposal:      &bundleProposal,
		Stakers:             stakers,
		TotalSelfDelegation: totalSelfDelegation,
		TotalDelegation:     totalDelegation,
		Status:              k.GetPoolStatus(ctx, pool),
	}
}
