package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Delegator returns all delegation info
func (k Keeper) Delegator(goCtx context.Context, req *types.QueryDelegatorRequest) (*types.QueryDelegatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryDelegatorResponse{}

	// Load pool
	_, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	// Load delegator
	delegator, found := k.GetDelegator(ctx, req.PoolId, req.Staker, req.Delegator)

	if !found {
		response.Delegator = &types.StakerDelegatorResponse{
			Delegator: "",
			CurrentReward: 0,
			DelegationAmount: 0,
			Staker: "",
		}

		return &response, nil
	}

	f1 := F1Distribution{
		k:                k,
		ctx:              ctx,
		poolId:           req.PoolId,
		stakerAddress:    delegator.Staker,
		delegatorAddress: delegator.Delegator,
	}

	response.Delegator = &types.StakerDelegatorResponse{
		Delegator: delegator.Delegator,
		CurrentReward: f1.getCurrentReward(),
		DelegationAmount: delegator.DelegationAmount,
		Staker: delegator.Staker,
	}

	return &response, nil
}
