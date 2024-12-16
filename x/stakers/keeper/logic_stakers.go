package keeper

import (
	m "math"
	"sort"

	"cosmossdk.io/math"

	"github.com/KYVENetwork/chain/util"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/KYVENetwork/chain/x/stakers/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// getEffectiveValidatorStakes returns a map for all pool stakers their effective stake. Effective stake
// is the actual stake at risk for a validator, it can be lower than the provided stake because
// of the max voting power per pool parameter, but it can never be higher. Thus, if the effective
// stake is lower than the provided one the slash
// TODO: explain or link to formula
func (k Keeper) getEffectiveValidatorStakes(ctx sdk.Context, poolId uint64, mustIncludeStakers ...string) (effectiveStakes map[string]uint64) {
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
		validator, _ := k.GetValidator(ctx, address)
		valaccount, _ := k.GetValaccount(ctx, poolId, address)
		stake := uint64(valaccount.StakeFraction.MulInt(validator.GetBondedTokens()).TruncateInt64())

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
		return make(map[string]uint64)
	}

	if math.LegacyOneDec().Quo(maxVotingPower).GT(math.LegacyNewDec(int64(len(addresses)))) {
		return make(map[string]uint64)
	}

	for i, validator := range validators {
		if math.LegacyNewDec(int64(effectiveStakes[validator.Address])).GT(maxVotingPower.MulInt64(totalStake)) {
			amount := math.LegacyNewDec(int64(effectiveStakes[validator.Address])).Sub(maxVotingPower.MulInt64(totalStake)).TruncateInt64()

			totalStakeRemainder -= int64(validator.Stake)
			effectiveStakes[validator.Address] -= uint64(amount)

			if totalStakeRemainder > 0 {
				for _, v := range validators[i+1:] {
					effectiveStakes[v.Address] += uint64(math.LegacyNewDec(int64(v.Stake)).QuoInt64(totalStakeRemainder).MulInt64(amount).TruncateInt64())
				}
			}
		}
	}

	scaleFactor := math.LegacyZeroDec()

	// get lowest validator with effective stake still bigger than zero
	for i := len(validators) - 1; i >= 0; i-- {
		if effectiveStakes[validators[i].Address] > 0 {
			scaleFactor = math.LegacyNewDec(int64(validators[i].Stake)).QuoInt64(int64(effectiveStakes[validators[i].Address]))
			break
		}
	}

	// scale all effective stakes down to scale factor
	for _, validator := range validators {
		// TODO: is rounding here fine?
		effectiveStakes[validator.Address] = uint64(scaleFactor.MulInt64(int64(effectiveStakes[validator.Address])).RoundInt64())
	}

	return
}

// getLowestStaker returns the staker with the lowest total stake
// (self-delegation + delegation) of a given pool.
// If all pool slots are taken, this is the staker who then
// gets kicked out.
func (k Keeper) getLowestStaker(ctx sdk.Context, poolId uint64) (val stakingTypes.Validator, found bool) {
	var minAmount uint64 = m.MaxUint64

	for _, staker := range k.getAllStakersOfPool(ctx, poolId) {
		stake := k.GetValidatorPoolStake(ctx, util.MustAccountAddressFromValAddress(staker.OperatorAddress), poolId)
		if stake < minAmount {
			minAmount = stake
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
