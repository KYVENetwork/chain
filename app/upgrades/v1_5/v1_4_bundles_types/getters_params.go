package v1_4_bundles_types

import (
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/bundles module parameters.
func GetParams(ctx sdk.Context, storeKey storeTypes.StoreKey, cdc codec.Codec) (params Params) {
	store := ctx.KVStore(storeKey)

	bz := store.Get(types.ParamsPrefix)
	if bz == nil {
		return params
	}

	cdc.MustUnmarshal(bz, &params)
	return params
}
