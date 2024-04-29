package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The `DelegationData` stores general aggregated variables for each existing staker
// as well as necessary variables needed for the F1 distribution algorithm
// Look at the proto-file for detailed explanation of the variables.
// Every staker with at least one delegator has this entry.

// AddAmountToDelegationRewards adds the specified amount to the current delegationData object.
// This is needed by the F1-algorithm to calculate to outstanding rewards
func (k Keeper) AddAmountToDelegationRewards(ctx sdk.Context, stakerAddress string, amount sdk.Coins) {
	delegationData, found := k.GetDelegationData(ctx, stakerAddress)
	if found {
		delegationData.CurrentRewards = delegationData.CurrentRewards.Add(amount...)
		k.SetDelegationData(ctx, delegationData)
	}
}

// SetDelegationData set a specific delegationPoolData in the store from its index
func (k Keeper) SetDelegationData(ctx sdk.Context, delegationData types.DelegationData) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.DelegationDataKeyPrefix)
	b := k.cdc.MustMarshal(&delegationData)
	store.Set(types.DelegationDataKey(delegationData.Staker), b)
}

// GetDelegationData returns a delegationData entry for a specific staker
// with `stakerAddress`
func (k Keeper) GetDelegationData(ctx sdk.Context, stakerAddress string) (val types.DelegationData, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.DelegationDataKeyPrefix)

	b := store.Get(types.DelegationDataKey(stakerAddress))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// DoesDelegationDataExist check if the staker with `stakerAddress` has
// a delegation data entry. This is the case if the staker as at least one delegator.
func (k Keeper) DoesDelegationDataExist(ctx sdk.Context, stakerAddress string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.DelegationDataKeyPrefix)
	return store.Has(types.DelegationDataKey(stakerAddress))
}

// RemoveDelegationData removes a delegationData entry from the pool
func (k Keeper) RemoveDelegationData(ctx sdk.Context, stakerAddress string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.DelegationDataKeyPrefix)
	store.Delete(types.DelegationDataKey(stakerAddress))
}

// GetAllDelegationData returns all delegationData entries
func (k Keeper) GetAllDelegationData(ctx sdk.Context) (list []types.DelegationData) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.DelegationDataKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DelegationData
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
