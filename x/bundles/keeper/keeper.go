package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	"github.com/KYVENetwork/chain/util"

	storetypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/bundles/types"
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

		accountKeeper    types.AccountKeeper
		bankKeeper       util.BankKeeper
		distrkeeper      types.DistrKeeper
		poolKeeper       types.PoolKeeper
		stakerKeeper     types.StakerKeeper
		delegationKeeper types.DelegationKeeper
		fundersKeeper    types.FundersKeeper
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
	distrkeeper types.DistrKeeper,
	poolKeeper types.PoolKeeper,
	stakerKeeper types.StakerKeeper,
	delegationKeeper types.DelegationKeeper,
	fundersKeeper types.FundersKeeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		memService:   memService,
		logger:       logger,

		authority: authority,

		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		distrkeeper:      distrkeeper,
		poolKeeper:       poolKeeper,
		stakerKeeper:     stakerKeeper,
		delegationKeeper: delegationKeeper,
		fundersKeeper:    fundersKeeper,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// A mem-store initialization needs to be performed in the begin-block.
// After a node restarts it will use the first begin-block which happens
// to rebuild the mem-store. After that `memStoreInitialized` indicates
// that the mem store was already built.
var memStoreInitialized = false

func (k Keeper) InitMemStore(gasCtx sdk.Context) {
	if !memStoreInitialized {

		// Update mem index
		noGasCtx := gasCtx.WithBlockGasMeter(storetypes.NewInfiniteGasMeter())
		for _, entry := range k.GetAllFinalizedBundles(noGasCtx) {
			k.SetFinalizedBundleIndexes(noGasCtx, entry)
		}

		memStoreInitialized = true
	}
}
