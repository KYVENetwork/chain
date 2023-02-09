package keeper

import (
	"encoding/binary"
	"strings"

	"github.com/KYVENetwork/chain/x/pool/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPoolCount get the total number of pools
func (k Keeper) GetPoolCount(ctx sdk.Context) uint64 {
	bz := ctx.KVStore(k.storeKey).Get(types.PoolCountKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// SetPoolCount sets the total number of pools
func (k Keeper) SetPoolCount(ctx sdk.Context, count uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	ctx.KVStore(k.storeKey).Set(types.PoolCountKey, bz)
}

// AppendPool appends a pool in the store with a new id and updates the count
func (k Keeper) AppendPool(ctx sdk.Context, pool types.Pool) uint64 {
	count := k.GetPoolCount(ctx)
	// Set the ID of the appended value
	pool.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)
	appendedValue := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolKeyPrefix(pool.Id), appendedValue)

	// Update pool count
	k.SetPoolCount(ctx, count+1)

	return count
}

// SetPool sets a specific pool in the store
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolKeyPrefix(pool.Id), b)
}

// GetPool returns a pool from its ID
func (k Keeper) GetPool(ctx sdk.Context, id uint64) (val types.Pool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)
	b := store.Get(types.PoolKeyPrefix(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)
	store.Delete(types.PoolKeyPrefix(id))
}

// GetAllPools returns all pools
func (k Keeper) GetAllPools(ctx sdk.Context) (list []types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPaginatedPoolsQuery performs a full search on all pools with the given parameters.
func (k Keeper) GetPaginatedPoolsQuery(
	ctx sdk.Context,
	pagination *query.PageRequest,
	search string,
	runtime string,
	disabled bool,
	storageProviderId uint32,
) ([]types.Pool, *query.PageResponse, error) {
	var pools []types.Pool

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)

	pageRes, err := query.FilteredPaginate(store, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var pool types.Pool
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return false, err
		}

		// filter search
		if !strings.Contains(strings.ToLower(pool.Name), strings.ToLower(search)) {
			return false, nil
		}

		// filter runtime
		if runtime != "" && runtime != pool.Runtime {
			return false, nil
		}

		// filter disabled
		if disabled != pool.Disabled {
			return false, nil
		}

		// filter storage provider id
		if storageProviderId != 0 && storageProviderId != pool.CurrentStorageProviderId {
			return false, nil
		}

		if accumulate {
			pools = append(pools, pool)
		}

		return true, nil
	})
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	return pools, pageRes, nil
}
