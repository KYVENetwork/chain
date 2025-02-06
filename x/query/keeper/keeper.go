package keeper

import (
	"fmt"

	"github.com/KYVENetwork/chain/util"

	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"

	"cosmossdk.io/log"
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	poolkeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	stakerskeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc    codec.BinaryCodec
		logger log.Logger

		accountKeeper authkeeper.AccountKeeper
		bankKeeper    bankkeeper.Keeper
		distrkeeper   util.DistributionKeeper
		poolKeeper    *poolkeeper.Keeper
		// TODO: rename to stakersKeeper
		stakerKeeper *stakerskeeper.Keeper
		// TODO: rename to bundlesKeeper
		bundleKeeper  types.BundlesKeeper
		globalKeeper  globalKeeper.Keeper
		govKeeper     *govkeeper.Keeper
		teamKeeper    teamKeeper.Keeper
		fundersKeeper fundersKeeper.Keeper
		stakingKeeper util.StakingKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	logger log.Logger,

	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	distrkeeper util.DistributionKeeper,
	poolKeeper *poolkeeper.Keeper,
	stakerKeeper *stakerskeeper.Keeper,
	bundleKeeper types.BundlesKeeper,
	globalKeeper globalKeeper.Keeper,
	govKeeper *govkeeper.Keeper,
	teamKeeper teamKeeper.Keeper,
	fundersKeeper fundersKeeper.Keeper,
	stakingKeeper util.StakingKeeper,
) Keeper {
	return Keeper{
		cdc:    cdc,
		logger: logger,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrkeeper:   distrkeeper,
		poolKeeper:    poolKeeper,
		stakerKeeper:  stakerKeeper,
		bundleKeeper:  bundleKeeper,
		globalKeeper:  globalKeeper,
		govKeeper:     govKeeper,
		teamKeeper:    teamKeeper,
		fundersKeeper: fundersKeeper,
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
