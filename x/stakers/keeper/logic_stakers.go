package keeper

import (
	"math"

	"github.com/KYVENetwork/chain/util"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/KYVENetwork/chain/x/stakers/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// getLowestStaker returns the staker with the lowest total stake
// (self-delegation + delegation) of a given pool.
// If all pool slots are taken, this is the staker who then
// gets kicked out.
func (k Keeper) getLowestStaker(ctx sdk.Context, poolId uint64) (val stakingTypes.Validator, found bool) {
	var minAmount uint64 = math.MaxUint64

	for _, staker := range k.getAllStakersOfPool(ctx, poolId) {
		delegationAmount := k.GetDelegationAmount(ctx, util.MustAccountAddressFromValAddress(staker.OperatorAddress))
		if delegationAmount < minAmount {
			minAmount = delegationAmount
			val = staker
		}
	}

	return
}

// ensureFreeSlot makes sure that a staker can join a given pool.
// If this is not possible an appropriate error is returned.
// A pool has a fixed amount of slots. If there are still free slots
// a staker can just join (even with the smallest stake possible).
// If all slots are taken, it checks if the new staker has more stake
// than the current lowest staker in that pool.
// If so, the lowest staker gets removed from the pool, so that the
// new staker can join.
func (k Keeper) ensureFreeSlot(ctx sdk.Context, poolId uint64, stakerAddress string) error {
	// check if slots are still available
	if k.GetStakerCountOfPool(ctx, poolId) >= types.MaxStakers {
		// if not - get lowest staker
		lowestStaker, _ := k.getLowestStaker(ctx, poolId)
		lowestStakerAddress := util.MustAccountAddressFromValAddress(lowestStaker.OperatorAddress)

		// if new pool joiner has more stake than lowest staker kick him out
		newAmount := k.GetDelegationAmount(ctx, stakerAddress)
		lowestAmount := k.GetDelegationAmount(ctx, lowestStakerAddress)
		if newAmount > lowestAmount {
			// remove lowest staker from pool
			k.LeavePool(ctx, lowestStakerAddress, poolId)
		} else {
			return errors.Wrapf(errorsTypes.ErrLogic, types.ErrStakeTooLow.Error(), lowestAmount)
		}
	}

	return nil
}
