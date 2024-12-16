package keeper

import (
	"context"
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

func (k Keeper) GetPaginatedStakersByDelegation(ctx sdk.Context, pagination *query.PageRequest, accumulator func(staker string, accumulate bool) bool) (*query.PageResponse, error) {
	validators, err := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		return nil, err
	}

	for _, validator := range validators {
		accumulator(util.MustAccountAddressFromValAddress(validator.OperatorAddress), true)
	}

	res := query.PageResponse{}
	return &res, nil
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

// GetDelegationOfPool returns the amount of how many ukyve users have delegated
// to stakers that are participating in the given pool
func (k Keeper) GetDelegationOfPool(ctx sdk.Context, poolId uint64) uint64 {
	totalDelegation := uint64(0)
	for _, address := range k.GetAllStakerAddressesOfPool(ctx, poolId) {
		totalDelegation += k.GetValidatorPoolStake(ctx, address, poolId)
	}
	return totalDelegation
}

// GetTotalAndHighestDelegationOfPool iterates all validators of a given pool and returns the stake of the validator
// with the highest stake and the sum of all stakes.
func (k Keeper) GetTotalAndHighestDelegationOfPool(ctx sdk.Context, poolId uint64) (totalDelegation, highestDelegation uint64) {
	for _, address := range k.GetAllStakerAddressesOfPool(ctx, poolId) {
		delegation := k.GetValidatorPoolStake(ctx, address, poolId)
		totalDelegation += delegation

		if delegation > highestDelegation {
			highestDelegation = delegation
		}
	}

	return totalDelegation, highestDelegation
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

// GetValidatorPoolStakeFraction returns the stake fraction a validator has inside the pool
func (k Keeper) GetValidatorPoolStakeFraction(ctx sdk.Context, staker string, poolId uint64) math.LegacyDec {
	valaccount, _ := k.GetValaccount(ctx, poolId, staker)
	return valaccount.StakeFraction
}

// GetValidatorPoolStake returns stake a validator has inside the pool
func (k Keeper) GetValidatorPoolStake(ctx sdk.Context, staker string, poolId uint64) uint64 {
	validator, found := k.GetValidator(ctx, staker)
	if !found {
		return 0
	}

	stakeFraction := k.GetValidatorPoolStakeFraction(ctx, staker, poolId)
	return uint64(math.LegacyNewDecFromInt(validator.BondedTokens()).Mul(stakeFraction).TruncateInt64())
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

// GetEffectiveValidatorStakes returns a map for all pool stakers their effective stake. Effective stake
// is the actual stake at risk for a validator, it can be lower than the provided stake because
// of the max voting power per pool parameter, but it can never be higher. Thus, if the effective
// stake is lower than the provided one the slash
// TODO: explain or link to formula
func (k Keeper) GetEffectiveValidatorStakes(ctx sdk.Context, poolId uint64, mustIncludeStakers ...string) (effectiveStakes map[string]uint64) {
	type ValidatorStake struct {
		Address string
		Stake   uint64
	}

	addresses := util.RemoveDuplicateStrings(append(k.GetAllStakerAddressesOfPool(ctx, poolId), mustIncludeStakers...))
	maxVotingPower := k.poolKeeper.GetMaxVotingPowerPerPool(ctx)
	effectiveStakes = make(map[string]uint64)
	validators := make([]ValidatorStake, 0)
	totalStake := int64(0)

	for _, address := range addresses {
		stake := k.GetValidatorPoolStake(ctx, address, poolId)
		effectiveStakes[address] = stake
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

	if totalStake == 0 || maxVotingPower.IsZero() {
		return
	}

	for i, validator := range validators {
		if math.LegacyNewDec(int64(effectiveStakes[validator.Address])).QuoInt64(totalStake).GT(maxVotingPower) {
			amount := math.LegacyNewDec(int64(effectiveStakes[validator.Address])).Sub(maxVotingPower.MulInt64(totalStake)).TruncateInt64()

			totalStakeRemainder -= int64(validator.Stake)
			effectiveStakes[validator.Address] -= uint64(amount)

			if totalStakeRemainder > 0 {
				for _, v := range validators[i+1:] {
					effectiveStakes[v.Address] += uint64(math.LegacyNewDec(int64(effectiveStakes[v.Address])).QuoInt64(totalStakeRemainder).MulInt64(amount).TruncateInt64())
				}
			}
		}
	}

	scaleFactor := math.LegacyZeroDec()

	// get lowest validator with effective stake still bigger than zero
	for i := len(validators) - 1; i >= 0; i-- {
		if effectiveStakes[validators[i].Address] > 0 {
			scaleFactor = math.LegacyNewDec(int64(validators[i].Stake)).QuoInt64(int64(effectiveStakes[validators[i].Address]))
		}
	}

	// scale all effective stakes down to scale factor
	for _, validator := range validators {
		effectiveStakes[validator.Address] = uint64(scaleFactor.MulInt64(int64(effectiveStakes[validator.Address])).TruncateInt64())
	}

	return
}

// Slash reduces the delegation of all delegators of `staker` by fraction. The slash itself is handled by the cosmos-sdk
func (k Keeper) Slash(ctx sdk.Context, poolId uint64, staker string, slashType stakertypes.SlashType) {
	validator, found := k.GetValidator(ctx, staker)
	if !found {
		return
	}

	consAddrBytes, _ := validator.GetConsAddr()

	// get real stake fraction from the effective stake which is truly at risk for getting slashed
	// by dividing it with the total bonded stake from the validator
	stakeFraction := math.LegacyZeroDec()

	if !validator.GetBondedTokens().IsZero() {
		stakeFraction = math.LegacyNewDec(int64(k.GetEffectiveValidatorStakes(ctx, poolId, staker)[staker])).QuoInt(validator.GetBondedTokens())
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
