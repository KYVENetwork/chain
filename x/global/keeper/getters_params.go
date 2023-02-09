package keeper

import (
	"github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/global module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetMinGasPrice returns the MinGasPrice param.
func (k Keeper) GetMinGasPrice(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).MinGasPrice
}

// GetBurnRatio returns the BurnRatio param.
func (k Keeper) GetBurnRatio(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).BurnRatio
}

// GetGasAdjustments returns the GasAdjustments param.
func (k Keeper) GetGasAdjustments(ctx sdk.Context) (res []types.GasAdjustment) {
	return k.GetParams(ctx).GasAdjustments
}

// GetGasRefunds returns the GasRefunds param.
func (k Keeper) GetGasRefunds(ctx sdk.Context) (res []types.GasRefund) {
	return k.GetParams(ctx).GasRefunds
}

// GetMinInitialDepositRatio returns the MinInitialDepositRatio param.
func (k Keeper) GetMinInitialDepositRatio(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).MinInitialDepositRatio
}

// SetParams sets the x/global module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
