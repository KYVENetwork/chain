package types

import (
	"cosmossdk.io/collections"
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "multi_coin_rewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	MultiCoinRewardsRedistributionAccountName = "multi_coin_rewards_distribution"
)

var (
	// ParamsKey is the prefix for all module params defined in params.proto
	ParamsKey = collections.NewPrefix(0)

	// MultiCoinRewardsEnabledKey is the key prefix for storing all users who have opted in for multi-coin-rewards.
	MultiCoinRewardsEnabledKey = collections.NewPrefix(1)

	// MultiCoinDistributionPolicyKey
	MultiCoinDistributionPolicyKey = collections.NewPrefix(2)

	// MultiCoinPendingRewardsEntryKeyPrefix | <index>
	MultiCoinPendingRewardsEntryKeyPrefix = []byte{3, 0}
	// MultiCoinPendingRewardsEntryKeyPrefixIndex2 | <address> | <poolId>
	MultiCoinPendingRewardsEntryKeyPrefixIndex2 = []byte{3, 1}
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
