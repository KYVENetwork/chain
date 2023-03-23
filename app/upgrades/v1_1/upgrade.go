package v1_1

import (
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingExported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
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

		MigrateStakerLogo(ctx, stakerKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// AdjustInvestorVesting correctly adjusts the vesting schedules of investors
// from our second funding round. In genesis, the accounts were set up with an
// 18-month cliff instead of a 6-month cliff.
func AdjustInvestorVesting(ctx sdk.Context, accountKeeper authKeeper.AccountKeeper, address sdk.AccAddress) {
	rawAccount := accountKeeper.GetAccount(ctx, address)
	account := rawAccount.(vestingExported.VestingAccount)

	baseAccount := authTypes.NewBaseAccount(
		account.GetAddress(), account.GetPubKey(), account.GetAccountNumber(), account.GetSequence(),
	)
	updatedAccount := vestingTypes.NewContinuousVestingAccount(
		baseAccount, account.GetOriginalVesting(), StartTime, EndTime,
	)

	accountKeeper.SetAccount(ctx, updatedAccount)
}

// MigrateStakerLogo migrates all existing staker metadata. The `Logo` field got replaced by `Identity`
// and must be a valid hex string. Therefore, all Logo URLs are reset to empty strings and the stakers
// have to use the UpdateMetadata message to set their identity.
func MigrateStakerLogo(ctx sdk.Context, keeper stakersKeeper.Keeper) {
	for _, staker := range keeper.GetAllStakers(ctx) {
		keeper.UpdateStakerMetadata(ctx, staker.Address, staker.Moniker, staker.Website, "", "", "")
	}
}
