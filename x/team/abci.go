package team

import (
	"fmt"

	"github.com/KYVENetwork/chain/util"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Mint
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	// Team
	"github.com/KYVENetwork/chain/x/team/keeper"
	// Upgrade
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
)

func DistributeTeamInflation(ctx sdk.Context, bk bankKeeper.Keeper, mk mintKeeper.Keeper, tk keeper.Keeper, uk upgradeKeeper.Keeper) {
	// Compute team allocation of minted coins.
	minter := mk.GetMinter(ctx)
	params := mk.GetParams(ctx)

	// get total inflation rewards for current block
	blockProvision := minter.BlockProvision(params)

	// calculate theoretical team balance. We don't use team module balance because a third party could skew
	// the team inflation rewards by simply transferring funds to the team module account
	teamBalance := tk.GetTeamInfo(ctx).RequiredModuleBalance

	// calculate total inflation rewards for team module.
	// We subtract current inflation because it was already applied to the total supply because BeginBlocker
	// x/mint runs before this method
	totalSupply := bk.GetSupply(ctx, blockProvision.Denom).Amount.Int64() - blockProvision.Amount.Int64()
	teamModuleRewardsShare := sdk.NewDec(int64(teamBalance)).Quo(sdk.NewDec(totalSupply))

	// if team module balance is greater than total supply panic
	if teamModuleRewardsShare.GT(sdk.NewDec(int64(1))) {
		util.PanicHalt(uk, ctx, fmt.Sprintf("team module balance %v is higher than total supply %v", teamBalance, totalSupply))
	}

	// calculate the total reward in $KYVE the entire team module receives this block
	teamModuleRewards := uint64(teamModuleRewardsShare.Mul(sdk.NewDec(blockProvision.Amount.Int64())).TruncateInt64())

	// count total account rewards
	totalAccountRewards := uint64(0)

	// distribute team module rewards between vesting accounts based on their vesting progress
	for _, account := range tk.GetTeamVestingAccounts(ctx) {
		// get current vesting progress
		status := teamKeeper.GetVestingStatus(account, uint64(ctx.BlockTime().Unix()))
		// calculate reward share of account
		accountShare := sdk.NewDec(int64(status.TotalVestedAmount - account.UnlockedClaimed)).Quo(sdk.NewDec(int64(types.TEAM_ALLOCATION)))
		// calculate total inflation rewards for account for this block
		accountRewards := uint64(sdk.NewDec(int64(teamModuleRewards)).Mul(accountShare).TruncateInt64())

		// save inflation rewards to account
		account.TotalRewards += accountRewards
		tk.SetTeamVestingAccount(ctx, account)

		// count total inflation rewards for team module
		totalAccountRewards += accountRewards
	}

	// panic if total account rewards are higher than team module rewards
	if totalAccountRewards > teamModuleRewards {
		util.PanicHalt(uk, ctx, fmt.Sprintf("account rewards %v are higher than entire team module rewards %v", totalAccountRewards, teamModuleRewards))
	}

	// track total authority inflation rewards
	authority := tk.GetAuthority(ctx)
	authority.TotalRewards += teamModuleRewards - totalAccountRewards
	tk.SetAuthority(ctx, authority)

	// distribute part of block provision to team module
	if err := util.TransferFromModuleToModule(bk, ctx, authTypes.FeeCollectorName, types.ModuleName, teamModuleRewards); err != nil {
		util.PanicHalt(uk, ctx, err.Error())
	}

	tk.Logger(ctx).Info("distributed portion of minted coins", "amount", teamModuleRewards)
}
