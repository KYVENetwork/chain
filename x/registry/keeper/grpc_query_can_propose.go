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

	// Check if minimum stake is reached
	if pool.TotalStake < pool.MinStake {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Not enough stake in pool",
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

	// Check if pool is upgrading
	if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Pool is upgrading",
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

	// Check if from_height matches
	if pool.BundleProposal.ToHeight != req.FromHeight {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Invalid from_height",
		}, nil
	}

	// Check if designated uploader
	if pool.BundleProposal.NextUploader != req.Proposer {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Not designated uploader",
		}, nil
	}

	// Check if upload interval has been surpassed
	if uint64(ctx.BlockTime().Unix()) < (pool.BundleProposal.CreatedAt + pool.UploadInterval) {
		return &types.QueryCanProposeResponse{
			Possible: false,
			Reason:   "Upload interval not surpassed",
		}, nil
	}

	return &types.QueryCanProposeResponse{
		Possible: true,
		Reason:   "",
	}, nil
}
