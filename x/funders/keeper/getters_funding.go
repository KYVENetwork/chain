package keeper

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DoesFundingExist checks if the funding exists
func (k Keeper) doesFundingExist(ctx sdk.Context, funderAddress string, poolId uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FundingKeyPrefixByFunder)
	return store.Has(types.FundingKeyByFunder(funderAddress, poolId))
}

// GetFunding returns the funding
func (k Keeper) getFunding(ctx sdk.Context, funderAddress string, poolId uint64) (funding *types.Funding, found bool) {
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
