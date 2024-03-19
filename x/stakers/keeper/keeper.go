package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"

	storetypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		authority string

		accountKeeper    types.AccountKeeper
		bankKeeper       types.BankKeeper
		distrkeeper      types.DistrKeeper
		poolKeeper       types.PoolKeeper
		upgradeKeeper    types.UpgradeKeeper
		delegationKeeper delegationKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,

	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrkeeper types.DistrKeeper,
	poolKeeper types.PoolKeeper,
	upgradeKeeper types.UpgradeKeeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,

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

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}
