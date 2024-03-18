package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		authority string

		accountKeeper types.AccountKeeper
		bankKeeper    util.BankKeeper
		poolKeeper    types.PoolKeeper
		upgradeKeeper util.UpgradeKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,

	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper util.BankKeeper,
	poolKeeper types.PoolKeeper,
	upgradeKeeper util.UpgradeKeeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		poolKeeper:    poolKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}
