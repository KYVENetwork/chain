package keeper

import (
	"context"
	"sort"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakersByPool(c context.Context, req *types.QueryStakersByPoolRequest) (*types.QueryStakersByPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if _, found := k.poolKeeper.GetPool(ctx, req.PoolId); !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	stakers := make([]types.FullStaker, 0)

	valaccounts := k.stakerKeeper.GetAllValaccountsOfPool(ctx, req.PoolId)
	for _, valaccount := range valaccounts {
		stakers = append(stakers, *k.GetFullStaker(ctx, valaccount.Staker))
	}

	stakes := k.stakerKeeper.GetValidatorPoolStakes(ctx, req.PoolId)

	sort.SliceStable(stakers, func(i, j int) bool {
		return stakes[stakers[i].Address] > stakes[stakers[j].Address]
	})

	return &types.QueryStakersByPoolResponse{Stakers: stakers}, nil
}
