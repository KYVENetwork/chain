package types

import (
	"fmt"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	// Staker
	stakerLeaving := make(map[string]bool)

	// Pool accounts
	poolAccountsMap := make(map[string]struct{})
	for _, elem := range gs.PoolAccountList {
		index := string(PoolAccountKey(elem.PoolId, elem.Staker))
		if _, ok := poolAccountsMap[index]; ok {
			return fmt.Errorf("duplicated index for pool account %v", elem)
		}
		poolAccountsMap[index] = struct{}{}
		stakerLeaving[index] = elem.IsLeaving
	}

	// Commission Change
	commissionChangeMap := make(map[string]struct{})

	for _, elem := range gs.CommissionChangeEntries {
		index := string(CommissionChangeEntryKey(elem.Index))
		if _, ok := commissionChangeMap[index]; ok {
			return fmt.Errorf("duplicated index for commission change entry %v", elem)
		}
		if elem.Index > gs.QueueStateCommission.HighIndex {
			return fmt.Errorf("commission change entry index too high: %v", elem)
		}
		if elem.Index < gs.QueueStateCommission.LowIndex {
			return fmt.Errorf("commission change entry index too low: %v", elem)
		}

		commissionChangeMap[index] = struct{}{}
	}

	// Leave Pool
	for _, elem := range gs.LeavePoolEntries {
		if elem.Index > gs.QueueStateLeave.HighIndex {
			return fmt.Errorf("unbonding stake entry index too high: %v", elem)
		}
		if elem.Index < gs.QueueStateLeave.LowIndex {
			return fmt.Errorf("unbonding stake entry index too low: %v", elem)
		}
		if !stakerLeaving[string(PoolAccountKey(elem.PoolId, elem.Staker))] {
			return fmt.Errorf("inconsistent staker leave: %v", elem)
		}
		stakerLeaving[string(PoolAccountKey(elem.PoolId, elem.Staker))] = false
	}

	for staker, isLeaving := range stakerLeaving {
		if isLeaving {
			return fmt.Errorf("inconsistent staker leave: %v", staker)
		}
	}

	// Stake Fraction Change
	stakeFractionChangeMap := make(map[string]struct{})

	for _, elem := range gs.StakeFractionChangeEntries {
		index := string(StakeFractionChangeEntryKey(elem.Index))
		if _, ok := stakeFractionChangeMap[index]; ok {
			return fmt.Errorf("duplicated index for stake fraction change entry %v", elem)
		}
		if elem.Index > gs.QueueStateStakeFraction.HighIndex {
			return fmt.Errorf("stake fraction change entry index too high: %v", elem)
		}
		if elem.Index < gs.QueueStateStakeFraction.LowIndex {
			return fmt.Errorf("stake fraction change entry index too low: %v", elem)
		}

		stakeFractionChangeMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
