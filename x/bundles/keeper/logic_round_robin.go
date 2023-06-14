package keeper

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

// Round Robin implementation

// RoundRobinValidatorPower contains the total delegation of a validator. It is used as a cache
// because the calculation of the total delegation requires gas and needs to access
// the KV Store.
type RoundRobinValidatorPower struct {
	Address string
	Power   int64
}

// RoundRobinValidatorSet is the in memory-object for working with the round-robin state
// It can not be stored to the KV-Store as go map iteration is non-deterministic
type RoundRobinValidatorSet struct {
	Validators []RoundRobinValidatorPower
	Progress   map[string]int64
}

// LoadRoundRobinValidatorSet initialises a validator set out of the current stored progress.
func (k Keeper) LoadRoundRobinValidatorSet(ctx sdk.Context, poolId uint64) RoundRobinValidatorSet {
	vs := RoundRobinValidatorSet{}
	vs.Progress = make(map[string]int64, 0)
	// Add all current stakers to round-robin set
	for _, address := range k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		vs.Validators = append(
			vs.Validators, RoundRobinValidatorPower{
				Address: address,
				Power:   int64(k.delegationKeeper.GetDelegationAmount(ctx, address)),
			},
		)
		vs.Progress[address] = 0
	}

	roundRobinProgress, _ := k.GetRoundRobinProgress(ctx, poolId)
	// Now add existing progress.
	for _, progress := range roundRobinProgress.ProgressList {
		_, ok := vs.Progress[progress.Address]
		// If the address is not found it means that the staker left the pool.
		// Therefore, this entry can be ignored.
		if ok {
			vs.Progress[progress.Address] += progress.Progress
		}
	}
	return vs
}

// SaveRoundRobinValidatorSet saves the current round-robin progress to the KV-Store
func (k Keeper) SaveRoundRobinValidatorSet(ctx sdk.Context, poolId uint64, vs RoundRobinValidatorSet) {
	roundRobinProgress := types.RoundRobinProgress{
		PoolId:       poolId,
		ProgressList: vs.GetRoundRobinProgress(),
	}
	k.SetRoundRobinProgress(ctx, roundRobinProgress)
}

func (vs *RoundRobinValidatorSet) GetRoundRobinProgress() []*types.RoundRobinSingleValidatorProgress {
	result := make([]*types.RoundRobinSingleValidatorProgress, 0)
	for address, progress := range vs.Progress {
		singleProgress := types.RoundRobinSingleValidatorProgress{
			Address:  address,
			Progress: progress,
		}
		result = append(result, &singleProgress)
	}
	sort.Slice(
		result, func(i, j int) bool {
			if result[i].Progress == result[j].Progress {
				return result[i].Address < result[j].Address
			}
			return result[i].Progress < result[j].Progress
		},
	)
	return result
}

func (vs *RoundRobinValidatorSet) getTotalDelegation() (total int64) {
	for _, vp := range vs.Validators {
		total += vp.Power
	}
	return
}

func (vs *RoundRobinValidatorSet) getTotalProgress() (total int64) {
	for _, vp := range vs.Progress {
		total += vp
	}
	return
}

func (vs *RoundRobinValidatorSet) getMinMaxDifference() int64 {
	if len(vs.Validators) == 0 {
		return 0
	}
	max := vs.Progress[vs.Validators[0].Address]
	min := max
	for _, p := range vs.Progress {
		if p > max {
			max = p
		}
		if p < min {
			min = p
		}
	}
	return max - min
}

func (vs *RoundRobinValidatorSet) size() int64 {
	return int64(len(vs.Validators))
}

func (vs *RoundRobinValidatorSet) normalize() {

	diff := vs.getMinMaxDifference()

	totalProgress := vs.getTotalProgress()
	threshold := 2 * vs.getTotalDelegation()
	if diff > threshold {

		totalProgress = 0
		for _, val := range vs.Validators {
			vs.Progress[val.Address] = vs.Progress[val.Address] * threshold / diff
			totalProgress += vs.Progress[val.Address]
		}
	}

	// center priorities around zero and update
	avg := totalProgress / vs.size()
	for key := range vs.Progress {
		vs.Progress[key] -= avg
	}
}

func (vs *RoundRobinValidatorSet) NextProposer(excludedAddresses ...string) string {
	vs.normalize()

	// If all addresses are excluded, then no address should be excluded
	if len(excludedAddresses) == len(vs.Validators) {
		excludedAddresses = make([]string, 0)
	}

	mapExcludedAddresses := make(map[string]bool)
	for _, excluded := range excludedAddresses {
		mapExcludedAddresses[excluded] = true
	}

	// update
	excludedPower := int64(0)
	for _, validator := range vs.Validators {
		if !mapExcludedAddresses[validator.Address] {
			vs.Progress[validator.Address] += validator.Power
		} else {
			excludedPower += validator.Power
		}
	}

	currentMaxValidator := vs.Validators[0].Address
	for _, validator := range vs.Validators {
		if !mapExcludedAddresses[validator.Address] {
			if vs.Progress[validator.Address] > vs.Progress[currentMaxValidator] {
				currentMaxValidator = validator.Address
			}
		}
	}

	vs.Progress[currentMaxValidator] -= vs.getTotalDelegation() - excludedPower

	return currentMaxValidator
}
