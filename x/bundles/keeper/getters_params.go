package keeper

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the current x/bundles module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	params, err := k.BundlesParams.Get(ctx)
	if err != nil {
		return types.DefaultParams()
	}
	return params
}

// GetUploadTimeout returns the UploadTimeout param
func (k Keeper) GetUploadTimeout(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).UploadTimeout
}

// GetStorageCost returns the StorageCost param
func (k Keeper) GetStorageCost(ctx sdk.Context, storageProviderId uint32) (res math.LegacyDec) {
	storageCosts := k.GetParams(ctx).StorageCosts
	for _, storageCost := range storageCosts {
		if storageCost.StorageProviderId == storageProviderId {
			return storageCost.Cost
		}
	}
	return math.LegacyZeroDec()
}

// GetNetworkFee returns the NetworkFee param
func (k Keeper) GetNetworkFee(ctx sdk.Context) (res math.LegacyDec) {
	return k.GetParams(ctx).NetworkFee
}

// GetMaxPoints returns the MaxPoints param
func (k Keeper) GetMaxPoints(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).MaxPoints
}

// SetParams sets the x/bundles module parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	_ = k.BundlesParams.Set(ctx, params)
}
