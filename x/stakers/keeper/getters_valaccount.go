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
func (k Keeper) IncrementPoints(ctx sdk.Context, poolId uint64, stakerAddress string) uint64 {
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)
	if found {
		valaccount.Points += 1
		k.SetValaccount(ctx, valaccount)
	}
	return valaccount.Points
}

// ResetPoints sets the point count for the staker in the given pool back to zero.
// Returns the amount of points the staker had before the reset.
func (k Keeper) ResetPoints(ctx sdk.Context, poolId uint64, stakerAddress string) (previousPoints uint64) {
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)
	if found {
		previousPoints = valaccount.Points
		valaccount.Points = 0
		k.SetValaccount(ctx, valaccount)
	}
	return
}

// GetAllValaccountsOfPool returns a list of all valaccount
func (k Keeper) GetAllValaccountsOfPool(ctx sdk.Context, poolId uint64) (val []*types.Valaccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ValaccountPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, util.GetByteKey(poolId))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		valaccount := types.Valaccount{}
		k.cdc.MustUnmarshal(iterator.Value(), &valaccount)

		if valaccount.Valaddress != "" {
			val = append(val, &valaccount)
		}
	}

	return
}

// GetValaccountsFromStaker returns all pools the staker has valaccounts in
func (k Keeper) GetValaccountsFromStaker(ctx sdk.Context, stakerAddress string) (val []*types.Valaccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	storeIndex2 := prefix.NewStore(storeAdapter, types.ValaccountPrefixIndex2)

	iterator := storeTypes.KVStorePrefixIterator(storeIndex2, util.GetByteKey(stakerAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolId := binary.BigEndian.Uint64(iterator.Key()[len(stakerAddress) : len(stakerAddress)+8])
		valaccount, active := k.GetValaccount(ctx, poolId, stakerAddress)

		if active {
			val = append(val, &valaccount)
		}
	}

	return val
}

// GetPoolCount returns the number of pools the current staker is
// currently participating.
func (k Keeper) GetPoolCount(ctx sdk.Context, stakerAddress string) (poolCount uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	storeIndex2 := prefix.NewStore(storeAdapter, types.ValaccountPrefixIndex2)
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

// SetValaccount set a specific Valaccount in the store from its index
func (k Keeper) SetValaccount(ctx sdk.Context, valaccount types.Valaccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ValaccountPrefix)
	b := k.cdc.MustMarshal(&valaccount)
	store.Set(types.ValaccountKey(
		valaccount.PoolId,
		valaccount.Staker,
	), b)

	storeIndex2 := prefix.NewStore(storeAdapter, types.ValaccountPrefixIndex2)
	storeIndex2.Set(types.ValaccountKeyIndex2(
		valaccount.Staker,
		valaccount.PoolId,
	), []byte{})
}

// GetValaccount returns a Valaccount from its index
func (k Keeper) GetValaccount(ctx sdk.Context, poolId uint64, stakerAddress string) (val types.Valaccount, active bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ValaccountPrefix)

	b := store.Get(types.ValaccountKey(
		poolId,
		stakerAddress,
	))
	if b == nil {
		val.Commission = math.LegacyZeroDec()
		val.StakeFraction = math.LegacyZeroDec()
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, val.Valaddress != ""
}

// GetAllValaccounts returns all active valaccounts
func (k Keeper) GetAllValaccounts(ctx sdk.Context) (list []types.Valaccount) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ValaccountPrefix)

	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Valaccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.Valaddress != "" {
			list = append(list, val)
		}
	}

	return
}
