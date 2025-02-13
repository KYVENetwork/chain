package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetMultiCoinRewardDistributionPolicy(goCtx context.Context, policy *types.MsgSetMultiCoinRewardsDistributionPolicy) (*types.MsgSetMultiCoinRewardsDistributionPolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	if params.MultiCoinDistributionPolicyAdminAddress != policy.Creator {
		return nil, types.ErrMultiCoinDistributionPolicyInvalidAdminAddress
	}

	if err := k.MultiCoinDistributionPolicy.Set(ctx, *policy.Policy); err != nil {
		return nil, errors.Wrap(err, types.ErrMultiCoinDistributionPolicyInvalid.Error())
	}

	return &types.MsgSetMultiCoinRewardsDistributionPolicyResponse{}, nil
}
