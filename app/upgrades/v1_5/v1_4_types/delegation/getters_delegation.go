package delegation

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetAllDelegationEntries(ctx sdk.Context, cdc codec.Codec, storeKey storeTypes.StoreKey) (list []DelegationEntry) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.DelegationEntriesKeyPrefix)
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
func GetAllDelegationData(ctx sdk.Context, cdc codec.Codec, storeKey storeTypes.StoreKey) (list []DelegationData) {
	store := ctx.KVStore(storeKey)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val DelegationData
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
