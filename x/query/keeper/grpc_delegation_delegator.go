package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The Delegator query returns the outstanding rewards and the total delegation amount of a
// delegator for its staker.
// If the delegator is not a staker both amounts will be zero.
// The request does not error.
func (k Keeper) Delegator(goCtx context.Context, req *types.QueryDelegatorRequest) (*types.QueryDelegatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryDelegatorResponse{}
	response.Delegator = &types.StakerDelegatorResponse{
		Delegator:        req.Delegator,
		CurrentRewards:   k.delegationKeeper.GetOutstandingRewards(ctx, req.Staker, req.Delegator),
		DelegationAmount: k.delegationKeeper.GetDelegationAmountOfDelegator(ctx, req.Staker, req.Delegator),
		Staker:           req.Staker,
	}

	return &response, nil
}
