package keeper

import (
	"fmt"

	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"

	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	delegationkeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	poolkeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	stakerskeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	"cosmossdk.io/log"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc    codec.BinaryCodec
		logger log.Logger

		accountKeeper    authkeeper.AccountKeeper
		bankKeeper       bankkeeper.Keeper
		distrkeeper      distrkeeper.Keeper
		poolKeeper       poolkeeper.Keeper
		stakerKeeper     stakerskeeper.Keeper
		delegationKeeper delegationkeeper.Keeper
		// TODO: rename to bundlesKeeper
		bundleKeeper  types.BundlesKeeper
		globalKeeper  globalKeeper.Keeper
		govKeeper     govkeeper.Keeper
		teamKeeper    teamKeeper.Keeper
		fundersKeeper fundersKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	logger log.Logger,

	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	distrkeeper distrkeeper.Keeper,
	poolKeeper poolkeeper.Keeper,
	stakerKeeper stakerskeeper.Keeper,
	delegationKeeper delegationkeeper.Keeper,
	bundleKeeper types.BundlesKeeper,
	globalKeeper globalKeeper.Keeper,
	govKeeper govkeeper.Keeper,
	teamKeeper teamKeeper.Keeper,
	fundersKeeper fundersKeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:    cdc,
		logger: logger,

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

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetBundlesKeeper(bundlesKeeper types.BundlesKeeper) {
	k.bundleKeeper = bundlesKeeper
}
