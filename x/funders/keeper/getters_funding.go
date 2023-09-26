package keeper

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DoesFundingExist checks if the funding exists
func (k Keeper) DoesFundingExist(ctx sdk.Context, funderAddress string, poolId uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)
	return store.Has(types.FundingKeyByFunder(funderAddress, poolId))
}

// GetFunding returns the funding
func (k Keeper) GetFunding(ctx sdk.Context, funderAddress string, poolId uint64) (funding *types.Funding, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)

	b := store.Get(types.FundingKeyByFunder(
		funderAddress,
		poolId,
	))
	if b == nil {
		return funding, false
	}

	k.cdc.MustUnmarshal(b, funding)
	return funding, true
}

// GetFundingsOfPool returns all fundings of a pool
func (k Keeper) GetFundingsOfPool(ctx sdk.Context, poolId uint64) (fundings []*types.Funding) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByPool)

	iterator := sdk.KVStorePrefixIterator(store, types.FundingKeyByPoolOnly(poolId))
	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var funding types.Funding
		k.cdc.MustUnmarshal(iterator.Value(), &funding)
		fundings = append(fundings, &funding)
	}
	return fundings
}

// SetFunding sets a specific funding in the store from its index
func (k Keeper) setFunding(ctx sdk.Context, funding *types.Funding) {
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
