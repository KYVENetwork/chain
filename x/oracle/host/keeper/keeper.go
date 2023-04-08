package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	// Bundles
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	// Oracle Host
	"github.com/KYVENetwork/chain/x/oracle/host/types"
	// Pool
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storeTypes.StoreKey

		bundlesKeeper bundlesKeeper.Keeper
		poolKeeper    poolKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storeTypes.StoreKey,
	bundlesKeeper bundlesKeeper.Keeper,
	poolKeeper poolKeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		bundlesKeeper: bundlesKeeper,
		poolKeeper:    poolKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
