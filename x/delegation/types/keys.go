package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "delegation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_delegation"
)

var ParamsKey = []byte{0x00}

var StakerIndexKeyPrefix = []byte{1} // StakerIndexPoolCountKeyPrefix = []byte{1}

var (
	// DelegatorKeyPrefix is the prefix to retrieve all Delegator entries
	DelegatorKeyPrefix = []byte{1, 0}

	// DelegatorKeyPrefixIndex2 is the prefix for a different key order for the DelegatorKeyPrefix
	DelegatorKeyPrefixIndex2 = []byte{1, 1}

	// DelegationEntriesKeyPrefix is the prefix to retrieve all DelegationEntries
	DelegationEntriesKeyPrefix = []byte{2}

	// DelegationDataKeyPrefix ...
	DelegationDataKeyPrefix = []byte{3}

	// DelegationSlashEntriesKeyPrefix ...
	DelegationSlashEntriesKeyPrefix = []byte{4}

	// QueueKey ...
	QueueKey = []byte{5}

	// UndelegationQueueKeyPrefix ...
	UndelegationQueueKeyPrefix = []byte{6, 0}

	// UndelegationQueueKeyPrefixIndex2 ...
	UndelegationQueueKeyPrefixIndex2 = []byte{6, 1}

	// RedelegationCooldownPrefix ...
	RedelegationCooldownPrefix = []byte{7}
)

// DelegatorKey returns the store Key to retrieve a Delegator from the index fields
func DelegatorKey(stakerAddress string, delegatorAddress string) []byte {
	return util.GetByteKey(stakerAddress, delegatorAddress)
}

// DelegatorKeyIndex2 returns the store Key to retrieve a Delegator from the index fields
func DelegatorKeyIndex2(delegatorAddress string, stakerAddress string) []byte {
	return util.GetByteKey(delegatorAddress, stakerAddress)
}

// DelegationEntriesKey returns the store Key to retrieve a DelegationEntries from the index fields
func DelegationEntriesKey(stakerAddress string, kIndex uint64) []byte {
	return util.GetByteKey(stakerAddress, kIndex)
}

// DelegationDataKey returns the store Key to retrieve a DelegationPoolData from the index fields
func DelegationDataKey(stakerAddress string) []byte {
	return util.GetByteKey(stakerAddress)
}

func UndelegationQueueKey(kIndex uint64) []byte {
	return util.GetByteKey(kIndex)
}

func UndelegationQueueKeyIndex2(stakerAddress string, kIndex uint64) []byte {
	return util.GetByteKey(stakerAddress, kIndex)
}

func RedelegationCooldownKey(delegator string, block uint64) []byte {
	return util.GetByteKey(delegator, block)
}

func DelegationSlashEntriesKey(stakerAddress string, kIndex uint64) []byte {
	return util.GetByteKey(stakerAddress, kIndex)
}

func StakerIndexKey(amount uint64, stakerAddress string) []byte {
	return util.GetByteKey(amount, stakerAddress)
}

//func StakerIndexByPoolCountKey(poolCount uint64, amount uint64, stakerAddress string) []byte {
//	return util.GetByteKey(poolCount, amount, stakerAddress)
//}
