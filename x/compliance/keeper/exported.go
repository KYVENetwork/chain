package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/compliance/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// HandleMultiCoinRewards is called by the distribution module to check the rewards for non-compliant tokens.
// If a user has not enabled multi-coin rewards, all tokens (except the native denom) are added to a queue.
// The user then has `MultiCoinRefundPendingTime` seconds to enable Multi-Coin rewards and claim the pending rewards.
// Otherwise, these tokens will get redistributed to the pools.
func (k Keeper) HandleMultiCoinRewards(goCtx context.Context, withdrawAddress sdk.AccAddress, coins sdk.Coins) sdk.Coins {
	ctx := sdk.UnwrapSDKContext(goCtx)

	has, err := k.MultiCoinRewardsEnabled.Has(ctx, withdrawAddress)
	if err == nil && has {
		// User has enabled multi-coin rewards, all coins are compliant
		return coins
	}

	// MultiCoinRewards is not enabled
	compliantRewards := sdk.NewCoins(sdk.NewCoin(globalTypes.Denom, coins.AmountOf(globalTypes.Denom)))
	nonCompliantRewards := coins.Sub(compliantRewards...)

	// Transfer non-compliant rewards
	if err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, distributionTypes.ModuleName, types.ModuleName, nonCompliantRewards); err != nil {
		panic(err)
	}
	// Add Pending-Queue entry
	if !nonCompliantRewards.Empty() {
		k.addPendingComplianceRewards(ctx, withdrawAddress.String(), nonCompliantRewards)
	}

	return compliantRewards
}
