package v1_3

import (
	"fmt"
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
	account, _ := ak.GetAccount(ctx, address).(vestingExported.VestingAccount)

	delegations := sk.GetAllDelegatorDelegations(ctx, address)
	totalDelegation := sdk.NewInt(0)

	for _, delegation := range delegations {
		// TODO: We assume a 1:1 ratio of shares to tokens as investors couldn't
		// perform any actions on delegations post v1.1 upgrade.
		totalDelegation.Add(delegation.GetShares().TruncateInt())
	}

	delegatedVesting := account.GetDelegatedVesting().AmountOf(denom)

	if totalDelegation.GT(delegatedVesting) {
		diff := sdk.NewCoins().Add(sdk.NewCoin(denom, totalDelegation.Sub(delegatedVesting)))
		account.TrackDelegation(ctx.BlockTime(), account.GetOriginalVesting(), diff)

		logger.Info(fmt.Sprintf("tracked delegation of %s with difference %s", address.String(), diff.String()))
	}

	if totalDelegation.LT(delegatedVesting) {
		diff := sdk.NewCoins().Add(sdk.NewCoin(denom, delegatedVesting.Sub(totalDelegation)))
		account.TrackUndelegation(diff)

		logger.Info(fmt.Sprintf("tracked undelegation of %s with difference %s", address.String(), diff.String()))
	}

	ak.SetAccount(ctx, account)
}
