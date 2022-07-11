package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.VoteSlash(ctx),
		k.UploadSlash(ctx),
		k.TimeoutSlash(ctx),
		k.UploadTimeout(ctx),
		k.StorageCost(ctx),
		k.NetworkFee(ctx),
		k.MaxPoints(ctx),
		k.UnbondingStakingTime(ctx),
		k.UnbondingStakingTime(ctx),
		k.RedelegationCooldown(ctx),
		k.RedelegationMaxAmount(ctx),
		k.CommissionChangeTime(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// VoteSlash returns the VoteSlash param
func (k Keeper) VoteSlash(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyVoteSlash, &res)
	return
}

// UploadSlash returns the UploadSlash param
func (k Keeper) UploadSlash(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyUploadSlash, &res)
	return
}

// TimeoutSlash returns the TimeoutSlash param
func (k Keeper) TimeoutSlash(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyTimeoutSlash, &res)
	return
}

// UploadTimeout returns the UploadTimeout param
func (k Keeper) UploadTimeout(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyUploadTimeout, &res)
	return
}

// StorageCost returns the StorageCost param
func (k Keeper) StorageCost(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyStorageCost, &res)
	return
}

// NetworkFee returns the NetworkFee param
func (k Keeper) NetworkFee(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyNetworkFee, &res)
	return
}

// MaxPoints returns the MaxPoints param
func (k Keeper) MaxPoints(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxPoints, &res)
	return
}

// UnbondingStakingTime ...
func (k Keeper) UnbondingStakingTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyUnbondingStakingTime, &res)
	return
}

// UnbondingDelegationTime ...
func (k Keeper) UnbondingDelegationTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyUnbondingDelegationTime, &res)
	return
}

// RedelegationCooldown ...
func (k Keeper) RedelegationCooldown(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRedelegationCooldown, &res)
	return
}

// RedelegationMaxAmount ...
func (k Keeper) RedelegationMaxAmount(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRedelegationMaxAmount, &res)
	return
}

// CommissionChangeTime ...
func (k Keeper) CommissionChangeTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyCommissionChangeTime, &res)
	return
}

// ParamStore ...
func (k Keeper) ParamStore() (paramStore paramtypes.Subspace) {
	return k.paramstore
}
