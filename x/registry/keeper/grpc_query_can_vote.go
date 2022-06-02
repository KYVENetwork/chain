package keeper

import (
	"context"
	"strings"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CanVote(goCtx context.Context, req *types.QueryCanVoteRequest) (*types.QueryCanVoteResponse, error) {
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
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Not enough nodes online",
		}, nil
	}

	// Check if pool has funds
	if pool.TotalFunds == 0 {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Pool has run out of funds",
		}, nil
	}

	// Check if pool is paused
	if pool.Paused {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Pool is paused",
		}, nil
	}

	// Check if pool is upgrading
	if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Pool is upgrading",
		}, nil
	}

	// Check if sender is a staker in pool
	_, isStaker := k.GetStaker(ctx, req.Voter, req.PoolId)
	if !isStaker {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Voter is no staker",
		}, nil
	}

	// Check if dropped bundle
	if pool.BundleProposal.BundleId == "" {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Can not vote on dropped bundle",
		}, nil
	}

	// Check if empty bundle
	if strings.HasPrefix(pool.BundleProposal.BundleId, types.KYVE_NO_DATA_BUNDLE) {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Can not vote on NO_DATA_BUNDLE",
		}, nil
	}

	// Check if tx matches current bundleProposal
	if req.BundleId != pool.BundleProposal.BundleId {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Provided bundleId does not match current one",
		}, nil
	}

	// check if voter is not uploader
	if pool.BundleProposal.Uploader == req.Voter {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Voter is uploader",
		}, nil
	}

	// Check if sender has not voted yet
	hasVotedValid, hasVotedInvalid, hasVotedAbstain := false, false, false

	for _, voter := range pool.BundleProposal.VotersValid {
		if voter == req.Voter {
			hasVotedValid = true
		}
	}

	for _, voter := range pool.BundleProposal.VotersInvalid {
		if voter == req.Voter {
			hasVotedInvalid = true
		}
	}

	for _, voter := range pool.BundleProposal.VotersAbstain {
		if voter == req.Voter {
			hasVotedAbstain = true
		}
	}

	if hasVotedValid || hasVotedInvalid {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Voter already voted",
		}, nil
	}

	if hasVotedAbstain {
		return &types.QueryCanVoteResponse{
			Possible: true,
			Reason:   "KYVE_VOTE_NO_ABSTAIN_ALLOWED",
		}, nil
	}

	return &types.QueryCanVoteResponse{
		Possible: true,
		Reason:   "",
	}, nil
}
