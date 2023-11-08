package types

import "github.com/KYVENetwork/chain/util"

const (
	// ModuleName defines the module name
	ModuleName = "funders"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_funders"
)

const (
	// MaxFunders which are allowed per pool
	MaxFunders = 50
)

var (
	// ParamsKey is the prefix for all module params defined in params.proto
	ParamsKey = []byte{0x00}

	// FunderKeyPrefix is indexed by the funder address
	// and contains all funders regardless of the pool
	// key -> FunderKeyPrefix | <funderAddr>
	FunderKeyPrefix = []byte{1}

	// FundingKeyPrefixByFunder stores funding for each funder and pool by funder
	// FundingKeyPrefixByFunder | <funder> | <poolId>
	FundingKeyPrefixByFunder = []byte{2, 0}

	// FundingKeyPrefixByPool stores funding for each funder and pool by pool
	// FundingKeyPrefixByPool | <poolId> | <funder>
	FundingKeyPrefixByPool = []byte{2, 1}

	// FundingStateKeyPrefix stores funding state for each pool
	// FundingStateKeyPrefix | <poolId> | <funder>
	FundingStateKeyPrefix = []byte{3, 0}
)

func FunderKey(funderAddress string) []byte {
	return util.GetByteKey(funderAddress)
}

func FundingKeyByFunder(funderAddress string, poolId uint64) []byte {
	return util.GetByteKey(funderAddress, poolId)
}

func FundingKeyByPool(funderAddress string, poolId uint64) []byte {
	return util.GetByteKey(poolId, funderAddress)
}

// FundingKeyByFunderIter is used to query all fundings for a funder
func FundingKeyByFunderIter(funderAddress string) []byte {
	return util.GetByteKey(funderAddress)
}

// FundingKeyByPoolIter is used to query all fundings for a pool
func FundingKeyByPoolIter(poolId uint64) []byte {
	return util.GetByteKey(poolId)
}

func FundingStateKey(poolId uint64) []byte {
	return util.GetByteKey(poolId)
}
