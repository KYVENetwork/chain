package v1_3

import (
	"fmt"

	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/log"

	// Pool
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	// Upgrade
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingExported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	// Staking
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	poolKeeper poolKeeper.Keeper,
	accountKeeper authKeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	stakingKeeper stakingKeeper.Keeper,
	bundlesKeeper bundlesKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		CheckPoolAccounts(ctx, logger, poolKeeper)
		UpdateBundlesVersionMap(bundlesKeeper, ctx)

		SetBundleParams(ctx, poolKeeper)

		if ctx.ChainID() == MainnetChainID {
			for _, address := range InvestorAccounts {
				TrackInvestorDelegation(ctx, logger, sdk.MustAccAddressFromBech32(address), accountKeeper, bankKeeper, stakingKeeper)
			}
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// TrackInvestorDelegation performs a correction of the delegation tracking inside the vesting account.
// The correction is done by performing a full untracking and then tracking the actual total delegated amount
// (including slashed amounts).
func TrackInvestorDelegation(ctx sdk.Context, logger log.Logger, address sdk.AccAddress, ak authKeeper.AccountKeeper, bk bankKeeper.Keeper, sk stakingKeeper.Keeper) {
	denom := sk.BondDenom(ctx)
	account, _ := ak.GetAccount(ctx, address).(vestingExported.VestingAccount)

	// Obtain total delegation of address
	totalDelegation := sdk.NewInt(0)
	for _, delegation := range sk.GetAllDelegatorDelegations(ctx, address) {
		// We take the shares as the total delegation as this is the amount which is
		// tracked inside the vesting account. (slashes are ignored, which is correct)
		totalDelegation = totalDelegation.Add(delegation.GetShares().TruncateInt())
	}

	// Fetch current balance.
	balanceCoin := bk.GetBalance(ctx, address, denom)

	// This is the balance a user would have if all tokens are unbonded (even the ones which got slashed).
	maxPossibleBalance := balanceCoin.Amount.Add(totalDelegation)
	maxPossibleBalanceCoins := sdk.NewCoins().Add(sdk.NewCoin(denom, maxPossibleBalance))

	if totalDelegation.GT(sdk.ZeroInt()) {

		// Untrack entire vesting delegation using maximum amount. This will set both `delegated_free`
		// and `delegated_vesting` back to zero.
		account.TrackUndelegation(sdk.NewCoins(sdk.NewCoin("ukyve", maxPossibleBalance)))

		// Track the delegation using the total delegation
		account.TrackDelegation(ctx.BlockTime(), maxPossibleBalanceCoins, sdk.NewCoins(sdk.NewCoin("ukyve", totalDelegation)))

		logger.Info(fmt.Sprintf("tracked delegation of %s with %s", address.String(), totalDelegation.String()))
		ak.SetAccount(ctx, account)
	}
}

// CheckPoolAccounts ensures that each pool account exists post upgrade.
func CheckPoolAccounts(ctx sdk.Context, logger log.Logger, keeper poolKeeper.Keeper) {
	pools := keeper.GetAllPools(ctx)

	for _, pool := range pools {
		keeper.EnsurePoolAccount(ctx, pool.Id)

		name := fmt.Sprintf("%s/%d", poolTypes.ModuleName, pool.Id)
		logger.Info("successfully initialised pool account", "name", name)
	}
}

func UpdateBundlesVersionMap(keeper bundlesKeeper.Keeper, ctx sdk.Context) {
	keeper.SetBundleVersionMap(ctx, bundlesTypes.BundleVersionMap{
		Versions: []*bundlesTypes.BundleVersionEntry{
			{
				Height:  uint64(ctx.BlockHeight()),
				Version: 2,
			},
		},
	})
}

func SetBundleParams(ctx sdk.Context, keeper poolKeeper.Keeper) {
	keeper.SetParams(ctx, poolTypes.Params{
		ProtocolInflationShare:  sdk.MustNewDecFromStr("0.08"),
		PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
	})
}
