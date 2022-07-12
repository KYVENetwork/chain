package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "registry"

	// StoreKey defines the primary module store Key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing Key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store Key
	MemStoreKey = "mem_registry"
)

// registry constants
const (
	MaxFunders          = 50 // maximum amount of funders which are allowed
	MaxStakers          = 50 // maximum amount of stakers which are allowed
	DefaultCommission   = "0.9"
	KYVE_NO_DATA_BUNDLE = "KYVE_NO_DATA_BUNDLE"
)

// ============ KV-STORE ===============

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	PoolKey      = "Pool-value-" // []byte{0,1}
	PoolCountKey = "Pool-count-" // []byte{0,2}

	// UnbondingStakingQueueStateKey ...
	UnbondingStakingQueueStateKey = []byte{0, 3}

	// UnbondingDelegationQueueStateKey ...
	UnbondingDelegationQueueStateKey = []byte{0, 4}

	// CommissionChangeQueueStateKey ...
	CommissionChangeQueueStateKey = []byte{0, 5}
)

var (
	// StakerKeyPrefix is the prefix to retrieve all Staker
	StakerKeyPrefix = "Staker/value/" // []byte{0x01}

	// FunderKeyPrefix is the prefix to retrieve all Funder
	FunderKeyPrefix = "Funder/value/" // []byte{0x02}

	// DelegatorKeyPrefix is the prefix to retrieve all Delegator
	DelegatorKeyPrefix = "Delegator/value/" // []byte{0x03}
	// DelegatorKeyPrefixIndex2 is the prefix for a different key order for the DelegatorKeyPrefix
	DelegatorKeyPrefixIndex2 = []byte{4}

	// DelegationEntriesKeyPrefix is the prefix to retrieve all DelegationEntries
	DelegationEntriesKeyPrefix = "DelegationEntries/value/" // []byte{0x05}

	// DelegationPoolDataKeyPrefix is the prefix to retrieve all DelegationPoolData
	DelegationPoolDataKeyPrefix = "DelegationPoolData/value/" // []byte{0x06}

	// ProposalKeyPrefix is the prefix to retrieve all Proposal
	ProposalKeyPrefix = "Proposal/value/" // []byte{0x07}
	// ProposalKeyPrefixIndex2 is the prefix for a different key order for the DelegatorKeyPrefix
	ProposalKeyPrefixIndex2 = []byte{8, 0}
	// ProposalKeyPrefixIndex3 is the prefix for a different key order for the DelegatorKeyPrefix
	ProposalKeyPrefixIndex3 = []byte{8, 1}

	// UnbondingStakingQueueEntryKeyPrefix ...
	UnbondingStakingQueueEntryKeyPrefix = []byte{9}
	// UnbondingStakingQueueEntryKeyPrefixIndex2 ...
	UnbondingStakingQueueEntryKeyPrefixIndex2 = []byte{10}
	// UnbondingStakerKeyPrefix ...
	UnbondingStakerKeyPrefix = []byte{11}

	// UnbondingDelegationQueueEntryKeyPrefix ...
	UnbondingDelegationQueueEntryKeyPrefix = []byte{12}
	// UnbondingDelegationQueueEntryKeyPrefixIndex2 ...
	UnbondingDelegationQueueEntryKeyPrefixIndex2 = []byte{13}

	// RedelegationCooldownPrefix ...
	RedelegationCooldownPrefix = []byte{14}

	// CommissionChangeQueueEntryKeyPrefix ...
	CommissionChangeQueueEntryKeyPrefix = []byte{15}
	// CommissionChangeQueueEntryKeyPrefixIndex2 ...
	CommissionChangeQueueEntryKeyPrefixIndex2 = []byte{16}
)

// StakerKey returns the store Key to retrieve a Staker from the index fields
func StakerKey(staker string, poolId uint64) []byte {
	return KeyPrefixBuilder{}.AString(staker).AInt(poolId).Key
}

// FunderKey returns the store Key to retrieve a Funder from the index fields
func FunderKey(funder string, poolId uint64) []byte {
	return KeyPrefixBuilder{}.AString(funder).AInt(poolId).Key
}

// === DELEGATION ===

// DelegatorKey returns the store Key to retrieve a Delegator from the index fields
func DelegatorKey(poolId uint64, stakerAddress string, delegatorAddress string) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AString(stakerAddress).AString(delegatorAddress).Key
}

// DelegatorKeyIndex2 returns the store Key to retrieve a Delegator from the index fields
func DelegatorKeyIndex2(delegatorAddress string, poolId uint64, stakerAddress string) []byte {
	return KeyPrefixBuilder{}.AString(delegatorAddress).AInt(poolId).AString(stakerAddress).Key
}

// DelegationEntriesKey returns the store Key to retrieve a DelegationEntries from the index fields
func DelegationEntriesKey(poolId uint64, stakerAddress string, kIndex uint64) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AString(stakerAddress).AInt(kIndex).Key
}

// DelegationPoolDataKey returns the store Key to retrieve a DelegationPoolData from the index fields
func DelegationPoolDataKey(poolId uint64, stakerAddress string) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AString(stakerAddress).Key
}

// === PROPOSALS ===

// ProposalKey returns the store Key to retrieve a Proposal from the index fields
func ProposalKey(storageId string) []byte {
	return KeyPrefixBuilder{}.AString(storageId).Key
}

// ProposalKeyIndex2 ...
func ProposalKeyIndex2(poolId uint64, id uint64) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AInt(id).Key
}

// ProposalKeyIndex3 ...
func ProposalKeyIndex3(poolId uint64, finalizedAt uint64) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AInt(finalizedAt).Key
}

// === UNBONDING ===

func UnbondingStakingQueueEntryKey(index uint64) []byte {
	return KeyPrefixBuilder{}.AInt(index).Key
}
func UnbondingStakingQueueEntryKeyIndex2(staker string, index uint64) []byte {
	return KeyPrefixBuilder{}.AString(staker).AInt(index).Key
}

func UnbondingStakerKey(poolId uint64, staker string) []byte {
	return KeyPrefixBuilder{}.AString(staker).AInt(poolId).Key
}

func UnbondingDelegationQueueEntryKey(index uint64) []byte {
	return KeyPrefixBuilder{}.AInt(index).Key
}
func UnbondingDelegationQueueEntryKeyIndex2(delegator string, index uint64) []byte {
	return KeyPrefixBuilder{}.AString(delegator).AInt(index).Key
}

func RedelegationCooldownKey(delegator string, block uint64) []byte {
	return KeyPrefixBuilder{}.AString(delegator).AInt(block).Key
}

func CommissionChangeQueueEntryKey(index uint64) []byte {
	return KeyPrefixBuilder{}.AInt(index).Key
}

// Important: only one queue entry per staker+poolId is allowed at a time.
func CommissionChangeQueueEntryKeyIndex2(staker string, poolId uint64) []byte {
	return KeyPrefixBuilder{}.AString(staker).AInt(poolId).Key
}

type KeyPrefixBuilder struct {
	Key []byte
}

func (k KeyPrefixBuilder) AInt(n uint64) KeyPrefixBuilder {
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, n)
	k.Key = append(k.Key, indexBytes...)
	k.Key = append(k.Key, []byte("/")...)
	return k
}

func (k KeyPrefixBuilder) AString(s string) KeyPrefixBuilder {
	k.Key = append(k.Key, []byte(s)...)
	k.Key = append(k.Key, []byte("/")...)
	return k
}
