package v1_1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingExported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	// IBC Transfer
	transferKeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	transferTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
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
	transferKeeper transferKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if ctx.ChainID() == MainnetChainID {
			for _, address := range InvestorAccounts {
				AdjustInvestorVesting(ctx, accountKeeper, sdk.MustAccAddressFromBech32(address))
			}
		}

		EnableIBCTransfers(ctx, transferKeeper)
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

// EnableIBCTransfers updates the parameters of the IBC Transfer module to
// allow both sending and receiving of IBC tokens. Since the default parameters
// of the module have everything enabled, we simply switch to the defaults.
func EnableIBCTransfers(ctx sdk.Context, keeper transferKeeper.Keeper) {
	params := transferTypes.DefaultParams()
	keeper.SetParams(ctx, params)
}

// MigrateStakerMetadata migrates all existing staker metadata. The `Logo`
// field has been deprecated and replaced by the `Identity` field. This new
// field must be a valid hex string; therefore, must be set to empty for now.
func MigrateStakerMetadata(ctx sdk.Context, keeper stakersKeeper.Keeper) {
	for _, staker := range keeper.GetAllStakers(ctx) {
		keeper.UpdateStakerMetadata(ctx, staker.Address, staker.Moniker, staker.Website, "", "", "")
	}
}
