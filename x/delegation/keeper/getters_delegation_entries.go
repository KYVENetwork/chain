package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The `DelegationEntry` stores the quotient of the collected rewards
// and the total delegation of every period. A period is a phase
// where the total delegation was unchanged and just rewards were
// paid out. More details can be found in the specs of this module

// SetDelegationEntry set a specific delegationEntry in the store for the staker
// and a given index
func (k Keeper) SetDelegationEntry(ctx sdk.Context, delegationEntries types.DelegationEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationEntriesKeyPrefix)
	b := k.cdc.MustMarshal(&delegationEntries)
	store.Set(types.DelegationEntriesKey(
		delegationEntries.Staker,
		delegationEntries.KIndex,
	), b)
}

// GetDelegationEntry returns a delegationEntry from its index
func (k Keeper) GetDelegationEntry(
	ctx sdk.Context,
	stakerAddress string,
	kIndex uint64,
) (val types.DelegationEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationEntriesKeyPrefix)

	b := store.Get(types.DelegationEntriesKey(
		stakerAddress,
		kIndex,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDelegationEntry removes a delegationEntry for the given staker with the
// given index from the store
func (k Keeper) RemoveDelegationEntry(
	ctx sdk.Context,
	stakerAddress string,
	kIndex uint64,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationEntriesKeyPrefix)
	store.Delete(types.DelegationEntriesKey(
		stakerAddress,
		kIndex,
	))
}

// GetAllDelegationEntries returns all delegationEntries (of all stakers)
func (k Keeper) GetAllDelegationEntries(ctx sdk.Context) (list []types.DelegationEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationEntriesKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DelegationEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
