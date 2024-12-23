package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/math"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	storeTypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddPoolAccountToPool adds a pool account to a pool.
// If pool account already active in the to pool nothing happens.
func (k Keeper) AddPoolAccountToPool(ctx sdk.Context, stakerAddress string, poolId uint64, poolAddress string, commission, stakeFraction math.LegacyDec) {
	if _, validatorExists := k.GetValidator(ctx, stakerAddress); validatorExists {
		if _, active := k.GetPoolAccount(ctx, stakerAddress, poolId); !active {
			k.SetPoolAccount(ctx, types.PoolAccount{
				PoolId:        poolId,
				Staker:        stakerAddress,
				PoolAddress:   poolAddress,
				Commission:    commission,
				StakeFraction: stakeFraction,
			})
			k.AddOneToCount(ctx, poolId)
			k.AddActiveStaker(ctx, stakerAddress)
		}
	}
}

// RemovePoolAccountFromPool removes a pool account from a given pool and updates
// all aggregated variables. If the pool account is not in the pool nothing happens.
func (k Keeper) RemovePoolAccountFromPool(ctx sdk.Context, stakerAddress string, poolId uint64) {
	if poolAccount, active := k.GetPoolAccount(ctx, stakerAddress, poolId); active {
		// remove pool account from pool by setting pool address to zero address
		poolAccount.PoolAddress = ""
		poolAccount.Points = 0
		poolAccount.IsLeaving = false
		k.SetPoolAccount(ctx, poolAccount)
		k.subtractOneFromCount(ctx, poolId)
		k.removeActiveStaker(ctx, stakerAddress)
	}
}

// #############################
// #  Raw KV-Store operations  #
// #############################

func (k Keeper) getAllStakersOfPool(ctx sdk.Context, poolId uint64) []stakingTypes.Validator {
	poolAccounts := k.GetAllPoolAccountsOfPool(ctx, poolId)

	stakers := make([]stakingTypes.Validator, 0)

	for _, poolAccount := range poolAccounts {
		staker, _ := k.GetValidator(ctx, poolAccount.Staker)
		stakers = append(stakers, staker)
	}

	return stakers
}

// GetLegacyStaker returns a staker from its index
func (k Keeper) GetLegacyStaker(
	ctx sdk.Context,
	staker string,
) (val types.Staker, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakerKeyPrefix)

	b := store.Get(types.StakerKey(
		staker,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetPaginatedLegacyStakerQuery(
	ctx sdk.Context,
	pagination *query.PageRequest,
	accumulator func(staker types.Staker),
) (*query.PageResponse, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakerKeyPrefix)

	pageRes, err := query.FilteredPaginate(store, pagination, func(
		key []byte,
		value []byte,
		accumulate bool,
	) (bool, error) {
		if accumulate {
			var staker types.Staker
			if err := k.cdc.Unmarshal(value, &staker); err != nil {
				return false, err
			}
			accumulator(staker)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return pageRes, nil
}

// GetAllStakers returns all staker
func (k Keeper) GetAllLegacyStakers(ctx sdk.Context) (list []types.Staker) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.StakerKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Staker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// #############################
// #     Aggregation Data      #
// #############################

func (k Keeper) GetStakerCountOfPool(ctx sdk.Context, poolId uint64) uint64 {
	return k.getStat(ctx, poolId, types.STAKER_STATS_COUNT)
}

func (k Keeper) AddOneToCount(ctx sdk.Context, poolId uint64) {
	count := k.getStat(ctx, poolId, types.STAKER_STATS_COUNT)
	k.setStat(ctx, poolId, types.STAKER_STATS_COUNT, count+1)
}

func (k Keeper) subtractOneFromCount(ctx sdk.Context, poolId uint64) {
	count := k.getStat(ctx, poolId, types.STAKER_STATS_COUNT)
	k.setStat(ctx, poolId, types.STAKER_STATS_COUNT, count-1)
}

// getStat get the total number of pool
func (k Keeper) getStat(ctx sdk.Context, poolId uint64, statType types.STAKER_STATS) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(util.GetByteKey(string(statType), poolId))
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// setStat set the total number of pool
func (k Keeper) setStat(ctx sdk.Context, poolId uint64, statType types.STAKER_STATS, count uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(util.GetByteKey(string(statType), poolId), bz)
}

// #############################
// #      Active Staker        #
// #############################
// Active Staker stores all stakers that are at least in one pool

// AddActiveStaker increases the active-staker-count of the given staker by one.
// The amount tracks the number of pools the staker is in. It also allows
// to determine that a given staker is at least in one pool.
func (k Keeper) AddActiveStaker(ctx sdk.Context, staker string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ActiveStakerIndex)
	// Get current count
	count := uint64(0)
	storeBytes := store.Get(types.ActiveStakerKeyIndex(staker))
	bytes := make([]byte, 8)
	copy(bytes, storeBytes)
	if len(bytes) == 8 {
		count = binary.BigEndian.Uint64(bytes)
	} else {
		bytes = make([]byte, 8)
	}
	// Count represents in how many pools the current staker is active
	count += 1

	// Encode and store
	binary.BigEndian.PutUint64(bytes, count)
	store.Set(types.ActiveStakerKeyIndex(staker), bytes)
}

// removeActiveStaker decrements the active-staker-count of the given staker
// by one. If the amount drop to zero the staker is removed from the set.
// Therefore, one can be sure, that only stakers which are participating in
// at least one pool are part of the set
func (k Keeper) removeActiveStaker(ctx sdk.Context, staker string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ActiveStakerIndex)
	// Get current count
	count := uint64(0)
	storeBytes := store.Get(types.ActiveStakerKeyIndex(staker))
	bytes := make([]byte, 8)
	copy(bytes, storeBytes)

	if len(bytes) == 8 {
		count = binary.BigEndian.Uint64(bytes)
	} else {
		bytes = make([]byte, 8)
	}

	if count == 0 || count == 1 {
		store.Delete(types.ActiveStakerKeyIndex(staker))
		return
	}

	// Count represents in how many pools the current staker is active
	count -= 1

	// Encode and store
	binary.BigEndian.PutUint64(bytes, count)
	store.Set(types.ActiveStakerKeyIndex(staker), bytes)
}

// getAllActiveStakers returns all active stakers, i.e. every staker
// that is member of at least one pool.
func (k Keeper) getAllActiveStakers(ctx sdk.Context) (list []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.ActiveStakerIndex)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		list = append(list, string(iterator.Key()))
	}

	return
}

// arrayPaginationAccumulator helps to parse the query.PageRequest for an array
// instead of a KV-Store.
func arrayPaginationAccumulator(slice []string, pagination *query.PageRequest, accumulator func(address string, accumulate bool) bool) (*query.PageResponse, error) {
	if pagination != nil && pagination.Key != nil {
		return nil, fmt.Errorf("key pagination not supported")
	}

	page, limit, err := query.ParsePagination(pagination)
	if err != nil {
		return nil, err
	}

	count := 0
	minIndex := (page - 1) * limit
	maxIndex := (page) * limit

	for i := 0; i < len(slice); i++ {
		if accumulator(slice[i], count >= minIndex && count < maxIndex) {
			count++
		}
	}

	pageRes := &query.PageResponse{
		NextKey: nil,
		Total:   uint64(count),
	}

	return pageRes, nil
}
