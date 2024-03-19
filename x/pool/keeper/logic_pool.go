package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/KYVENetwork/chain/x/pool/types"
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

// ChargeInflationPool charges the inflation pool and transfers the funds to the pool module
// so the payout can be performed
func (k Keeper) ChargeInflationPool(ctx sdk.Context, poolId uint64) (payout uint64, err error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return payout, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrPoolNotFound.Error(), poolId)
	}

	account := pool.GetPoolAccount()
	balance := k.bankKeeper.GetBalance(ctx, account, globalTypes.Denom).Amount.Int64()

	// charge X percent from current pool balance and use it as payout
	payout = uint64(math.LegacyNewDec(balance).Mul(k.GetPoolInflationPayoutRate(ctx)).TruncateInt64())

	// transfer funds to pool module account so bundle reward can be paid out from there
	if err := util.TransferFromAddressToModule(k.bankKeeper, ctx, account.String(), types.ModuleName, payout); err != nil {
		util.PanicHalt(k.upgradeKeeper, ctx, err.Error())
	}

	return payout, nil
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
