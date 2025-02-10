package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/KYVENetwork/chain/util"
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

		accountKeeper types.AccountKeeper
		bankKeeper    util.BankKeeper
		distrkeeper   util.DistributionKeeper
		poolKeeper    types.PoolKeeper
		stakerKeeper  types.StakerKeeper
		fundersKeeper types.FundersKeeper

		Schema        collections.Schema
		BundlesParams collections.Item[types.Params]
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
	distrKeeper util.DistributionKeeper,
	poolKeeper types.PoolKeeper,
	stakerKeeper types.StakerKeeper,
	fundersKeeper types.FundersKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		memService:   memService,
		logger:       logger,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrkeeper:   distrKeeper,
		poolKeeper:    poolKeeper,
		stakerKeeper:  stakerKeeper,
		fundersKeeper: fundersKeeper,

		BundlesParams: collections.NewItem(sb, types.ParamsPrefix, "params", codec.CollValue[types.Params](cdc)),
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

// TODO: remove after v2 migration
func (k Keeper) Migration_GetStoreService() store.KVStoreService {
	return k.storeService
}

// TODO: remove after v2 migration
func (k Keeper) Migration_GetCodec() codec.BinaryCodec {
	return k.cdc
}
