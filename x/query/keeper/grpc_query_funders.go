package keeper

import (
	"context"

	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		fundings := k.fundersKeeper.GetFundingsOfFunder(ctx, funder.Address)
		data = append(data, k.parseFunder(&funder, fundings))
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
		return nil, status.Error(codes.NotFound, "funder not found")
	}

	fundings := k.fundersKeeper.GetFundingsOfFunder(ctx, funder.Address)

	funderData := k.parseFunder(&funder, fundings)
	fundingsData := k.parseFundings(fundings, req.WithInactiveFundings)

	return &types.QueryFunderResponse{
		Funder:   &funderData,
		Fundings: fundingsData,
	}, nil
}

func (k Keeper) parseFunder(funder *fundersTypes.Funder, fundings []fundersTypes.Funding) types.Funder {
	totalUsedFunds := uint64(0)
	totalAllocatedFunds := uint64(0)
	totalAmountPerBundle := uint64(0)
	poolsFunded := make([]uint64, 0)

	for _, funding := range fundings {
		// Only count active fundings for totalAmountPerBundle
		if funding.Amount > 0 {
			totalAmountPerBundle += funding.AmountPerBundle
		}

		totalUsedFunds += funding.TotalFunded
		totalAllocatedFunds += funding.Amount

		poolsFunded = append(poolsFunded, funding.PoolId)
	}

	return types.Funder{
		Address:     funder.Address,
		Moniker:     funder.Moniker,
		Identity:    funder.Identity,
		Website:     funder.Website,
		Contact:     funder.Contact,
		Description: funder.Description,
		Stats: &types.FundingStats{
			TotalUsedFunds:       totalUsedFunds,
			TotalAllocatedFunds:  totalAllocatedFunds,
			TotalAmountPerBundle: totalAmountPerBundle,
			PoolsFunded:          poolsFunded,
		},
	}
}
