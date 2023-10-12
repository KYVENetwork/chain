package v1_3_types

import (
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	"github.com/KYVENetwork/chain/x/pool/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllPools returns all pools
func GetAllPools(ctx sdk.Context, poolKeeper poolKeeper.Keeper, cdc codec.BinaryCodec) (list []Pool, err error) {
	store := prefix.NewStore(ctx.KVStore(poolKeeper.StoreKey()), types.PoolKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val Pool
		err = cdc.Unmarshal(iterator.Value(), &val)
		if err != nil {
			return
		}
		list = append(list, val)
	}

	return
}
