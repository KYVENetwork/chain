package types

import (
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Pool
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
)

type DelegationKeeper interface {
	util.DelegationKeeper

	Delegate(sdk.Context, string, string, uint64) error
}

type PoolKeeper interface {
	util.PoolKeeper

	GetPoolWithError(sdk.Context, uint64) (poolTypes.Pool, error)
}
