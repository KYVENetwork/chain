package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CanValidate(c context.Context, req *types.QueryCanValidateRequest) (*types.QueryCanValidateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if _, err := k.poolKeeper.GetPoolWithError(ctx, req.PoolId); err != nil {
		return &types.QueryCanValidateResponse{
			Possible: false,
			Reason:   err.Error(),
		}, nil
	}

	var staker string

	// Check if valaddress has a valaccount in pool
	for _, poolAccount := range k.stakerKeeper.GetAllPoolAccountsOfPool(ctx, req.PoolId) {
		if poolAccount.PoolAddress == req.PoolAddress {
			staker = poolAccount.Staker
			break
		}
	}

	if staker == "" {
		return &types.QueryCanValidateResponse{
			Possible: false,
			Reason:   "no valaccount found",
		}, nil
	}

	return &types.QueryCanValidateResponse{
		Possible: true,
		Reason:   staker,
	}, nil
}
