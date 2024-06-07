package keeper

import (
	"context"

	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) FundingsByFunder(c context.Context, req *types.QueryFundingsByFunderRequest) (*types.QueryFundingsByFunderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	fundings, pageRes, err := k.fundersKeeper.GetPaginatedFundingQuery(ctx, req.Pagination, &req.Address, nil, req.Status)
	if err != nil {
		return nil, err
	}

	whitelist := k.fundersKeeper.GetCoinWhitelistMap(ctx)
	data := k.parseFundings(fundings, whitelist)

	return &types.QueryFundingsByFunderResponse{Fundings: data, Pagination: pageRes}, nil
}

func (k Keeper) FundingsByPool(c context.Context, req *types.QueryFundingsByPoolRequest) (*types.QueryFundingsByPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	fundings, pageRes, err := k.fundersKeeper.GetPaginatedFundingQuery(ctx, req.Pagination, nil, &req.PoolId, req.Status)
	if err != nil {
		return nil, err
	}

	whitelist := k.fundersKeeper.GetCoinWhitelistMap(ctx)
	data := k.parseFundings(fundings, whitelist)

	return &types.QueryFundingsByPoolResponse{Fundings: data, Pagination: pageRes}, nil
}

func (k Keeper) parseFundings(fundings []fundersTypes.Funding, whitelist map[string]fundersTypes.WhitelistCoinEntry) []types.Funding {
	fundingsData := make([]types.Funding, 0)
	for _, funding := range fundings {
		fundingsData = append(fundingsData, types.Funding{
			FunderAddress:    funding.FunderAddress,
			PoolId:           funding.PoolId,
			Amounts:          funding.Amounts,
			AmountsPerBundle: funding.AmountsPerBundle,
			TotalFunded:      funding.TotalFunded,
			Score:            funding.GetScore(whitelist),
		})
	}
	return fundingsData
}
