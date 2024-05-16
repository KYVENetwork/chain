package stakers

import (
	"cosmossdk.io/store/prefix"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"

	storeTypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllStakers returns all staker
func GetAllStakers(ctx sdk.Context, cdc codec.Codec) (list []Staker) {
	storeService := runtime.NewKVStoreService(storeTypes.NewKVStoreKey(stakersTypes.StoreKey))
	storeAdapter := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, stakersTypes.StakerKeyPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val Staker
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
