package keeper

import (
	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/multi_coin_rewards module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetMultiCoinDistributionPendingTime returns the MultiCoinDistributionPendingTime param
func (k Keeper) GetMultiCoinDistributionPendingTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).MultiCoinDistributionPendingTime
}

// GetMultiCoinDistributionPolicyAdminAddress returns the admin address which is allowed to update the coin weights
// distribution policy
func (k Keeper) GetMultiCoinDistributionPolicyAdminAddress(ctx sdk.Context) (res string) {
	return k.GetParams(ctx).MultiCoinDistributionPolicyAdminAddress
}

// SetParams sets the x/multi_coin_rewards module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
