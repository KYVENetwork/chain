package v2_0

import (
	"context"
	"fmt"

	"github.com/KYVENetwork/chain/x/stakers/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

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
	UpgradeName = "v2.0.0"
)

var logger log.Logger

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	delegationKeeper delegationkeeper.Keeper,
	stakersKeeper *stakersKeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger = sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// Run cosmos migrations
		migratedVersionMap, err := mm.RunMigrations(ctx, configurator, fromVM)

		// Run KYVE migrations
		migrateProtocolStakers(sdkCtx, delegationKeeper, stakersKeeper, stakingKeeper, bankKeeper)

		// TODO: migrate slash params, commission change queues

		return migratedVersionMap, err
	}
}

func migrateProtocolStakers(ctx sdk.Context, delegationKeeper delegationkeeper.Keeper,
	stakersKeeper *stakersKeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) {
	// Process current unbonding queue
	delegationKeeper.FullyProcessDelegatorUnbondingQueue(ctx)

	totalMigratedStake := uint64(0)
	for _, mapping := range ValidatorMappings {
		// TODO check if staker exists
		delegators := delegationKeeper.GetDelegatorsByStaker(ctx, mapping.ProtocolAddress)

		for _, delegator := range delegators {
			amount := delegationKeeper.PerformFullUndelegation(ctx, mapping.ProtocolAddress, delegator)
			totalMigratedStake += amount

			stakingServer := stakingkeeper.NewMsgServerImpl(stakingKeeper)
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
	}
	logger.Info("migrated %d ukyve from protocol to chain", totalMigratedStake)

	// Undelegate Remaining
	totalReturnedStake := uint64(0)
	for _, delegator := range delegationKeeper.GetAllDelegators(ctx) {
		amount := delegationKeeper.PerformFullUndelegation(ctx, delegator.Staker, delegator.Delegator)
		totalReturnedStake += amount
	}
	logger.Info("returned %d ukyve from protocol to users", totalReturnedStake)

	// Withdraw Pending Commissions
	for _, staker := range stakersKeeper.GetAllLegacyStakers(ctx) {
		if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(staker.Address), staker.CommissionRewards); err != nil {
			panic(err)
		}
		_ = ctx.EventManager().EmitTypedEvent(&types.EventClaimCommissionRewards{
			Staker:  staker.Address,
			Amounts: staker.CommissionRewards.String(),
		})
	}

	// TODO delete staker object
}
