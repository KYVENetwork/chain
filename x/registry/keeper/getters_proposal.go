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
		proposal.BundleId,
	), b)

	// Insert bundle id for second index
	storeIndex := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalKeyPrefixIndex2)
	storeIndex.Set(types.ProposalKeyIndex2(proposal.PoolId, proposal.FromHeight), []byte(proposal.BundleId))
}

// GetProposal returns a proposal from its index
func (k Keeper) GetProposal(
	ctx sdk.Context,
	bundleId string,
) (val types.Proposal, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProposalKeyPrefix))

	b := store.Get(types.ProposalKey(
		bundleId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
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
