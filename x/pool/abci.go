package pool

import (
	"fmt"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Mint
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/keeper"
	// Upgrade
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
)

func SplitInflation(ctx sdk.Context, bk bankKeeper.Keeper, mk mintKeeper.Keeper, tk teamKeeper.Keeper, pk keeper.Keeper, uk upgradeKeeper.Keeper) {
	// Compute team allocation of minted coins.
	minter := mk.GetMinter(ctx)
	params := mk.GetParams(ctx)

	// protocol_split TODO: move to pool params
	protocolSplit, err := sdk.NewDecFromStr("0.04")
	if err != nil {
		panic(err)
	}

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

	// calculate the remaining block provision after x/team took its share
	remainingBlockProvision := blockProvision.Amount.Int64() - teamModuleRewardsShare.Mul(sdk.NewDec(blockProvision.Amount.Int64())).TruncateInt64()
	protocolAmount := sdk.NewDec(remainingBlockProvision).Mul(protocolSplit)

	totalOperatingCost := uint64(0)

	for _, pool := range pk.GetAllPools(ctx) {
		totalOperatingCost += pool.OperatingCost
	}

	if totalOperatingCost == 0 {
		return
	}

	for _, pool := range pk.GetAllPools(ctx) {
		amount := uint64(sdk.NewDec(int64(pool.OperatingCost)).Quo(sdk.NewDec(int64(totalOperatingCost))).Mul(protocolAmount).TruncateInt64())

		if err := util.TransferFromModuleToAddress(bk, ctx, authTypes.FeeCollectorName, pool.GetPoolAccount(), amount); err != nil {
			util.PanicHalt(uk, ctx, err.Error())
		}
	}

	pk.Logger(ctx).Info("split portion of minted coins to protocol", "amount", protocolAmount.TruncateInt64())
}
