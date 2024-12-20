package keeper

import (
	"encoding/binary"

	storeTypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetCommissionChangeEntry ...
func (k Keeper) SetCommissionChangeEntry(ctx sdk.Context, commissionChangeEntry types.CommissionChangeEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CommissionChangeEntryKeyPrefix)
	b := k.cdc.MustMarshal(&commissionChangeEntry)
	store.Set(types.CommissionChangeEntryKey(commissionChangeEntry.Index), b)

	// Insert the same entry with a different key prefix for query lookup
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, commissionChangeEntry.Index)

	indexStore := prefix.NewStore(storeAdapter, types.CommissionChangeEntryKeyPrefixIndex2)
	indexStore.Set(types.CommissionChangeEntryKeyIndex2(
		commissionChangeEntry.Staker,
		commissionChangeEntry.PoolId,
	), indexBytes)
}

// GetCommissionChangeEntry ...
func (k Keeper) GetCommissionChangeEntry(ctx sdk.Context, index uint64) (val types.CommissionChangeEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CommissionChangeEntryKeyPrefix)

	b := store.Get(types.CommissionChangeEntryKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetCommissionChangeEntryByIndex2 returns a pending commission change entry by staker address (if there is one)
func (k Keeper) GetCommissionChangeEntryByIndex2(ctx sdk.Context, staker string, poolId uint64) (val types.CommissionChangeEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CommissionChangeEntryKeyPrefixIndex2)

	b := store.Get(types.CommissionChangeEntryKeyIndex2(staker, poolId))
	if b == nil {
		return val, false
	}

	index := binary.BigEndian.Uint64(b)

	return k.GetCommissionChangeEntry(ctx, index)
}

// RemoveCommissionChangeEntry ...
func (k Keeper) RemoveCommissionChangeEntry(ctx sdk.Context, commissionChangeEntry *types.CommissionChangeEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CommissionChangeEntryKeyPrefix)
	store.Delete(types.CommissionChangeEntryKey(commissionChangeEntry.Index))

	indexStore := prefix.NewStore(storeAdapter, types.CommissionChangeEntryKeyPrefixIndex2)
	indexStore.Delete(types.CommissionChangeEntryKeyIndex2(
		commissionChangeEntry.Staker,
		commissionChangeEntry.PoolId,
	))
}

// GetAllCommissionChangeEntries returns all pending commission change entries of all stakers
func (k Keeper) GetAllCommissionChangeEntries(ctx sdk.Context) (list []types.CommissionChangeEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.CommissionChangeEntryKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CommissionChangeEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
