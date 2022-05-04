package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CanPropose(goCtx context.Context, req *types.QueryCanProposeRequest) (*types.QueryCanProposeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Load pool
	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	// Check if enough nodes are online
	if len(pool.Stakers) < 2 {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Not enough nodes online",
		}, nil
	}

	// Check if pool has funds
	if pool.TotalFunds == 0 {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Pool has run out of funds",
		}, nil
	}

	// Check if pool is paused
	if pool.Paused {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Pool is paused",
		}, nil
	}

	// Check if sender is a staker in pool
	_, isStaker := k.GetStaker(ctx, req.Proposer, req.PoolId)
	if !isStaker {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Proposer is no staker",
		}, nil
	}

	// Check if upload interval has been surpassed
	if uint64(ctx.BlockTime().Unix()) < (pool.BundleProposal.CreatedAt + pool.UploadInterval) {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Upload interval not surpassed",
		}, nil
	}

	// Check if designated uploader
	if pool.BundleProposal.NextUploader != req.Proposer {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Not designated uploader",
		}, nil
	}

	// Check if consensus has already been reached.
	valid := false
	invalid := false

	if len(pool.Stakers) > 1 {
		// subtract one because of uploader
		valid = len(pool.BundleProposal.VotersValid)*2 > (len(pool.Stakers) - 1)
		invalid = len(pool.BundleProposal.VotersInvalid)*2 >= (len(pool.Stakers) - 1)
	}

	if !valid && !invalid {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Quorum not reached yet",
		}, nil
	}

	return &types.QueryCanProposeResponse{
		Possible: true,
		Reason:   "",
	}, nil
}
