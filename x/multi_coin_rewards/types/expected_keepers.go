package types

import (
	"cosmossdk.io/math"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolKeeper interface {
	GetMaxVotingPowerPerPool(ctx sdk.Context) (res math.LegacyDec)
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
	GetPoolWithError(ctx sdk.Context, poolId uint64) (poolTypes.Pool, error)
	EnsurePoolAccount(ctx sdk.Context, id uint64)
}
