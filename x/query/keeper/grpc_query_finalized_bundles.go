package keeper

import (
	"context"
	"strconv"

	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FinalizedBundlesQuery(c context.Context, req *types.QueryFinalizedBundlesRequest) (*types.QueryFinalizedBundlesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if req.Index != "" {
		index, err := strconv.ParseUint(req.Index, 10, 64)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "index needs to be an unsigned integer")
		}
		bundle, found := k.bundleKeeper.GetFinalizedBundleByIndex(ctx, req.PoolId, index)
		data := make([]types.FinalizedBundle, 0)
		if found {
			data = append(data, bundle)
		}
		return &types.QueryFinalizedBundlesResponse{FinalizedBundles: data, Pagination: nil}, nil
	} else {
		finalizedBundles, pageRes, err := k.bundleKeeper.GetPaginatedFinalizedBundleQuery(ctx, req.Pagination, req.PoolId)
		return &types.QueryFinalizedBundlesResponse{FinalizedBundles: finalizedBundles, Pagination: pageRes}, err
	}
}

func (k Keeper) FinalizedBundleQuery(c context.Context, req *types.QueryFinalizedBundleRequest) (*types.FinalizedBundle, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	finalizedBundle, found := k.bundleKeeper.GetFinalizedBundle(ctx, req.PoolId, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	versionMap := k.bundleKeeper.GetBundleVersionMap(ctx).GetMap()
	response := bundlesKeeper.RawBundleToQueryBundle(finalizedBundle, versionMap)
	return &response, nil
}
