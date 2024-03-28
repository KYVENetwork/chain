package types

import (
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
	GetPoolWithError(ctx sdk.Context, poolId uint64) (poolTypes.Pool, error)
}
