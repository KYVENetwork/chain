package v1_1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingExported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	// Stakers
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	accountKeeper authKeeper.AccountKeeper,
	stakerKeeper stakersKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if ctx.ChainID() == MainnetChainID {
			for _, address := range InvestorAccounts {
				AdjustInvestorVesting(ctx, accountKeeper, sdk.MustAccAddressFromBech32(address))
			}
		}

		MigrateStakerMetadata(ctx, stakerKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// AdjustInvestorVesting correctly adjusts the vesting schedules of investors
// from our second funding round. In genesis, the accounts were set up with an
// 18-month cliff instead of a 6-month cliff.
func AdjustInvestorVesting(ctx sdk.Context, keeper authKeeper.AccountKeeper, address sdk.AccAddress) {
	rawAccount := keeper.GetAccount(ctx, address)
	account := rawAccount.(vestingExported.VestingAccount)

	baseAccount := authTypes.NewBaseAccount(
		account.GetAddress(), account.GetPubKey(), account.GetAccountNumber(), account.GetSequence(),
	)
	updatedAccount := vestingTypes.NewContinuousVestingAccount(
		baseAccount, account.GetOriginalVesting(), StartTime, EndTime,
	)

	keeper.SetAccount(ctx, updatedAccount)
}

// MigrateStakerMetadata migrates all existing staker metadata. The `Logo`
// field has been deprecated and replaced by the `Identity` field. This new
// field must be a valid hex string; therefore, must be set to empty for now.
func MigrateStakerMetadata(ctx sdk.Context, keeper stakersKeeper.Keeper) {
	for _, staker := range keeper.GetAllStakers(ctx) {
		keeper.UpdateStakerMetadata(ctx, staker.Address, staker.Moniker, staker.Website, "", "", "")
	}
}
