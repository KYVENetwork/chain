package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"

	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AssertPoolCanRun checks whether the given pool fulfils all
// technical/formal requirements to produce bundles
func (k Keeper) AssertPoolCanRun(ctx sdk.Context, poolId uint64) error {
	pool, poolErr := k.poolKeeper.GetPoolWithError(ctx, poolId)
	if poolErr != nil {
		return poolErr
	}

	// Error if the pool is upgrading.
	if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
		return types.ErrPoolCurrentlyUpgrading
	}

	// Error if the pool is disabled.
	if pool.Disabled {
		return types.ErrPoolDisabled
	}

	// Error if the end key is reached. The pool will simply halt if this is the case,
	// it is the responsibility of the protocol nodes to reach final consensus and that
	// a bundle does not exceed the end_key
	if pool.EndKey != "" && pool.CurrentKey == pool.EndKey {
		return types.ErrEndKeyReached
	}

	// Get the total and the highest delegation of a single validator in the pool
	totalDelegation, highestDelegation := k.delegationKeeper.GetTotalAndHighestDelegationOfPool(ctx, poolId)

	// Error if min delegation is not reached
	if totalDelegation < pool.MinDelegation {
		return types.ErrMinDelegationNotReached
	}

	maxVotingPower := k.poolKeeper.GetMaxVotingPowerPerPool(ctx)
	maxDelegation := uint64(maxVotingPower.MulInt64(int64(totalDelegation)).TruncateInt64())

	// Error if highest delegation exceeds max voting power
	if highestDelegation > maxDelegation {
		return types.ErrVotingPowerTooHigh
	}

	return nil
}

// AssertCanVote checks whether a participant in the network can vote on
// a bundle proposal in a storage pool
func (k Keeper) AssertCanVote(ctx sdk.Context, poolId uint64, staker string, voter string, storageId string) error {
	// Check basic pool configs
	if err := k.AssertPoolCanRun(ctx, poolId); err != nil {
		return err
	}

	// Check if sender is a staker in pool
	if err := k.stakerKeeper.AssertValaccountAuthorized(ctx, poolId, staker, voter); err != nil {
		return err
	}

	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// Check if dropped bundle
	if bundleProposal.StorageId == "" {
		return types.ErrBundleDropped
	}

	// Check if tx matches current bundleProposal
	if storageId != bundleProposal.StorageId {
		return types.ErrInvalidStorageId
	}

	// Check if the sender has already voted on the bundle.
	hasVotedValid := util.ContainsString(bundleProposal.VotersValid, staker)
	hasVotedInvalid := util.ContainsString(bundleProposal.VotersInvalid, staker)

	if hasVotedValid {
		return types.ErrAlreadyVotedValid
	}

	if hasVotedInvalid {
		return types.ErrAlreadyVotedInvalid
	}

	return nil
}

// AssertCanPropose checks whether a participant can submit the next bundle
// proposal in a storage pool
func (k Keeper) AssertCanPropose(ctx sdk.Context, poolId uint64, staker string, proposer string, fromIndex uint64) error {
	// Check basic pool configs
	if err := k.AssertPoolCanRun(ctx, poolId); err != nil {
		return err
	}

	// Check if sender is a staker in pool
	if err := k.stakerKeeper.AssertValaccountAuthorized(ctx, poolId, staker, proposer); err != nil {
		return err
	}

	pool, _ := k.poolKeeper.GetPoolWithError(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// Check if designated uploader
	if bundleProposal.NextUploader != staker {
		return errors.Wrapf(types.ErrNotDesignatedUploader, "expected %v received %v", bundleProposal.NextUploader, staker)
	}

	// Check if upload interval has been surpassed
	if uint64(ctx.BlockTime().Unix()) < (bundleProposal.UpdatedAt + pool.UploadInterval) {
		return errors.Wrapf(types.ErrUploadInterval, "expected %v < %v", ctx.BlockTime().Unix(), bundleProposal.UpdatedAt+pool.UploadInterval)
	}

	// Check if from_index matches
	if pool.CurrentIndex+bundleProposal.BundleSize != fromIndex {
		return errors.Wrapf(types.ErrFromIndex, "expected %v received %v", pool.CurrentIndex+bundleProposal.BundleSize, fromIndex)
	}

	return nil
}

// validateSubmitBundleArgs validates various bundle proposal metadata for correctness and
// fails if at least one requirement is not met
func (k Keeper) validateSubmitBundleArgs(ctx sdk.Context, bundleProposal *types.BundleProposal, msg *types.MsgSubmitBundleProposal) error {
	pool, err := k.poolKeeper.GetPoolWithError(ctx, msg.PoolId)
	if err != nil {
		return err
	}

	// Validate storage id
	if msg.StorageId == "" {
		return types.ErrInvalidArgs
	}

	// Validate from index
	if pool.CurrentIndex+bundleProposal.BundleSize != msg.FromIndex {
		return errors.Wrapf(types.ErrFromIndex, "expected %v received %v", pool.CurrentIndex+bundleProposal.BundleSize, msg.FromIndex)
	}

	// Validate if bundle is bigger than zero
	if msg.BundleSize == 0 {
		return types.ErrInvalidArgs
	}

	// Validate if bundle is not too big
	if msg.BundleSize > pool.MaxBundleSize {
		return errors.Wrapf(types.ErrMaxBundleSize, "expected %v received %v", pool.MaxBundleSize, msg.BundleSize)
	}

	// Validate key values
	if msg.FromKey == "" || msg.ToKey == "" {
		return types.ErrInvalidArgs
	}

	return nil
}

// slashDelegatorsAndRemoveStaker slashes a staker with a certain slashType and all including
// delegators and removes him from the storage pool
func (k Keeper) slashDelegatorsAndRemoveStaker(ctx sdk.Context, poolId uint64, stakerAddress string, slashType delegationTypes.SlashType) {
	k.delegationKeeper.SlashDelegators(ctx, poolId, stakerAddress, slashType)
	k.stakerKeeper.LeavePool(ctx, stakerAddress, poolId)
}

// resetPoints resets the points from a valaccount to zero
func (k Keeper) resetPoints(ctx sdk.Context, poolId uint64, stakerAddress string) {
	previousPoints := k.stakerKeeper.ResetPoints(ctx, poolId, stakerAddress)

	// only reset points if valaccount has at least a point
	if previousPoints > 0 {
		_ = ctx.EventManager().EmitTypedEvent(&types.EventPointsReset{
			PoolId: poolId,
			Staker: stakerAddress,
		})
	}
}

// addPoint increases the points of a valaccount with one and automatically
// slashes and removes the staker once he reaches max points
func (k Keeper) addPoint(ctx sdk.Context, poolId uint64, stakerAddress string) {
	// Add one point to staker in given pool
	points := k.stakerKeeper.IncrementPoints(ctx, poolId, stakerAddress)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventPointIncreased{
		PoolId:        poolId,
		Staker:        stakerAddress,
		CurrentPoints: points,
	})

	if points >= k.GetMaxPoints(ctx) {
		// slash all delegators with a timeout slash and remove staker from pool.
		// points are reset due to the valaccount being deleted while leaving the pool
		k.slashDelegatorsAndRemoveStaker(ctx, poolId, stakerAddress, delegationTypes.SLASH_TYPE_TIMEOUT)
	}
}

// handleNonVoters checks if stakers in a pool voted on the current bundle proposal
// if a staker did not vote at all on a bundle proposal he received points
// if a staker receives a certain number of points he receives a timeout slash and gets
// kicked out of a pool
func (k Keeper) handleNonVoters(ctx sdk.Context, poolId uint64) {
	voters := map[string]bool{}
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	for _, address := range bundleProposal.VotersValid {
		voters[address] = true
	}

	for _, address := range bundleProposal.VotersInvalid {
		voters[address] = true
	}

	for _, address := range bundleProposal.VotersAbstain {
		voters[address] = true
	}

	for _, staker := range k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		if !voters[staker] {
			k.addPoint(ctx, poolId, staker)
		}
	}
}

// calculatePayouts calculates the different payouts to treasury, uploader and delegators from the total payout
// the pool module provides for this bundle round
func (k Keeper) calculatePayouts(ctx sdk.Context, poolId uint64, totalPayout sdk.Coins) (bundleReward types.BundleReward) {
	// This method first subtracts the network fee from the total payout dedicated for this bundle.
	// After that the uploader receives the storage rewards which are based on the byte size of the bundle.
	// If the total payout does not cover the storage rewards we pay out the remains, the commission and
	// delegation rewards will be empty  in this case. After the payout of the storage rewards the remains
	// are divided between uploader and its delegators based on the uploader's commission.
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// Should not happen, if so make no payouts
	if !k.stakerKeeper.DoesStakerExist(ctx, bundleProposal.Uploader) {
		return
	}

	// Already return if there are no payouts
	if totalPayout.IsZero() {
		return
	}

	bundleReward.Total = totalPayout

	// calculate share of treasury from total payout
	treasuryPayout, _ := sdk.NewDecCoinsFromCoins(totalPayout...).MulDec(k.GetNetworkFee(ctx)).TruncateDecimal()
	bundleReward.Treasury = treasuryPayout

	// subtract treasury payout from total payout, so we can continue splitting up the rewards from here
	totalPayout = totalPayout.Sub(bundleReward.Treasury...)
	if totalPayout.IsZero() {
		return
	}

	// subtract storage cost from remaining total payout. We split the storage cost between all coins and charge
	// the amount per coin, the idea is that every coin should contribute the same USD value to the total storage
	// reward. This is done by defining the storage cost as USD / byte and the coin weights as USD / coin denom.
	//
	// If there is not enough of a coin available to cover the storage reward per coin we simply charge what is left,
	// so there can be the case that the storageRewards are less than what we actually wanted to pay out. This is
	// acceptable because this case is very rare, usually the minFundingAmount ensures that there are always enough
	// funds left of each coin, and in the case there are not enough the coins are removed and therefore for the
	// next bundle we split between the other remaining coins.
	whitelist := k.fundersKeeper.GetCoinWhitelistMap(ctx)
	// wantedStorageRewards are the amounts based on the current storage cost we want to pay out, this can be more
	// than we have available in totalPayout
	wantedStorageRewards := sdk.NewCoins()
	// storageCostPerCoin is the storage cost in $USD for each coin. This implies that each coin contributes the same
	// amount of value to the storage rewards
	storageCostPerCoin := k.GetStorageCost(ctx, bundleProposal.StorageProviderId).MulInt64(int64(bundleProposal.DataSize)).QuoInt64(int64(totalPayout.Len()))
	for _, coin := range totalPayout {
		weight := whitelist[coin.Denom].CoinWeight
		if weight.IsZero() {
			continue
		}

		// currencyUnit is the amount of base denoms of the currency
		currencyUnit := math.LegacyNewDec(10).Power(uint64(whitelist[coin.Denom].CoinDecimals))
		// amount is the value of storageCostPerCoin in the base denomination of the currency. We calculate this
		// by multiplying first with the amount of base denoms of the currency and then divide this by the $USD
		// value per currency unit which is the weight.
		amount := storageCostPerCoin.Mul(currencyUnit).Quo(weight).TruncateInt()
		wantedStorageRewards = wantedStorageRewards.Add(sdk.NewCoin(coin.Denom, amount))
	}

	// we take the min here since there can be the case where we want to charge more coins for the storage
	// reward than we have left in the total payout
	bundleReward.UploaderStorageCost = totalPayout.Min(wantedStorageRewards)

	// the remaining total payout is split between the uploader and his delegators.
	totalPayout = totalPayout.Sub(bundleReward.UploaderStorageCost...)
	if totalPayout.IsZero() {
		return
	}

	commission := k.stakerKeeper.GetCommission(ctx, bundleProposal.Uploader)
	commissionRewards, _ := sdk.NewDecCoinsFromCoins(totalPayout...).MulDec(commission).TruncateDecimal()
	bundleReward.UploaderCommission = commissionRewards

	// the remaining total payout belongs to the delegators
	totalPayout = totalPayout.Sub(bundleReward.UploaderCommission...)
	if totalPayout.IsZero() {
		return
	}

	// if the uploader has no delegators he receives the entire remaining amount
	if k.delegationKeeper.GetDelegationAmount(ctx, bundleProposal.Uploader) > 0 {
		bundleReward.Delegation = totalPayout
	} else {
		bundleReward.UploaderCommission = bundleReward.UploaderCommission.Add(totalPayout...)
	}

	return
}

// registerBundleProposalFromUploader handles the registration of the new bundle proposal
// an uploader has just submitted. With this new bundle proposal other participants
// can vote on it.
func (k Keeper) registerBundleProposalFromUploader(ctx sdk.Context, msg *types.MsgSubmitBundleProposal, nextUploader string) {
	pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)

	bundleProposal := types.BundleProposal{
		PoolId:            msg.PoolId,
		Uploader:          msg.Staker,
		NextUploader:      nextUploader,
		StorageId:         msg.StorageId,
		DataSize:          msg.DataSize,
		BundleSize:        msg.BundleSize,
		UpdatedAt:         uint64(ctx.BlockTime().Unix()),
		VotersValid:       append(make([]string, 0), msg.Staker),
		FromKey:           msg.FromKey,
		ToKey:             msg.ToKey,
		BundleSummary:     msg.BundleSummary,
		DataHash:          msg.DataHash,
		StorageProviderId: pool.CurrentStorageProviderId,
		CompressionId:     pool.CurrentCompressionId,
	}

	k.SetBundleProposal(ctx, bundleProposal)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleProposed{
		PoolId:            bundleProposal.PoolId,
		Id:                pool.TotalBundles,
		StorageId:         bundleProposal.StorageId,
		Uploader:          bundleProposal.Uploader,
		DataSize:          bundleProposal.DataSize,
		FromIndex:         pool.CurrentIndex,
		BundleSize:        bundleProposal.BundleSize,
		FromKey:           bundleProposal.FromKey,
		ToKey:             bundleProposal.ToKey,
		BundleSummary:     bundleProposal.BundleSummary,
		DataHash:          bundleProposal.DataHash,
		ProposedAt:        uint64(ctx.BlockTime().Unix()),
		StorageProviderId: bundleProposal.StorageProviderId,
		CompressionId:     bundleProposal.CompressionId,
	})

	// Emit a vote event. Uploader automatically votes valid on their bundle.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleVote{
		PoolId:    msg.PoolId,
		Staker:    msg.Staker,
		StorageId: msg.StorageId,
		Vote:      types.VOTE_TYPE_VALID,
	})
}

// finalizeCurrentBundleProposal takes the data of the current evaluated proposal
// and stores it as a finalized proposal. This only happens if the network
// reached quorum on the proposal's validity.
func (k Keeper) finalizeCurrentBundleProposal(ctx sdk.Context, poolId uint64, voteDistribution types.VoteDistribution, fundersPayout sdk.Coins, inflationPayout uint64, bundleReward types.BundleReward, nextUploader string) {
	pool, _ := k.poolKeeper.GetPool(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// save finalized bundle
	finalizedAt := types.FinalizedAt{
		Height:    uint64(ctx.BlockHeight()),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	}
	finalizedBundle := types.FinalizedBundle{
		StorageId:         bundleProposal.StorageId,
		PoolId:            pool.Id,
		Id:                pool.TotalBundles,
		Uploader:          bundleProposal.Uploader,
		FromIndex:         pool.CurrentIndex,
		ToIndex:           pool.CurrentIndex + bundleProposal.BundleSize,
		FinalizedAt:       &finalizedAt,
		FromKey:           bundleProposal.FromKey,
		ToKey:             bundleProposal.ToKey,
		BundleSummary:     bundleProposal.BundleSummary,
		DataHash:          bundleProposal.DataHash,
		StorageProviderId: bundleProposal.StorageProviderId,
		CompressionId:     bundleProposal.CompressionId,
		StakeSecurity: &types.StakeSecurity{
			ValidVotePower: voteDistribution.Valid,
			TotalVotePower: voteDistribution.Total,
		},
	}

	k.SetFinalizedBundle(ctx, finalizedBundle)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalized{
		PoolId:                    finalizedBundle.PoolId,
		Id:                        finalizedBundle.Id,
		Valid:                     voteDistribution.Valid,
		Invalid:                   voteDistribution.Invalid,
		Abstain:                   voteDistribution.Abstain,
		Total:                     voteDistribution.Total,
		Status:                    voteDistribution.Status,
		FundersPayout:             fundersPayout.String(),
		InflationPayout:           inflationPayout,
		RewardTreasury:            bundleReward.Treasury.String(),
		RewardUploader:            bundleReward.UploaderCommission.Add(bundleReward.UploaderStorageCost...).String(),
		RewardDelegation:          bundleReward.Delegation.String(),
		RewardTotal:               bundleReward.Total.String(),
		FinalizedAt:               uint64(ctx.BlockTime().Unix()),
		Uploader:                  bundleProposal.Uploader,
		NextUploader:              nextUploader,
		RewardUploaderCommission:  bundleReward.UploaderCommission.String(),
		RewardUploaderStorageCost: bundleReward.UploaderStorageCost.String(),
	})

	// Finalize the proposal, saving useful information.
	k.poolKeeper.IncrementBundleInformation(ctx, pool.Id, pool.CurrentIndex+bundleProposal.BundleSize, bundleProposal.ToKey, bundleProposal.BundleSummary)
}

// dropCurrentBundleProposal removes the current proposal due to not reaching
// a required quorum on the validity of the data. When the proposal is dropped
// the same next uploader as before can submit his proposal since it is not his
// fault, that the last one did not reach any quorum.
func (k Keeper) dropCurrentBundleProposal(ctx sdk.Context, poolId uint64, voteDistribution types.VoteDistribution, nextUploader string) {
	pool, _ := k.poolKeeper.GetPool(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalized{
		PoolId:      pool.Id,
		Id:          pool.TotalBundles,
		Valid:       voteDistribution.Valid,
		Invalid:     voteDistribution.Invalid,
		Abstain:     voteDistribution.Abstain,
		Total:       voteDistribution.Total,
		Status:      voteDistribution.Status,
		FinalizedAt: uint64(ctx.BlockTime().Unix()),
		Uploader:    bundleProposal.Uploader,
	})

	// drop bundle
	bundleProposal = types.BundleProposal{
		PoolId:       pool.Id,
		NextUploader: nextUploader,
		UpdatedAt:    uint64(ctx.BlockTime().Unix()),
	}

	k.SetBundleProposal(ctx, bundleProposal)
}

// calculateVotingPower calculates the voting power one staker has in a
// storage pool based only on the total delegation this staker has
func (k Keeper) calculateVotingPower(delegation uint64) (votingPower uint64) {
	// voting power is linear
	votingPower = delegation
	return
}

// chooseNextUploader selects the next uploader based on a fixed set of stakers in a pool.
// It is guaranteed that someone is chosen deterministically if the round-robin set itself is not empty.
func (k Keeper) chooseNextUploader(ctx sdk.Context, poolId uint64, excluded ...string) (nextUploader string) {
	vs := k.LoadRoundRobinValidatorSet(ctx, poolId)
	nextUploader = vs.NextProposer(excluded...)
	k.SaveRoundRobinValidatorSet(ctx, vs)
	return
}

// chooseNextUploader selects the next uploader based on a fixed set of stakers in a pool.
// It is guaranteed that someone is chosen deterministically if the round-robin set itself is not empty.
func (k Keeper) chooseNextUploaderFromList(ctx sdk.Context, poolId uint64, included []string) (nextUploader string) {
	vs := k.LoadRoundRobinValidatorSet(ctx, poolId)

	// Calculate set difference to obtain excluded
	includedMap := make(map[string]bool)
	for _, entry := range included {
		includedMap[entry] = true
	}
	excluded := make([]string, 0)
	for _, entry := range vs.Validators {
		if !includedMap[entry.Address] {
			excluded = append(excluded, entry.Address)
		}
	}

	nextUploader = vs.NextProposer(excluded...)
	k.SaveRoundRobinValidatorSet(ctx, vs)
	return
}

// GetVoteDistribution is an internal function evaluates the quorum status
// based on the voting power of the current bundle proposal.
func (k Keeper) GetVoteDistribution(ctx sdk.Context, poolId uint64) (voteDistribution types.VoteDistribution) {
	bundleProposal, found := k.GetBundleProposal(ctx, poolId)
	if !found {
		return
	}

	// get voting power for valid
	for _, voter := range bundleProposal.VotersValid {
		// valaccount was found the voter is active in the pool
		if k.stakerKeeper.DoesValaccountExist(ctx, poolId, voter) {
			delegation := k.delegationKeeper.GetDelegationAmount(ctx, voter)
			voteDistribution.Valid += k.calculateVotingPower(delegation)
		}
	}

	// get voting power for invalid
	for _, voter := range bundleProposal.VotersInvalid {
		// valaccount was found the voter is active in the pool
		if k.stakerKeeper.DoesValaccountExist(ctx, poolId, voter) {
			delegation := k.delegationKeeper.GetDelegationAmount(ctx, voter)
			voteDistribution.Invalid += k.calculateVotingPower(delegation)
		}
	}

	// get voting power for abstain
	for _, voter := range bundleProposal.VotersAbstain {
		// valaccount was found the voter is active in the pool
		if k.stakerKeeper.DoesValaccountExist(ctx, poolId, voter) {
			delegation := k.delegationKeeper.GetDelegationAmount(ctx, voter)
			voteDistribution.Abstain += k.calculateVotingPower(delegation)
		}
	}

	// get total voting power
	for _, staker := range k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		delegation := k.delegationKeeper.GetDelegationAmount(ctx, staker)
		voteDistribution.Total += k.calculateVotingPower(delegation)
	}

	if voteDistribution.Total == 0 {
		// if total voting power is zero no quorum can be reached
		voteDistribution.Status = types.BUNDLE_STATUS_NO_QUORUM
	} else if voteDistribution.Valid*2 > voteDistribution.Total {
		// if more than 50% voted for valid quorum is reached
		voteDistribution.Status = types.BUNDLE_STATUS_VALID
	} else if voteDistribution.Invalid*2 >= voteDistribution.Total {
		// if more or equal than 50% voted for invalid quorum is reached
		voteDistribution.Status = types.BUNDLE_STATUS_INVALID
	} else {
		// if neither valid nor invalid reached 50% no quorum was reached
		voteDistribution.Status = types.BUNDLE_STATUS_NO_QUORUM
	}

	return
}

// tallyBundleProposal evaluates the votes of a bundle proposal and determines the outcome
func (k Keeper) tallyBundleProposal(ctx sdk.Context, bundleProposal types.BundleProposal, poolId uint64) (types.TallyResult, error) {
	// Increase points of stakers who did not vote at all + slash + remove if there is a bundle proposal.
	// The protocol requires everybody to stay always active.
	if bundleProposal.StorageId != "" {
		k.handleNonVoters(ctx, poolId)
	}

	// evaluate all votes and determine status based on the votes weighted with stake + delegation
	voteDistribution := k.GetVoteDistribution(ctx, poolId)

	// Handle tally outcome
	switch voteDistribution.Status {
	case types.BUNDLE_STATUS_VALID:
		// charge the funders of the pool
		fundersPayout, err := k.fundersKeeper.ChargeFundersOfPool(ctx, poolId, poolTypes.ModuleName)
		if err != nil {
			return types.TallyResult{}, err
		}

		// charge the inflation pool
		inflationPayout, err := k.poolKeeper.ChargeInflationPool(ctx, poolId)
		if err != nil {
			return types.TallyResult{}, err
		}

		// combine funders payout with inflation payout to calculate the rewards for the different stakeholders
		// like treasury, uploader and delegators
		totalPayout := fundersPayout.Add(sdk.NewInt64Coin(globalTypes.Denom, int64(inflationPayout)))
		bundleReward := k.calculatePayouts(ctx, poolId, totalPayout)

		// payout rewards to treasury
		if err := k.distrkeeper.FundCommunityPool(ctx, bundleReward.Treasury, k.accountKeeper.GetModuleAddress(poolTypes.ModuleName)); err != nil {
			return types.TallyResult{}, err
		}

		// payout rewards to uploader through commission rewards
		uploaderReward := bundleReward.UploaderCommission.Add(bundleReward.UploaderStorageCost...)
		if err := k.stakerKeeper.IncreaseStakerCommissionRewards(ctx, bundleProposal.Uploader, poolTypes.ModuleName, uploaderReward); err != nil {
			return types.TallyResult{}, err
		}

		// payout rewards to delegators through delegation rewards
		if err := k.delegationKeeper.PayoutRewards(ctx, bundleProposal.Uploader, bundleReward.Delegation, poolTypes.ModuleName); err != nil {
			return types.TallyResult{}, err
		}

		// slash stakers who voted incorrectly
		for _, voter := range bundleProposal.VotersInvalid {
			k.slashDelegatorsAndRemoveStaker(ctx, poolId, voter, delegationTypes.SLASH_TYPE_VOTE)
		}

		return types.TallyResult{
			Status:           types.TallyResultValid,
			VoteDistribution: voteDistribution,
			FundersPayout:    fundersPayout,
			InflationPayout:  inflationPayout,
			BundleReward:     bundleReward,
		}, nil
	case types.BUNDLE_STATUS_INVALID:
		// If the bundles is invalid, everybody who voted incorrectly gets slashed.
		// The bundle provided by the message-sender is of no mean, because the previous bundle
		// turned out to be incorrect.
		// There this round needs to start again and the message-sender stays uploader.

		// slash stakers who voted incorrectly - uploader receives upload slash
		for _, voter := range bundleProposal.VotersValid {
			if voter == bundleProposal.Uploader {
				k.slashDelegatorsAndRemoveStaker(ctx, poolId, voter, delegationTypes.SLASH_TYPE_UPLOAD)
			} else {
				k.slashDelegatorsAndRemoveStaker(ctx, poolId, voter, delegationTypes.SLASH_TYPE_VOTE)
			}
		}

		return types.TallyResult{
			Status:           types.TallyResultInvalid,
			VoteDistribution: voteDistribution,
			FundersPayout:    sdk.NewCoins(),
			InflationPayout:  0,
			BundleReward:     types.BundleReward{},
		}, nil
	default:
		// If the bundle is neither valid nor invalid the quorum has not been reached yet.
		return types.TallyResult{
			Status:           types.TallyResultNoQuorum,
			VoteDistribution: voteDistribution,
			FundersPayout:    sdk.NewCoins(),
			InflationPayout:  0,
			BundleReward:     types.BundleReward{},
		}, nil
	}
}
