package bundles

import (
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/bundles/types"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/bundles module parameters.
func GetParams(ctx sdk.Context, cdc codec.Codec) (params Params) {
	storeService := runtime.NewKVStoreService(storeTypes.NewKVStoreKey(bundlesTypes.StoreKey))
	store := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))

	bz := store.Get(types.ParamsPrefix)
	if bz == nil {
		return params
	}

	cdc.MustUnmarshal(bz, &params)
	return params
}
