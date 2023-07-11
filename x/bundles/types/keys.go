package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "bundles"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_bundles"
)

var ParamsKey = []byte{0x00}

var (
	// BundleKeyPrefix ...
	BundleKeyPrefix = []byte{1}
	// FinalizedBundlePrefix ...
	FinalizedBundlePrefix = []byte{2}
	// FinalizedBundleVersionMapKey ...
	FinalizedBundleVersionMapKey = []byte{3}

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

// FinalizedBundleByIndexKey ...
func FinalizedBundleByIndexKey(poolId uint64, height uint64) []byte {
	return util.GetByteKey(poolId, height)
}
