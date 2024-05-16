package funders

import (
	storeTypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllFundings returns all fundings
func GetAllFundings(ctx sdk.Context, cdc codec.Codec, storeKey storeTypes.StoreKey) (fundings []Funding) {
	store := ctx.KVStore(storeKey)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	//goland:noinspection GoUnhandledErrorResult
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val Funding
		cdc.MustUnmarshal(iterator.Value(), &val)
		fundings = append(fundings, val)
	}

	return fundings
}
