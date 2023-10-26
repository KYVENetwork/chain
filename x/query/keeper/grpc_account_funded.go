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

	pools := k.poolKeeper.GetAllPools(ctx)

	for _, pool := range pools {
		funding, found := k.fundersKeeper.GetFunding(ctx, req.Address, pool.Id)
		if !found {
			return nil, status.Error(codes.Internal, "funding not found")
		}
		fundingState, found := k.fundersKeeper.GetFundingState(ctx, pool.Id)
		if !found {
			return nil, status.Error(codes.Internal, "funding state not found")
		}

		if funding.Amount > 0 {
			funded = append(funded, types.Funded{
				Amount: funding.Amount,
				Pool: &types.BasicPool{
					Id:                   pool.Id,
					Name:                 pool.Name,
					Runtime:              pool.Runtime,
					Logo:                 pool.Logo,
					InflationShareWeight: pool.InflationShareWeight,
					UploadInterval:       pool.UploadInterval,
					TotalFunds:           k.fundersKeeper.GetTotalActiveFunding(ctx, pool.Id),
					TotalDelegation:      k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id),
					Status:               k.GetPoolStatus(ctx, &pool, &fundingState),
				},
			})
		}
	}

	return &types.QueryAccountFundedListResponse{
		Funded: funded,
	}, nil
}
