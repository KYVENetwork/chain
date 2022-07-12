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
		// Set pool status
		if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
			pool.Status = types.POOL_STATUS_UPGRADING
		} else if pool.Paused {
			pool.Status = types.POOL_STATUS_PAUSED
		} else if len(pool.Stakers) < 2 {
			pool.Status = types.POOL_STATUS_NOT_ENOUGH_VALIDATORS
		} else if pool.TotalStake < pool.MinStake {
			pool.Status = types.POOL_STATUS_NOT_ENOUGH_STAKE
		} else if pool.TotalFunds == 0 {
			pool.Status = types.POOL_STATUS_NO_FUNDS
		} else {
			pool.Status = types.POOL_STATUS_ACTIVE
		}

		// Remove next uploader if pool is not active
		if pool.Status != types.POOL_STATUS_ACTIVE {
			pool.BundleProposal.NextUploader = ""
		}

		// Update status
		k.SetPool(ctx, pool)

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

		// Skip if we haven't reached the upload interval.
		if uint64(ctx.BlockTime().Unix()) < (pool.BundleProposal.CreatedAt + pool.UploadInterval) {
			continue
		}

		// Check if bundle needs to be dropped
		if pool.BundleProposal.StorageId != "" && !strings.HasPrefix(pool.BundleProposal.StorageId, types.KYVE_NO_DATA_BUNDLE) {
			// check if the quorum was actually reached
			valid, invalid, abstain, total := k.getVoteDistribution(ctx, &pool)
			quorum := k.getQuorumStatus(valid, invalid, abstain, total)

			if quorum == types.BUNDLE_STATUS_NO_QUORUM {
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
				ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalised{
					PoolId:       pool.Id,
					StorageId:    pool.BundleProposal.StorageId,
					ByteSize:     pool.BundleProposal.ByteSize,
					Uploader:     pool.BundleProposal.Uploader,
					NextUploader: pool.BundleProposal.NextUploader,
					Reward:       0,
					Valid:        valid,
					Invalid:      invalid,
					FromHeight:   pool.BundleProposal.FromHeight,
					ToHeight:     pool.BundleProposal.ToHeight,
					Status:       types.BUNDLE_STATUS_NO_QUORUM,
					ToKey:        pool.BundleProposal.ToKey,
					ToValue:      pool.BundleProposal.ToValue,
					Id:           0,
					Abstain: abstain,
					Total: total,
				})

				pool.BundleProposal = &types.BundleProposal{
					NextUploader: nextUploader,
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
			ctx.EventManager().EmitTypedEvent(&types.EventSlash{
				PoolId:    pool.Id,
				Address:   staker.Account,
				Amount:    slashAmount,
				SlashType: types.SLASH_TYPE_TIMEOUT,
			})

			staker, foundStaker = k.GetStaker(ctx, pool.BundleProposal.NextUploader, pool.Id)

			// check if next uploader is still there or already removed
			if foundStaker {

				deactivateStaker(&pool, &staker)
				k.SetStaker(ctx, staker)

				ctx.EventManager().EmitTypedEvent(&types.EventStakerStatusChanged{
					PoolId:  pool.Id,
					Address: staker.Account,
					Status:  types.STAKER_STATUS_INACTIVE,
				})
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
