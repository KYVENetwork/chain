package pool

import (
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
