package keeper

import (
	"fmt"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/KYVENetwork/chain/x/pool/types"
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

		stakersKeeper types.StakersKeeper
		accountKeeper authkeeper.AccountKeeper
		bankKeeper    bankkeeper.Keeper
		distrkeeper   distrkeeper.Keeper
		upgradeKeeper types.UpgradeKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,

	authority string,

	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	upgradeKeeper types.UpgradeKeeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrkeeper:   distrKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func SetStakersKeeper(k *Keeper, stakersKeeper types.StakersKeeper) {
	k.stakersKeeper = stakersKeeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}
