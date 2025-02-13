package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) MultiCoinDistributionPolicyQuery(ctx context.Context, request *types.QueryMultiCoinDistributionPolicyRequest) (*types.QueryMultiCoinDistributionPolicyResponse, error) {
	policy, err := k.MultiCoinDistributionPolicy.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryMultiCoinDistributionPolicyResponse{Policy: policy}, nil
}

func (k Keeper) MultiCoinStatus(ctx context.Context, request *types.QueryMultiCoinStatusRequest) (*types.QueryMultiCoinStatusResponse, error) {
	account, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}

	has, err := k.MultiCoinRewardsEnabled.Has(ctx, account)
	if err != nil {
		return nil, err
	}

	entries, _ := k.GetMultiCoinPendingRewardsEntriesByIndex2(sdk.UnwrapSDKContext(ctx), request.Address)

	pendingRewards := sdk.NewCoins()

	for _, entry := range entries {
		pendingRewards = pendingRewards.Add(entry.Rewards...)
	}

	return &types.QueryMultiCoinStatusResponse{
		Enabled:                 has,
		PendingMultiCoinRewards: pendingRewards,
	}, nil
}
