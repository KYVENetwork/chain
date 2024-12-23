package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "stakers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	// ParamsKey is the prefix for all module params defined in params.proto
	ParamsKey = []byte{0x00}

	// StakerKeyPrefix is indexed by the staker address
	// and contains all stakers regardless of the pool
	// key -> StakerKeyPrefix | <stakerAddr>
	StakerKeyPrefix = []byte{1}

	// PoolAccountPrefix stores pool account for each staker and pool
	// PoolAccountPrefix | <poolId> | <staker>
	PoolAccountPrefix = []byte{2, 0}
	// PoolAccountPrefixIndex2 | <staker> | <poolId>
	PoolAccountPrefixIndex2 = []byte{2, 1}

	// CommissionChangeEntryKeyPrefix | <index>
	CommissionChangeEntryKeyPrefix = []byte{4, 0}
	// CommissionChangeEntryKeyPrefixIndex2 | <staker> | <poolId>
	CommissionChangeEntryKeyPrefixIndex2 = []byte{4, 1}

	// LeavePoolEntryKeyPrefix | <index>
	LeavePoolEntryKeyPrefix = []byte{5, 0}
	// LeavePoolEntryKeyPrefixIndex2 | <staker> | <poolId>
	LeavePoolEntryKeyPrefixIndex2 = []byte{5, 1}

	ActiveStakerIndex = []byte{6}

	// StakeFractionChangeEntryKeyPrefix | <index>
	StakeFractionChangeEntryKeyPrefix = []byte{7, 0}
	// StakeFractionChangeKeyPrefixIndex2 | <staker> | <poolId>
	StakeFractionChangeKeyPrefixIndex2 = []byte{7, 1}
)

// ENUM aggregated data types
type STAKER_STATS string

var STAKER_STATS_COUNT STAKER_STATS = "total_stakers"

// ENUM queue types identifiers
type QUEUE_IDENTIFIER []byte

var (
	QUEUE_IDENTIFIER_COMMISSION     QUEUE_IDENTIFIER = []byte{30, 2}
	QUEUE_IDENTIFIER_LEAVE          QUEUE_IDENTIFIER = []byte{30, 3}
	QUEUE_IDENTIFIER_STAKE_FRACTION QUEUE_IDENTIFIER = []byte{30, 4}
)

const MaxStakers = 50

func PoolAccountKey(poolId uint64, staker string) []byte {
	return util.GetByteKey(poolId, staker)
}

func PoolAccountKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}

func CommissionChangeEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}

// Important: only one queue entry per staker is allowed at a time.
func CommissionChangeEntryKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}

func LeavePoolEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}

func LeavePoolEntryKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}

func StakeFractionChangeEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}

func StakeFractionChangeEntryKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}
