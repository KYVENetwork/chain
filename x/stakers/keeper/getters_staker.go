package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"

	"cosmossdk.io/math"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	storeTypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
