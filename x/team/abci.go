package team

import (
	"fmt"

	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	// Team
	"github.com/KYVENetwork/chain/x/team/keeper"
	"github.com/KYVENetwork/chain/x/team/types"
	// Upgrade
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
)

func DistributeTeamInflation(ctx sdk.Context, bk bankKeeper.Keeper, tk keeper.Keeper, uk upgradeKeeper.Keeper) {
	// get the total team reward the module is eligible for in this block
	teamModuleRewards := tk.GetTeamBlockProvision(ctx)

	// count total account rewards
	totalAccountRewards := uint64(0)

	// distribute team module rewards between vesting accounts based on their vesting progress
	for _, account := range tk.GetTeamVestingAccounts(ctx) {
		// get current vesting progress
		status := keeper.GetVestingStatus(account, uint64(ctx.BlockTime().Unix()))
		// calculate reward share of account
		accountShare := sdk.NewDec(int64(status.TotalVestedAmount - account.UnlockedClaimed)).Quo(sdk.NewDec(int64(types.TEAM_ALLOCATION)))
		// calculate total inflation rewards for account for this block
		accountRewards := uint64(sdk.NewDec(teamModuleRewards).Mul(accountShare).TruncateInt64())

		// save inflation rewards to account
		account.TotalRewards += accountRewards
		tk.SetTeamVestingAccount(ctx, account)

		// count total inflation rewards for team module
		totalAccountRewards += accountRewards
	}

	// panic if total account rewards are higher than team module rewards
	if totalAccountRewards > uint64(teamModuleRewards) {
		util.PanicHalt(uk, ctx, fmt.Sprintf("account rewards %v are higher than entire team module rewards %v", totalAccountRewards, teamModuleRewards))
	}

	// track total authority inflation rewards
	authority := tk.GetAuthority(ctx)
	authority.TotalRewards += uint64(teamModuleRewards) - totalAccountRewards
	tk.SetAuthority(ctx, authority)

	// distribute part of block provision to team module
	if err := util.TransferFromModuleToModule(bk, ctx, authTypes.FeeCollectorName, types.ModuleName, uint64(teamModuleRewards)); err != nil {
		util.PanicHalt(uk, ctx, err.Error())
	}

	tk.Logger(ctx).Info("distributed portion of minted coins", "amount", teamModuleRewards)
}
