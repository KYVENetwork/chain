package funders

import (
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/funders/types"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func GetParams(ctx sdk.Context, cdc codec.Codec) (params Params) {
	storeService := runtime.NewKVStoreService(storeTypes.NewKVStoreKey(fundersTypes.StoreKey))
	store := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	cdc.MustUnmarshal(bz, &params)
	return params
}
