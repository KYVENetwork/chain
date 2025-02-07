package keeper

import (
	"context"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// HandleMultiCoinRewards is called by the distribution module to check if the user has enabled multi-coin-rewards.
// If a user has not enabled multi-coin rewards, all tokens (except the native denom) are added to a queue.
// The user then has `MultiCoinDistributionPendingTime` seconds to enable Multi-Coin rewards and claim the pending rewards.
// Otherwise, these tokens will get redistributed to the pools.
func (k Keeper) HandleMultiCoinRewards(goCtx context.Context, withdrawAddress sdk.AccAddress, coins sdk.Coins) sdk.Coins {
	ctx := sdk.UnwrapSDKContext(goCtx)

	has, err := k.MultiCoinRewardsEnabled.Has(ctx, withdrawAddress)
	if err == nil && has {
		// User has enabled multi-coin rewards, all coins are available
		return coins
	}

	// MultiCoinRewards is not enabled
	enabledRewards := sdk.NewCoins(sdk.NewCoin(globalTypes.Denom, coins.AmountOf(globalTypes.Denom)))
	disabledRewards := coins.Sub(enabledRewards...)

	// Transfer non-enabled rewards
	if err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, distributionTypes.ModuleName, types.ModuleName, disabledRewards); err != nil {
		panic(err)
	}
	// Add Pending-Queue entry
	if !disabledRewards.Empty() {
		k.addPendingRewards(ctx, withdrawAddress.String(), disabledRewards)
	}

	return enabledRewards
}
