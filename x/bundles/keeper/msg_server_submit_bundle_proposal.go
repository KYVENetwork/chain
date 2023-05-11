package keeper

import (
	"context"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"

	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitBundleProposal handles the logic of an SDK message that allows protocol nodes to submit a new bundle proposal.
func (k msgServer) SubmitBundleProposal(
	goCtx context.Context, msg *types.MsgSubmitBundleProposal,
) (*types.MsgSubmitBundleProposalResponse, error) {
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
		nextUploader := k.chooseNextUploaderFromAllStakers(ctx, msg.PoolId)

		k.registerBundleProposalFromUploader(ctx, msg, nextUploader)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

	// Previous round contains a bundle which needs to be validated now.

	// Increase points of stakers who did not vote at all + slash + remove if necessary.
	// The protocol requires everybody to stay always active.
	k.handleNonVoters(ctx, msg.PoolId)

	// evaluate all votes and determine status based on the votes weighted with stake + delegation
	voteDistribution := k.GetVoteDistribution(ctx, msg.PoolId)

	// Handle tally outcome
	switch voteDistribution.Status {

	case types.BUNDLE_STATUS_VALID:
		// If a bundle is valid the following things happen:
		// 1. A reward is paid out to the uploader, its delegators and the treasury
		//    The appropriate funds are deducted from the total pool funds
		// 2. The next uploader is randomly selected based on everybody who
		//    voted valid on this bundle.
		// 3. The bundle is finalized by added it permanently to the state.
		// 4. The sender immediately starts the next round by registering
		//    his new bundle proposal.

		// Calculate the total reward for the bundle, and individual payouts.
		bundleReward := k.calculatePayouts(ctx, msg.PoolId)

		if err := k.poolKeeper.ChargeFundersOfPool(ctx, msg.PoolId, bundleReward.Total); err != nil {
			// update the latest time on bundle to indicate that the bundle is still active
			// protocol nodes use this to determine the upload timeout
			bundleProposal.UpdatedAt = uint64(ctx.BlockTime().Unix())
			k.SetBundleProposal(ctx, bundleProposal)

			// emit event which indicates that pool has run out of funds
			_ = ctx.EventManager().EmitTypedEvent(&pooltypes.EventPoolOutOfFunds{
				PoolId: msg.PoolId,
			})

			return &types.MsgSubmitBundleProposalResponse{}, nil
		}

		pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)
		bundleProposal, _ := k.GetBundleProposal(ctx, msg.PoolId)

		uploaderPayout := bundleReward.Uploader

		delegationPayoutSuccessful := k.delegationKeeper.PayoutRewards(ctx, bundleProposal.Uploader, bundleReward.Delegation, pooltypes.ModuleName)
		// If staker has no delegators add all delegation rewards to the staker rewards
		if !delegationPayoutSuccessful {
			uploaderPayout += bundleReward.Delegation
		}

		// transfer funds from pool to stakers module
		if err := util.TransferFromModuleToModule(k.bankKeeper, ctx, pooltypes.ModuleName, stakertypes.ModuleName, uploaderPayout); err != nil {
			return nil, err
		}

		// increase commission rewards of uploader
		k.stakerKeeper.IncreaseStakerCommissionRewards(ctx, bundleProposal.Uploader, uploaderPayout)

		// send network fee to treasury
		if err := util.TransferFromModuleToTreasury(k.accountKeeper, k.distrkeeper, ctx, pooltypes.ModuleName, bundleReward.Treasury); err != nil {
			return nil, err
		}

		// slash stakers who voted incorrectly
		for _, voter := range bundleProposal.VotersInvalid {
			k.slashDelegatorsAndRemoveStaker(ctx, msg.PoolId, voter, delegationTypes.SLASH_TYPE_VOTE)
		}

		// Determine next uploader and register next bundle

		// Get next uploader from stakers who voted `valid` and are still active
		activeVoters := make([]string, 0)
		nextUploader := ""
		for _, voter := range bundleProposal.VotersValid {
			if k.stakerKeeper.DoesValaccountExist(ctx, msg.PoolId, voter) {
				activeVoters = append(activeVoters, voter)
			}
		}

		if len(activeVoters) > 0 {
			nextUploader = k.chooseNextUploaderFromSelectedStakers(ctx, msg.PoolId, activeVoters)
		} else {
			nextUploader = k.chooseNextUploaderFromAllStakers(ctx, msg.PoolId)
		}

		k.finalizeCurrentBundleProposal(ctx, pool.Id, voteDistribution, bundleReward, nextUploader)

		// Register the provided bundle as a new proposal for the next round
		k.registerBundleProposalFromUploader(ctx, msg, nextUploader)

		return &types.MsgSubmitBundleProposalResponse{}, nil

	case types.BUNDLE_STATUS_INVALID:
		// If the bundles is invalid, everybody who voted incorrectly gets slashed.
		// The bundle provided by the message-sender is of no mean, because the previous bundle
		// turned out to be incorrect.
		// There this round needs to start again and the message-sender stays uploader.

		// slash stakers who voted incorrectly - uploader receives upload slash
		for _, voter := range bundleProposal.VotersValid {
			if voter == bundleProposal.Uploader {
				k.slashDelegatorsAndRemoveStaker(ctx, msg.PoolId, voter, delegationTypes.SLASH_TYPE_UPLOAD)
			} else {
				k.slashDelegatorsAndRemoveStaker(ctx, msg.PoolId, voter, delegationTypes.SLASH_TYPE_VOTE)
			}
		}

		// Drop current bundle. Can't register the provided bundle because the previous bundles
		// needs to be resubmitted first.
		k.dropCurrentBundleProposal(ctx, msg.PoolId, voteDistribution, bundleProposal.NextUploader)

		return &types.MsgSubmitBundleProposalResponse{}, nil

	default:
		// If the bundle is neither valid nor invalid the quorum has not been reached yet.
		return nil, types.ErrQuorumNotReached
	}
}
