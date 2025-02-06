package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/compliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ToggleMultiCoinRewards(ctx context.Context, compliance *types.MsgToggleMultiCoinRewards) (*types.MsgToggleMultiCoinRewardsResponse, error) {
	accountAddress, err := sdk.AccAddressFromBech32(compliance.Creator)
	if err != nil {
		return nil, err
	}

	// If entry exists, it means that the user has currently multi-coin rewards enabled
	rewardsCurrentlyEnabled, err := k.MultiCoinRewardsEnabled.Has(ctx, accountAddress)
	if err != nil {
		return nil, err
	}

	if compliance.Enabled {
		// User wants to enable multi-coin rewards

		if rewardsCurrentlyEnabled {
			return nil, types.ErrMultiCoinRewardsAlreadyEnabled
		}

		rewards, _ := k.GetMultiCoinPendingRewardsEntriesByIndex2(sdk.UnwrapSDKContext(ctx), accountAddress.String())
		totalRewards := sdk.NewCoins()
		for _, reward := range rewards {
			totalRewards = totalRewards.Add(reward.Rewards...)
		}

		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accountAddress, totalRewards)
		if err != nil {
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

	return &types.MsgToggleMultiCoinRewardsResponse{}, nil
}

func (k msgServer) SetMultiCoinRewardRefundPolicy(goCtx context.Context, policy *types.MsgSetMultiCoinRewardsRefundPolicy) (*types.MsgSetMultiCoinRewardsRefundPolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	if params.MultiCoinRefundPolicyAdminAddress != policy.Creator {
		return nil, types.ErrMultiCoinRefundPolicyInvalidAdminAddress
	}

	if err := k.MultiCoinRefundPolicy.Set(ctx, *policy.Policy); err != nil {
		return nil, errors.Wrap(err, types.ErrMultiCoinRefundPolicyInvalid.Error())
	}

	return &types.MsgSetMultiCoinRewardsRefundPolicyResponse{}, nil
}
