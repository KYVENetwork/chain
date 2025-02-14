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
		if err := SetMultiCoinRewardsPolicy(sdkCtx, multiCoinRewardsKeeper); err != nil {
			return migratedVersionMap, err
		}

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

func SetMultiCoinRewardsPolicy(ctx sdk.Context, multiCoinRewardsKeeper multicoinrewardskeeper.Keeper) error {
	if ctx.ChainID() == "kyve-1" {
		return multiCoinRewardsKeeper.MultiCoinDistributionPolicy.Set(ctx, multicoinrewardstypes.MultiCoinDistributionPolicy{
			Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{
				{
					Denom: "ibc/A59C9E368C043E72968615DE82D4AD4BC88E34E6F353262B6769781C07390E8A", // andromeda
					PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 14,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					},
				},
				{
					Denom: "ibc/F4E5517A3BA2E77906A0847014EBD39E010E28BEB4181378278144D22442DB91", // source
					PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 11,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
						{
							PoolId: 12,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					},
				},
				{
					Denom: "ibc/D0C5DCA29836D2FD5937714B21206DD8243E5E76B1D0F180741CCB43DCAC1584", // dydx
					PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 13,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					},
				},
				{
					Denom: "ibc/506478E08FB0A2D3B12D493E3B182572A3B0D7BD5DCBE71610D2F393DEDDF4CA", // xion
					PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 16,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
						{
							PoolId: 17,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					},
				},
				{
					Denom: "ibc/7D5A9AE91948931279BA58A04FBEB9BF4F7CA059F7D4BDFAC6C3C43705973E1E", // lava
					PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 18,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					},
				},
			},
		})
	}

	return nil
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
			Type:     "/kyve.bundles.v1beta1.MsgSubmitBundleProposal",
			Fraction: math.LegacyMustNewDecFromStr("0.99"),
		},
		{
			Type:     "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
			Fraction: math.LegacyMustNewDecFromStr("0.99"),
		},
		{
			Type:     "/kyve.bundles.v1beta1.MsgSkipUploaderRole",
			Fraction: math.LegacyMustNewDecFromStr("0.99"),
		},
		{
			Type:     "/kyve.bundles.v1beta1.MsgClaimUploaderRole",
			Fraction: math.LegacyMustNewDecFromStr("0.99"),
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
