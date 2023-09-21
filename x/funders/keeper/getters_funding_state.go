package keeper

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DoesFundingStateExist checks if the FundingState exists
func (k Keeper) doesFundingStateExist(ctx sdk.Context, poolId uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingStateKeyPrefix)
	return store.Has(types.FundingStateKey(poolId))
}

// GetFundingState returns the FundingState
func (k Keeper) getFundingState(ctx sdk.Context, poolId uint64) (fundingState types.FundingState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingStateKeyPrefix)

	b := store.Get(types.FundingStateKey(
		poolId,
	))
	if b == nil {
		return fundingState, false
	}

	k.cdc.MustUnmarshal(b, &fundingState)
	return fundingState, true
}

// SetFundingState sets a specific FundingState in the store from its index
func (k Keeper) setFundingState(ctx sdk.Context, fundingState types.FundingState) {
	b := k.cdc.MustMarshal(&fundingState)
	storeByFunder := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingStateKeyPrefix)
	storeByFunder.Set(types.FundingStateKey(
		fundingState.PoolId,
	), b)
}
