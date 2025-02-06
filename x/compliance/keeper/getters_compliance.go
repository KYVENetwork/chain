package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/compliance/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// #####################
// === QUEUE ENTRIES ===
// #####################

// SetMultiCoinPendingRewardsEntry ...
func (k Keeper) SetMultiCoinPendingRewardsEntry(ctx sdk.Context, compliancePendingRewards types.MultiCoinPendingRewardsEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.MultiCoinPendingRewardsEntryKeyPrefix)
	b := k.cdc.MustMarshal(&compliancePendingRewards)
	store.Set(types.MultiCoinPendingRewardsKeyEntry(
		compliancePendingRewards.Index,
	), b)

	// Insert the same entry with a different key prefix for query lookup
	indexStore := prefix.NewStore(storeAdapter, types.MultiCoinPendingRewardsEntryKeyPrefixIndex2)
	indexStore.Set(types.MultiCoinPendingRewardsKeyEntryIndex2(
		compliancePendingRewards.Address,
		compliancePendingRewards.Index,
	), []byte{1})
}

// GetMultiCoinPendingRewardsEntry ...
func (k Keeper) GetMultiCoinPendingRewardsEntry(ctx sdk.Context, index uint64) (val types.MultiCoinPendingRewardsEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.MultiCoinPendingRewardsEntryKeyPrefix)

	b := store.Get(types.MultiCoinPendingRewardsKeyEntry(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetMultiCoinPendingRewardsEntriesByIndex2 ...
func (k Keeper) GetMultiCoinPendingRewardsEntriesByIndex2(ctx sdk.Context, address string) (list []types.MultiCoinPendingRewardsEntry, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.MultiCoinPendingRewardsEntryKeyPrefixIndex2)

	iterator := storeTypes.KVStorePrefixIterator(store, util.GetByteKey(address))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MultiCoinPendingRewardsEntry
		println("key")
		println(string(iterator.Key()))
		println(string(iterator.Value()))
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// RemoveMultiCoinPendingRewardsEntry ...
func (k Keeper) RemoveMultiCoinPendingRewardsEntry(ctx sdk.Context, entry *types.MultiCoinPendingRewardsEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.MultiCoinPendingRewardsEntryKeyPrefix)
	store.Delete(types.MultiCoinPendingRewardsKeyEntry(entry.Index))

	indexStore := prefix.NewStore(storeAdapter, types.MultiCoinPendingRewardsEntryKeyPrefixIndex2)
	indexStore.Delete(types.MultiCoinPendingRewardsKeyEntryIndex2(
		entry.Address,
		entry.Index,
	))
}

// GetAllMultiCoinPendingRewardsEntries ...
func (k Keeper) GetAllMultiCoinPendingRewardsEntries(ctx sdk.Context) (list []types.MultiCoinPendingRewardsEntry) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.MultiCoinPendingRewardsEntryKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MultiCoinPendingRewardsEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetAllEnabledMultiCoinAddresses(ctx sdk.Context) []string {
	addresses := make([]string, 0)
	if iter, err := k.MultiCoinRewardsEnabled.Iterate(ctx, nil); err == nil {
		if accounts, err := iter.Keys(); err == nil {
			for _, account := range accounts {
				addresses = append(addresses, account.String())
			}
		}
	}
	return addresses
}
