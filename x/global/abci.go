package global

import (
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Global
	"github.com/KYVENetwork/chain/x/global/keeper"
	// Upgrade
	upgradeKeeper "cosmossdk.io/x/upgrade/keeper"
)

// EndBlocker handles the fee burning if it is configured
func EndBlocker(ctx sdk.Context, ak authKeeper.AccountKeeper, bk bankKeeper.Keeper, gk keeper.Keeper, uk upgradeKeeper.Keeper) {
	// Since no fees are paid in the genesis block, skip.
	// NOTE: This is Tendermint specific.
	if ctx.BlockHeight() == 1 {
		return
	}

	burnRatio := gk.GetBurnRatio(ctx)
	if burnRatio.IsZero() {
		return
	}

	// Obtain all collected fees.
	feeCoinsInt := bk.GetAllBalances(ctx, ak.GetModuleAddress(authTypes.FeeCollectorName))
	feeCoins := sdk.NewDecCoinsFromCoins(feeCoinsInt...)
	if feeCoins.IsZero() {
		return
	}

	// Sum burn ratio amount.
	burnCoins := sdk.NewCoins()
	for _, coin := range feeCoins {
		amount := coin.Amount.Mul(burnRatio)
		burnCoins = burnCoins.Add(sdk.NewCoin(coin.Denom, amount.TruncateInt()))
	}

	err := bk.BurnCoins(ctx, authTypes.FeeCollectorName, burnCoins)
	if err != nil {
		util.PanicHalt(uk, ctx, err.Error())
	}
}
