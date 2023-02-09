package keeper

import (
	"fmt"
	"math"
	"sort"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// arrayPagination helps to parse the query.PageRequest for an array
// instead of a KV-Store.
func arrayPagination(slice []string, pagination *query.PageRequest) ([]string, *query.PageResponse, error) {
	if pagination != nil && pagination.Key != nil {
		return nil, nil, fmt.Errorf("key pagination not supported")
	}

	page, limit, err := query.ParsePagination(pagination)
	if err != nil {
		return nil, nil, err
	}

	resultLength := util.MinInt(limit, len(slice)-page*limit)
	result := make([]string, resultLength)

	for i := 0; i < resultLength; i++ {
		result[i] = slice[page*limit+i]
	}

	pageRes := &query.PageResponse{
		NextKey: nil,
		Total:   uint64(len(slice)),
	}

	return result, pageRes, nil
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

// SetStakerIndex sets and Index-entry which sorts all stakers (active and passive)
// by its total delegation
func (k Keeper) SetStakerIndex(ctx sdk.Context, staker string) {
	amount := k.GetDelegationAmount(ctx, staker)
	store := prefix.NewStore(ctx.KVStore(k.memKey), types.StakerIndexKeyPrefix)
	store.Set(types.StakerIndexKey(math.MaxUint64-amount, staker), []byte{0})
}

// RemoveStakerIndex deletes and Index-entry which sorts all stakers (active and passive)
// by its total delegation
func (k Keeper) RemoveStakerIndex(ctx sdk.Context, staker string) {
	amount := k.GetDelegationAmount(ctx, staker)
	store := prefix.NewStore(ctx.KVStore(k.memKey), types.StakerIndexKeyPrefix)
	store.Delete(types.StakerIndexKey(math.MaxUint64-amount, staker))
}

// GetPaginatedStakersByDelegation returns all stakers (active and inactive)
// sorted by its current total delegation. It supports the cosmos query.PageRequest pagination.
func (k Keeper) GetPaginatedStakersByDelegation(ctx sdk.Context, pagination *query.PageRequest, accumulator func(staker string, accumulate bool) bool) (*query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.memKey), types.StakerIndexKeyPrefix)

	pageRes, err := query.FilteredPaginate(store, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		address := string(key[8 : 8+43])
		return accumulator(address, accumulate), nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return pageRes, nil
}

// GetPaginatedActiveStakersByDelegation returns all active stakers
// sorted by its current total delegation. It supports the cosmos query.PageRequest pagination.
func (k Keeper) GetPaginatedActiveStakersByDelegation(ctx sdk.Context, pagination *query.PageRequest, accumulator func(staker string, accumulate bool) bool) (*query.PageResponse, error) {
	activeStakers := k.stakersKeeper.GetActiveStakers(ctx)

	sort.Slice(activeStakers, func(i, j int) bool {
		return k.GetDelegationAmount(ctx, activeStakers[i]) > k.GetDelegationAmount(ctx, activeStakers[j])
	})

	pageRes, err := arrayPaginationAccumulator(activeStakers, pagination, accumulator)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return pageRes, nil
}

// GetPaginatedInactiveStakersByDelegation returns all inactive stakers
// sorted by its current total delegation. It supports the cosmos query.PageRequest pagination.
func (k Keeper) GetPaginatedInactiveStakersByDelegation(ctx sdk.Context, pagination *query.PageRequest, accumulator func(staker string, accumulate bool) bool) (*query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.memKey), types.StakerIndexKeyPrefix)

	pageRes, err := query.FilteredPaginate(store, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		address := string(key[8 : 8+43])
		if k.stakersKeeper.GetPoolCount(ctx, address) > 0 {
			return false, nil
		}
		return accumulator(address, accumulate), nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return pageRes, nil
}

// GetPaginatedActiveStakersByPoolCountAndDelegation returns all active stakers
// sorted by the amount of pools they are participating. If the poolCount is equal
// they are sorted by current total delegation. It supports the cosmos query.PageRequest pagination.
func (k Keeper) GetPaginatedActiveStakersByPoolCountAndDelegation(ctx sdk.Context, pagination *query.PageRequest) ([]string, *query.PageResponse, error) {
	activeStakers := k.stakersKeeper.GetActiveStakers(ctx)
	sort.Slice(activeStakers, func(i, j int) bool {
		pc_i := k.stakersKeeper.GetPoolCount(ctx, activeStakers[i])
		pc_j := k.stakersKeeper.GetPoolCount(ctx, activeStakers[j])

		if pc_i == pc_j {
			return k.GetDelegationAmount(ctx, activeStakers[i]) > k.GetDelegationAmount(ctx, activeStakers[j])
		}
		return pc_i > pc_j
	})

	return arrayPagination(activeStakers, pagination)
}
