package funders

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllFundings returns all fundings
func GetAllFundings(ctx sdk.Context, cdc codec.Codec, storeKey storeTypes.StoreKey) (fundings []Funding) {
	store := prefix.NewStore(ctx.KVStore(storeKey), fundersTypes.FundingKeyPrefixByFunder)
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
