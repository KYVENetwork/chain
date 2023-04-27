package v1_2

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Consensus
	consensusKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	// Params
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	consensusKeeper consensusKeeper.Keeper,
	paramsKeeper paramsKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Migrate consensus parameters from x/params to dedicated x/consensus module.
		baseAppSubspace := paramsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramsTypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppSubspace, &consensusKeeper)

		// TODO: Migrate MinInitialDepositRatio from x/global to x/gov.

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
