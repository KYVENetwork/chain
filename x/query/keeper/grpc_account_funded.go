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

	for i := range pools {
		pool := pools[i]
		// TODO(rapha): fix this
		//amount := pool.GetFunderAmount(req.Address)
		amount := uint64(0)

		if amount > 0 {
			funded = append(funded, types.Funded{
				Amount: amount,
				Pool: &types.BasicPool{
					Id:             pool.Id,
					Name:           pool.Name,
					Runtime:        pool.Runtime,
					Logo:           pool.Logo,
					OperatingCost:  pool.OperatingCost,
					UploadInterval: pool.UploadInterval,
					// TODO(rapha): fix this
					TotalFunds: 0,
					//TotalFunds:      pool.TotalFunds,
					TotalDelegation: k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id),
					Status:          k.GetPoolStatus(ctx, &pool),
				},
			})
		}
	}

	return &types.QueryAccountFundedListResponse{
		Funded: funded,
	}, nil
}
