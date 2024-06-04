package pool

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/pool/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the x/pool params from state.
func GetParams(ctx sdk.Context, cdc codec.Codec, storeKey storeTypes.StoreKey) (params types.Params) {
	bz := ctx.KVStore(storeKey).Get(types.ParamsKey)
	if bz != nil {
		cdc.MustUnmarshal(bz, &params)
	}

	return
}

// GetAllPools returns all pools
func GetAllPools(ctx sdk.Context, storeKey storeTypes.StoreKey, cdc codec.Codec) (list []Pool) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.PoolKey)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val Pool
		cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
