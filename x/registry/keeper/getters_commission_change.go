package keeper

import (
	"encoding/binary"
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// #####################
// === QUEUE ENTRIES ===
// #####################

// SetCommissionChangeQueueEntry ...
func (k Keeper) SetCommissionChangeQueueEntry(ctx sdk.Context, commissionChangeQueueEntry types.CommissionChangeQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CommissionChangeQueueEntryKeyPrefix)
	b := k.cdc.MustMarshal(&commissionChangeQueueEntry)
	store.Set(types.CommissionChangeQueueEntryKey(commissionChangeQueueEntry.Index), b)

	// Insert the same entry with a different key prefix for query lookup
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, commissionChangeQueueEntry.Index)

	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.CommissionChangeQueueEntryKeyPrefixIndex2)
	indexStore.Set(types.CommissionChangeQueueEntryKeyIndex2(
		commissionChangeQueueEntry.Staker,
		commissionChangeQueueEntry.PoolId,
	), indexBytes)
}

// GetCommissionChangeQueueEntry ...
func (k Keeper) GetCommissionChangeQueueEntry(ctx sdk.Context, index uint64) (val types.CommissionChangeQueueEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CommissionChangeQueueEntryKeyPrefix)

	b := store.Get(types.CommissionChangeQueueEntryKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetCommissionChangeQueueEntryByIndex2 ...
func (k Keeper) GetCommissionChangeQueueEntryByIndex2(ctx sdk.Context, staker string, poolId uint64) (val types.CommissionChangeQueueEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CommissionChangeQueueEntryKeyPrefixIndex2)

	b := store.Get(types.CommissionChangeQueueEntryKeyIndex2(staker, poolId))
	if b == nil {
		return val, false
	}

	index := binary.BigEndian.Uint64(b)

	return k.GetCommissionChangeQueueEntry(ctx, index)
}

// RemoveCommissionChangeQueueEntry ...
func (k Keeper) RemoveCommissionChangeQueueEntry(ctx sdk.Context, commissionChangeQueueEntry *types.CommissionChangeQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CommissionChangeQueueEntryKeyPrefix)
	store.Delete(types.CommissionChangeQueueEntryKey(commissionChangeQueueEntry.Index))

	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.CommissionChangeQueueEntryKeyPrefixIndex2)
	indexStore.Delete(types.CommissionChangeQueueEntryKeyIndex2(
		commissionChangeQueueEntry.Staker,
		commissionChangeQueueEntry.PoolId,
	))
}

// GetAllCommissionChangeQueueEntries ...
func (k Keeper) GetAllCommissionChangeQueueEntries(ctx sdk.Context) (list []types.CommissionChangeQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CommissionChangeQueueEntryKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CommissionChangeQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ###################
// === QUEUE STATE ===
// ###################

// GetCommissionChangeQueueState ...
func (k Keeper) GetCommissionChangeQueueState(ctx sdk.Context) (state types.CommissionChangeQueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := store.Get(types.CommissionChangeQueueStateKey)

	if b == nil {
		return state
	}

	k.cdc.MustUnmarshal(b, &state)
	return
}

// SetCommissionChangeQueueState ...
func (k Keeper) SetCommissionChangeQueueState(ctx sdk.Context, state types.CommissionChangeQueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := k.cdc.MustMarshal(&state)
	store.Set(types.CommissionChangeQueueStateKey, b)
}
