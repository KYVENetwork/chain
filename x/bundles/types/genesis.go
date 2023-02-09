package types

import (
	"fmt"
	"sort"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	// Bundle proposal
	bundleProposalKey := make(map[string]struct{})

	for _, elem := range gs.BundleProposalList {
		index := string(BundleProposalKey(elem.PoolId))
		if _, ok := bundleProposalKey[index]; ok {
			return fmt.Errorf("duplicated pool-id for bundle proposal %v", elem)
		}
		bundleProposalKey[index] = struct{}{}
	}

	// Finalized bundles
	finalizedBundleProposals := make(map[string]struct{})
	previousIndexPerPool := make(map[uint64]uint64)

	sort.Slice(gs.FinalizedBundleList, func(i, j int) bool {
		return gs.FinalizedBundleList[i].Id < gs.FinalizedBundleList[j].Id
	})

	for _, elem := range gs.FinalizedBundleList {
		index := string(FinalizedBundleKey(elem.PoolId, elem.Id))
		if _, ok := finalizedBundleProposals[index]; ok {
			return fmt.Errorf("duplicated index for finalized bundle %v", elem)
		}
		finalizedBundleProposals[index] = struct{}{}

		if previousIndexPerPool[elem.PoolId] == elem.Id {
			previousIndexPerPool[elem.PoolId] += 1
		} else {
			return fmt.Errorf("missing finalized bundle %v", elem)
		}
	}

	return gs.Params.Validate()
}
