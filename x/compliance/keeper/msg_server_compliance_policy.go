package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/compliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
