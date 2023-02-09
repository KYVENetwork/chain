package types

import "fmt"

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in account entries
	accountsIndexMap := make(map[string]struct{})

	for _, elem := range gs.AccountList {
		index := string(TeamVestingAccountKeyPrefix(elem.Id))
		if _, ok := accountsIndexMap[index]; ok {
			return fmt.Errorf("duplicated account id %v", elem)
		}
		accountsIndexMap[index] = struct{}{}
		if elem.Id >= gs.AccountCount {
			return fmt.Errorf("account id higher than account count %v", elem)
		}
	}

	if gs.Authority.RewardsClaimed > gs.Authority.TotalRewards {
		return fmt.Errorf("claimed is greater than total rewards %#v", gs.Authority)
	}

	return nil
}
