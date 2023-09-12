package v1_4

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibcTmMigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"

	// Consensus
	consensusKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	// Global
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	// Governance
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	// IBC Core
	ibcKeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

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
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("Run v1.4 upgrade")

		var err error

		// Migrate consensus parameters from x/params to dedicated x/consensus module.
		baseAppSubspace := paramsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramsTypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppSubspace, &consensusKeeper)

		// ibc-go v7.0 to v7.1 upgrade
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := ibcKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		ibcKeeper.ClientKeeper.SetParams(ctx, params)

		// Run module migrations.
		vm, err = mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Prune expired Tendermint consensus states.
		_, err = ibcTmMigrations.PruneExpiredConsensusStates(ctx, cdc, ibcKeeper.ClientKeeper)
		if err != nil {
			return vm, err
		}

		// Migrate initial deposit ratio.
		err = migrateInitialDepositRatio(ctx, globalKeeper, govKeeper)
		if err != nil {
			return vm, err
		}

		return vm, nil
	}
}

// migrateInitialDepositRatio migrates the MinInitialDepositRatio parameter from
// our custom x/global module to the x/gov module.
func migrateInitialDepositRatio(
	ctx sdk.Context,
	globalKeeper globalKeeper.Keeper,
	govKeeper govKeeper.Keeper,
) error {
	minInitialDepositRatio := globalKeeper.GetMinInitialDepositRatio(ctx)

	params := govKeeper.GetParams(ctx)
	params.MinInitialDepositRatio = minInitialDepositRatio.String()

	return govKeeper.SetParams(ctx, params)
}
