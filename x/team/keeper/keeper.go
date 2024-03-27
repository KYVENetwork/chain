package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	"github.com/KYVENetwork/chain/util"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

	"github.com/cosmos/cosmos-sdk/codec"

	// Team
	"github.com/KYVENetwork/chain/x/team/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		accountKeeper util.AccountKeeper
		bankKeeper    types.BankKeeper
		mintKeeper    mintKeeper.Keeper
		upgradeKeeper util.UpgradeKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,

	accountKeeper util.AccountKeeper,
	bankKeeper types.BankKeeper,
	mintKeeper mintKeeper.Keeper,
	upgradeKeeper util.UpgradeKeeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		mintKeeper:    mintKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
