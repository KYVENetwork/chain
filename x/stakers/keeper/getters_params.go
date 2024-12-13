package keeper

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/stakers module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetCommissionChangeTime returns the CommissionChangeTime param
func (k Keeper) GetCommissionChangeTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).CommissionChangeTime
}

// GetLeavePoolTime returns the LeavePoolTime param
func (k Keeper) GetLeavePoolTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).LeavePoolTime
}

// GetStakeFractionChangeTime returns the StakeFractionChangeTime param
func (k Keeper) GetStakeFractionChangeTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).StakeFractionChangeTime
}

// GetVoteSlash returns the VoteSlash param
func (k Keeper) GetVoteSlash(ctx sdk.Context) (res math.LegacyDec) {
	return k.GetParams(ctx).VoteSlash
}

// GetUploadSlash returns the UploadSlash param
func (k Keeper) GetUploadSlash(ctx sdk.Context) (res math.LegacyDec) {
	return k.GetParams(ctx).UploadSlash
}

// GetTimeoutSlash returns the TimeoutSlash param
func (k Keeper) GetTimeoutSlash(ctx sdk.Context) (res math.LegacyDec) {
	return k.GetParams(ctx).TimeoutSlash
}

func (k Keeper) getSlashFraction(ctx sdk.Context, slashType types.SlashType) (slashAmountRatio math.LegacyDec) {
	// Retrieve slash fraction from params
	switch slashType {
	case types.SLASH_TYPE_TIMEOUT:
		slashAmountRatio = k.GetTimeoutSlash(ctx)
	case types.SLASH_TYPE_VOTE:
		slashAmountRatio = k.GetVoteSlash(ctx)
	case types.SLASH_TYPE_UPLOAD:
		slashAmountRatio = k.GetUploadSlash(ctx)
	}
	return
}

// SetParams sets the x/stakers module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
