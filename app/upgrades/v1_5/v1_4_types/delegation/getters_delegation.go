package delegation

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetAllDelegationEntries(ctx sdk.Context, cdc codec.Codec) (list []DelegationEntry) {
	storeService := runtime.NewKVStoreService(storeTypes.NewKVStoreKey(delegationTypes.StoreKey))
	storeAdapter := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, delegationTypes.DelegationEntriesKeyPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val DelegationEntry
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllDelegationData returns all delegationData entries
func GetAllDelegationData(ctx sdk.Context, cdc codec.Codec) (list []DelegationData) {
	storeService := runtime.NewKVStoreService(storeTypes.NewKVStoreKey(delegationTypes.StoreKey))
	storeAdapter := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, delegationTypes.DelegationDataKeyPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val DelegationData
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
