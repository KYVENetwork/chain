package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitBundleProposal handles the logic of an SDK message that allows protocol nodes to submit a new bundle proposal.
func (k msgServer) SubmitBundleProposal(goCtx context.Context, msg *types.MsgSubmitBundleProposal) (*types.MsgSubmitBundleProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.AssertCanPropose(ctx, msg.PoolId, msg.Staker, msg.Creator, msg.FromIndex); err != nil {
		return nil, err
	}

	bundleProposal, _ := k.GetBundleProposal(ctx, msg.PoolId)

	// Validate submit bundle args.
	if err := k.validateSubmitBundleArgs(ctx, &bundleProposal, msg); err != nil {
		return nil, err
	}

	// Reset points of uploader as node has proven to be active.
	k.resetPoints(ctx, msg.PoolId, msg.Staker)

	// If previous bundle was dropped just register the new bundle.
	// No previous round needs to be evaluated
	if bundleProposal.StorageId == "" {
		nextUploader := k.chooseNextUploader(ctx, msg.PoolId)

		k.registerBundleProposalFromUploader(ctx, msg, nextUploader)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

	// Previous round contains a bundle which needs to be validated now.
	result, err := k.tallyBundleProposal(ctx, bundleProposal, msg.PoolId)
	if err != nil {
		return nil, err
	}

	switch result.Status {
	case types.TallyResultValid:
		// Get next uploader from stakers who voted `valid`
		nextUploader := k.chooseNextUploaderFromList(ctx, msg.PoolId, bundleProposal.VotersValid)

		// Bundle is finalized by adding it to the store
		k.finalizeCurrentBundleProposal(ctx, msg.PoolId, result.VoteDistribution, result.FundersPayout, result.InflationPayout, result.BundleReward, nextUploader)

		// Register the provided bundle as a new proposal for the next round
		k.registerBundleProposalFromUploader(ctx, msg, nextUploader)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	case types.TallyResultInvalid:
		// Drop current bundle. Can't register the provided bundle because the previous bundles
		// needs to be resubmitted first.
		k.dropCurrentBundleProposal(ctx, msg.PoolId, result.VoteDistribution, bundleProposal.NextUploader)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	case types.TallyResultNoQuorum:
		return nil, types.ErrQuorumNotReached
	}
	return nil, types.ErrQuorumNotReached
}
