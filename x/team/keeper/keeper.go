package keeper

import (
	"fmt"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Team
	"github.com/KYVENetwork/chain/x/team/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storeTypes.StoreKey

		accountKeeper authKeeper.AccountKeeper
		bankKeeper    bankKeeper.Keeper
		mintKeeper    mintKeeper.Keeper
		upgradeKeeper upgradeKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storeTypes.StoreKey,
	accountKeeper authKeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	mintKeeper mintKeeper.Keeper,
	upgradeKeeper upgradeKeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		mintKeeper:    mintKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storeTypes.StoreKey {
	return k.storeKey
}
