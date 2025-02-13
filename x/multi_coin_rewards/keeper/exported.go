package keeper

import (
	"context"

	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMultiCoinRewards is called by the distribution module to check if the user has enabled multi-coin-rewards.
// If a user has not enabled multi-coin rewards, all tokens (except the native denom) are added to a queue.
// The user then has `MultiCoinDistributionPendingTime` seconds to enable Multi-Coin rewards and claim the pending rewards.
// Otherwise, these tokens will get redistributed to the pools.
func (k Keeper) HandleMultiCoinRewards(goCtx context.Context, withdrawAddress sdk.AccAddress, coins sdk.Coins) (sdk.Coins, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	has, err := k.MultiCoinRewardsEnabled.Has(ctx, withdrawAddress)
	if err != nil {
		return nil, err
	}

	if has {
		// User has enabled multi-coin rewards, all coins are available
		return coins, nil
	}

	// MultiCoinRewards is not enabled
	enabledRewards := sdk.NewCoins(sdk.NewCoin(globalTypes.Denom, coins.AmountOf(globalTypes.Denom)))
	disabledRewards := coins.Sub(enabledRewards...)

	// Add Pending-Queue entry
	if !disabledRewards.Empty() {
		// Transfer non-enabled rewards
		if err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, distributionTypes.ModuleName, types.ModuleName, disabledRewards); err != nil {
			return nil, err
		}
		k.addPendingRewards(ctx, withdrawAddress.String(), disabledRewards)
	}

	return enabledRewards, nil
}
