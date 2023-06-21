package keeper

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRoundRobinProgress stores the round-robin progress for a pool
func (k Keeper) SetRoundRobinProgress(ctx sdk.Context, roundRobinProgress types.RoundRobinProgress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RoundRobinProgressPrefix)
	b := k.cdc.MustMarshal(&roundRobinProgress)
	store.Set(types.RoundRobinProgressKey(roundRobinProgress.PoolId), b)
}

// GetRoundRobinProgress returns the round-robin progress for a pool
func (k Keeper) GetRoundRobinProgress(ctx sdk.Context, poolId uint64) (val types.RoundRobinProgress, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RoundRobinProgressPrefix)

	b := store.Get(types.RoundRobinProgressKey(poolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllRoundRobinProgress returns the round-robin progress of all pools
func (k Keeper) GetAllRoundRobinProgress(ctx sdk.Context) (list []types.RoundRobinProgress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RoundRobinProgressPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	for ; iterator.Valid(); iterator.Next() {
		var val types.RoundRobinProgress
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
