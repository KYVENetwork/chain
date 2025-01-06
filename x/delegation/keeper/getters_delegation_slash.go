package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The `DelegationSlash` entry stores every slash that happened to a staker.
// This is needed by the F1-Fee algorithm to correctly calculate the
// remaining delegation of delegators whose staker got slashed.

// GetAllDelegationSlashEntries returns all delegation slash entries (of all stakers)
func (k Keeper) GetAllDelegationSlashEntries(ctx sdk.Context) (list []types.DelegationSlash) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.DelegationSlashEntriesKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DelegationSlash
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllDelegationSlashesBetween returns all Slashes that happened between the given periods
// `start` and `end` are both inclusive.
func (k Keeper) GetAllDelegationSlashesBetween(ctx sdk.Context, staker string, start uint64, end uint64) (list []types.DelegationSlash) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.DelegationSlashEntriesKeyPrefix)

	// use iterator with end+1 because the end of the iterator is exclusive
	iterator := store.Iterator(util.GetByteKey(staker, start), util.GetByteKey(staker, end+1))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DelegationSlash
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
