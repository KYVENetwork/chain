package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleUploadTimeout is an end block hook that triggers an upload timeout for every pool (if applicable).
func (k Keeper) HandleUploadTimeout(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Iterate over all pool Ids.
	for _, pool := range k.poolKeeper.GetAllPools(ctx) {
		err := k.AssertPoolCanRun(ctx, pool.Id)
		bundleProposal, _ := k.GetBundleProposal(ctx, pool.Id)

		// Check if pool is active
		if err != nil {
			// if pool was disabled we drop the current bundle. We only drop
			// if there is an ongoing bundle proposal. Else we just remove the next
			// uploader
			if err == types.ErrPoolDisabled && bundleProposal.StorageId != "" {
				k.dropCurrentBundleProposal(ctx, pool.Id, types.VoteDistribution{
					Valid:   0,
					Invalid: 0,
					Abstain: 0,
					Total:   0,
					Status:  types.BUNDLE_STATUS_DISABLED,
				}, "")
			} else if bundleProposal.NextUploader != "" {
				bundleProposal.NextUploader = ""
				k.SetBundleProposal(ctx, bundleProposal)
			}

			// since a paused or disabled pool can not produce any bundles
			// we continue because timeout slashes don't apply in this case
			continue
		}

		// Skip if we haven't reached the upload interval.
		if uint64(ctx.BlockTime().Unix()) < (bundleProposal.UpdatedAt + pool.UploadInterval) {
			continue
		}

		// Check if bundle needs to be dropped
		if bundleProposal.StorageId != "" {
			// check if the quorum was actually reached
			voteDistribution := k.GetVoteDistribution(ctx, pool.Id)

			if voteDistribution.Status == types.BUNDLE_STATUS_NO_QUORUM {
				// handle stakers who did not vote at all
				k.handleNonVoters(ctx, pool.Id)

				// Get next uploader from all pool stakers
				nextUploader := k.chooseNextUploader(ctx, pool.Id)

				// If consensus wasn't reached, we drop the bundle and emit an event.
				k.dropCurrentBundleProposal(ctx, pool.Id, voteDistribution, nextUploader)
				continue
			}
		}

		// Skip if we haven't reached the upload timeout.
		if uint64(ctx.BlockTime().Unix()) < (bundleProposal.UpdatedAt + pool.UploadInterval + k.GetUploadTimeout(ctx)) {
			continue
		}

		// We now know that the pool is active and the upload timeout has been reached.

		timedoutUploader := bundleProposal.NextUploader

		// Check if we have a bundle proposal to validate.
		if bundleProposal.StorageId != "" {
			// Previous round contains a bundle which needs to be validated now.
			result, err := k.tallyBundleProposal(ctx, bundleProposal, pool.Id)
			if err != nil {
				// If we have an error here we might have an inconsistent state.
				continue
			}

			switch result.Status {
			case types.TallyResultValid:
				// Get next uploader from stakers who voted `valid`
				nextUploader := k.chooseNextUploaderFromList(ctx, pool.Id, bundleProposal.VotersValid)

				// Finalize bundle by adding it to the store
				k.finalizeCurrentBundleProposal(ctx, pool.Id, result.VoteDistribution, result.FundersPayout, result.InflationPayout, result.BundleReward, nextUploader)

				// Register empty bundle with next uploader
				bundleProposal = types.BundleProposal{
					PoolId:       pool.Id,
					NextUploader: nextUploader,
					UpdatedAt:    uint64(ctx.BlockTime().Unix()),
				}
				k.SetBundleProposal(ctx, bundleProposal)
			default:
				// In every other case the bundle is dropped.

				// Get next uploader from all pool stakers
				nextUploader := k.chooseNextUploader(ctx, pool.Id)

				// Drop current bundle and set next uploader
				k.dropCurrentBundleProposal(ctx, pool.Id, result.VoteDistribution, nextUploader)
			}
		} else {
			// Update bundle proposal and choose next uploader
			bundleProposal.NextUploader = k.chooseNextUploader(ctx, pool.Id)
			bundleProposal.UpdatedAt = uint64(ctx.BlockTime().Unix())
			k.SetBundleProposal(ctx, bundleProposal)
		}

		// Now we increase the points of the valaccount
		// (if he is still participating in the pool)
		if k.stakerKeeper.DoesValaccountExist(ctx, pool.Id, timedoutUploader) {
			k.addPoint(ctx, pool.Id, timedoutUploader)
		}
	}
}
