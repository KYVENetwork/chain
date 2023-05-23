package keeper

import (
	"fmt"

	"github.com/KYVENetwork/chain/util"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		authority string

		accountKeeper util.AccountKeeper
		bankKeeper    util.BankKeeper
		distrKeeper   util.DistributionKeeper
		upgradeKeeper util.UpgradeKeeper
		stakersKeeper types.StakersKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,

	authority string,

	accountKeeper util.AccountKeeper,
	bankKeeper util.BankKeeper,
	distrkeeper util.DistributionKeeper,
	upgradeKeeper util.UpgradeKeeper,
	stakersKeeper types.StakersKeeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrkeeper,
		upgradeKeeper: upgradeKeeper,
		stakersKeeper: stakersKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}

var memStoreInitialized = false

func (k Keeper) InitMemStore(gasCtx sdk.Context) {
	if !memStoreInitialized {

		// Update mem index
		noGasCtx := gasCtx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())
		for _, entry := range k.GetAllDelegationData(noGasCtx) {
			k.SetStakerIndex(noGasCtx, entry.Staker)
		}

		memStoreInitialized = true
	}
}
