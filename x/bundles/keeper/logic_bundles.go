package keeper

import (
	"encoding/binary"
	"math/rand"
	"sort"

	"cosmossdk.io/errors"

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

	// Error if the pool has no funds.
	if len(pool.Funders) == 0 {
		return types.ErrPoolOutOfFunds
	}

	// Error if min delegation is not reached
	if k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id) < pool.MinDelegation {
		return types.ErrMinDelegationNotReached
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

// calculatePayouts deducts the network fee from the rewards and splits the remaining amount
// between the staker and its delegators. If there are no delegators, the entire amount is
// awarded to the staker.
func (k Keeper) calculatePayouts(ctx sdk.Context, poolId uint64) (bundleReward types.BundleReward) {
	pool, _ := k.poolKeeper.GetPoolWithError(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// Should not happen, if so move everything to the treasury
	if !k.stakerKeeper.DoesStakerExist(ctx, bundleProposal.Uploader) {
		bundleReward.Treasury = bundleReward.Total

		return
	}

	// formula for calculating the rewards
	bundleReward.Total = pool.OperatingCost + uint64(k.GetStorageCost(ctx).MulInt64(int64(bundleProposal.DataSize)).TruncateInt64())

	// Add fee to treasury
	bundleReward.Treasury = uint64(sdk.NewDec(int64(bundleReward.Total)).Mul(k.GetNetworkFee(ctx)).TruncateInt64())

	// Remaining rewards to be split between staker and its delegators
	totalNodeReward := bundleReward.Total - bundleReward.Treasury

	// Payout delegators
	if k.delegationKeeper.GetDelegationAmount(ctx, bundleProposal.Uploader) > 0 {
		commission := k.stakerKeeper.GetCommission(ctx, bundleProposal.Uploader)

		bundleReward.Uploader = uint64(sdk.NewDec(int64(totalNodeReward)).Mul(commission).TruncateInt64())
		bundleReward.Delegation = totalNodeReward - bundleReward.Uploader
	} else {
		bundleReward.Uploader = totalNodeReward
		bundleReward.Delegation = 0
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

	// Emit a vote event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleVote{
		PoolId:    msg.PoolId,
		Staker:    msg.Staker,
		StorageId: msg.StorageId,
		Vote:      types.VOTE_TYPE_VALID,
	})

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
}

// finalizeCurrentBundleProposal takes the data of the current evaluated proposal
// and stores it as a finalized proposal. This only happens if the network
// reached quorum on the proposal's validity.
func (k Keeper) finalizeCurrentBundleProposal(ctx sdk.Context, poolId uint64, voteDistribution types.VoteDistribution, bundleReward types.BundleReward, nextUploader string) {
	pool, _ := k.poolKeeper.GetPool(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// save finalized bundle
	finalizedBundle := types.FinalizedBundle{
		StorageId:         bundleProposal.StorageId,
		PoolId:            pool.Id,
		Id:                pool.TotalBundles,
		Uploader:          bundleProposal.Uploader,
		FromIndex:         pool.CurrentIndex,
		ToIndex:           pool.CurrentIndex + bundleProposal.BundleSize,
		FinalizedAt:       uint64(ctx.BlockHeight()),
		FromKey:           bundleProposal.FromKey,
		ToKey:             bundleProposal.ToKey,
		BundleSummary:     bundleProposal.BundleSummary,
		DataHash:          bundleProposal.DataHash,
		StorageProviderId: bundleProposal.StorageProviderId,
		CompressionId:     bundleProposal.CompressionId,
	}

	k.SetFinalizedBundle(ctx, finalizedBundle)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalized{
		PoolId:           finalizedBundle.PoolId,
		Id:               finalizedBundle.Id,
		Valid:            voteDistribution.Valid,
		Invalid:          voteDistribution.Invalid,
		Abstain:          voteDistribution.Abstain,
		Total:            voteDistribution.Total,
		Status:           voteDistribution.Status,
		RewardTreasury:   bundleReward.Treasury,
		RewardUploader:   bundleReward.Uploader,
		RewardDelegation: bundleReward.Delegation,
		RewardTotal:      bundleReward.Total,
		FinalizedAt:      uint64(ctx.BlockTime().Unix()),
		Uploader:         bundleProposal.Uploader,
		NextUploader:     nextUploader,
	})

	// Finalize the proposal, saving useful information.
	k.poolKeeper.IncrementBundleInformation(ctx, pool.Id, pool.CurrentIndex+bundleProposal.BundleSize, bundleProposal.ToKey, bundleProposal.BundleSummary)
}

// dropCurrentBundleProposal removes the current proposal due to not reaching
// a required quorum on the validity of the data. When the proposal is dropped
// the same next uploader as before can submit his proposal since it is not his
// fault, that the last one did not reach any quorum.
func (k Keeper) dropCurrentBundleProposal(
	ctx sdk.Context,
	poolId uint64,
	voteDistribution types.VoteDistribution,
	nextUploader string,
) {
	pool, _ := k.poolKeeper.GetPool(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalized{
		PoolId:           pool.Id,
		Id:               pool.TotalBundles,
		Valid:            voteDistribution.Valid,
		Invalid:          voteDistribution.Invalid,
		Abstain:          voteDistribution.Abstain,
		Total:            voteDistribution.Total,
		Status:           voteDistribution.Status,
		RewardTreasury:   0,
		RewardUploader:   0,
		RewardDelegation: 0,
		RewardTotal:      0,
		FinalizedAt:      uint64(ctx.BlockTime().Unix()),
		Uploader:         bundleProposal.Uploader,
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

// RandomChoiceCandidate holds the voting power of a candidate for the
// next uploader selection
type RandomChoiceCandidate struct {
	Account     string
	VotingPower uint64
}

// getWeightedRandomChoice is an internal function that returns a weighted random
// selection out of a list of candidates based on their voting power.
func (k Keeper) getWeightedRandomChoice(candidates []RandomChoiceCandidate, seed int64) string {
	type WeightedRandomChoice struct {
		Elements    []string
		Weights     []uint64
		TotalWeight uint64
	}

	wrc := WeightedRandomChoice{}

	for _, candidate := range candidates {
		i := sort.Search(len(wrc.Weights), func(i int) bool { return wrc.Weights[i] > candidate.VotingPower })
		wrc.Weights = append(wrc.Weights, 0)
		wrc.Elements = append(wrc.Elements, "")
		copy(wrc.Weights[i+1:], wrc.Weights[i:])
		copy(wrc.Elements[i+1:], wrc.Elements[i:])
		wrc.Weights[i] = candidate.VotingPower
		wrc.Elements[i] = candidate.Account
		wrc.TotalWeight += candidate.VotingPower
	}

	if wrc.TotalWeight == 0 {
		return ""
	}

	value := rand.New(rand.NewSource(seed)).Uint64() % wrc.TotalWeight

	for key, weight := range wrc.Weights {
		if weight > value {
			return wrc.Elements[key]
		}

		value -= weight
	}

	return ""
}

// chooseNextUploaderFromSelectedStakers selects the next uploader based on a
// fixed set of stakers in a pool. It is guaranteed that someone is chosen
// deterministically
func (k Keeper) chooseNextUploaderFromSelectedStakers(ctx sdk.Context, poolId uint64, addresses []string) (nextUploader string) {
	var _candidates []RandomChoiceCandidate

	if len(addresses) == 0 {
		return ""
	}

	for _, s := range addresses {
		if k.stakerKeeper.DoesValaccountExist(ctx, poolId, s) {
			delegation := k.delegationKeeper.GetDelegationAmount(ctx, s)

			_candidates = append(_candidates, RandomChoiceCandidate{
				Account:     s,
				VotingPower: k.calculateVotingPower(delegation),
			})
		}
	}

	seed := int64(binary.BigEndian.Uint64(ctx.BlockHeader().AppHash))
	return k.getWeightedRandomChoice(_candidates, seed)
}

// chooseNextUploaderFromAllStakers selects the next uploader based on all
// stakers in a pool. It is guaranteed that someone is chosen
// deterministically
func (k Keeper) chooseNextUploaderFromAllStakers(ctx sdk.Context, poolId uint64) (nextUploader string) {
	stakers := k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, poolId)
	return k.chooseNextUploaderFromSelectedStakers(ctx, poolId, stakers)
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
