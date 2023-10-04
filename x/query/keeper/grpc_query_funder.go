package keeper

import (
	"context"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KYVENetwork/chain/x/query/types"
)

func (k Keeper) Funders(c context.Context, req *types.QueryFundersRequest) (*types.QueryFundersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	funders, pageRes, err := k.fundersKeeper.GetPaginatedFundersQuery(ctx, req.Pagination, req.Search)
	if err != nil {
		return nil, err
	}

	data := make([]types.Funder, 0)
	for _, funder := range funders {
		data = append(data, k.parseFunder(&funder))
	}

	return &types.QueryFundersResponse{Funders: data, Pagination: pageRes}, nil
}

func (k Keeper) Funder(c context.Context, req *types.QueryFunderRequest) (*types.QueryFunderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	funder, found := k.fundersKeeper.GetFunder(ctx, req.Address)
	if !found {
		return nil, errorsTypes.ErrKeyNotFound
	}
	funderData := k.parseFunder(&funder)

	totalUsedFunds := uint64(0)
	totalAllocatedFunds := uint64(0)
	poolsFunded := make([]uint64, 0)

	fundings := k.fundersKeeper.GetFundingsOfFunder(ctx, req.Address)
	fundingsData := make([]types.Funding, 0)
	for _, funding := range fundings {
		if funding.Amount > 0 || req.WithInactiveFundings {
			fundingsData = append(fundingsData, types.Funding{
				PoolId:          funding.PoolId,
				Amount:          funding.Amount,
				AmountPerBundle: funding.AmountPerBundle,
				TotalFunded:     funding.TotalFunded,
			})
		}
		totalUsedFunds += funding.TotalFunded
		totalAllocatedFunds += funding.Amount
		poolsFunded = append(poolsFunded, funding.PoolId)
	}

	statsData := &types.FunderStats{
		TotalUsedFunds:      totalUsedFunds,
		TotalAllocatedFunds: totalAllocatedFunds,
		PoolsFunded:         poolsFunded,
	}

	return &types.QueryFunderResponse{
		Funder:   &funderData,
		Fundings: fundingsData,
		Stats:    statsData,
	}, nil
}

func (k Keeper) parseFunder(funder *fundersTypes.Funder) types.Funder {
	return types.Funder{
		Address:     funder.Address,
		Moniker:     funder.Moniker,
		Identity:    funder.Identity,
		Logo:        funder.Logo,
		Website:     funder.Website,
		Contact:     funder.Contact,
		Description: funder.Description,
	}
}
