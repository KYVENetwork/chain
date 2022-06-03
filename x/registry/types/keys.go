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

// ========== EVENTS ===================
// general event props
const (
	EventName    = "EventName"
	EventPoolId  = "PoolId"
	EventCreator = "Creator"
	EventAmount  = "Amount"
)

// voting
const (
	VoteEventKey      = "Voted"
	VoteEventBundleId = "BundleId"
	VoteEventVote     = "Vote"
)

// slashing
const (
	SlashEventKey = "ReceivedSlash"
	SlashAccount  = "Account"
)

// Activity
const (
	ProposalEventKey          = "ProposalEnded"
	ProposalEventBundleId     = "BundleId"
	ProposalEventByteSize     = "ByteSize"
	ProposalEventUploader     = "Uploader"
	ProposalEventNextUploader = "NextUploader"
	ProposalEventReward       = "BundleReward"
	ProposalEventValid        = "Valid"
	ProposalEventInvalid      = "Invalid"
	ProposalEventFromHeight   = "FromHeight"
	ProposalEventToHeight     = "ToHeight"
	ProposalEventStatus       = "Status"
)

const (
	UpdateMetadataEventKey = "UpdateMetadata"
	UpdateMetadataCommission = "Commission"
	UpdateMetadataMoniker = "Moniker"
	UpdateMetadataWebsite = "Website"
	UpdateMetadataLogo = "Logo"
)

// ============ KV-STORE ===============

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	PoolKey           = "Pool-value-"
	PoolCountKey      = "Pool-count-"
	UnbondingStateKey = "UnbondingState-value-"
)

const (
	// DelegationEntriesKeyPrefix is the prefix to retrieve all DelegationEntries
	// refactor to []byte{0x01}
	DelegationEntriesKeyPrefix = "DelegationEntries/value/"
	// DelegationPoolDataKeyPrefix is the prefix to retrieve all DelegationPoolData
	// refactor to []byte{0x02}
	DelegationPoolDataKeyPrefix = "DelegationPoolData/value/"
	// DelegatorKeyPrefix is the prefix to retrieve all Delegator
	// refactor to []byte{0x03}
	DelegatorKeyPrefix = "Delegator/value/"
	// FunderKeyPrefix is the prefix to retrieve all Funder
	// refactor to []byte{0x05}
	FunderKeyPrefix = "Funder/value/"
	// ProposalKeyPrefix is the prefix to retrieve all Proposal
	// refactor to []byte{0x06}
	ProposalKeyPrefix = "Proposal/value/"
	// StakerKeyPrefix is the prefix to retrieve all Staker
	// refactor to []byte{0x08}
	StakerKeyPrefix = "Staker/value/"
	// UnbondingEntriesKeyPrefix is the prefix to retrieve all UnbondingEntries
	// refactor to []byte{0x09}
	UnbondingEntriesKeyPrefix = "UnbondingEntries/value/"
	// UnbondingEntriesKeyPrefixByDelegator is the prefix to retrieve all UnbondingEntries
	// refactor to []byte{0x10}
	UnbondingEntriesKeyPrefixByDelegator = "UnbondingEntriesByDelegator/value/"
)

var (
	// DelegatorKeyPrefixIndex2 is the prefix for a different key order for the DelegatorKeyPrefix
	DelegatorKeyPrefixIndex2 = []byte{0x04}

	// ProposalKeyPrefixIndex2 is the prefix for a different key order for the DelegatorKeyPrefix
	ProposalKeyPrefixIndex2 = []byte{0x08}
)

// DelegationEntriesKey returns the store Key to retrieve a DelegationEntries from the index fields
func DelegationEntriesKey(poolId uint64, stakerAddress string, kIndex uint64) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AString(stakerAddress).AInt(kIndex).Key
}

// DelegationPoolDataKey returns the store Key to retrieve a DelegationPoolData from the index fields
func DelegationPoolDataKey(poolId uint64, stakerAddress string) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AString(stakerAddress).Key
}

// DelegatorKey returns the store Key to retrieve a Delegator from the index fields
func DelegatorKey(poolId uint64, stakerAddress string, delegatorAddress string) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AString(stakerAddress).AString(delegatorAddress).Key
}

// DelegatorKeyIndex2 returns the store Key to retrieve a Delegator from the index fields
func DelegatorKeyIndex2(delegatorAddress string, poolId uint64, stakerAddress string) []byte {
	return KeyPrefixBuilder{}.AString(delegatorAddress).AInt(poolId).AString(stakerAddress).Key
}

// FunderKey returns the store Key to retrieve a Funder from the index fields
func FunderKey(funder string, poolId uint64) []byte {
	return KeyPrefixBuilder{}.AString(funder).AInt(poolId).Key
}

// ProposalKey returns the store Key to retrieve a Proposal from the index fields
func ProposalKey(bundleId string) []byte {
	return KeyPrefixBuilder{}.AString(bundleId).Key
}

// ProposalKey returns the store Key to retrieve a Proposal from the index fields
func ProposalKeyIndex2(poolId uint64, fromId uint64) []byte {
	return KeyPrefixBuilder{}.AInt(poolId).AInt(fromId).Key
}

// StakerKey returns the store Key to retrieve a Staker from the index fields
func StakerKey(staker string, poolId uint64) []byte {
	return KeyPrefixBuilder{}.AString(staker).AInt(poolId).Key
}

// UnbondingEntriesKey returns the store Key to retrieve a UnbondingEntries from the index fields
func UnbondingEntriesKey(index uint64) []byte {
	return KeyPrefixBuilder{}.AInt(index).Key
}

// UnbondingEntriesByDelegatorKey returns the store Key to retrieve a UnbondingEntries from the index fields
// Index is still needed to make Key unique
func UnbondingEntriesByDelegatorKey(delegator string, index uint64) []byte {
	return KeyPrefixBuilder{}.AString(delegator).AInt(index).Key
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
