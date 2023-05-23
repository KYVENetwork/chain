package types

import (
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Delegation
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Pool
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
)

type DelegationKeeper interface {
	util.DelegationKeeper

	SlashDelegators(ctx sdk.Context, poolId uint64, staker string, slashType delegationTypes.SlashType)
}

type PoolKeeper interface {
	util.PoolKeeper

	GetAllPools(ctx sdk.Context) []poolTypes.Pool
	GetPool(sdk.Context, uint64) (poolTypes.Pool, bool)
	GetPoolWithError(sdk.Context, uint64) (poolTypes.Pool, error)
}
