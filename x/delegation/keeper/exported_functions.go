package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// GetOutstandingRewards calculates the current rewards a delegator has collected for
// the given staker.
func (k Keeper) GetOutstandingRewards(ctx sdk.Context, staker string, delegator string) sdk.Coins {
	return k.f1GetOutstandingRewards(ctx, staker, delegator)
}
