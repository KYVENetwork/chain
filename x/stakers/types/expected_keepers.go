package types

import (
	"github.com/KYVENetwork/chain/util"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolKeeper interface {
	util.PoolKeeper

	GetPoolWithError(sdk.Context, uint64) (poolTypes.Pool, error)
}
