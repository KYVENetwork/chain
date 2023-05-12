package v1_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Team
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	teamKeeper teamKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if ctx.ChainID() == MainnetChainID {
			AdjustTeamVesting(ctx, teamKeeper)
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// AdjustTeamVesting adjusts the commencement (vesting start) and clawback
// (vesting end) dates of multiple members belonging to the KYVE Core Team.
func AdjustTeamVesting(ctx sdk.Context, keeper teamKeeper.Keeper) {
	// ----- Account ID = 8 -----
	accountEight, _ := keeper.GetTeamVestingAccount(ctx, 8)
	accountEight.Commencement = 1614556800 // Mar 1st, 2021.

	keeper.SetTeamVestingAccount(ctx, accountEight)

	// ----- Account ID = 11 -----
	accountEleven, _ := keeper.GetTeamVestingAccount(ctx, 11)
	accountEleven.Commencement = 1639094400 // Dec 10th, 2021.
	accountEleven.Clawback = 1682899200     // May 1st, 2023.

	keeper.SetTeamVestingAccount(ctx, accountEleven)
}
