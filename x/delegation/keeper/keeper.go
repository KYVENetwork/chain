package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/KYVENetwork/chain/util"

	storetypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		memService   store.MemoryStoreService
		logger       log.Logger

		authority string

		accountKeeper types.AccountKeeper
		bankKeeper    util.BankKeeper
		distrKeeper   util.DistributionKeeper
		poolKeeper    types.PoolKeeper
		upgradeKeeper util.UpgradeKeeper
		stakersKeeper types.StakersKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	memService store.MemoryStoreService,
	logger log.Logger,

	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper util.BankKeeper,
	distrkeeper util.DistributionKeeper,
	poolKeeper types.PoolKeeper,
	upgradeKeeper util.UpgradeKeeper,
	stakersKeeper types.StakersKeeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		memService:   memService,
		logger:       logger,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrkeeper,
		poolKeeper:    poolKeeper,
		upgradeKeeper: upgradeKeeper,
		stakersKeeper: stakersKeeper,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

var memStoreInitialized = false

func (k Keeper) InitMemStore(gasCtx sdk.Context) {
	if !memStoreInitialized {

		// Update mem index
		noGasCtx := gasCtx.WithBlockGasMeter(storetypes.NewInfiniteGasMeter())
		for _, entry := range k.GetAllDelegationData(noGasCtx) {
			k.SetStakerIndex(noGasCtx, entry.Staker)
		}

		memStoreInitialized = true
	}
}
