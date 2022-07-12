package keeper

import (
	"context"
	"strings"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// VoteProposal handles the logic of an SDK message that allows protocol nodes to vote on a pool's bundle proposal.
func (k msgServer) VoteProposal(
	goCtx context.Context, msg *types.MsgVoteProposal,
) (*types.MsgVoteProposalResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if enough nodes are online
	if len(pool.Stakers) < 2 {
		return nil, types.ErrNotEnoughNodesOnline
	}

	// Check if minimum stake is reached
	if pool.TotalStake < pool.MinStake {
		return nil, types.ErrNotEnoughStake
	}

	// Error if the pool has no funds.
	if len(pool.Funders) == 0 {
		return nil, sdkErrors.Wrap(sdkErrors.ErrInsufficientFunds, types.ErrFundsTooLow.Error())
	}

	// Error if the pool is paused.
	if pool.Paused {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrPoolPaused.Error())
	}

	// Error if the pool is upgrading.
	if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrPoolCurrentlyUpgrading.Error())
	}

	// Check if the sender is a protocol node (aka has staked into this pool).
	staker, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Check if the sender is also the bundle's uploader.
	if pool.BundleProposal.Uploader == msg.Creator {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrVoterIsUploader.Error())
	}

	// Check if bundle is not dropped or NO_DATA_BUNDLE
	if pool.BundleProposal.StorageId == "" || strings.HasPrefix(pool.BundleProposal.StorageId, types.KYVE_NO_DATA_BUNDLE) {
		return nil, sdkErrors.Wrapf(
			sdkErrors.ErrNotFound, types.ErrInvalidStorageId.Error(), pool.BundleProposal.StorageId,
		)
	}

	// Check if the sender is voting on the same bundle.
	if msg.StorageId != pool.BundleProposal.StorageId {
		return nil, sdkErrors.Wrapf(
			sdkErrors.ErrNotFound, types.ErrInvalidStorageId.Error(), pool.BundleProposal.StorageId,
		)
	}

	// Check if the sender has already voted on the bundle.
	hasVotedValid, hasVotedInvalid, hasVotedAbstain := false, false, false

	for _, voter := range pool.BundleProposal.VotersValid {
		if voter == msg.Creator {
			hasVotedValid = true
		}
	}

	for _, voter := range pool.BundleProposal.VotersInvalid {
		if voter == msg.Creator {
			hasVotedInvalid = true
		}
	}

	for _, voter := range pool.BundleProposal.VotersAbstain {
		if voter == msg.Creator {
			hasVotedAbstain = true
		}
	}

	if hasVotedValid || hasVotedInvalid {
		return nil, sdkErrors.Wrapf(
			sdkErrors.ErrUnauthorized, types.ErrAlreadyVoted.Error(), pool.BundleProposal.StorageId,
		)
	}

	if hasVotedAbstain {
		if msg.Vote == types.VOTE_TYPE_ABSTAIN {
			return nil, sdkErrors.Wrapf(
				sdkErrors.ErrUnauthorized, types.ErrAlreadyVoted.Error(), pool.BundleProposal.StorageId,
			)
		}

		// remove voter from abstain votes
		pool.BundleProposal.VotersAbstain = removeStringFromList(pool.BundleProposal.VotersAbstain, msg.Creator)
	}

	// Update and return.
	if msg.Vote == types.VOTE_TYPE_YES {
		pool.BundleProposal.VotersValid = append(pool.BundleProposal.VotersValid, msg.Creator)
	} else if msg.Vote == types.VOTE_TYPE_NO {
		pool.BundleProposal.VotersInvalid = append(pool.BundleProposal.VotersInvalid, msg.Creator)
	} else if msg.Vote == types.VOTE_TYPE_ABSTAIN {
		pool.BundleProposal.VotersAbstain = append(pool.BundleProposal.VotersAbstain, msg.Creator)
	} else {
		return nil, sdkErrors.Wrapf(
			sdkErrors.ErrUnauthorized, types.ErrInvalidVote.Error(), msg.Vote,
		)
	}

	// reset points
	staker.Points = 0
	k.SetStaker(ctx, staker)

	// Emit a vote event.
	err := ctx.EventManager().EmitTypedEvent(&types.EventBundleVote{
		PoolId:   msg.Id,
		Address:  msg.Creator,
		StorageId: msg.StorageId,
		Vote:     msg.Vote,
	})
	if err != nil {
		return nil, err
	}

	k.SetPool(ctx, pool)

	return &types.MsgVoteProposalResponse{}, nil
}
