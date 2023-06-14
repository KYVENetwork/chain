package v1_3

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/log"

	// Pool
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	poolKeeper poolKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		CheckPoolAddresses(ctx, logger, poolKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// CheckPoolAddresses ensures that each pool account exists post upgrade.
func CheckPoolAddresses(ctx sdk.Context, logger log.Logger, keeper poolKeeper.Keeper) {
	pools := keeper.GetAllPools(ctx)

	for _, pool := range pools {
		keeper.EnsurePoolAccount(ctx, pool.Id)

		name := fmt.Sprintf("%s/%d", poolTypes.ModuleName, pool.Id)
		logger.Info("successfully initialised pool account", "name", name)
	}
}
