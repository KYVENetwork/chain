package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
}

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
}
