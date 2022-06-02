package keeper

import (
	"context"
	"strings"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleUploadTimeout is an end block hook that triggers an upload timeout for every pool (if applicable).
func (k Keeper) HandleUploadTimeout(goCtx context.Context) {
	// Unwrap context and fetch all pools.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pools := k.GetAllPool(ctx)

	// Iterate over all pools.
	for _, pool := range pools {
		// Check if there is an upcoming pool upgrade
		if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
			// Check if pool upgrade already has been applied
			if pool.Protocol.Version != pool.UpgradePlan.Version || pool.Protocol.Binaries != pool.UpgradePlan.Binaries {
				// perform pool upgrade
				pool.Protocol.Version = pool.UpgradePlan.Version
				pool.Protocol.Binaries = pool.UpgradePlan.Binaries
				pool.Protocol.LastUpgrade = pool.UpgradePlan.ScheduledAt
			}

			// Check if upgrade duration was reached
			if uint64(ctx.BlockTime().Unix()) >= (pool.UpgradePlan.ScheduledAt + pool.UpgradePlan.Duration) {
				// reset upgrade plan to default values
				pool.UpgradePlan = &types.UpgradePlan{}
			}
	
			k.SetPool(ctx, pool)
		}

		// Remove next uploader immediately if not enough nodes are online
		if len(pool.Stakers) < 2 && pool.BundleProposal.NextUploader != "" {
			pool.BundleProposal.NextUploader = ""
			k.SetPool(ctx, pool)
			continue
		}

		// Remove next uploader immediately if pool has no funds
		if pool.TotalFunds == 0 && pool.BundleProposal.NextUploader != "" {
			pool.BundleProposal.NextUploader = ""
			k.SetPool(ctx, pool)
			continue
		}

		// Remove next uploader immediately if pool is paused
		if pool.Paused && pool.BundleProposal.NextUploader != "" {
			pool.BundleProposal.NextUploader = ""
			k.SetPool(ctx, pool)
			continue
		}

		// Remove next uploader immediately if pool is upgrading
		if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt && pool.BundleProposal.NextUploader != "" {
			pool.BundleProposal.NextUploader = ""
			k.SetPool(ctx, pool)
			continue
		}

		// Skip if we haven't reached the upload interval.
		if uint64(ctx.BlockTime().Unix()) < (pool.BundleProposal.CreatedAt + pool.UploadInterval) {
			continue
		}

		// Check if bundle needs to be dropped
		if pool.BundleProposal.BundleId != "" && !strings.HasPrefix(pool.BundleProposal.BundleId, types.KYVE_NO_DATA_BUNDLE) {
			// Check if quorum has already been reached.
			valid := false
			invalid := false

			if len(pool.Stakers) > 1 {
				// subtract one because of uploader
				valid = len(pool.BundleProposal.VotersValid)*2 > (len(pool.Stakers) - 1)
				invalid = len(pool.BundleProposal.VotersInvalid)*2 >= (len(pool.Stakers) - 1)
			}

			// check if the quorum was actually reached
			if !valid && !invalid {
				// handle stakers who did not vote at all
				k.handleNonVoters(ctx, &pool)

				// Get next uploader
				voters := append(pool.BundleProposal.VotersValid, pool.BundleProposal.VotersInvalid...)
				nextUploader := ""

				if len(voters) > 0 {
					nextUploader = k.getNextUploaderByRandom(ctx, &pool, voters)
				} else {
					nextUploader = k.getNextUploaderByRandom(ctx, &pool, pool.Stakers)
				}

				// If consensus wasn't reached, we drop the bundle and emit an event.
				types.EmitBundleDroppedQuorumNotReachedEvent(ctx, &pool)

				pool.BundleProposal = &types.BundleProposal{
					NextUploader: nextUploader,
					FromHeight:   pool.BundleProposal.FromHeight,
					ToHeight:     pool.BundleProposal.FromHeight,
					CreatedAt:    uint64(ctx.BlockTime().Unix()),
				}

				k.SetPool(ctx, pool)
			}
		}

		// Skip if we haven't reached the upload timeout.
		if uint64(ctx.BlockTime().Unix()) < (pool.BundleProposal.CreatedAt + pool.UploadInterval + k.UploadTimeout(ctx)) {
			continue
		}

		// We now know that the pool is active and the upload timeout has been reached.
		// Now we slash and remove the current next_uploader and select a new one.

		staker, foundStaker := k.GetStaker(ctx, pool.BundleProposal.NextUploader, pool.Id)

		// skip timeout slash if staker is not found
		if foundStaker {
			// slash next_uploader for not uploading in time
			slashAmount := k.slashStaker(ctx, &pool, staker.Account, k.TimeoutSlash(ctx))

			// emit slashing event
			types.EmitSlashEvent(ctx, pool.Id, staker.Account, slashAmount)

			staker, foundStaker = k.GetStaker(ctx, pool.BundleProposal.NextUploader, pool.Id)

			// check if next uploader is still there or already removed
			if foundStaker {
				// remove current next_uploader
				k.removeStaker(ctx, &pool, &staker)

				// Transfer remaining stake to account.
				k.TransferToAddress(ctx, staker.Account, staker.Amount)

				// Emit an unstake event.
				types.EmitUnstakeEvent(ctx, pool.Id, staker.Account, staker.Amount)
			}

			// Update current lowest staker
			k.updateLowestStaker(ctx, &pool)
		}

		nextUploader := ""

		if len(pool.Stakers) > 0 {
			nextUploader = k.getNextUploaderByRandom(ctx, &pool, pool.Stakers)
		}

		// update bundle proposal
		pool.BundleProposal.NextUploader = nextUploader
		pool.BundleProposal.CreatedAt = uint64(ctx.BlockTime().Unix())

		k.SetPool(ctx, pool)
	}
}
