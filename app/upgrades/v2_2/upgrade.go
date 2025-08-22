package v2_2

import (
	"context"
	"fmt"

	liquidkeeper "github.com/KYVENetwork/chain/x/liquid/keeper"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v2.2.0"
)

var logger log.Logger

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	liquidStakingKeeper *liquidkeeper.Keeper,
	stakingKeeper *stakingKeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger = sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// Run cosmos migrations
		migratedVersionMap, err := mm.RunMigrations(ctx, configurator, fromVM)

		// Run KYVE migrations
		validators, _ := stakingKeeper.GetValidators(ctx, 1_000)
		for _, v := range validators {
			valAddress, _ := sdk.ValAddressFromBech32(v.OperatorAddress)
			err := liquidStakingKeeper.Hooks().AfterValidatorCreated(ctx, valAddress)
			if err != nil {
				return nil, err
			}
		}

		logger.Info(fmt.Sprintf("finished upgrade %v", UpgradeName))

		return migratedVersionMap, err
	}
}
