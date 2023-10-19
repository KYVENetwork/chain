package types

import "fmt"

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	fundersIndexMap := make(map[string]struct{})
	for _, funder := range gs.FunderList {
		index := FunderKey(funder.Address)
		if _, ok := fundersIndexMap[string(index)]; ok {
			return fmt.Errorf("duplicated funder id for %v", funder)
		}
	}

	fundingByFunderIndexMap := make(map[string]struct{})
	fundingByPoolIndexMap := make(map[string]struct{})
	for _, funding := range gs.FundingList {
		byFunderIndex := FundingKeyByFunder(funding.FunderAddress, funding.PoolId)
		if _, ok := fundingByFunderIndexMap[string(byFunderIndex)]; ok {
			return fmt.Errorf("duplicated funding id for %v", funding)
		}
		byPoolIndex := FundingKeyByPool(funding.FunderAddress, funding.PoolId)
		if _, ok := fundingByPoolIndexMap[string(byPoolIndex)]; ok {
			return fmt.Errorf("duplicated funding id for %v", funding)
		}
	}

	fundingStateIndexMap := make(map[string]struct{})
	for _, fundingState := range gs.FundingStateList {
		index := FundingStateKey(fundingState.PoolId)
		if _, ok := fundingStateIndexMap[string(index)]; ok {
			return fmt.Errorf("duplicated funding state id for %v", fundingState)
		}
	}
	return gs.Params.Validate()
}
