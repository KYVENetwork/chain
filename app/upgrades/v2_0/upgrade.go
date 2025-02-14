package v2_0

import (
	"context"
	"fmt"
	poolkeeper "github.com/KYVENetwork/chain/x/pool/keeper"

	multicoinrewardskeeper "github.com/KYVENetwork/chain/x/multi_coin_rewards/keeper"
	multicoinrewardstypes "github.com/KYVENetwork/chain/x/multi_coin_rewards/types"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	bundleskeeper "github.com/KYVENetwork/chain/x/bundles/keeper"

	"cosmossdk.io/math"

	globalkeeper "github.com/KYVENetwork/chain/x/global/keeper"

	"github.com/KYVENetwork/chain/x/stakers/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	delegationkeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	stakerskeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
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
	accountKeeper authkeeper.AccountKeeper,
	delegationKeeper delegationkeeper.Keeper,
	stakersKeeper *stakerskeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	bundlesKeeper bundleskeeper.Keeper,
	globalKeeper globalkeeper.Keeper,
	multiCoinRewardsKeeper multicoinrewardskeeper.Keeper,
	poolKeeper *poolkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger = sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// Run cosmos migrations
		migratedVersionMap, err := mm.RunMigrations(ctx, configurator, fromVM)

		// Run KYVE migrations
		migrateProtocolStakers(sdkCtx, delegationKeeper, stakersKeeper, stakingKeeper, bankKeeper)
		EnsureMultiCoinDistributionAccount(sdkCtx, accountKeeper, multicoinrewardstypes.ModuleName)
		EnsureMultiCoinDistributionAccount(sdkCtx, accountKeeper, multicoinrewardstypes.MultiCoinRewardsRedistributionAccountName)
		AdjustGasConfig(sdkCtx, globalKeeper)

		SetMultiCoinRewardsParams(sdkCtx, multiCoinRewardsKeeper)
		SetPoolParams(sdkCtx, poolKeeper)

		// TODO set withdraw address

		// Run Bundles Merkle Roots migrations
		bundlesKeeper.SetBundlesMigrationUpgradeHeight(sdkCtx, uint64(sdkCtx.BlockHeight()))

		logger.Info(fmt.Sprintf("finished upgrade %v", UpgradeName))

		return migratedVersionMap, err
	}
}

func SetMultiCoinRewardsParams(ctx sdk.Context, multiCoinRewardsKeeper multicoinrewardskeeper.Keeper) {
	params := multiCoinRewardsKeeper.GetParams(ctx)
	params.MultiCoinDistributionPendingTime = 60 * 60 * 24 * 14
	// KYVE Public Good Funding address
	params.MultiCoinDistributionPolicyAdminAddress = "kyve1t0uez3nn28ljnzlwndzxffyjuhean3edhtjee8"
	multiCoinRewardsKeeper.SetParams(ctx, params)
}

func SetPoolParams(ctx sdk.Context, poolKeeper *poolkeeper.Keeper) {
	params := poolKeeper.GetParams(ctx)

	// TODO: set new mainnet inflation split

	poolKeeper.SetParams(ctx, params)
}

func AdjustGasConfig(ctx sdk.Context, globalKeeper globalkeeper.Keeper) {
	params := globalKeeper.GetParams(ctx)
	params.MinGasPrice = math.LegacyMustNewDecFromStr("2")
	params.GasAdjustments = []globalTypes.GasAdjustment{
		{
			Type:   "/cosmos.staking.v1beta1.MsgCreateValidator",
			Amount: 50000000,
		},
		{
			Type:   "/kyve.funders.v1beta1.MsgCreateFunder",
			Amount: 50000000,
		},
	}
	params.GasRefunds = []globalTypes.GasRefund{
		{
			Type:     "kyve.bundles.v1beta1.MsgSubmitBundleProposal",
			Fraction: math.LegacyMustNewDecFromStr("0.95"),
		},
		{
			Type:     "kyve.bundles.v1beta1.MsgVoteBundleProposal",
			Fraction: math.LegacyMustNewDecFromStr("0.95"),
		},
		{
			Type:     "kyve.bundles.v1beta1.MsgSkipUploaderRole",
			Fraction: math.LegacyMustNewDecFromStr("0.95"),
		},
		{
			Type:     "kyve.bundles.v1beta1.MsgClaimUploaderRole",
			Fraction: math.LegacyMustNewDecFromStr("0.95"),
		},
	}
	globalKeeper.SetParams(ctx, params)
}

func EnsureMultiCoinDistributionAccount(ctx sdk.Context, ak authkeeper.AccountKeeper, name string) {
	address := authTypes.NewModuleAddress(name)
	account := ak.GetAccount(ctx, address)

	if account == nil {
		// account doesn't exist, initialise a new module account.
		account = authTypes.NewEmptyModuleAccount(name)
		ak.NewAccount(ctx, account)
	} else {
		// account exists, adjust it to a module account.
		baseAccount := authTypes.NewBaseAccount(address, nil, account.GetAccountNumber(), 0)
		account = authTypes.NewModuleAccount(baseAccount, name)
		ak.SetAccount(ctx, account)
	}
}

func migrateProtocolStakers(ctx sdk.Context, delegationKeeper delegationkeeper.Keeper,
	stakersKeeper *stakerskeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) {
	// Process current unbonding queue
	delegationKeeper.FullyProcessDelegatorUnbondingQueue(ctx)

	var validatorMapping []ValidatorMapping
	if ctx.ChainID() == "kyve-1" {
		validatorMapping = ValidatorMappingsMainnet
	} else if ctx.ChainID() == "kaon-1" {
		validatorMapping = ValidatorMappingsKaon
	} else if ctx.ChainID() == "korellia-2" {
		validatorMapping = ValidatorMappingsKorellia
	}

	totalMigratedStake := uint64(0)
	for _, mapping := range validatorMapping {
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
				panic(err)
			}
		}
	}
	logger.Info(fmt.Sprintf("migrated %d ukyve from protocol to chain", totalMigratedStake))

	// Undelegate Remaining
	totalReturnedStake := uint64(0)
	for _, delegator := range delegationKeeper.GetAllDelegators(ctx) {
		amount := delegationKeeper.PerformFullUndelegation(ctx, delegator.Staker, delegator.Delegator)
		totalReturnedStake += amount
	}
	logger.Info(fmt.Sprintf("returned %d ukyve from protocol to users", totalReturnedStake))

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

	// Delete all legacy stakers objects
	stakersKeeper.Migration_ResetOldState(ctx)

	// Migrate Params
	delegationParams := delegationKeeper.GetParams(ctx)
	stakersParams := stakersKeeper.GetParams(ctx)

	stakersParams.TimeoutSlash = delegationParams.TimeoutSlash
	stakersParams.UploadSlash = delegationParams.UploadSlash
	stakersParams.VoteSlash = delegationParams.VoteSlash

	stakersKeeper.SetParams(ctx, stakersParams)
}
