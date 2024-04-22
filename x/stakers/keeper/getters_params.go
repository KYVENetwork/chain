package keeper

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/stakers module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetCommissionChangeTime returns the CommissionChangeTime param
func (k Keeper) GetCommissionChangeTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).CommissionChangeTime
}

// GetLeavePoolTime returns the LeavePoolTime param
func (k Keeper) GetLeavePoolTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).LeavePoolTime
}

// SetParams sets the x/stakers module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
