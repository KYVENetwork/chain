package keeper

import (
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the x/pool params from state.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	bz := ctx.KVStore(k.storeKey).Get(types.ParamsKey)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &params)
	}

	return
}

// GetProtocolInflationShare returns the ProtocolInflationShare param.
func (k Keeper) GetProtocolInflationShare(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).ProtocolInflationShare
}

// GetPoolInflationPayoutRate returns the GetPoolInflationPayoutRate param
func (k Keeper) GetPoolInflationPayoutRate(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).PoolInflationPayoutRate
}

// SetParams stores the x/pool params in state.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	bz := k.cdc.MustMarshal(&params)
	ctx.KVStore(k.storeKey).Set(types.ParamsKey, bz)
}
