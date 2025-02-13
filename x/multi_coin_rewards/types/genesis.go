package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                      DefaultParams(),
		MultiCoinDistributionPolicy: &MultiCoinDistributionPolicy{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	// Multi Coin Pending Rewards
	multiCoinPendingRewardsMap := make(map[string]struct{})

	for _, elem := range gs.MultiCoinPendingRewardsEntries {
		index := string(MultiCoinPendingRewardsKeyEntry(elem.Index))
		if _, ok := multiCoinPendingRewardsMap[index]; ok {
			return fmt.Errorf("duplicated index for multi coin pending rewards entry %v", elem)
		}
		if elem.Index > gs.QueueStatePendingRewards.HighIndex {
			return fmt.Errorf(" multi coin pending rewards entry index too high: %v", elem)
		}
		if elem.Index < gs.QueueStatePendingRewards.LowIndex {
			return fmt.Errorf(" multi coin pending rewards entry index too low: %v", elem)
		}

		multiCoinPendingRewardsMap[index] = struct{}{}
	}

	if _, err := ParseAndNormalizeMultiCoinDistributionMap(*gs.MultiCoinDistributionPolicy); err != nil {
		return err
	}

	return gs.Params.Validate()
}
