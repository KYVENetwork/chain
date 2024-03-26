package keeper

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/delegation module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetUnbondingDelegationTime returns the UnbondingDelegationTime param
func (k Keeper) GetUnbondingDelegationTime(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).UnbondingDelegationTime
}

// GetRedelegationCooldown returns the RedelegationCooldown param
func (k Keeper) GetRedelegationCooldown(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).RedelegationCooldown
}

// GetRedelegationMaxAmount returns the RedelegationMaxAmount param
func (k Keeper) GetRedelegationMaxAmount(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).RedelegationMaxAmount
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

// SetParams sets the x/delegation module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}
