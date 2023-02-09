package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CurrentVoteStatus(c context.Context, req *types.QueryCurrentVoteStatusRequest) (*types.QueryCurrentVoteStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	_, found := k.poolKeeper.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.ErrKeyNotFound
	}

	voteDistribution := k.bundleKeeper.GetVoteDistribution(ctx, req.PoolId)

	return &types.QueryCurrentVoteStatusResponse{
		Valid:   voteDistribution.Valid,
		Invalid: voteDistribution.Invalid,
		Abstain: voteDistribution.Abstain,
		Total:   voteDistribution.Total,
	}, nil
}
