package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// VoteBundleProposal handles the logic of an SDK message that allows protocol nodes to vote on a pool's bundle proposal.
func (k msgServer) VoteBundleProposal(
	goCtx context.Context, msg *types.MsgVoteBundleProposal,
) (*types.MsgVoteBundleProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.AssertCanVote(ctx, msg.PoolId, msg.Staker, msg.Creator, msg.StorageId); err != nil {
		return nil, err
	}

	bundleProposal, _ := k.GetBundleProposal(ctx, msg.PoolId)
	hasVotedAbstain := util.ContainsString(bundleProposal.VotersAbstain, msg.Staker)

	if hasVotedAbstain {
		if msg.Vote == types.VOTE_TYPE_ABSTAIN {
			return nil, types.ErrAlreadyVotedAbstain
		}

		// remove voter from abstain votes
		bundleProposal.VotersAbstain, _ = util.RemoveFromStringArrayStable(bundleProposal.VotersAbstain, msg.Staker)
	}

	switch msg.Vote {
	case types.VOTE_TYPE_VALID:
		bundleProposal.VotersValid = append(bundleProposal.VotersValid, msg.Staker)
	case types.VOTE_TYPE_INVALID:
		bundleProposal.VotersInvalid = append(bundleProposal.VotersInvalid, msg.Staker)
	case types.VOTE_TYPE_ABSTAIN:
		bundleProposal.VotersAbstain = append(bundleProposal.VotersAbstain, msg.Staker)
	default:
		return nil, sdkErrors.Wrapf(sdkErrors.ErrUnauthorized, types.ErrInvalidVote.Error(), msg.Vote)
	}

	k.SetBundleProposal(ctx, bundleProposal)

	// reset points as user has now proven to be active
	k.resetPoints(ctx, msg.PoolId, msg.Staker)

	// Emit a vote event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleVote{
		PoolId:    msg.PoolId,
		Staker:    msg.Staker,
		StorageId: msg.StorageId,
		Vote:      msg.Vote,
	})

	return &types.MsgVoteBundleProposalResponse{}, nil
}
