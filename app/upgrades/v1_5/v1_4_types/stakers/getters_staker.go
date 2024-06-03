package stakers

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllStakers returns all staker
func GetAllStakers(ctx sdk.Context, cdc codec.Codec, storeKey storeTypes.StoreKey) (list []Staker) {
	store := prefix.NewStore(ctx.KVStore(storeKey), stakersTypes.StakerKeyPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val Staker
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
