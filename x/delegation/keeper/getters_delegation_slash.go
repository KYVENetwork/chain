package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The `DelegationSlash` entry stores every slash that happened to a staker.
// This is needed by the F1-Fee algorithm to correctly calculate the
// remaining delegation of delegators whose staker got slashed.

// SetDelegationSlashEntry for the affected staker with the index of the period
// the slash is starting.
func (k Keeper) SetDelegationSlashEntry(ctx sdk.Context, slashEntry types.DelegationSlash) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationSlashEntriesKeyPrefix)
	b := k.cdc.MustMarshal(&slashEntry)
	store.Set(types.DelegationEntriesKey(
		slashEntry.Staker,
		slashEntry.KIndex,
	), b)
}

// GetDelegationSlashEntry returns a DelegationSlash for the given staker and index.
func (k Keeper) GetDelegationSlashEntry(
	ctx sdk.Context,
	stakerAddress string,
	kIndex uint64,
) (val types.DelegationSlash, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationSlashEntriesKeyPrefix)

	b := store.Get(types.DelegationSlashEntriesKey(
		stakerAddress,
		kIndex,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDelegationSlashEntry removes an entry for a given staker and index
func (k Keeper) RemoveDelegationSlashEntry(
	ctx sdk.Context,
	stakerAddress string,
	kIndex uint64,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationSlashEntriesKeyPrefix)
	store.Delete(types.DelegationSlashEntriesKey(
		stakerAddress,
		kIndex,
	))
}

// GetAllDelegationSlashEntries returns all delegation slash entries (of all stakers)
func (k Keeper) GetAllDelegationSlashEntries(ctx sdk.Context) (list []types.DelegationSlash) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationSlashEntriesKeyPrefix)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationSlashEntriesKeyPrefix)

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
