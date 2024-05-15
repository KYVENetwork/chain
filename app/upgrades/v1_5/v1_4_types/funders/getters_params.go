package funders

import (
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func GetParams(ctx sdk.Context, storeKey storeTypes.StoreKey, cdc codec.Codec) (params Params) {
	store := ctx.KVStore(storeKey)

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	cdc.MustUnmarshal(bz, &params)
	return params
}
