package keeper

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetProposal set a specific proposal in the store from its index
func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))
	b := k.cdc.MustMarshal(&proposal)
	store.Set(types.ProposalKey(
		proposal.StorageId,
	), b)

	// Insert bundle id for second index
	storeIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	storeIndex.Set(types.ProposalKeyIndex2(proposal.PoolId, proposal.Id), []byte(proposal.StorageId))

	// Insert bundle id for second index
	storeIndex3 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex3)
	storeIndex3.Set(types.ProposalKeyIndex3(proposal.PoolId, proposal.FinalizedAt), []byte(proposal.StorageId))
}

// GetProposal returns a proposal from its index
func (k Keeper) GetProposal(ctx sdk.Context, storageId string) (val types.Proposal, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))

	b := store.Get(types.ProposalKey(
		storageId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetProposalByPoolIdAndBundleId returns a proposal from its index
func (k Keeper) GetProposalByPoolIdAndBundleId(ctx sdk.Context, poolId uint64, bundleId uint64) (val types.Proposal, found bool) {
	// Insert bundle id for second index
	storeIndex2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	storageIdBytes := storeIndex2.Get(types.ProposalKeyIndex2(poolId, bundleId))

	if storageIdBytes == nil {
		return val, false
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))

	b := store.Get(types.ProposalKey(string(storageIdBytes)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetProposalsByPoolIdSinceBundleId returns for a given pool all proposals that have
// an ID equal or higher to minBundleId
func (k Keeper) GetProposalsByPoolIdSinceBundleId(ctx sdk.Context, poolId uint64, minBundleId uint64) (proposals []types.Proposal) {
	proposalPrefixBuilder := types.KeyPrefixBuilder{Key: types.ProposalKeyPrefixIndex2}.AInt(poolId)
	proposalIndexStore := prefix.NewStore(ctx.KVStore(k.storeKey), proposalPrefixBuilder.Key)
	proposalIndexIterator := proposalIndexStore.Iterator(types.KeyPrefixBuilder{}.AInt(minBundleId).Key, nil)

	defer proposalIndexIterator.Close()

	for ; proposalIndexIterator.Valid() ; proposalIndexIterator.Next() {
		storageId := string(proposalIndexIterator.Value())
		proposal, _ := k.GetProposal(ctx, storageId)
		proposals = append(proposals, proposal)
	}

	return
}

// RemoveProposal removes a proposal from the store
func (k Keeper) RemoveProposal(ctx sdk.Context, proposal types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))
	store.Delete(types.ProposalKey(proposal.StorageId))

	indexStore2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	indexStore2.Delete(types.ProposalKeyIndex2(proposal.PoolId, proposal.Id))

	// Insert bundle id for second index
	storeIndex3 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex3)
	storeIndex3.Delete(types.ProposalKeyIndex3(proposal.PoolId, proposal.FinalizedAt))
}

// GetAllProposal returns all proposal
func (k Keeper) GetAllProposal(ctx sdk.Context) (list []types.Proposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Proposal
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// TODO delete this function after v0.6.0 migration
func (k Keeper) UpgradeHelperV060MigrateSecondIndex(ctx sdk.Context) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var keysToDelete [][]byte

	for ; iterator.Valid(); iterator.Next() {
		keyArray := iterator.Key()
		var key = make([]byte, len(keyArray))
		copy(key, keyArray)
		keysToDelete = append(keysToDelete, key)
	}

	println("Delete ", len(keysToDelete), " index keys")

	for _, key := range keysToDelete {
		store.Delete(key)
	}

	storeIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	var counter uint64 = 0
	for _, proposal := range k.GetAllProposal(ctx) {
		counter++
		// Insert bundle id for second index
		storeIndex.Set(types.ProposalKeyIndex2(proposal.PoolId, proposal.Id), []byte(proposal.StorageId))
	}
	println("Created ", counter, " index keys")

}
