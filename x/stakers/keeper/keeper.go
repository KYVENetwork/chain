package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/core/store"

	"github.com/KYVENetwork/chain/util"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"

	"cosmossdk.io/log"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		memService   store.MemoryStoreService
		logger       log.Logger

		authority string

		accountKeeper    util.AccountKeeper
		bankKeeper       util.BankKeeper
		distrkeeper      util.DistributionKeeper
		poolKeeper       types.PoolKeeper
		upgradeKeeper    util.UpgradeKeeper
		delegationKeeper delegationKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	memService store.MemoryStoreService,
	logger log.Logger,

	authority string,

	accountKeeper util.AccountKeeper,
	bankKeeper util.BankKeeper,
	distrkeeper util.DistributionKeeper,
	poolKeeper types.PoolKeeper,
	upgradeKeeper util.UpgradeKeeper,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		memService:   memService,
		logger:       logger,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrkeeper:   distrkeeper,
		poolKeeper:    poolKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func SetDelegationKeeper(k *Keeper, delegationKeeper delegationKeeper.Keeper) {
	k.delegationKeeper = delegationKeeper
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// TODO: remove after migration
func (k Keeper) Migration_SetStaker(ctx sdk.Context, staker types.Staker) {
	k.setStaker(ctx, staker)
}
