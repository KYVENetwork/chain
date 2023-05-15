package keeper

import (
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/stakers module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetGlobalMinDelegation returns the GlobalMinDelegation param
func (k Keeper) GetGlobalMinDelegation(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).GlobalMinDelegation
}

// SetParams sets the x/stakers module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
