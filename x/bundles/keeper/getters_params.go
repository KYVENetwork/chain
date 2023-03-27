package keeper

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/bundles module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetUploadTimeout returns the UploadTimeout param
func (k Keeper) GetUploadTimeout(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).UploadTimeout
}

// GetStorageCost returns the StorageCost param
func (k Keeper) GetStorageCost(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).StorageCost
}

// GetNetworkFee returns the NetworkFee param
func (k Keeper) GetNetworkFee(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).NetworkFee
}

// GetMaxPoints returns the MaxPoints param
func (k Keeper) GetMaxPoints(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).MaxPoints
}

// SetParams sets the x/bundles module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
