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

const (
	MaxFunders = 50 // maximum amount of funders which are allowed
)

var (
	PoolKey      = []byte{1}
	PoolCountKey = []byte{2}
)

func PoolKeyPrefix(poolId uint64) []byte {
	return util.GetByteKey(poolId)
}
