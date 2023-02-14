package keeper

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/pool/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetPoolWithError returns a pool by its poolId, if the pool does not exist,
// a types.ErrPoolNotFound error is returned
func (k Keeper) GetPoolWithError(ctx sdk.Context, poolId uint64) (types.Pool, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return types.Pool{}, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrPoolNotFound.Error(), poolId)
	}
	return pool, nil
}

// AssertPoolExists returns nil if the pool exists and types.ErrPoolNotFound if it does not.
func (k Keeper) AssertPoolExists(ctx sdk.Context, poolId uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)
	if store.Has(types.PoolKeyPrefix(poolId)) {
		return nil
	}
	return errors.Wrapf(errorsTypes.ErrNotFound, types.ErrPoolNotFound.Error(), poolId)
}

// IncrementBundleInformation updates the latest finalized bundle of a pool
func (k Keeper) IncrementBundleInformation(
	ctx sdk.Context,
	poolId uint64,
	currentIndex uint64,
	currentKey string,
	currentSummary string,
) {
	pool, found := k.GetPool(ctx, poolId)
	if found {
		pool.CurrentIndex = currentIndex
		pool.TotalBundles = pool.TotalBundles + 1
		pool.CurrentKey = currentKey
		pool.CurrentSummary = currentSummary
		k.SetPool(ctx, pool)
	}
}
