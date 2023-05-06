package v1_2

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Consensus
	consensusKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	// Global
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	// Governance
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	// IBC Core
	ibcKeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	// IBC Light Clients
	ibcTmMigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	// Params
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.BinaryCodec,
	consensusKeeper consensusKeeper.Keeper,
	globalKeeper globalKeeper.Keeper,
	govKeeper govKeeper.Keeper,
	ibcKeeper ibcKeeper.Keeper,
	paramsKeeper paramsKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		var err error

		// Migrate consensus parameters from x/params to dedicated x/consensus module.
		baseAppSubspace := paramsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramsTypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppSubspace, &consensusKeeper)

		// Prune expired Tendermint consensus states.
		_, err = ibcTmMigrations.PruneExpiredConsensusStates(ctx, cdc, ibcKeeper.ClientKeeper)
		if err != nil {
			return vm, err
		}

		// Run module migrations.
		vm, err = mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Migrate initial deposit ratio.
		err = MigrateInitialDepositRatio(ctx, globalKeeper, govKeeper)
		if err != nil {
			return vm, err
		}

		return vm, nil
	}
}

// MigrateInitialDepositRatio migrates the MinInitialDepositRatio parameter from
// our custom x/global module to the x/gov module.
func MigrateInitialDepositRatio(
	ctx sdk.Context,
	globalKeeper globalKeeper.Keeper,
	govKeeper govKeeper.Keeper,
) error {
	minInitialDepositRatio := globalKeeper.GetMinInitialDepositRatio(ctx)

	params := govKeeper.GetParams(ctx)
	params.MinInitialDepositRatio = minInitialDepositRatio.String()

	return govKeeper.SetParams(ctx, params)
}
