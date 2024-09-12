package v1_6

import (
	"context"
	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	_ "embed"
	"fmt"
	bundleskeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v1.6.0"
)

var logger log.Logger

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	upgradeHeight int64,
	bundlesKeeper bundleskeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger = sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// Run cosmos migrations
		migratedVersionMap, err := mm.RunMigrations(ctx, configurator, fromVM)

		// Run KYVE migrations
		bundlesKeeper.SetBundlesMigrationUpgradeHeight(sdkCtx, uint64(upgradeHeight))

		return migratedVersionMap, err
	}
}
