package types

import (
	"fmt"

	"github.com/KYVENetwork/chain/util"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                     DefaultParams(),
		DelegatorList:              []Delegator{},
		DelegationEntryList:        []DelegationEntry{},
		DelegationDataList:         []DelegationData{},
		DelegationSlashList:        []DelegationSlash{},
		UndelegationQueueEntryList: []UndelegationQueueEntry{},
		QueueStateUndelegation:     QueueState{},
		RedelegationCooldownList:   []RedelegationCooldown{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	if err := gs.validateF1(); err != nil {
		return err
	}

	if err := gs.validateUnbondingQueue(); err != nil {
		return err
	}

	if err := gs.validateRedelegation(); err != nil {
		return err
	}

	return gs.Params.Validate()
}

func (gs *GenesisState) validateF1() error {
	// Check for duplicated index in Delegator
	delegatorMap := make(map[string]struct{})
	delegatorKIndexMap := make(map[string]struct{})

	for _, elem := range gs.DelegatorList {
		index := string(DelegatorKey(elem.Staker, elem.Delegator))
		if _, ok := delegatorMap[index]; ok {
			return fmt.Errorf("duplicated index for delegator %v", elem)
		}
		delegatorMap[index] = struct{}{}

		kIndex := string(util.GetByteKey(elem.Staker, elem.KIndex))
		if _, ok := delegatorKIndexMap[kIndex]; ok {
			return fmt.Errorf("duplicated k-index for delegator %v", elem)
		}
		delegatorKIndexMap[kIndex] = struct{}{}
	}

	// Check for duplicated index in DelegationEntries
	delegationEntriesKIndexMap := make(map[string]struct{})
	for _, elem := range gs.DelegationEntryList {

		index := string(DelegationEntriesKey(elem.Staker, elem.KIndex))
		if _, ok := delegationEntriesKIndexMap[index]; ok {
			return fmt.Errorf("duplicated k-index for delegation entry %v", elem)
		}
		delegationEntriesKIndexMap[index] = struct{}{}
	}

	// Check for duplicated index in DelegationEntries
	delegationDataMap := make(map[string]struct{})

	for _, elem := range gs.DelegationDataList {
		index := string(DelegationDataKey(elem.Staker))
		if _, ok := delegationDataMap[index]; ok {
			return fmt.Errorf("duplicated index for delegation data %v", elem)
		}
		delegationDataMap[index] = struct{}{}
	}

	// Check for duplicated index in SlashEntries
	slashMap := make(map[string]struct{})

	for _, elem := range gs.DelegationSlashList {
		index := string(DelegationSlashEntriesKey(elem.Staker, elem.KIndex))
		if _, ok := slashMap[index]; ok {
			return fmt.Errorf("duplicated k-index for delegation entry %v", elem)
		}
		//nolint:all
		entryIndex := string(DelegationEntriesKey(elem.Staker, elem.KIndex))
		if _, ok := delegationEntriesKIndexMap[entryIndex]; !ok {
			return fmt.Errorf("slash entry pointing to non-existent delegation index: %v", elem)
		}

		slashMap[index] = struct{}{}
	}

	return nil
}

func (gs *GenesisState) validateUnbondingQueue() error {
	// Check undelegation queue
	unbondingMap := make(map[string]struct{})

	for _, elem := range gs.UndelegationQueueEntryList {
		index := string(UndelegationQueueKey(elem.Index))
		if _, ok := unbondingMap[index]; ok {
			return fmt.Errorf("duplicated index for unbonding entry %v", elem)
		}
		if elem.Index > gs.QueueStateUndelegation.HighIndex {
			return fmt.Errorf("unbonding entry index too high: %v", elem)
		}
		if elem.Index < gs.QueueStateUndelegation.LowIndex {
			return fmt.Errorf("unbonding entry index too low: %v", elem)
		}

		unbondingMap[index] = struct{}{}
	}
	return nil
}

func (gs *GenesisState) validateRedelegation() error {
	// Check undelegation queue
	redelegationMap := make(map[string]struct{})

	for _, elem := range gs.RedelegationCooldownList {
		index := string(RedelegationCooldownKey(elem.Address, elem.CreationDate))
		if _, ok := redelegationMap[index]; ok {
			return fmt.Errorf("duplicated index for redelegation entry %v", elem)
		}

		redelegationMap[index] = struct{}{}
	}
	return nil
}
