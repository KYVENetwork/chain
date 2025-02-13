package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ToggleMultiCoinRewards(ctx context.Context, toggle *types.MsgToggleMultiCoinRewards) (*types.MsgToggleMultiCoinRewardsResponse, error) {
	accountAddress, err := sdk.AccAddressFromBech32(toggle.Creator)
	if err != nil {
		return nil, err
	}

	// If entry exists, it means that the user has currently multi-coin rewards enabled
	rewardsCurrentlyEnabled, err := k.MultiCoinRewardsEnabled.Has(ctx, accountAddress)
	if err != nil {
		return nil, err
	}

	totalRewards := sdk.NewCoins()
	if toggle.Enabled {
		// User wants to enable multi-coin rewards

		if rewardsCurrentlyEnabled {
			return nil, types.ErrMultiCoinRewardsAlreadyEnabled
		}

		rewards, _ := k.GetMultiCoinPendingRewardsEntriesByIndex2(sdk.UnwrapSDKContext(ctx), accountAddress.String())
		for _, reward := range rewards {
			totalRewards = totalRewards.Add(reward.Rewards...)
		}

		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accountAddress, totalRewards); err != nil {
			return nil, err
		}

		if err := k.MultiCoinRewardsEnabled.Set(ctx, accountAddress); err != nil {
			return nil, err
		}
	} else {
		// User wants to disable multi-coin rewards

		if !rewardsCurrentlyEnabled {
			return nil, types.ErrMultiCoinRewardsAlreadyDisabled
		}

		if err := k.MultiCoinRewardsEnabled.Remove(ctx, accountAddress); err != nil {
			return nil, err
		}
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventToggleMultiCoinRewards{
		Address:               toggle.Creator,
		Enabled:               toggle.Enabled,
		PendingRewardsClaimed: totalRewards.String(),
	})

	return &types.MsgToggleMultiCoinRewardsResponse{}, nil
}
