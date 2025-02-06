package keeper

import (
	"github.com/KYVENetwork/chain/x/compliance/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/compliance module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetMultiCoinRefundPendingTime returns the MultiCoinRefundPendingTime param
func (k Keeper) GetMultiCoinRefundPendingTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).MultiCoinRefundPendingTime
}

// GetMultiCoinRefundPolicyAdminAddress returns the admin address which is allowed to update the coin weights
// refund policy
func (k Keeper) GetMultiCoinRefundPolicyAdminAddress(ctx sdk.Context) (res string) {
	return k.GetParams(ctx).MultiCoinRefundPolicyAdminAddress
}

// SetParams sets the x/compliance module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
