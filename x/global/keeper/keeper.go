package keeper

import (
	"fmt"

	"github.com/KYVENetwork/chain/x/global/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storeTypes.StoreKey

		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storeTypes.StoreKey,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		authority: authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
