package v2_0

import (
	"context"
	"fmt"

	poolTypes "github.com/KYVENetwork/chain/x/pool/types"

	poolkeeper "github.com/KYVENetwork/chain/x/pool/keeper"

	multicoinrewardskeeper "github.com/KYVENetwork/chain/x/multi_coin_rewards/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

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
	distrKeeper *distrkeeper.Keeper,
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

		// Run Bundles Merkle Roots migrations
		bundlesKeeper.SetBundlesMigrationUpgradeHeight(sdkCtx, uint64(sdkCtx.BlockHeight()))
		UpgradeRuntimes(sdkCtx, poolKeeper)
		UpdateUploadIntervals(sdkCtx, poolKeeper)

		// TODO: update coin weights and storage cost for mainnet

		// Set MultiCoinRewards and Withdraw address for the KYVE Foundation
		if sdkCtx.ChainID() == "kyve-1" {
			SetWithdrawAddressAndMultiCoinRewards(
				sdkCtx, multiCoinRewardsKeeper, accountKeeper, distrKeeper,
				"kyve1dur8kw9qh28p00urmjulmnxyt0m34k3j7veehz", "kyve173pnpz27lcn6zq4x37392n09y8mnz5vadjx2m9")
		}

		logger.Info(fmt.Sprintf("finished upgrade %v", UpgradeName))

		return migratedVersionMap, err
	}
}

// SetWithdrawAddressAndMultiCoinRewards sets a withdraw-address and enables multi-coin rewards for
// a given delegator
func SetWithdrawAddressAndMultiCoinRewards(
	ctx sdk.Context,
	multiCoinRewardsKeeper multicoinrewardskeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	distrKeeper *distrkeeper.Keeper,
	delegatorAddress string,
	withdrawAddress string,
) {
	delegatorAccAddress, err := accountKeeper.AddressCodec().StringToBytes(delegatorAddress)
	if err != nil {
		panic(err)
	}

	withdrawAccAddress, err := accountKeeper.AddressCodec().StringToBytes(withdrawAddress)
	if err != nil {
		panic(err)
	}

	if err = distrKeeper.SetWithdrawAddr(ctx, delegatorAccAddress, withdrawAccAddress); err != nil {
		panic(err)
	}

	if err := multiCoinRewardsKeeper.MultiCoinRewardsEnabled.Set(ctx, delegatorAccAddress); err != nil {
		panic(err)
	}
}

func SetMultiCoinRewardsParams(ctx sdk.Context, multiCoinRewardsKeeper multicoinrewardskeeper.Keeper) {
	params := multiCoinRewardsKeeper.GetParams(ctx)

	if ctx.ChainID() == "kyve-1" {
		params.MultiCoinDistributionPendingTime = 60 * 60 * 24 * 14
		// KYVE Public Good Funding address
		params.MultiCoinDistributionPolicyAdminAddress = "kyve1t0uez3nn28ljnzlwndzxffyjuhean3edhtjee8"
	} else if ctx.ChainID() == "kaon-1" {
		params.MultiCoinDistributionPendingTime = 60 * 60 * 24
		// Kaon Ecosystem
		params.MultiCoinDistributionPolicyAdminAddress = "kyve1z67jal9d9unjvmzsadps9jytzt9kx2m2vgc3wm"
	} else if ctx.ChainID() == "korellia-2" {
		params.MultiCoinDistributionPendingTime = 600
		// Korellia Ecosystem
		params.MultiCoinDistributionPolicyAdminAddress = "kyve1ygtqlzxhvp3t0wwcjd5lmq4zxl0qcck9g3mmgp"
	}

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

	return multiCoinRewardsKeeper.MultiCoinDistributionPolicy.Set(ctx, multicoinrewardstypes.MultiCoinDistributionPolicy{
		Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{},
	})
}

func SetPoolParams(ctx sdk.Context, poolKeeper *poolkeeper.Keeper) {
	params := poolKeeper.GetParams(ctx)

	if ctx.ChainID() == "kyve-1" {
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.4")
		params.MaxVotingPowerPerPool = math.LegacyNewDec(1).QuoInt64(3)
	}

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

func UpgradeRuntimes(sdkCtx sdk.Context, poolKeeper *poolkeeper.Keeper) {
	// Upgrade duration set to 10mins
	upgrades := []poolTypes.MsgScheduleRuntimeUpgrade{
		{
			Runtime:     "@kyvejs/tendermint",
			Version:     "1.3.0",
			ScheduledAt: uint64(sdkCtx.BlockTime().Unix()),
			Duration:    600,
			Binaries:    "{\"kyve-linux-arm64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint%401.3.0/kyve-linux-arm64.zip\",\"kyve-linux-x64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint%401.3.0/kyve-linux-x64.zip\",\"kyve-macos-x64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint%401.3.0/kyve-macos-x64.zip\"}",
		},
		{
			Runtime:     "@kyvejs/tendermint-bsync",
			Version:     "1.2.9",
			ScheduledAt: uint64(sdkCtx.BlockTime().Unix()),
			Duration:    600,
			Binaries:    "{\"kyve-linux-arm64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint-bsync%401.2.9/kyve-linux-arm64.zip\",\"kyve-linux-x64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint-bsync%401.2.9/kyve-linux-x64.zip\",\"kyve-macos-x64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint-bsync%401.2.9/kyve-macos-x64.zip\"}",
		},
		{
			Runtime:     "@kyvejs/tendermint-ssync",
			Version:     "1.3.0",
			ScheduledAt: uint64(sdkCtx.BlockTime().Unix()),
			Duration:    600,
			Binaries:    "{\"kyve-linux-arm64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint-ssync%401.3.0/kyve-linux-arm64.zip\",\"kyve-linux-x64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint-ssync%401.3.0/kyve-linux-x64.zip\",\"kyve-macos-x64\":\"https://github.com/KYVENetwork/kyvejs/releases/download/%40kyvejs%2Ftendermint-ssync%401.3.0/kyve-macos-x64.zip\"}",
		},
	}

	for _, upgrade := range upgrades {
		affectedPools := make([]uint64, 0)
		for _, pool := range poolKeeper.GetAllPools(sdkCtx) {
			// only schedule upgrade if the runtime matches
			if pool.Runtime != upgrade.Runtime {
				continue
			}

			// only schedule upgrade if there is no upgrade already
			if pool.UpgradePlan.ScheduledAt != 0 {
				continue
			}

			pool.UpgradePlan = &poolTypes.UpgradePlan{
				Version:     upgrade.Version,
				Binaries:    upgrade.Binaries,
				ScheduledAt: upgrade.ScheduledAt,
				Duration:    upgrade.Duration,
			}

			affectedPools = append(affectedPools, pool.Id)

			poolKeeper.SetPool(sdkCtx, pool)
		}

		_ = sdkCtx.EventManager().EmitTypedEvent(&poolTypes.EventRuntimeUpgradeScheduled{
			Runtime:       upgrade.Runtime,
			Version:       upgrade.Version,
			ScheduledAt:   upgrade.ScheduledAt,
			Duration:      upgrade.Duration,
			Binaries:      upgrade.Binaries,
			AffectedPools: affectedPools,
		})
	}
}

func UpdateUploadIntervals(sdkCtx sdk.Context, poolKeeper *poolkeeper.Keeper) {
	if sdkCtx.ChainID() == "kyve-1" {
		for _, pool := range poolKeeper.GetAllPools(sdkCtx) {
			if pool.Id == 4 || pool.Id == 8 || pool.Id == 10 {
				pool.UploadInterval = 120
				poolKeeper.SetPool(sdkCtx, pool)
			}
		}
	}
}
