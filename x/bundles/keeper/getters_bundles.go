package keeper

import (
	"encoding/binary"
	"fmt"

	cosmossdk_io_math "cosmossdk.io/math"

	queryTypes "github.com/KYVENetwork/chain/x/query/types"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

// SetBundleProposal stores a current bundle proposal in the KV-Store.
// There is only one bundle proposal per pool
func (k Keeper) SetBundleProposal(ctx sdk.Context, bundleProposal types.BundleProposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BundleKeyPrefix)
	b := k.cdc.MustMarshal(&bundleProposal)
	store.Set(types.BundleProposalKey(
		bundleProposal.PoolId,
	), b)
}

// GetBundleProposal returns the bundle proposal for the given pool with id `poolId`
func (k Keeper) GetBundleProposal(ctx sdk.Context, poolId uint64) (val types.BundleProposal, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BundleKeyPrefix)

	b := store.Get(types.BundleProposalKey(poolId))
	if b == nil {
		val.PoolId = poolId
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllBundleProposals returns all bundle proposals of all pools
func (k Keeper) GetAllBundleProposals(ctx sdk.Context) (list []types.BundleProposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BundleKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	for ; iterator.Valid(); iterator.Next() {
		var val types.BundleProposal
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// =====================
// = Finalized Bundles =
// =====================

// SetFinalizedBundle stores a finalized bundle identified by its `poolId` and `id`.
func (k Keeper) SetFinalizedBundle(ctx sdk.Context, finalizedBundle types.FinalizedBundle) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FinalizedBundlePrefix)
	b := k.cdc.MustMarshal(&finalizedBundle)
	store.Set(types.FinalizedBundleKey(
		finalizedBundle.PoolId,
		finalizedBundle.Id,
	), b)

	k.SetFinalizedBundleIndexes(ctx, finalizedBundle)
}

// SetFinalizedBundleIndexes sets an in-memory reference for every bundle sorted by pool/fromIndex
// to allow querying for specific bundle ranges.
func (k Keeper) SetFinalizedBundleIndexes(ctx sdk.Context, finalizedBundle types.FinalizedBundle) {
	indexByStorageIndex := prefix.NewStore(ctx.KVStore(k.memKey), types.FinalizedBundleByIndexPrefix)
	indexByStorageIndex.Set(
		types.FinalizedBundleByIndexKey(finalizedBundle.PoolId, finalizedBundle.FromIndex),
		util.GetByteKey(finalizedBundle.Id))
}

func (k Keeper) GetAllFinalizedBundles(ctx sdk.Context) (list []types.FinalizedBundle) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FinalizedBundlePrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	for ; iterator.Valid(); iterator.Next() {
		var val types.FinalizedBundle
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetFinalizedBundlesByPool(ctx sdk.Context, poolId uint64) (list []types.FinalizedBundle) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FinalizedBundlePrefix)
	iterator := sdk.KVStorePrefixIterator(store, util.GetByteKey(poolId))

	for ; iterator.Valid(); iterator.Next() {
		var val types.FinalizedBundle
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetFinalizedBundle returns a finalized bundle by its identifier
func (k Keeper) GetFinalizedBundle(ctx sdk.Context, poolId, id uint64) (val types.FinalizedBundle, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FinalizedBundlePrefix)

	b := store.Get(types.FinalizedBundleKey(poolId, id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func RawBundleToQueryBundle(rawFinalizedBundle types.FinalizedBundle, versionMap map[int32]uint64) (queryBundle queryTypes.FinalizedBundle) {
	finalizedHeight := cosmossdk_io_math.NewInt(int64(rawFinalizedBundle.FinalizedAt.Height))
	finalizedTimestamp := cosmossdk_io_math.NewInt(int64(rawFinalizedBundle.FinalizedAt.Timestamp))

	finalizedBundle := queryTypes.FinalizedBundle{
		PoolId:        rawFinalizedBundle.PoolId,
		Id:            rawFinalizedBundle.Id,
		StorageId:     rawFinalizedBundle.StorageId,
		Uploader:      rawFinalizedBundle.Uploader,
		FromIndex:     rawFinalizedBundle.FromIndex,
		ToIndex:       rawFinalizedBundle.ToIndex,
		ToKey:         rawFinalizedBundle.ToKey,
		BundleSummary: rawFinalizedBundle.BundleSummary,
		DataHash:      rawFinalizedBundle.DataHash,
		FinalizedAt: &queryTypes.FinalizedAt{
			Height:    &finalizedHeight,
			Timestamp: &finalizedTimestamp,
		},
		FromKey:           rawFinalizedBundle.FromKey,
		StorageProviderId: uint64(rawFinalizedBundle.StorageProviderId),
		CompressionId:     uint64(rawFinalizedBundle.CompressionId),
		StakeSecurity: &queryTypes.StakeSecurity{
			ValidVotePower: nil,
			TotalVotePower: nil,
		},
	}

	// Check for version 2
	if rawFinalizedBundle.FinalizedAt.Height >= versionMap[2] {
		validPower := cosmossdk_io_math.NewInt(int64(rawFinalizedBundle.StakeSecurity.ValidVotePower))
		totalPower := cosmossdk_io_math.NewInt(int64(rawFinalizedBundle.StakeSecurity.TotalVotePower))
		finalizedBundle.StakeSecurity.ValidVotePower = &validPower
		finalizedBundle.StakeSecurity.TotalVotePower = &totalPower
	}

	return finalizedBundle
}

// GetPaginatedFinalizedBundleQuery parses a paginated request and builds a valid response out of the
// raw finalized bundles. It uses the fact that the ID of a bundle increases incrementally (starting with 0)
// and allows therefore for efficient queries using `offset`.
func (k Keeper) GetPaginatedFinalizedBundleQuery(ctx sdk.Context, pagination *query.PageRequest, poolId uint64) ([]queryTypes.FinalizedBundle, *query.PageResponse, error) {
	// Parse basic pagination
	if pagination == nil {
		pagination = &query.PageRequest{CountTotal: true}
	}

	offset := pagination.Offset
	key := pagination.Key
	limit := pagination.Limit
	reverse := pagination.Reverse

	if limit == 0 {
		limit = query.DefaultLimit
	}

	pageResponse := query.PageResponse{}

	// user has to use either offset or key, not both
	if offset > 0 && key != nil {
		return nil, nil, fmt.Errorf("invalid request, either offset or key is expected, got both")
	}

	// Init Bundles Store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), util.GetByteKey(types.FinalizedBundlePrefix, poolId))

	// Get latest bundle id by obtaining last item from the iterator
	reverseIterator := store.ReverseIterator(nil, nil)
	if reverseIterator.Valid() {
		// Current bundle_id equals the total amount of bundles - 1
		bundleId := binary.BigEndian.Uint64(reverseIterator.Key())
		pageResponse.Total = bundleId + 1
	}
	_ = reverseIterator.Close()

	// Translate offset to next page keys
	if len(key) == 0 {
		if reverse {
			pagination.Key = util.GetByteKey(pageResponse.Total - offset)
		} else {
			pagination.Key = util.GetByteKey(offset)
		}
	}

	var iterator storeTypes.Iterator
	// Use correct iterator depending on the request
	if reverse {
		iterator = store.ReverseIterator(nil, pagination.Key)
	} else {
		iterator = store.Iterator(pagination.Key, nil)
	}

	var data []queryTypes.FinalizedBundle
	versionMap := k.GetBundleVersionMap(ctx).GetMap()

	// Iterate bundle store and build actual response
	for i := uint64(0); i < limit; i++ {
		if iterator.Valid() {
			var rawFinalizedBundle types.FinalizedBundle
			if err := k.cdc.Unmarshal(iterator.Value(), &rawFinalizedBundle); err != nil {
				return nil, nil, err
			}
			data = append(data, RawBundleToQueryBundle(rawFinalizedBundle, versionMap))
			pageResponse.NextKey = iterator.Key()
			iterator.Next()
		} else {
			break
		}
	}
	// Fetch next key (if there is one)
	if iterator.Valid() {
		if !reverse {
			pageResponse.NextKey = iterator.Key()
		}
	} else {
		pageResponse.NextKey = nil
	}
	_ = iterator.Close()

	return data, &pageResponse, nil
}

func (k Keeper) GetFinalizedBundleByIndex(ctx sdk.Context, poolId, index uint64) (val queryTypes.FinalizedBundle, found bool) {
	proposalIndexStore := prefix.NewStore(ctx.KVStore(k.memKey), util.GetByteKey(types.FinalizedBundleByIndexPrefix, poolId))
	proposalIndexIterator := proposalIndexStore.ReverseIterator(nil, util.GetByteKey(index+1))
	defer proposalIndexIterator.Close()

	if proposalIndexIterator.Valid() {
		bundleId := binary.BigEndian.Uint64(proposalIndexIterator.Value())

		bundle, bundleFound := k.GetFinalizedBundle(ctx, poolId, bundleId)
		if bundleFound {
			if bundle.FromIndex <= index && bundle.ToIndex > index {
				versionMap := k.GetBundleVersionMap(ctx).GetMap()
				return RawBundleToQueryBundle(bundle, versionMap), true
			}
		}
	}
	return
}

// Finalized Bundle Version Map

// SetBundleVersionMap stores the bundle version map
func (k Keeper) SetBundleVersionMap(ctx sdk.Context, bundleVersionMap types.BundleVersionMap) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&bundleVersionMap)
	store.Set(types.FinalizedBundleVersionMapKey, b)
}

// GetBundleVersionMap returns the bundle version map
func (k Keeper) GetBundleVersionMap(ctx sdk.Context) (val types.BundleVersionMap) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.FinalizedBundleVersionMapKey)
	if b == nil {
		val.Versions = make([]*types.BundleVersionEntry, 0)
		return val
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}
