package keeper

import (
	"fmt"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"

	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	bundlekeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	delegationkeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	poolkeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	stakerskeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper    authkeeper.AccountKeeper
		bankKeeper       bankkeeper.Keeper
		distrkeeper      distrkeeper.Keeper
		poolKeeper       poolkeeper.Keeper
		stakerKeeper     stakerskeeper.Keeper
		delegationKeeper delegationkeeper.Keeper
		bundleKeeper     bundlekeeper.Keeper
		globalKeeper     globalKeeper.Keeper
		govKeeper        govkeeper.Keeper
		teamKeeper       teamKeeper.Keeper
		fundersKeeper    fundersKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	distrkeeper distrkeeper.Keeper,
	poolKeeper poolkeeper.Keeper,
	stakerKeeper stakerskeeper.Keeper,
	delegationKeeper delegationkeeper.Keeper,
	bundleKeeper bundlekeeper.Keeper,
	globalKeeper globalKeeper.Keeper,
	govKeeper govkeeper.Keeper,
	teamKeeper teamKeeper.Keeper,
	fundersKeeper fundersKeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		distrkeeper:      distrkeeper,
		poolKeeper:       poolKeeper,
		stakerKeeper:     stakerKeeper,
		delegationKeeper: delegationKeeper,
		bundleKeeper:     bundleKeeper,
		globalKeeper:     globalKeeper,
		govKeeper:        govKeeper,
		teamKeeper:       teamKeeper,
		fundersKeeper:    fundersKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
