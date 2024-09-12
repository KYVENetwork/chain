package types

import (
	"cosmossdk.io/collections"
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "bundles"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var ParamsPrefix = collections.NewPrefix(0)

var (
	// BundleKeyPrefix ...
	BundleKeyPrefix = []byte{1}
	// FinalizedBundlePrefix ...
	FinalizedBundlePrefix = []byte{2}
	// FinalizedBundleVersionMapKey ...
	FinalizedBundleVersionMapKey = []byte{3}
	// RoundRobinProgressPrefix ...
	RoundRobinProgressPrefix = []byte{4}
	// BundlesMigrationHeightKey ...
	BundlesMigrationHeightKey = []byte{5}

	FinalizedBundleByIndexPrefix = []byte{11}
)

// BundleProposalKey ...
func BundleProposalKey(poolId uint64) []byte {
	return util.GetByteKey(poolId)
}

// FinalizedBundleKey ...
func FinalizedBundleKey(poolId uint64, id uint64) []byte {
	return util.GetByteKey(poolId, id)
}

// RoundRobinProgressKey ...
func RoundRobinProgressKey(poolId uint64) []byte {
	return util.GetByteKey(poolId)
}

// FinalizedBundleByIndexKey ...
func FinalizedBundleByIndexKey(poolId uint64, height uint64) []byte {
	return util.GetByteKey(poolId, height)
}
