package keeper

import (
	"encoding/binary"

	"cosmossdk.io/math"

	storeTypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IncrementPoints increments to Points for a staker in a given pool.
// Returns the amount of the current points (including the current incrementation)
func (k Keeper) IncrementPoints(ctx sdk.Context, stakerAddress string, poolId uint64) uint64 {
	poolAccount, found := k.GetPoolAccount(ctx, stakerAddress, poolId)
	if found {
		poolAccount.Points += 1
		k.SetPoolAccount(ctx, poolAccount)
	}
	return poolAccount.Points
}

// ResetPoints sets the point count for the staker in the given pool back to zero.
// Returns the amount of points the staker had before the reset.
func (k Keeper) ResetPoints(ctx sdk.Context, stakerAddress string, poolId uint64) (previousPoints uint64) {
	poolAccount, found := k.GetPoolAccount(ctx, stakerAddress, poolId)
	if found {
		previousPoints = poolAccount.Points
		poolAccount.Points = 0
		k.SetPoolAccount(ctx, poolAccount)
	}
	return
}

// GetAllPoolAccountsOfPool returns a list of all pool accounts
func (k Keeper) GetAllPoolAccountsOfPool(ctx sdk.Context, poolId uint64) (val []*types.PoolAccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PoolAccountPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, util.GetByteKey(poolId))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolAccount := types.PoolAccount{}
		k.cdc.MustUnmarshal(iterator.Value(), &poolAccount)

		if poolAccount.PoolAddress != "" {
			val = append(val, &poolAccount)
		}
	}

	return
}

// GetPoolAccountsFromStaker returns all pools the staker has pool accounts in
func (k Keeper) GetPoolAccountsFromStaker(ctx sdk.Context, stakerAddress string) (val []*types.PoolAccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	storeIndex2 := prefix.NewStore(storeAdapter, types.PoolAccountPrefixIndex2)

	iterator := storeTypes.KVStorePrefixIterator(storeIndex2, util.GetByteKey(stakerAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolId := binary.BigEndian.Uint64(iterator.Key()[len(stakerAddress) : len(stakerAddress)+8])
		poolAccount, active := k.GetPoolAccount(ctx, stakerAddress, poolId)

		if active {
			val = append(val, &poolAccount)
		}
	}

	return val
}

// GetPoolCount returns the number of pools the current staker is
// currently participating.
func (k Keeper) GetPoolCount(ctx sdk.Context, stakerAddress string) (poolCount uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	storeIndex2 := prefix.NewStore(storeAdapter, types.PoolAccountPrefixIndex2)
	iterator := storeTypes.KVStorePrefixIterator(storeIndex2, util.GetByteKey(stakerAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolCount += 1
	}
	return
}

// #############################
// #  Raw KV-Store operations  #
// #############################

// SetPoolAccount set a specific pool account in the store from its index
func (k Keeper) SetPoolAccount(ctx sdk.Context, poolAccount types.PoolAccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PoolAccountPrefix)
	b := k.cdc.MustMarshal(&poolAccount)
	store.Set(types.PoolAccountKey(
		poolAccount.PoolId,
		poolAccount.Staker,
	), b)

	storeIndex2 := prefix.NewStore(storeAdapter, types.PoolAccountPrefixIndex2)
	storeIndex2.Set(types.PoolAccountKeyIndex2(
		poolAccount.Staker,
		poolAccount.PoolId,
	), []byte{})
}

// GetPoolAccount returns a pool account from its index
func (k Keeper) GetPoolAccount(ctx sdk.Context, stakerAddress string, poolId uint64) (val types.PoolAccount, active bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PoolAccountPrefix)

	b := store.Get(types.PoolAccountKey(
		poolId,
		stakerAddress,
	))
	if b == nil {
		val.Commission = math.LegacyZeroDec()
		val.StakeFraction = math.LegacyZeroDec()
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, val.PoolAddress != ""
}

// GetAllPoolAccounts returns all active pool accounts
func (k Keeper) GetAllPoolAccounts(ctx sdk.Context) (list []types.PoolAccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.PoolAccountPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolAccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.PoolAddress != "" {
			list = append(list, val)
		}
	}

	return
}
