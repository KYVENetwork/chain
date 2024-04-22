package v1_5

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v1.5.0"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := sdk.UnwrapSDKContext(ctx).Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))
		// TODO: migrate gov params
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
