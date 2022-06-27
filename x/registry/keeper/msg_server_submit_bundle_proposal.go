package keeper

import (
	"context"
	"strings"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SubmitBundleProposal handles the logic of an SDK message that allows protocol nodes to submit a new bundle proposal.
func (k msgServer) SubmitBundleProposal(
	goCtx context.Context, msg *types.MsgSubmitBundleProposal,
) (*types.MsgSubmitBundleProposalResponse, error) {
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
	_, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Validate bundle id.
	if msg.BundleId == "" {
		return nil, types.ErrInvalidArgs
	}

	// Get current height from where the bundle proposal should resume
	current_height := pool.CurrentHeight

	if pool.BundleProposal.ToHeight != 0 {
		current_height = pool.BundleProposal.ToHeight
	}

	// Validate from height
	if msg.FromHeight != current_height {
		return nil, types.ErrFromHeight
	}

	// Validate to height
	if msg.ToHeight < current_height {
		return nil, types.ErrToHeight
	}

	if msg.ToHeight-current_height > pool.MaxBundleSize {
		return nil, types.ErrMaxBundleSize
	}

	current_key := pool.CurrentKey

	if pool.BundleProposal.ToKey != "" {
		current_key = pool.BundleProposal.ToKey
	}

	// Validate from key
	if msg.FromKey != current_key {
		return nil, types.ErrFromKey
	}

	// Check if the sender is the designated uploader.
	if pool.BundleProposal.NextUploader != msg.Creator {
		return nil, types.ErrNotDesignatedUploader
	}

	// Check if upload_interval has been surpassed
	if uint64(ctx.BlockTime().Unix()) < (pool.BundleProposal.CreatedAt + pool.UploadInterval) {
		return nil, types.ErrUploadInterval
	}

	// EVALUATE PREVIOUS ROUND

	// Check if quorum has already been reached.
	valid := false
	invalid := false

	if len(pool.Stakers) > 1 {
		// subtract one because of uploader
		valid = len(pool.BundleProposal.VotersValid)*2 > (len(pool.Stakers) - 1)
		invalid = len(pool.BundleProposal.VotersInvalid)*2 >= (len(pool.Stakers) - 1)
	}

	// Check args of bundle types
	if strings.HasPrefix(msg.BundleId, types.KYVE_NO_DATA_BUNDLE) {
		// Validate bundle args
		if msg.ToHeight != current_height || msg.ByteSize != 0 {
			return nil, types.ErrInvalidArgs
		}

		// Validate key values
		if msg.ToKey != "" || msg.ToValue != "" {
			return nil, types.ErrInvalidArgs
		}
	} else {
		if msg.ToHeight <= current_height || msg.ByteSize == 0 {
			return nil, types.ErrInvalidArgs
		}

		// Validate key values
		if msg.ToKey == "" || msg.ToValue == "" {
			return nil, types.ErrInvalidArgs
		}
	}

	// If bundle was dropped or is of type KYVE_NO_DATA_BUNDLE just register new bundle.
	if pool.BundleProposal.BundleId == "" || strings.HasPrefix(pool.BundleProposal.BundleId, types.KYVE_NO_DATA_BUNDLE) {
		pool.BundleProposal = &types.BundleProposal{
			Uploader:     msg.Creator,
			NextUploader: k.getNextUploaderByRandom(ctx, &pool, pool.Stakers),
			BundleId:     msg.BundleId,
			ByteSize:     msg.ByteSize,
			ToHeight:     msg.ToHeight,
			CreatedAt:    uint64(ctx.BlockTime().Unix()),
			ToKey:        msg.ToKey,
			ToValue:      msg.ToValue,
		}

		k.SetPool(ctx, pool)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

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

	// handle valid proposal
	if valid {
		// Calculate the total reward for the bundle, and individual payouts.
		bundleReward := pool.OperatingCost + (pool.BundleProposal.ByteSize * k.StorageCost(ctx))

		// load and parse network fee
		networkFee, err := sdk.NewDecFromStr(k.NetworkFee(ctx))
		if err != nil {
			k.PanicHalt(ctx, "Invalid value for params: "+err.Error())
		}

		treasuryPayout := uint64(sdk.NewDec(int64(bundleReward)).Mul(networkFee).RoundInt64())
		uploaderPayout := bundleReward - treasuryPayout

		// Calculate the delegation rewards for the uploader.
		uploader, foundUploader := k.GetStaker(ctx, pool.BundleProposal.Uploader, pool.Id)
		uploaderDelegation, foundUploaderDelegation := k.GetDelegationPoolData(ctx, pool.Id, pool.BundleProposal.Uploader)

		if foundUploader && foundUploaderDelegation {
			// If the uploader has no delegators, it keeps the delegation reward.

			if uploaderDelegation.DelegatorCount > 0 {
				// Calculate the reward, factoring in the node commission, and subtract from the uploader payout.
				commission, _ := sdk.NewDecFromStr(uploader.Commission)
				delegationReward := uint64(
					sdk.NewDec(int64(uploaderPayout)).Mul(sdk.NewDec(1).Sub(commission)).RoundInt64(),
				)

				uploaderPayout -= delegationReward
				uploaderDelegation.CurrentRewards += delegationReward

				k.SetDelegationPoolData(ctx, uploaderDelegation)
			}
		}

		// Calculate the individual cost for each pool funder.
		// NOTE: Because of integer division, it is possible that there is a small remainder.
		// This remainder is in worst case MaxFundersAmount(tkyve) and is charged to the lowest funder.
		fundersCost := bundleReward / uint64(len(pool.Funders))
		fundersCostRemainder := bundleReward - (uint64(len(pool.Funders)) * fundersCost)

		// Fetch the lowest funder, and find a new one if the current one isn't found.
		lowestFunder, foundLowestFunder := k.GetFunder(ctx, pool.LowestFunder, pool.Id)

		if !foundLowestFunder {
			k.updateLowestFunder(ctx, &pool)
			lowestFunder, _ = k.GetFunder(ctx, pool.LowestFunder, pool.Id)
		}

		slashedFunds := uint64(0)

		// Remove every funder who can't afford the funder cost.
		for fundersCost+fundersCostRemainder > lowestFunder.Amount {
			// Now, let's remove all other funders who have run out of funds.
			for _, account := range pool.Funders {
				funder, _ := k.GetFunder(ctx, account, pool.Id)

				if funder.Amount < fundersCost {
					// remove funder
					k.removeFunder(ctx, &pool, &funder)

					// transfer amount to treasury
					slashedFunds += funder.Amount

					// Emit a defund event.
					errEmit := ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
						PoolId:  msg.Id,
						Address: funder.Account,
						Amount:  funder.Amount,
					})
					if errEmit != nil {
						return nil, errEmit
					}
				}
			}

			if pool.TotalFunds > 0 {
				fundersCost = bundleReward / uint64(len(pool.Funders))
				fundersCostRemainder = bundleReward - (uint64(len(pool.Funders)) * fundersCost)

				k.updateLowestFunder(ctx, &pool)
				lowestFunder, _ = k.GetFunder(ctx, pool.LowestFunder, pool.Id)
			} else {
				// Recalculate the lowest funder, update, and return.
				k.updateLowestFunder(ctx, &pool)

				if slashedFunds > 0 {
					// transfer slashed funds to treasury
					err := k.transferToTreasury(ctx, slashedFunds)
					if err != nil {
						return nil, err
					}
				}

				pool.BundleProposal = &types.BundleProposal{
					Uploader:      pool.BundleProposal.Uploader,
					NextUploader:  pool.BundleProposal.NextUploader,
					BundleId:      pool.BundleProposal.BundleId,
					ByteSize:      pool.BundleProposal.ByteSize,
					ToHeight:      pool.BundleProposal.ToHeight,
					CreatedAt:     uint64(ctx.BlockTime().Unix()),
					VotersValid:   pool.BundleProposal.VotersValid,
					VotersInvalid: pool.BundleProposal.VotersInvalid,
					ToKey:         pool.BundleProposal.ToKey,
					ToValue:       pool.BundleProposal.ToValue,
				}

				k.SetPool(ctx, pool)

				// Emit a bundle dropped event because of insufficient funds.
				errEmit := ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalised{
					PoolId:       pool.Id,
					BundleId:     pool.BundleProposal.BundleId,
					ByteSize:     pool.BundleProposal.ByteSize,
					Uploader:     pool.BundleProposal.Uploader,
					NextUploader: pool.BundleProposal.NextUploader,
					Reward:       0,
					Valid:        uint64(len(pool.BundleProposal.VotersValid)),
					Invalid:      uint64(len(pool.BundleProposal.VotersInvalid)),
					FromHeight:   pool.BundleProposal.FromHeight,
					ToHeight:     pool.BundleProposal.ToHeight,
					Status:       types.BUNDLE_STATUS_NO_FUNDS,
					ToKey:        pool.BundleProposal.ToKey,
					ToValue:      pool.BundleProposal.ToValue,
					Id:           0,
				})
				if errEmit != nil {
					return nil, errEmit
				}

				return &types.MsgSubmitBundleProposalResponse{}, nil
			}
		}

		if slashedFunds > 0 {
			// transfer slashed funds to treasury
			err := k.transferToTreasury(ctx, slashedFunds)
			if err != nil {
				return nil, err
			}
		}

		// Charge every funder equally.
		for _, account := range pool.Funders {
			funder, _ := k.GetFunder(ctx, account, pool.Id)

			if funder.Amount >= fundersCost {
				funder.Amount -= fundersCost
			}

			k.SetFunder(ctx, funder)
		}

		// Remove any remainder cost from the lowest funder.
		lowestFunder, _ = k.GetFunder(ctx, pool.LowestFunder, pool.Id)

		if lowestFunder.Amount >= fundersCostRemainder {
			lowestFunder.Amount -= fundersCostRemainder
		}

		k.SetFunder(ctx, lowestFunder)

		// Subtract bundle reward from the pool's total funds.
		pool.TotalFunds -= bundleReward

		// Partially slash all nodes who voted incorrectly.
		for _, voter := range pool.BundleProposal.VotersInvalid {
			slashAmount := k.slashStaker(ctx, &pool, voter, k.VoteSlash(ctx))

			errEmit := ctx.EventManager().EmitTypedEvent(&types.EventSlash{
				PoolId:    pool.Id,
				Address:   voter,
				Amount:    slashAmount,
				SlashType: types.SLASH_TYPE_VOTE,
			})
			if errEmit != nil {
				return nil, errEmit
			}
		}

		// Send payout to treasury.
		errTreasury := k.transferToTreasury(ctx, treasuryPayout)
		if errTreasury != nil {
			return nil, errTreasury
		}

		// Send payout to uploader.
		errTransfer := k.TransferToAddress(ctx, pool.BundleProposal.Uploader, uploaderPayout)
		if errTransfer != nil {
			return nil, errTransfer
		}

		// save valid bundle
		k.SetProposal(ctx, types.Proposal{
			BundleId:    pool.BundleProposal.BundleId,
			PoolId:      pool.Id,
			Id:          pool.TotalBundles,
			Uploader:    pool.BundleProposal.Uploader,
			FromHeight:  pool.CurrentHeight,
			ToHeight:    pool.BundleProposal.ToHeight,
			FinalizedAt: uint64(ctx.BlockHeight()),
			Key:         pool.BundleProposal.ToKey,
			Value:       pool.BundleProposal.ToValue,
		})

		// Finalise the proposal, saving useful information.
		pool.CurrentHeight = pool.BundleProposal.ToHeight
		pool.TotalBytes = pool.TotalBytes + pool.BundleProposal.ByteSize
		pool.TotalBundles = pool.TotalBundles + 1
		pool.TotalBundleRewards = pool.TotalBundleRewards + bundleReward
		pool.CurrentKey = pool.BundleProposal.ToKey
		pool.CurrentValue = pool.BundleProposal.ToValue

		// Emit a valid bundle event.
		errEmit := ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalised{
			PoolId:       pool.Id,
			BundleId:     pool.BundleProposal.BundleId,
			ByteSize:     pool.BundleProposal.ByteSize,
			Uploader:     pool.BundleProposal.Uploader,
			NextUploader: pool.BundleProposal.NextUploader,
			Reward:       bundleReward,
			Valid:        uint64(len(pool.BundleProposal.VotersValid)),
			Invalid:      uint64(len(pool.BundleProposal.VotersInvalid)),
			FromHeight:   pool.BundleProposal.FromHeight,
			ToHeight:     pool.BundleProposal.ToHeight,
			Status:       types.BUNDLE_STATUS_VALID,
			ToKey:        pool.BundleProposal.ToKey,
			ToValue:      pool.BundleProposal.ToValue,
			Id:           pool.TotalBundles - 1,
		})
		if errEmit != nil {
			return nil, errEmit
		}

		// Set submitted bundle as new bundle proposal and select new next_uploader
		pool.BundleProposal = &types.BundleProposal{
			Uploader:     msg.Creator,
			NextUploader: nextUploader,
			BundleId:     msg.BundleId,
			ByteSize:     msg.ByteSize,
			ToHeight:     msg.ToHeight,
			CreatedAt:    uint64(ctx.BlockTime().Unix()),
			ToKey:        msg.ToKey,
			ToValue:      msg.ToValue,
		}

		k.SetPool(ctx, pool)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	} else if invalid {
		// Partially slash all nodes who voted incorrectly.
		for _, voter := range pool.BundleProposal.VotersValid {
			slashAmount := k.slashStaker(ctx, &pool, voter, k.VoteSlash(ctx))

			errEmit := ctx.EventManager().EmitTypedEvent(&types.EventSlash{
				PoolId:    pool.Id,
				Address:   voter,
				Amount:    slashAmount,
				SlashType: types.SLASH_TYPE_VOTE,
			})
			if errEmit != nil {
				return nil, errEmit
			}
		}

		// Partially slash the uploader.
		slashAmount := k.slashStaker(ctx, &pool, pool.BundleProposal.Uploader, k.UploadSlash(ctx))

		// emit slash event
		errEmit := ctx.EventManager().EmitTypedEvent(&types.EventSlash{
			PoolId:    pool.Id,
			Address:   pool.BundleProposal.Uploader,
			Amount:    slashAmount,
			SlashType: types.SLASH_TYPE_UPLOAD,
		})
		if errEmit != nil {
			return nil, errEmit
		}

		// Update the current lowest staker.
		k.updateLowestStaker(ctx, &pool)

		// Emit an invalid bundle event.
		errEmit = ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalised{
			PoolId:       pool.Id,
			BundleId:     pool.BundleProposal.BundleId,
			ByteSize:     pool.BundleProposal.ByteSize,
			Uploader:     pool.BundleProposal.Uploader,
			NextUploader: pool.BundleProposal.NextUploader,
			Reward:       0,
			Valid:        uint64(len(pool.BundleProposal.VotersValid)),
			Invalid:      uint64(len(pool.BundleProposal.VotersInvalid)),
			FromHeight:   pool.BundleProposal.FromHeight,
			ToHeight:     pool.BundleProposal.ToHeight,
			Status:       types.BUNDLE_STATUS_INVALID,
			ToKey:        pool.BundleProposal.ToKey,
			ToValue:      pool.BundleProposal.ToValue,
			Id:           0,
		})
		if errEmit != nil {
			return nil, errEmit
		}

		// Update and return.
		pool.BundleProposal = &types.BundleProposal{
			NextUploader: pool.BundleProposal.NextUploader,
			CreatedAt:    uint64(ctx.BlockTime().Unix()),
		}

		k.SetPool(ctx, pool)

		return &types.MsgSubmitBundleProposalResponse{}, nil
	} else {
		return nil, types.ErrQuorumNotReached
	}
}
