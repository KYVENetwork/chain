package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "compliance"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	MultiCoinRewardsRedistributionAccountName = "multi_coin_rewards_redistribution"
)

var (
	// ParamsKey is the prefix for all module params defined in params.proto
	ParamsKey = []byte{0x00}

	// MultiCoinPendingRewardsEntryKeyPrefix | <index>
	MultiCoinPendingRewardsEntryKeyPrefix = []byte{8, 0}
	// MultiCoinPendingRewardsEntryKeyPrefixIndex2 | <staker> | <poolId>
	MultiCoinPendingRewardsEntryKeyPrefixIndex2 = []byte{8, 1}

	MultiCoinRewardsEnabledKeyPrefix = []byte{9, 2}

	MultiCoinRefundPolicyKeyPrefix = []byte{9, 3}
)

// ENUM queue types identifiers
type QUEUE_IDENTIFIER []byte

var QUEUE_IDENTIFIER_MULTI_COIN_REWARDS QUEUE_IDENTIFIER = []byte{30, 5}

func MultiCoinPendingRewardsKeyEntry(index uint64) []byte {
	return util.GetByteKey(index)
}

func MultiCoinPendingRewardsKeyEntryIndex2(address string, index uint64) []byte {
	return util.GetByteKey(address, index)
}
