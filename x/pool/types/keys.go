package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "pool"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_pool"
)

var (
	// ParamsKey is the prefix for all module params defined in params.proto
	ParamsKey = []byte{0}

	// PoolKey is the prefix for all pools defined in pool.proto
	PoolKey = []byte{1}

	// PoolCountKey is the prefix for the pool counter defined in pool.proto
	PoolCountKey = []byte{2}
)

func PoolKeyPrefix(poolId uint64) []byte {
	return util.GetByteKey(poolId)
}
