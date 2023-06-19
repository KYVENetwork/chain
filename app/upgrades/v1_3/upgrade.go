package v1_3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/log"

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
	accountKeeper authKeeper.AccountKeeper,
	stakingKeeper stakingKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		if ctx.ChainID() == MainnetChainID {
			for _, address := range InvestorAccounts {
				TrackInvestorDelegation(ctx, logger, sdk.MustAccAddressFromBech32(address), accountKeeper, stakingKeeper)
			}
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// TrackInvestorDelegation ...
func TrackInvestorDelegation(ctx sdk.Context, logger log.Logger, address sdk.AccAddress, ak authKeeper.AccountKeeper, sk stakingKeeper.Keeper) {
	denom := sk.BondDenom(ctx)
	rawAccount := ak.GetAccount(ctx, address)
	account, _ := rawAccount.(vestingExported.VestingAccount)

	delegations := sk.GetAllDelegatorDelegations(ctx, address)
	totalDelegation := sdk.NewCoins()

	for _, delegation := range delegations {
		// TODO: We assume a 1:1 ratio of shares to tokens as investors couldn't
		// perform any actions on delegations post v1.1 upgrade.
		totalDelegation.Add(sdk.NewCoin(denom, delegation.GetShares().TruncateInt()))
	}

	trackedDifference := totalDelegation.Sub(account.GetDelegatedVesting()...)
	if !trackedDifference.IsZero() {
		// TODO: We assume that the usable balance is the total vesting amount
		// as the investor cliff is still ongoing.
		account.TrackDelegation(ctx.BlockTime(), account.GetOriginalVesting(), trackedDifference)

		ak.SetAccount(ctx, account)
		logger.Info("fixed vesting account tracked delegation", "difference", trackedDifference.String())
	}
}
