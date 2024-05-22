package v1_4_pool_types

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllPools returns all pools
func GetAllPools(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.Codec) (list []Pool) {
	store := ctx.KVStore(storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val Pool
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
