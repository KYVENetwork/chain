package keeper

import (
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/pool module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetProtocolInflationShare returns the ProtocolInflationShare param
func (k Keeper) GetProtocolInflationShare(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).ProtocolInflationShare
}

// GetPoolInflationPayoutRate returns the GetPoolInflationPayoutRate param
func (k Keeper) GetPoolInflationPayoutRate(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).PoolInflationPayoutRate
}

// SetParams sets the x/pool module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
