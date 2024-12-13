package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// Stakers
	"github.com/KYVENetwork/chain/x/stakers/types"
)

// These functions are meant to be called from external modules.
// For now this is the bundles module and the delegation module
// which need to interact with the stakers module.

// LeavePool removes a staker from a pool and emits the corresponding event.
// The staker is no longer able to participate in the given pool.
// All points the staker had in that pool are deleted.
func (k Keeper) LeavePool(ctx sdk.Context, staker string, poolId uint64) {
	k.RemoveValaccountFromPool(ctx, poolId, staker)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventLeavePool{
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
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)

	if !found {
		return types.ErrValaccountUnauthorized
	}

	if valaccount.Valaddress != valaddress {
		return types.ErrValaccountUnauthorized
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
		totalDelegation += k.GetDelegationAmount(ctx, address)
	}
	return totalDelegation
}

// GetTotalAndHighestDelegationOfPool iterates all validators of a given pool and returns the stake of the validator
// with the highest stake and the sum of all stakes.
func (k Keeper) GetTotalAndHighestDelegationOfPool(ctx sdk.Context, poolId uint64) (totalDelegation, highestDelegation uint64) {
	for _, address := range k.GetAllStakerAddressesOfPool(ctx, poolId) {
		delegation := k.GetDelegationAmount(ctx, address)
		totalDelegation += delegation

		if delegation > highestDelegation {
			highestDelegation = delegation
		}
	}

	return totalDelegation, highestDelegation
}

// GetDelegationAmount returns the stake of a given validator in ukyve
func (k Keeper) GetDelegationAmount(ctx sdk.Context, validator string) uint64 {
	validatorAddress, err := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(validator))
	if err != nil {
		return 0
	}
	chainValidator, err := k.stakingKeeper.GetValidator(ctx, validatorAddress)
	if err != nil {
		return 0
	}

	return chainValidator.GetBondedTokens().Uint64()
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
func (k Keeper) Slash(ctx sdk.Context, poolId uint64, staker string, slashFraction math.LegacyDec) {
	validator, found := k.GetValidator(ctx, staker)
	if !found {
		return
	}

	consAddrBytes, _ := validator.GetConsAddr()

	amount, err := k.stakingKeeper.Slash(ctx, consAddrBytes, ctx.BlockHeight(), validator.GetConsensusPower(math.NewInt(1000000)), slashFraction)
	if err != nil {
		return
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventSlash{
		PoolId: poolId,
		Staker: staker,
		Amount: amount.Uint64(),
		// SlashType: slashType, TODO add slash type, once migrated away from delegation module
	})
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

	err := k.distKeeper.AllocateTokensToValidator(ctx, validator, sdk.NewDecCoinsFromCoins(amount...))
	if err != nil {
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
