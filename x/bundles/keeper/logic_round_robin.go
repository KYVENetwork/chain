package keeper

import (
	"sort"

	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*

Weighted Round Robin Uploader Selection

This file implements all necessary logic for a weighted round-robin algorithm. An example introduction can be found
[here](https://learnblockchain.cn/docs/tendermint/spec/reactors/consensus/proposer-selection.html).

Our implementation additionally has the option of temporarily skipping participants for single rounds.
They do not advance in the round-robin progress and can not be selected as an uploader.
The frequencies of uploader selection with respect to the excluded ones can be described as follows.

Let $R$ denote the number of total rounds and $r$ the index of the current round.
Let $N$ denote the number of total validators and $n$ the index of the n-th validator.

The stake (+ delegation) of each validator for each round is given by
    $s(n, r)$

Then the total stake for round r is given by
    $S(r) = \sum_{i=1}^N s(i, r)$

Ignoring the existing progress, the likeliness of being selected in the next round is given by
    $p(n, r) = s(n, r) / S(r)$

Using this value one can obtain the frequencies for uploader selection over all rounds, which is
   $P(n) = 1/R * \sum_{r=1}^R p(n, r)$

Except for rounding errors $P(n)$ is independent from $R$ if $p(n, r)$ is constant.
If validators $i$ is excluded for round $k$, this is denoted by $s(i, k) = 0$. So in general $S(r)$ is
dependent on validator exclusions and validators set changes.

*/

// RoundRobinValidatorPower contains the total delegation of a validator. It is used as a cache
// because the calculation of the total delegation needs to access the KV-Store and therefore
// consumes gas everytime it is called.
// This value is only stored for the current round and only lives inside the memory.
type RoundRobinValidatorPower struct {
	Address string
	Power   int64
}

// RoundRobinValidatorSet is the in memory-object for working with the round-robin state
// It can not be stored to the KV-Store as go map iteration is non-deterministic.
// To obtain a deterministic state of the current state call GetRoundRobinProgress().
type RoundRobinValidatorSet struct {
	PoolId     uint64
	Validators []RoundRobinValidatorPower
	Progress   map[string]int64
}

// LoadRoundRobinValidatorSet initialises a validator set for the given pool id.
// If available it fetches the current round-robin state. Then it iterates all current pool
// validators and initialises the set accordingly.
// If a validator left the pool, the progress will be ignored.
// If new validators joined the pool, their progress will be zero.
func (k Keeper) LoadRoundRobinValidatorSet(ctx sdk.Context, poolId uint64) RoundRobinValidatorSet {
	vs := RoundRobinValidatorSet{}
	vs.PoolId = poolId
	vs.Progress = make(map[string]int64, 0)
	totalDelegation := int64(0)
	// Used for calculating the set difference of active validators and existing round-robin set
	newValidators := make(map[string]bool, 0)
	// Add all current validators to round-robin set
	for _, address := range k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		delegation := k.delegationKeeper.GetDelegationAmount(ctx, address)
		if delegation > 0 {
			// If a validator has no delegation do not add to round-robin set. Validator is basically non-existent.
			vs.Validators = append(vs.Validators, RoundRobinValidatorPower{
				Address: address,
				Power:   int64(delegation),
			})
			vs.Progress[address] = 0
			totalDelegation += int64(delegation)
			newValidators[address] = true
		}
	}

	// Fetch stored progress
	roundRobinProgress, _ := k.GetRoundRobinProgress(ctx, poolId)
	for _, progress := range roundRobinProgress.ProgressList {
		_, ok := vs.Progress[progress.Address]
		// If the address is not found it means that the validator left the pool.
		// Therefore, this entry can be ignored.
		if ok {
			vs.Progress[progress.Address] += progress.Progress
		}

		newValidators[progress.Address] = false
	}

	for newAddress, isNew := range newValidators {
		if isNew {
			vs.Progress[newAddress] = sdk.MustNewDecFromStr("-1.125").MulInt64(totalDelegation).TruncateInt64()
		}
	}

	vs.normalize()
	return vs
}

// SaveRoundRobinValidatorSet saves the current round-robin progress for the given poolId to the KV-Store
func (k Keeper) SaveRoundRobinValidatorSet(ctx sdk.Context, vs RoundRobinValidatorSet) {
	roundRobinProgress := types.RoundRobinProgress{
		PoolId:       vs.PoolId,
		ProgressList: vs.GetRoundRobinProgress(),
	}
	k.SetRoundRobinProgress(ctx, roundRobinProgress)
}

// GetRoundRobinProgress returns a deterministic (sorted) list of the current round-robin progress.
// Due to the fact that maps have no order in Go, we must introduce one by ourselves.
// This is done by sorting all entries by their addresses alphabetically.
func (vs *RoundRobinValidatorSet) GetRoundRobinProgress() []*types.RoundRobinSingleValidatorProgress {
	// Convert map to list
	result := make([]*types.RoundRobinSingleValidatorProgress, 0)
	for address, progress := range vs.Progress {
		singleProgress := types.RoundRobinSingleValidatorProgress{
			Address:  address,
			Progress: progress,
		}
		result = append(result, &singleProgress)
	}
	// Sort addresses alphabetically
	sort.Slice(
		result, func(i, j int) bool {
			return result[i].Address < result[j].Address
		},
	)
	return result
}

// getTotalDelegation returns the total delegation power of the current set.
func (vs *RoundRobinValidatorSet) getTotalDelegation() (total int64) {
	for _, vp := range vs.Validators {
		total += vp.Power
	}
	return
}

// getTotalDelegation returns the sum of all progresses of each validator.
// This value is supposed to always be zero. However, if new validators join or leave this
// value might no longer be 0. Have a look at normalize().
func (vs *RoundRobinValidatorSet) getTotalProgress() (total int64) {
	for _, vp := range vs.Progress {
		total += vp
	}
	return
}

// getMinMaxDifference returns the difference of progresses of the two validators with the
// most and least progress.
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

// size returns the number of participants in the current round-robin set.
func (vs *RoundRobinValidatorSet) size() int64 {
	return int64(len(vs.Validators))
}

func (vs *RoundRobinValidatorSet) normalize() {
	if vs.size() == 0 {
		return
	}

	diff := vs.getMinMaxDifference()

	totalProgress := vs.getTotalProgress()
	threshold := 2 * vs.getTotalDelegation()
	if diff > threshold {

		totalProgress = 0
		for _, val := range vs.Validators {
			decProgress := sdk.NewDec(vs.Progress[val.Address])
			vs.Progress[val.Address] = decProgress.MulInt64(threshold).QuoInt64(diff).TruncateInt64()
			totalProgress += vs.Progress[val.Address]
		}
	}

	// center priorities around zero and update
	avg := sdk.NewDec(totalProgress).QuoInt64(vs.size()).TruncateInt64()
	for key := range vs.Progress {
		vs.Progress[key] -= avg
	}
}

func (vs *RoundRobinValidatorSet) NextProposer(excludedAddresses ...string) string {
	if vs.size() == 0 {
		return ""
	}

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
