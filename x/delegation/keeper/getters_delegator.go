package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// `Delegator` is created for every delegator (address) that delegates
// to a staker. It stores the initial amount delegated and the index
// of the F1-period where the user started to become a delegator.
// When the user performs a redelegation this object is recreated.
// To query the current delegation use `GetDelegationAmountOfDelegator()`
// as the `initialAmount` does not consider slashes.

// SetDelegator set a specific delegator in the store from its index
func (k Keeper) SetDelegator(ctx sdk.Context, delegator types.Delegator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefix)
	b := k.cdc.MustMarshal(&delegator)
	store.Set(types.DelegatorKey(
		delegator.Staker,
		delegator.Delegator,
	), b)

	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefixIndex2)
	indexStore.Set(types.DelegatorKeyIndex2(
		delegator.Delegator,
		delegator.Staker,
	), []byte{1})
}

// GetDelegator returns a delegator from its index
func (k Keeper) GetDelegator(
	ctx sdk.Context,
	stakerAddress string,
	delegatorAddress string,
) (val types.Delegator, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefix)
	b := store.Get(types.DelegatorKey(stakerAddress, delegatorAddress))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// DoesDelegatorExist checks if the key exists in the KV-store
func (k Keeper) DoesDelegatorExist(
	ctx sdk.Context,
	stakerAddress string,
	delegatorAddress string,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefix)
	return store.Has(types.DelegatorKey(stakerAddress, delegatorAddress))
}

// RemoveDelegator removes a delegator from the store
func (k Keeper) RemoveDelegator(
	ctx sdk.Context,
	stakerAddress string,
	delegatorAddress string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefix)
	store.Delete(types.DelegatorKey(
		stakerAddress,
		delegatorAddress,
	))
	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefixIndex2)
	indexStore.Delete(types.DelegatorKeyIndex2(
		delegatorAddress,
		stakerAddress,
	))
}

// GetAllDelegators returns all delegators (of all stakers)
func (k Keeper) GetAllDelegators(ctx sdk.Context) (list []types.Delegator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Delegator
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetStakersByDelegator(ctx sdk.Context, delegator string) (list []string) {
	delegatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegatorKeyPrefixIndex2)
	iterator := sdk.KVStorePrefixIterator(delegatorStore, util.GetByteKey(delegator))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		staker := string(iterator.Key()[43 : 43+43])
		list = append(list, staker)
	}
	return
}
