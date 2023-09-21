package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
	//GetPoolWithError(ctx sdk.Context, poolId uint64) (pooltypes.Pool, error)
}
