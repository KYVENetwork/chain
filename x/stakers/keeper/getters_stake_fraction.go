package keeper

import (
	"encoding/binary"

	storeTypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetStakeFractionChangeEntry ...
func (k Keeper) SetStakeFractionChangeEntry(ctx sdk.Context, stakeFractionChangeEntry types.StakeFractionChangeEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakeFractionChangeEntryKeyPrefix)
	b := k.cdc.MustMarshal(&stakeFractionChangeEntry)
	store.Set(types.StakeFractionChangeEntryKey(stakeFractionChangeEntry.Index), b)

	// Insert the same entry with a different key prefix for query lookup
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, stakeFractionChangeEntry.Index)

	indexStore := prefix.NewStore(storeAdapter, types.StakeFractionChangeKeyPrefixIndex2)
	indexStore.Set(types.StakeFractionChangeEntryKeyIndex2(
		stakeFractionChangeEntry.Staker,
		stakeFractionChangeEntry.PoolId,
	), indexBytes)
}

// GetStakeFractionChangeEntry ...
func (k Keeper) GetStakeFractionChangeEntry(ctx sdk.Context, index uint64) (val types.StakeFractionChangeEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakeFractionChangeEntryKeyPrefix)

	b := store.Get(types.StakeFractionChangeEntryKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetStakeFractionChangeEntryByIndex2 returns a pending stake fraction change entry by staker address (if there is one)
func (k Keeper) GetStakeFractionChangeEntryByIndex2(ctx sdk.Context, staker string, poolId uint64) (val types.StakeFractionChangeEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakeFractionChangeKeyPrefixIndex2)

	b := store.Get(types.StakeFractionChangeEntryKeyIndex2(staker, poolId))
	if b == nil {
		return val, false
	}

	index := binary.BigEndian.Uint64(b)

	return k.GetStakeFractionChangeEntry(ctx, index)
}

// RemoveStakeFractionEntry ...
func (k Keeper) RemoveStakeFractionEntry(ctx sdk.Context, stakeFractionChangeEntry *types.StakeFractionChangeEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakeFractionChangeEntryKeyPrefix)
	store.Delete(types.StakeFractionChangeEntryKey(stakeFractionChangeEntry.Index))

	indexStore := prefix.NewStore(storeAdapter, types.StakeFractionChangeKeyPrefixIndex2)
	indexStore.Delete(types.StakeFractionChangeEntryKeyIndex2(
		stakeFractionChangeEntry.Staker,
		stakeFractionChangeEntry.PoolId,
	))
}

// GetAllStakeFractionChangeEntries returns all pending stake fraction change entries of all stakers
func (k Keeper) GetAllStakeFractionChangeEntries(ctx sdk.Context) (list []types.StakeFractionChangeEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakeFractionChangeEntryKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.StakeFractionChangeEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
