package keeper

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DoesFunderExist checks if the funding exists
func (k Keeper) DoesFunderExist(ctx sdk.Context, funderAddress string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FunderKeyPrefix)
	return store.Has(types.FunderKey(funderAddress))
}

// GetFunder returns the funder
func (k Keeper) GetFunder(ctx sdk.Context, funderAddress string) (funder types.Funder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FunderKeyPrefix)

	b := store.Get(types.FunderKey(
		funderAddress,
	))
	if b == nil {
		return funder, false
	}

	k.cdc.MustUnmarshal(b, &funder)
	return funder, true
}

// SetFunder sets a specific funder in the store from its index
func (k Keeper) setFunder(ctx sdk.Context, funder *types.Funder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FunderKeyPrefix)
	b := k.cdc.MustMarshal(funder)
	store.Set(types.FunderKey(
		funder.Address,
	), b)
}
