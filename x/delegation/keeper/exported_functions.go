package keeper

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// These functions are meant to be called from external modules
// For now this is the bundles module which needs to interact
// with the delegation module.
// All these functions are safe in the way that they do not return errors
// and every edge case is handled within the function itself.

// GetDelegationAmount returns the sum of all delegations for a specific staker.
// If the staker does not exist, it returns zero as the staker has zero delegations
func (k Keeper) GetDelegationAmount(ctx sdk.Context, staker string) uint64 {
	delegationData, found := k.GetDelegationData(ctx, staker)

	if found {
		return delegationData.TotalDelegation
	}

	return 0
}

// GetDelegationAmountOfDelegator returns the amount of how many $KYVE `delegatorAddress`
// has delegated to `stakerAddress`. If one of the addresses does not exist, it returns zero.
func (k Keeper) GetDelegationAmountOfDelegator(ctx sdk.Context, stakerAddress string, delegatorAddress string) uint64 {
	return k.f1GetCurrentDelegation(ctx, stakerAddress, delegatorAddress)
}

// GetDelegationOfPool returns the amount of how many $KYVE users have delegated
// to stakers that are participating in the given pool
func (k Keeper) GetDelegationOfPool(ctx sdk.Context, poolId uint64) uint64 {
	totalDelegation := uint64(0)
	for _, address := range k.stakersKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		totalDelegation += k.GetDelegationAmount(ctx, address)
	}
	return totalDelegation
}

// GetTotalAndHighestDelegationOfPool returns the total delegation amount of all validators in the given pool and
// the highest total delegation amount of a single validator in a pool
func (k Keeper) GetTotalAndHighestDelegationOfPool(ctx sdk.Context, poolId uint64) (totalDelegation, highestDelegation uint64) {
	// Get the total delegation and the highest delegation of a staker in the pool
	for _, address := range k.stakersKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		delegation := k.GetDelegationAmount(ctx, address)
		totalDelegation += delegation

		if delegation > highestDelegation {
			highestDelegation = delegation
		}
	}

	return
}

// PayoutRewards transfers `amount` $nKYVE from the `payerModuleName`-module to the delegation module.
// It then awards these tokens internally to all delegators of staker `staker`.
// Delegators can then receive these rewards if they call the `withdraw`-transaction.
// If the staker has no delegators or the module to module transfer fails the method fails and
// returns the error.
func (k Keeper) PayoutRewards(ctx sdk.Context, staker string, amount uint64, payerModuleName string) error {
	// Assert there is an amount
	if amount == 0 {
		return nil
	}

	// Assert there are delegators
	if !k.DoesDelegationDataExist(ctx, staker) {
		return errors.Wrapf(sdkErrors.ErrLogic, "Staker has no delegators which can be paid out.")
	}

	// Add amount to the rewards pool
	k.AddAmountToDelegationRewards(ctx, staker, amount)

	// Transfer tokens to the delegation module
	if err := util.TransferFromModuleToModule(k.bankKeeper, ctx, payerModuleName, types.ModuleName, amount); err != nil {
		return err
	}

	return nil
}

// SlashDelegators reduces the delegation of all delegators of `staker` by fraction
// and transfers the amount to the Treasury.
func (k Keeper) SlashDelegators(ctx sdk.Context, poolId uint64, staker string, slashType types.SlashType) {
	// Only slash if staker has delegators
	if k.DoesDelegationDataExist(ctx, staker) {

		// Update in-memory staker index for efficient queries
		k.RemoveStakerIndex(ctx, staker)
		defer k.SetStakerIndex(ctx, staker)

		// Perform F1-slash and get slashed amount in ukyve
		slashedAmount := k.f1Slash(ctx, staker, k.getSlashFraction(ctx, slashType))

		// Transfer tokens to the Treasury
		if err := util.TransferFromModuleToTreasury(k.accountKeeper, k.distrKeeper, ctx, types.ModuleName, slashedAmount); err != nil {
			util.PanicHalt(k.upgradeKeeper, ctx, "Not enough tokens in module")
		}

		// Emit slash event
		_ = ctx.EventManager().EmitTypedEvent(&types.EventSlash{
			PoolId:    poolId,
			Staker:    staker,
			Amount:    slashedAmount,
			SlashType: slashType,
		})
	}
}

// GetOutstandingRewards calculates the current rewards a delegator has collected for
// the given staker.
func (k Keeper) GetOutstandingRewards(ctx sdk.Context, staker string, delegator string) uint64 {
	return k.f1GetOutstandingRewards(ctx, staker, delegator)
}
