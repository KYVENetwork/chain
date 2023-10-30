package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountFundedList(goCtx context.Context, req *types.QueryAccountFundedListRequest) (*types.QueryAccountFundedListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var funded []types.Funded

	fundings := k.fundersKeeper.GetFundingsOfFunder(ctx, req.Address)

	for _, funding := range fundings {
		if funding.Amount > 0 {
			pool, found := k.poolKeeper.GetPool(ctx, funding.PoolId)
			if !found {
				return nil, status.Error(codes.Internal, "pool not found")
			}
			funded = append(funded, types.Funded{
				Amount: funding.Amount,
				Pool: &types.BasicPool{
					Id:                   funding.PoolId,
					Name:                 pool.Name,
					Runtime:              pool.Runtime,
					Logo:                 pool.Logo,
					InflationShareWeight: pool.InflationShareWeight,
					UploadInterval:       pool.UploadInterval,
					TotalFunds:           k.fundersKeeper.GetTotalActiveFunding(ctx, pool.Id),
					TotalDelegation:      k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id),
					Status:               k.GetPoolStatus(ctx, &pool),
				},
			})
		}
	}

	return &types.QueryAccountFundedListResponse{
		Funded: funded,
	}, nil
}
