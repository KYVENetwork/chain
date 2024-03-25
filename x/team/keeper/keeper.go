package keeper

import (
	"cosmossdk.io/log"
	"fmt"

	upgradeKeeper "cosmossdk.io/x/upgrade/keeper"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

	storeTypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
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
		logger   log.Logger

		accountKeeper authKeeper.AccountKeeper
		bankKeeper    bankKeeper.Keeper
		mintKeeper    mintKeeper.Keeper
		upgradeKeeper upgradeKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storeTypes.StoreKey,
	logger log.Logger,

	accountKeeper authKeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	mintKeeper mintKeeper.Keeper,
	upgradeKeeper upgradeKeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		logger:   logger,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		mintKeeper:    mintKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storeTypes.StoreKey {
	return k.storeKey
}
