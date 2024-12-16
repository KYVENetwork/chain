package keeper

import (
	m "math"

	"cosmossdk.io/math"

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
	var minAmount uint64 = m.MaxUint64

	for _, staker := range k.getAllStakersOfPool(ctx, poolId) {
		delegationAmount := k.GetValidatorPoolStake(ctx, util.MustAccountAddressFromValAddress(staker.OperatorAddress), poolId)
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
func (k Keeper) ensureFreeSlot(ctx sdk.Context, poolId uint64, stakerAddress string, stakeFraction math.LegacyDec) error {
	// check if slots are still available
	if k.GetStakerCountOfPool(ctx, poolId) >= types.MaxStakers {
		// if not - get lowest staker
		lowestStaker, _ := k.getLowestStaker(ctx, poolId)
		lowestStakerAddress := util.MustAccountAddressFromValAddress(lowestStaker.OperatorAddress)

		// if new pool joiner would have more stake than lowest staker kick him out
		newStaker, _ := k.GetValidator(ctx, stakerAddress)
		newAmount := uint64(math.LegacyNewDecFromInt(newStaker.GetBondedTokens()).Mul(stakeFraction).TruncateInt64())
		lowestAmount := k.GetValidatorPoolStake(ctx, lowestStakerAddress, poolId)
		if newAmount > lowestAmount {
			// remove lowest staker from pool
			k.LeavePool(ctx, lowestStakerAddress, poolId)
		} else {
			return errors.Wrapf(errorsTypes.ErrLogic, types.ErrStakeTooLow.Error(), lowestAmount)
		}
	}

	return nil
}
