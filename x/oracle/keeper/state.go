package keeper

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	bz := ctx.KVStore(k.storeKey).Get(types.ParamsKey)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &params)
	}

	return
}

func (k Keeper) GetPricePerByte(ctx sdk.Context) math.LegacyDec {
	return k.GetParams(ctx).PricePerByte
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	bz := k.cdc.MustMarshal(&params)
	ctx.KVStore(k.storeKey).Set(types.ParamsKey, bz)
}
