package keeper

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// Stakers
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
)

// These functions are meant to be called from external modules.
// For now this is the bundles module and the delegation module
// which need to interact with the stakers module.

// LeavePool removes a staker from a pool and emits the corresponding event.
// The staker is no longer able to participate in the given pool.
// All points the staker had in that pool are deleted.
func (k Keeper) LeavePool(ctx sdk.Context, staker string, poolId uint64) {
	k.RemoveValaccountFromPool(ctx, poolId, staker)

	_ = ctx.EventManager().EmitTypedEvent(&stakertypes.EventLeavePool{
		PoolId: poolId,
		Staker: staker,
	})
}

// GetAllStakerAddressesOfPool returns a list of all stakers
// which have currently a valaccount registered for the given pool
// and are therefore allowed to participate in that pool.
func (k Keeper) GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string) {
	for _, valaccount := range k.GetAllValaccountsOfPool(ctx, poolId) {
		stakers = append(stakers, valaccount.Staker)
	}

	return stakers
}

func (k Keeper) GetPaginatedStakersByPoolStake(ctx sdk.Context, pagination *query.PageRequest, accumulator func(validator string, accumulate bool) bool) (*query.PageResponse, error) {
	validators, err := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: sort by pool stake
	addresses := make([]string, 0)

	for _, validator := range validators {
		addresses = append(addresses, util.MustAccountAddressFromValAddress(validator.OperatorAddress))
	}

	pageRes, err := arrayPaginationAccumulator(addresses, pagination, accumulator)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return pageRes, nil
}

// AssertValaccountAuthorized checks if the given `valaddress` is allowed to vote in pool
// with id `poolId` to vote in favor of `stakerAddress`.
// If the valaddress is not authorized the appropriate error is returned.
// Otherwise, it returns `nil`
func (k Keeper) AssertValaccountAuthorized(ctx sdk.Context, poolId uint64, stakerAddress string, valaddress string) error {
	valaccount, active := k.GetValaccount(ctx, poolId, stakerAddress)
	if !active {
		return stakertypes.ErrValaccountUnauthorized
	}

	if valaccount.Valaddress != valaddress {
		return stakertypes.ErrValaccountUnauthorized
	}

	return nil
}

// GetActiveStakers returns all staker-addresses that are
// currently participating in at least one pool.
func (k Keeper) GetActiveStakers(ctx sdk.Context) []string {
	return k.getAllActiveStakers(ctx)
}

// GetTotalStakeOfPool returns the amount in uykve which actively secures
// the given pool
func (k Keeper) GetTotalStakeOfPool(ctx sdk.Context, poolId uint64) (totalStake uint64) {
	effectiveStakes := k.GetValidatorPoolStakes(ctx, poolId)
	for _, stake := range effectiveStakes {
		totalStake += stake
	}
	return
}

// IsVotingPowerTooHigh returns whether there are enough validators in a pool
// to successfully stay below the max voting power
func (k Keeper) IsVotingPowerTooHigh(ctx sdk.Context, poolId uint64) bool {
	addresses := int64(len(k.GetAllStakerAddressesOfPool(ctx, poolId)))
	maxVotingPower := k.poolKeeper.GetMaxVotingPowerPerPool(ctx)

	if maxVotingPower.IsZero() {
		return true
	}

	return math.LegacyOneDec().Quo(maxVotingPower).GT(math.LegacyNewDec(addresses))
}

// GetValidator returns the Cosmos-validator for a given kyve-address.
func (k Keeper) GetValidator(ctx sdk.Context, staker string) (stakingTypes.Validator, bool) {
	valAddress, err := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(staker))
	if err != nil {
		return stakingTypes.Validator{}, false
	}
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return stakingTypes.Validator{}, false
	}

	return validator, true
}

// GetValidatorPoolCommission returns the commission a validator has inside the pool
func (k Keeper) GetValidatorPoolCommission(ctx sdk.Context, staker string, poolId uint64) math.LegacyDec {
	valaccount, _ := k.GetValaccount(ctx, poolId, staker)
	return valaccount.Commission
}

// GetValidatorPoolStake returns stake a validator has actively and at risk inside the pool
func (k Keeper) GetValidatorPoolStake(ctx sdk.Context, staker string, poolId uint64) uint64 {
	return k.GetValidatorPoolStakes(ctx, poolId, staker)[staker]
}

// GetValidatorPoolStakes returns a map for all pool validators with their effective stake. Effective stake
// is the actual amount which determines the validator's voting power and is the actual amount at risk for
// slashing. The effective stake can be lower (never higher) than the specified stake by the validators by
// his stake fraction because of the maximum voting power which limits the voting power a validator can have
// in a pool.
//
// We limit the voting power by first sorting all validators based on their stake and start at the highest. From
// there we check if this validator exceeds the voting power with his stake, if he does we cut off an amount
// from his stake and redistribute it to all the validators below him based on their stakes. Because we are simply
// redistributing exact amounts we do not need to worry about changing the total stake which would mess up voting
// powers of other validators again. After this has been repeated for every validator that exceeds the max voting
// power we finally scale down all stakes to the lowest validator's stake. This is because the top validators lost
// stake and the lowest validators gained stake, but because it is not allowed to risk more than specified so we scale
// everything down accordingly. This results in a stake distribution where every validator who was below the max
// voting power has the same effective stake as his dedicated stake and those validators who where above it have
// less effective stake.
func (k Keeper) GetValidatorPoolStakes(ctx sdk.Context, poolId uint64, mustIncludeStakers ...string) map[string]uint64 {
	type ValidatorStake struct {
		Address string
		Stake   uint64
	}

	validators := make([]ValidatorStake, 0)
	stakes := make(map[string]uint64)

	// we include given stakers since in some instances the staker we are looking for has been already kicked
	// out of a pool, but we still need to payout rewards or slash him afterward depending on his last action
	// right before leaving the pool
	addresses := util.RemoveDuplicateStrings(append(k.GetAllStakerAddressesOfPool(ctx, poolId), mustIncludeStakers...))
	maxVotingPower := k.poolKeeper.GetMaxVotingPowerPerPool(ctx)

	// it is impossible regardless how many validators are in a pool to have a max voting power of 0%,
	// therefore we return here
	if maxVotingPower.IsZero() {
		return stakes
	}

	// if there are not enough validators in a pool so that the max voting power is always
	// exceeded by at least one validator we return
	if math.LegacyOneDec().Quo(maxVotingPower).GT(math.LegacyNewDec(int64(len(addresses)))) {
		return stakes
	}

	totalStake := int64(0)

	for _, address := range addresses {
		validator, _ := k.GetValidator(ctx, address)
		valaccount, _ := k.GetValaccount(ctx, poolId, address)

		// calculate the stake the validator has specifically chosen for this pool
		// with his stake fraction
		stake := uint64(valaccount.StakeFraction.MulInt(validator.GetBondedTokens()).TruncateInt64())

		stakes[address] = stake
		validators = append(validators, ValidatorStake{
			Address: address,
			Stake:   stake,
		})
		totalStake += int64(stake)
	}

	totalStakeRemainder := totalStake

	// sort descending based on stake
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].Stake > validators[j].Stake
	})

	// if the total stake of the pool is zero we return
	if totalStake == 0 {
		return stakes
	}

	for i, validator := range validators {
		// check if the validator has a higher stake than allowed by the max voting power
		if math.LegacyNewDec(int64(stakes[validator.Address])).GT(maxVotingPower.MulInt64(totalStake)) {
			// if the validator got a stake which would give him a higher voting power than the maximum allowed
			// one we cut off the exact amount from his stake so that he is just below the max voting power
			cutoffAmount := math.LegacyNewDec(int64(stakes[validator.Address])).Sub(maxVotingPower.MulInt64(totalStake)).TruncateInt64()

			totalStakeRemainder -= int64(validator.Stake)
			stakes[validator.Address] -= uint64(cutoffAmount)

			// we take the cutoff amount and distribute it on the remaining validators down the list
			// who all have less voting power than the current one. We distribute the cutoff amount
			// based on the validator's stake
			if totalStakeRemainder > 0 {
				for _, v := range validators[i+1:] {
					stakes[v.Address] += uint64(math.LegacyNewDec(int64(v.Stake)).QuoInt64(totalStakeRemainder).MulInt64(cutoffAmount).TruncateInt64())
				}
			}
		}
	}

	// after we have redistributed all cutoff amounts so that no validator exceeds the maximum voting power
	// based on their remaining effective stake we now scale the stakes to get the true effective staking amount.
	// This is because while the top validators who got their voting power reduced the lower validators have actually
	// gained voting power relatively. But because their effective stake is now bigger than their allocated stake
	// we have to scale it down again, because we are not allowed to slash more than the validator has allocated in
	// this particular pool. Therefore, we take the stake of the lowest validator (with a bigger stake than 0) and scale
	// down all stakes of all other validators to that accordingly
	scaleFactor := math.LegacyZeroDec()

	// get the lowest validator with effective stake still bigger than zero and determine the scale factor
	for i := len(validators) - 1; i >= 0; i-- {
		if stakes[validators[i].Address] > 0 {
			scaleFactor = math.LegacyNewDec(int64(validators[i].Stake)).QuoInt64(int64(stakes[validators[i].Address]))
			break
		}
	}

	// scale all effective stakes down to scale factor
	for _, validator := range validators {
		stakes[validator.Address] = uint64(scaleFactor.MulInt64(int64(stakes[validator.Address])).RoundInt64())
	}

	// the result is a map which contains the effective stake for every validator in a pool. The effective stake
	// can not be higher than the dedicated stake specified by the validator by his stake fraction, therefore it
	// represents the true amount of $KYVE which is at risk for slashing
	return stakes
}

// GetOutstandingCommissionRewards returns the outstanding commission rewards for a given validator
func (k Keeper) GetOutstandingCommissionRewards(ctx sdk.Context, staker string) sdk.Coins {
	valAddress, err := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(staker))
	if err != nil {
		return sdk.NewCoins()
	}

	commission, err := k.distKeeper.GetValidatorAccumulatedCommission(ctx, valAddress)
	if err != nil {
		return sdk.NewCoins()
	}

	return util.TruncateDecCoins(sdk.NewDecCoins(commission.Commission...))
}

// GetOutstandingRewards returns the outstanding delegation rewards for a given delegator
func (k Keeper) GetOutstandingRewards(orgCtx sdk.Context, staker string, delegator string) sdk.Coins {
	valAdr, err := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(staker))
	if err != nil {
		return sdk.NewCoins()
	}

	delAdr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdk.NewCoins()
	}

	ctx, _ := orgCtx.CacheContext()

	del, err := k.stakingKeeper.Delegation(ctx, delAdr, valAdr)
	if err != nil {
		return sdk.NewCoins()
	}

	val, found := k.GetValidator(ctx, staker)
	if !found {
		return sdk.NewCoins()
	}

	endingPeriod, err := k.distKeeper.IncrementValidatorPeriod(ctx, val)
	if err != nil {
		return sdk.NewCoins()
	}

	rewards, err := k.distKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
	if err != nil {
		return sdk.NewCoins()
	}

	return util.TruncateDecCoins(rewards)
}

// Slash reduces the delegation of all delegators of `staker` by fraction. The slash itself is handled by the cosmos-sdk
func (k Keeper) Slash(ctx sdk.Context, poolId uint64, staker string, slashType stakertypes.SlashType) {
	validator, found := k.GetValidator(ctx, staker)
	if !found {
		return
	}

	consAddrBytes, _ := validator.GetConsAddr()

	// here the stake fraction can be actually different from the stake fraction the validator has specified
	// for this pool because with the max voting power in place his effective stake in a pool could have been
	// reduced because his original stake was too high. Therefore, we determine the true stake fraction
	// by dividing his effective stake with his bonded amount.
	stakeFraction := math.LegacyZeroDec()

	if !validator.GetBondedTokens().IsZero() {
		stakeFraction = math.LegacyNewDec(int64(k.GetValidatorPoolStake(ctx, staker, poolId))).QuoInt(validator.GetBondedTokens())
	}

	// the validator can only be slashed for his effective stake fraction in a pool, therefore we
	// update the slash fraction accordingly
	slashFraction := k.getSlashFraction(ctx, slashType).Mul(stakeFraction)

	amount, err := k.stakingKeeper.Slash(
		ctx,
		consAddrBytes,
		ctx.BlockHeight(),
		validator.GetConsensusPower(math.NewInt(1000000)),
		slashFraction,
	)
	if err != nil {
		return
	}

	_ = ctx.EventManager().EmitTypedEvent(&stakertypes.EventSlash{
		PoolId:        poolId,
		Staker:        staker,
		Amount:        amount.Uint64(),
		SlashType:     slashType,
		StakeFraction: stakeFraction,
	})
}

// GOVERNANCE - BONDING
// The next functions are used in our custom fork of the cosmos-sdk
// which includes protocol staking into the governance.
// The behavior is exactly the same as with normal cosmos-validators.

// TotalBondedTokens returns all tokens which are currently bonded by the protocol
// I.e. the sum of all delegation of all stakers that are currently participating
// in at least one pool
func (k Keeper) TotalBondedTokens(ctx context.Context) math.Int {
	return math.ZeroInt()
}

// GetActiveValidators returns all protocol-node information which
// are needed by the governance to calculate the voting powers.
// The interface needs to correspond to github.com/cosmos/cosmos-sdk/x/gov/types/v1.ValidatorGovInfo
// But as there is no direct dependency in the cosmos-sdk-fork this value is passed as an interface{}
func (k Keeper) GetActiveValidators(ctx context.Context) (validators []interface{}) {
	return
}

// GetDelegations returns the address and the delegation amount of all active protocol-stakers the
// delegator as delegated to. This is used to calculate the vote weight each delegator has.
func (k Keeper) GetDelegations(ctx context.Context, delegator string) (validators []string, amounts []math.LegacyDec) {
	return
}

func (k Keeper) DoesStakerExist(ctx sdk.Context, staker string) bool {
	// ToDo remove after Delegation module got deleted
	_, found := k.GetValidator(ctx, staker)
	return found
}

// PayoutRewards transfers `amount` from the `payerModuleName`-module to the delegation module.
// It then awards these tokens internally to all delegators of staker `staker`.
// Delegators can then receive these rewards if they call the `withdraw`-transaction.
// If the staker has no delegators or the module to module transfer fails, the method fails and
// returns an error.
func (k Keeper) PayoutRewards(ctx sdk.Context, staker string, amount sdk.Coins, payerModuleName string) error {
	// Assert there is an amount
	if amount.Empty() {
		return nil
	}

	validator, _ := k.GetValidator(ctx, staker)
	valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
	if err != nil {
		return err
	}

	currentRewards, err := k.distKeeper.GetValidatorCurrentRewards(ctx, valBz)
	if err != nil {
		return err
	}

	currentRewards.Rewards = currentRewards.Rewards.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	if err := k.distKeeper.SetValidatorCurrentRewards(ctx, valBz, currentRewards); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			distrtypes.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(distrtypes.AttributeKeyValidator, validator.GetOperator()),
		),
	)

	outstanding, err := k.distKeeper.GetValidatorOutstandingRewards(ctx, valBz)
	if err != nil {
		return err
	}

	outstanding.Rewards = outstanding.Rewards.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	if err := k.distKeeper.SetValidatorOutstandingRewards(ctx, valBz, outstanding); err != nil {
		return err
	}

	// Transfer tokens to the delegation module
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, payerModuleName, distrtypes.ModuleName, amount); err != nil {
		return err
	}

	return nil
}

// PayoutAdditionalCommissionRewards pays out some additional tokens to the validator.
func (k Keeper) PayoutAdditionalCommissionRewards(ctx sdk.Context, validator string, payerModuleName string, amount sdk.Coins) error {
	// Assert there is an amount
	if amount.Empty() {
		return nil
	}

	// Assert the staker exists
	if _, found := k.GetValidator(ctx, validator); !found {
		return errors.Wrapf(sdkErrors.ErrNotFound, "staker does not exist")
	}

	// transfer funds from pool to distribution module
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, payerModuleName, distrtypes.ModuleName, amount); err != nil {
		return err
	}

	valAcc, err := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(validator))
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			distrtypes.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(distrtypes.AttributeKeyValidator, util.MustValaddressFromOperatorAddress(validator)),
		),
	)

	currentCommission, err := k.distKeeper.GetValidatorAccumulatedCommission(ctx, valAcc)
	if err != nil {
		return err
	}

	currentCommission.Commission = currentCommission.Commission.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	err = k.distKeeper.SetValidatorAccumulatedCommission(ctx, valAcc, currentCommission)
	if err != nil {
		return err
	}

	outstanding, err := k.distKeeper.GetValidatorOutstandingRewards(ctx, valAcc)
	if err != nil {
		return err
	}

	outstanding.Rewards = outstanding.Rewards.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	err = k.distKeeper.SetValidatorOutstandingRewards(ctx, valAcc, outstanding)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetDelegationAmountOfDelegator(ctx sdk.Context, validator, delegator string) uint64 {
	address, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		panic(err)
	}

	valAddress, err := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(validator))
	if err != nil {
		panic(err)
	}

	delegation, err := k.stakingKeeper.Delegation(ctx, address, valAddress)
	if err != nil {
		return 0
	}

	val, _ := k.stakingKeeper.GetValidator(ctx, valAddress)

	return uint64(val.TokensFromSharesTruncated(delegation.GetShares()).RoundInt64())
}
