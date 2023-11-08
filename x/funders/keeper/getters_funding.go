package keeper

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	queryTypes "github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoesFundingExist checks if the funding exists
func (k Keeper) DoesFundingExist(ctx sdk.Context, funderAddress string, poolId uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)
	return store.Has(types.FundingKeyByFunder(funderAddress, poolId))
}

// GetFunding returns the funding
func (k Keeper) GetFunding(ctx sdk.Context, funderAddress string, poolId uint64) (funding types.Funding, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)

	b := store.Get(types.FundingKeyByFunder(
		funderAddress,
		poolId,
	))
	if b == nil {
		return funding, false
	}

	k.cdc.MustUnmarshal(b, &funding)
	return funding, true
}

// GetFundingsOfFunder returns all fundings of a funder
func (k Keeper) GetFundingsOfFunder(ctx sdk.Context, funderAddress string) (fundings []types.Funding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)

	iterator := sdk.KVStorePrefixIterator(store, types.FundingKeyByFunderIter(funderAddress))
	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var funding types.Funding
		k.cdc.MustUnmarshal(iterator.Value(), &funding)
		fundings = append(fundings, funding)
	}
	return fundings
}

// GetFundingsOfPool returns all fundings of a pool
func (k Keeper) GetFundingsOfPool(ctx sdk.Context, poolId uint64) (fundings []types.Funding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByPool)

	iterator := sdk.KVStorePrefixIterator(store, types.FundingKeyByPoolIter(poolId))
	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var funding types.Funding
		k.cdc.MustUnmarshal(iterator.Value(), &funding)
		fundings = append(fundings, funding)
	}
	return fundings
}

// GetAllFundings returns all fundings
func (k Keeper) GetAllFundings(ctx sdk.Context) (fundings []types.Funding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Funding
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		fundings = append(fundings, val)
	}

	return fundings
}

// SetFunding sets a specific funding in the store from its index
func (k Keeper) SetFunding(ctx sdk.Context, funding *types.Funding) {
	b := k.cdc.MustMarshal(funding)

	storeByFunder := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)
	storeByFunder.Set(types.FundingKeyByFunder(
		funding.FunderAddress,
		funding.PoolId,
	), b)

	storeByPool := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByPool)
	storeByPool.Set(types.FundingKeyByPool(
		funding.FunderAddress,
		funding.PoolId,
	), b)
}

// GetPaginatedFundingQuery performs a full search on all fundings with the given parameters.
// Requires either funderAddress or poolId to be provided.
func (k Keeper) GetPaginatedFundingQuery(
	ctx sdk.Context,
	pagination *query.PageRequest,
	funderAddress *string,
	poolId *uint64,
	fundingStatus queryTypes.FundingStatus,
) ([]types.Funding, *query.PageResponse, error) {
	if funderAddress == nil && poolId == nil {
		return nil, nil, status.Error(codes.InvalidArgument, "either funderAddress or poolId must be provided")
	}
	keyPrefix := types.FundingKeyPrefixByFunder
	if funderAddress == nil {
		keyPrefix = types.FundingKeyPrefixByPool
	}

	var fundings []types.Funding
	store := prefix.NewStore(ctx.KVStore(k.storeKey), keyPrefix)

	pageRes, err := query.FilteredPaginate(store, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var funding types.Funding
		if err := k.cdc.Unmarshal(value, &funding); err != nil {
			return false, err
		}

		if funderAddress != nil && *funderAddress != funding.FunderAddress {
			return false, nil
		}

		if poolId != nil && *poolId != funding.PoolId {
			return false, nil
		}

		if fundingStatus == queryTypes.FUNDING_STATUS_ACTIVE && funding.IsInactive() {
			return false, nil
		}

		if fundingStatus == queryTypes.FUNDING_STATUS_INACTIVE && funding.IsActive() {
			return false, nil
		}

		if accumulate {
			fundings = append(fundings, funding)
		}

		return true, nil
	})
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	return fundings, pageRes, nil
}
