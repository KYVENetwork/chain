package keeper

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DoesFundingStateExist checks if the FundingState exists
func (k Keeper) DoesFundingStateExist(ctx sdk.Context, poolId uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingStateKeyPrefix)
	return store.Has(types.FundingStateKey(poolId))
}

// GetFundingState returns the FundingState
func (k Keeper) GetFundingState(ctx sdk.Context, poolId uint64) (fundingState types.FundingState, found bool) {
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

// GetAllFundingStates returns all FundingStates
func (k Keeper) GetAllFundingStates(ctx sdk.Context) (fundingStates []types.FundingState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingStateKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.FundingState
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		fundingStates = append(fundingStates, val)
	}

	return fundingStates
}

// SetFundingState sets a specific FundingState in the store from its index
func (k Keeper) SetFundingState(ctx sdk.Context, fundingState *types.FundingState) {
	b := k.cdc.MustMarshal(fundingState)
	storeByFunder := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingStateKeyPrefix)
	storeByFunder.Set(types.FundingStateKey(
		fundingState.PoolId,
	), b)
}

func (k Keeper) GetActiveFundings(ctx sdk.Context, fundingState types.FundingState) (fundings []types.Funding) {
	for _, funder := range fundingState.ActiveFunderAddresses {
		funding, found := k.GetFunding(ctx, funder, fundingState.PoolId)
		if found {
			fundings = append(fundings, funding)
		} // else should never happen or we have a corrupted state
	}
	return fundings
}
