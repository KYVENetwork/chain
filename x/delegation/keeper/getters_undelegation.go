package keeper

import (
	storeTypes "cosmossdk.io/store/types"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// #####################
// === QUEUE ENTRIES ===
// #####################

// SetUndelegationQueueEntry ...
func (k Keeper) SetUndelegationQueueEntry(ctx sdk.Context, undelegationQueueEntry types.UndelegationQueueEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.UndelegationQueueKeyPrefix)
	b := k.cdc.MustMarshal(&undelegationQueueEntry)
	store.Set(types.UndelegationQueueKey(
		undelegationQueueEntry.Index,
	), b)

	// Insert the same entry with a different key prefix for query lookup
	indexStore := prefix.NewStore(storeAdapter, types.UndelegationQueueKeyPrefixIndex2)
	indexStore.Set(types.UndelegationQueueKeyIndex2(
		undelegationQueueEntry.Delegator,
		undelegationQueueEntry.Index,
	), []byte{})
}

// GetUndelegationQueueEntry ...
func (k Keeper) GetUndelegationQueueEntry(ctx sdk.Context, index uint64) (val types.UndelegationQueueEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.UndelegationQueueKeyPrefix)

	b := store.Get(types.UndelegationQueueKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUndelegationQueueEntry ...
func (k Keeper) RemoveUndelegationQueueEntry(ctx sdk.Context, undelegationQueueEntry *types.UndelegationQueueEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.UndelegationQueueKeyPrefix)
	store.Delete(types.UndelegationQueueKey(undelegationQueueEntry.Index))

	indexStore := prefix.NewStore(storeAdapter, types.UndelegationQueueKeyPrefixIndex2)
	indexStore.Delete(types.UndelegationQueueKeyIndex2(
		undelegationQueueEntry.Delegator,
		undelegationQueueEntry.Index,
	))
}

// GetAllUnbondingDelegationQueueEntries returns all delegator unbondings
func (k Keeper) GetAllUnbondingDelegationQueueEntries(ctx sdk.Context) (list []types.UndelegationQueueEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.UndelegationQueueKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UndelegationQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllUnbondingDelegationQueueEntriesOfDelegator returns all delegator unbondings of the given address
func (k Keeper) GetAllUnbondingDelegationQueueEntriesOfDelegator(ctx sdk.Context, address string) (list []types.UndelegationQueueEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, util.GetByteKey(types.UndelegationQueueKeyPrefixIndex2, address))
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		index := binary.BigEndian.Uint64(iterator.Key()[0:8])

		entry, _ := k.GetUndelegationQueueEntry(ctx, index)
		list = append(list, entry)
	}

	return
}

// ###################
// === QUEUE STATE ===
// ###################

// GetQueueState returns the state for the undelegation queue
func (k Keeper) GetQueueState(ctx sdk.Context) (state types.QueueState) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	b := store.Get(types.QueueKey)

	if b == nil {
		return state
	}

	k.cdc.MustUnmarshal(b, &state)
	return
}

// SetQueueState saves the undelegation queue state
func (k Keeper) SetQueueState(ctx sdk.Context, state types.QueueState) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	b := k.cdc.MustMarshal(&state)
	store.Set(types.QueueKey, b)
}
