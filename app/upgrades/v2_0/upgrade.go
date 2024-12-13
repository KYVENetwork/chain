package v2_0

import (
	"context"
	"fmt"
	delegationkeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"

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
	delegationKeeper delegationkeeper.Keeper,
	stakersKeeper stakersKeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger = sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// Run cosmos migrations
		migratedVersionMap, err := mm.RunMigrations(ctx, configurator, fromVM)

		// Run KYVE migrations
		migrateProtocolStakers(sdkCtx, delegationKeeper, stakersKeeper, stakingKeeper)

		return migratedVersionMap, err
	}
}

func migrateProtocolStakers(ctx sdk.Context, delegationKeeper delegationkeeper.Keeper,
	stakersKeeper stakersKeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper) {

	// Process current unbonding queue
	delegationKeeper.FullyProcessDelegatorUnbondingQueue(ctx)

	for _, mapping := range ValidatorMappings {
		// TODO check if staker exists
		delegators := delegationKeeper.GetDelegatorsByStaker(ctx, mapping.ProtocolAddress)

		for _, delegator := range delegators {
			amount := delegationKeeper.GetDelegationAmountOfDelegator(ctx, mapping.ProtocolAddress, delegator)
			delegationKeeper.PerformUndelegation(ctx, mapping.ProtocolAddress, delegator, amount)

			stakingServer := stakingkeeper.NewMsgServerImpl(&stakingKeeper)
			_, err := stakingServer.Delegate(ctx, stakingTypes.NewMsgDelegate(
				delegator,
				mapping.ConsensusAddress,
				sdk.NewInt64Coin(globalTypes.Denom, int64(amount)),
			))
			if err != nil {
				// TODO how to handle
				return
			}
		}

		_ = mapping
	}

}
