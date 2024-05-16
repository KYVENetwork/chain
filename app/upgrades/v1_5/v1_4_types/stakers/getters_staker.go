package stakers

import (
	"github.com/cosmos/cosmos-sdk/codec"

	storeTypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllStakers returns all staker
func GetAllStakers(ctx sdk.Context, cdc codec.Codec, storeKey storeTypes.StoreKey) (list []Staker) {
	store := ctx.KVStore(storeKey)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val Staker
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
