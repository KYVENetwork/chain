package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Delegation
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Pool
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
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

	// Increase points of stakers who did not vote at all + slash + remove if necessary.
	// The protocol requires everybody to stay always active.
	k.handleNonVoters(ctx, msg.PoolId)

	// evaluate all votes and determine status based on the votes weighted with stake + delegation
	voteDistribution := k.GetVoteDistribution(ctx, msg.PoolId)

	// Handle tally outcome
	switch voteDistribution.Status {

	case types.BUNDLE_STATUS_VALID:
		// If a bundle is valid the following things happen:
		// 1. Funders and Inflation Pool are charged. The total payout is divided
		//    between the uploader, its delegators and the treasury.
		//    The appropriate funds are deducted from the total pool funds
		// 2. The next uploader is randomly selected based on everybody who
		//    voted valid on this bundle.
		// 3. The bundle is finalized by added it permanently to the state.
		// 4. The sender immediately starts the next round by registering
		//    his new bundle proposal.

		pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)

		// charge the funders of the pool
		fundersPayout, err := k.fundersKeeper.ChargeFundersOfPool(ctx, msg.PoolId)
		if err != nil {
			return &types.MsgSubmitBundleProposalResponse{}, err
		}

		// charge the inflation pool
		inflationPayout, err := k.poolKeeper.ChargeInflationPool(ctx, msg.PoolId)
		if err != nil {
			return &types.MsgSubmitBundleProposalResponse{}, err
		}

		// calculate payouts to the different stakeholders like treasury, uploader and delegators
		bundleReward := k.calculatePayouts(ctx, msg.PoolId, fundersPayout+inflationPayout)

		// payout rewards to treasury
		if err := util.TransferFromModuleToTreasury(k.accountKeeper, k.distrkeeper, ctx, poolTypes.ModuleName, bundleReward.Treasury); err != nil {
			return nil, err
		}

		// payout rewards to uploader through commission rewards
		if err := k.stakerKeeper.IncreaseStakerCommissionRewards(ctx, bundleProposal.Uploader, bundleReward.Uploader); err != nil {
			return nil, err
		}

		// payout rewards to delegators through delegation rewards
		if err := k.delegationKeeper.PayoutRewards(ctx, bundleProposal.Uploader, bundleReward.Delegation, poolTypes.ModuleName); err != nil {
			return nil, err
		}

		// slash stakers who voted incorrectly
		for _, voter := range bundleProposal.VotersInvalid {
			k.slashDelegatorsAndRemoveStaker(ctx, msg.PoolId, voter, delegationTypes.SLASH_TYPE_VOTE)
		}

		// Determine next uploader and register next bundle

		// Get next uploader from stakers who voted `valid`
		nextUploader := k.chooseNextUploaderFromList(ctx, msg.PoolId, bundleProposal.VotersValid)

		k.finalizeCurrentBundleProposal(ctx, pool.Id, voteDistribution, fundersPayout, inflationPayout, bundleReward, nextUploader)

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
