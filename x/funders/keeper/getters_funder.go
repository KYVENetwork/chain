package keeper

import (
	storeTypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"strings"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoesFunderExist checks if the funding exists
func (k Keeper) DoesFunderExist(ctx sdk.Context, funderAddress string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.FunderKeyPrefix)
	return store.Has(types.FunderKey(funderAddress))
}

// GetFunder returns the funder
func (k Keeper) GetFunder(ctx sdk.Context, funderAddress string) (funder types.Funder, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.FunderKeyPrefix)

	b := store.Get(types.FunderKey(
		funderAddress,
	))
	if b == nil {
		return funder, false
	}

	k.cdc.MustUnmarshal(b, &funder)
	return funder, true
}

// GetAllFunders returns all funders
func (k Keeper) GetAllFunders(ctx sdk.Context) (funders []types.Funder) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.FunderKeyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Funder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		funders = append(funders, val)
	}

	return funders
}

// SetFunder sets a specific funder in the store from its index
func (k Keeper) SetFunder(ctx sdk.Context, funder *types.Funder) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.FunderKeyPrefix)
	b := k.cdc.MustMarshal(funder)
	store.Set(types.FunderKey(
		funder.Address,
	), b)
}

// GetPaginatedFundersQuery performs a full search on all funders with the given parameters.
func (k Keeper) GetPaginatedFundersQuery(
	ctx sdk.Context,
	pagination *query.PageRequest,
	search string,
) ([]types.Funder, *query.PageResponse, error) {
	var funders []types.Funder

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.FunderKeyPrefix)

	pageRes, err := query.FilteredPaginate(store, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var funder types.Funder
		if err := k.cdc.Unmarshal(value, &funder); err != nil {
			return false, err
		}

		// filter search
		if !strings.Contains(strings.ToLower(funder.Moniker), strings.ToLower(search)) {
			return false, nil
		}

		if accumulate {
			funders = append(funders, funder)
		}

		return true, nil
	})
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	return funders, pageRes, nil
}
