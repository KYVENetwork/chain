package types

import (
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "stakers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_stakers"
)

var (
	// ParamsKey is the prefix for all module params defined in params.proto
	ParamsKey = []byte{0x00}

	// StakerKeyPrefix is indexed by the staker address
	// and contains all stakers regardless of the pool
	// key -> StakerKeyPrefix | <stakerAddr>
	StakerKeyPrefix = []byte{1}

	// ValaccountPrefix stores valaccount for each staker and pool
	// ValaccountPrefix | <poolId> | <staker>
	ValaccountPrefix = []byte{2, 0}
	// ValaccountPrefixIndex2 | <staker> | <poolId>
	ValaccountPrefixIndex2 = []byte{2, 1}

	// CommissionChangeEntryKeyPrefix | <index>
	CommissionChangeEntryKeyPrefix = []byte{4, 0}
	// CommissionChangeEntryKeyPrefixIndex2 | <staker>
	CommissionChangeEntryKeyPrefixIndex2 = []byte{4, 1}

	// LeavePoolEntryKeyPrefix | <index>
	LeavePoolEntryKeyPrefix = []byte{5, 0}
	// LeavePoolEntryKeyPrefixIndex2 | <staker> | <poolId>
	LeavePoolEntryKeyPrefixIndex2 = []byte{5, 1}

	ActiveStakerIndex = []byte{6}
)

// ENUM aggregated data types
type STAKER_STATS string

var STAKER_STATS_COUNT STAKER_STATS = "total_stakers"

// ENUM queue types identifiers
type QUEUE_IDENTIFIER []byte

var (
	QUEUE_IDENTIFIER_COMMISSION QUEUE_IDENTIFIER = []byte{30, 2}
	QUEUE_IDENTIFIER_LEAVE      QUEUE_IDENTIFIER = []byte{30, 3}
)

const MaxStakers = 50

var DefaultCommission = sdk.MustNewDecFromStr("0.1")

// StakerKey returns the store Key to retrieve a Staker from the index fields
func StakerKey(staker string) []byte {
	return util.GetByteKey(staker)
}

func ValaccountKey(poolId uint64, staker string) []byte {
	return util.GetByteKey(poolId, staker)
}

func ValaccountKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}

func CommissionChangeEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}

// Important: only one queue entry per staker is allowed at a time.
func CommissionChangeEntryKeyIndex2(staker string) []byte {
	return util.GetByteKey(staker)
}

func LeavePoolEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}

func LeavePoolEntryKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}

func ActiveStakerKeyIndex(staker string) []byte {
	return util.GetByteKey(staker)
}
