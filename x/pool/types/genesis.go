package types

import (
	"fmt"
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
	// Check for duplicated index in DelegationEntries
	poolIndexMap := make(map[string]struct{})

	for _, elem := range gs.PoolList {
		index := string(PoolKeyPrefix(elem.Id))
		if _, ok := poolIndexMap[index]; ok {
			return fmt.Errorf("duplicated pool id %v", elem)
		}
		poolIndexMap[index] = struct{}{}
		if elem.Id >= gs.PoolCount {
			return fmt.Errorf("pool id higher than pool count %v", elem)
		}
	}

	return gs.Params.Validate()
}
