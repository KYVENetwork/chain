package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRoundRobinProgress stores the round-robin progress for a pool
func (k Keeper) SetRoundRobinProgress(ctx sdk.Context, roundRobinProgress types.RoundRobinProgress) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.RoundRobinProgressPrefix)

	b := k.cdc.MustMarshal(&roundRobinProgress)
	store.Set(types.RoundRobinProgressKey(roundRobinProgress.PoolId), b)
}

// GetRoundRobinProgress returns the round-robin progress for a pool
func (k Keeper) GetRoundRobinProgress(ctx sdk.Context, poolId uint64) (val types.RoundRobinProgress, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.RoundRobinProgressPrefix)

	b := store.Get(types.RoundRobinProgressKey(poolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllRoundRobinProgress returns the round-robin progress of all pools
func (k Keeper) GetAllRoundRobinProgress(ctx sdk.Context) (list []types.RoundRobinProgress) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.RoundRobinProgressPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	for ; iterator.Valid(); iterator.Next() {
		var val types.RoundRobinProgress
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
