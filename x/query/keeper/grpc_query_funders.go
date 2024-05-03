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

	params := k.fundersKeeper.GetParams(ctx)

	data := make([]types.Funder, 0)
	for _, funder := range funders {
		fundings := k.fundersKeeper.GetFundingsOfFunder(ctx, funder.Address)
		data = append(data, k.parseFunder(&funder, fundings, params.CoinWhitelist))
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

	allFundings := k.fundersKeeper.GetFundingsOfFunder(ctx, funder.Address)
	fundings := k.filterFundingsOnStatus(allFundings, req.Status)

	params := k.fundersKeeper.GetParams(ctx)
	funderData := k.parseFunder(&funder, allFundings, params.CoinWhitelist)
	fundingsData := k.parseFundings(fundings, params.CoinWhitelist)

	return &types.QueryFunderResponse{
		Funder:   &funderData,
		Fundings: fundingsData,
	}, nil
}

func (k Keeper) filterFundingsOnStatus(fundings []fundersTypes.Funding, fundingStatus types.FundingStatus) []fundersTypes.Funding {
	if fundingStatus == types.FUNDING_STATUS_UNSPECIFIED {
		return fundings
	}

	filtered := make([]fundersTypes.Funding, 0)
	for _, funding := range fundings {
		if fundingStatus == types.FUNDING_STATUS_ACTIVE && funding.Amounts.IsAllPositive() {
			filtered = append(filtered, funding)
		}
		if fundingStatus == types.FUNDING_STATUS_INACTIVE && funding.Amounts.IsZero() {
			filtered = append(filtered, funding)
		}
	}
	return filtered
}

func (k Keeper) parseFunder(funder *fundersTypes.Funder, fundings []fundersTypes.Funding, whitelist []*fundersTypes.WhitelistCoinEntry) types.Funder {
	stats := types.FundingStats{
		TotalUsedFunds:       sdk.NewCoins(),
		TotalAllocatedFunds:  sdk.NewCoins(),
		TotalAmountPerBundle: sdk.NewCoins(),
		PoolsFunded:          make([]uint64, 0),
		Score:                uint64(0),
	}

	for _, funding := range fundings {
		// Only count active fundings for totalAmountPerBundle
		if funding.IsActive() {
			stats.TotalAmountPerBundle = stats.TotalAmountPerBundle.Add(funding.AmountsPerBundle...)
		}

		stats.TotalUsedFunds = stats.TotalUsedFunds.Add(funding.TotalFunded...)
		stats.TotalAllocatedFunds = stats.TotalAllocatedFunds.Add(funding.Amounts...)
		stats.Score += funding.GetScore(whitelist)

		stats.PoolsFunded = append(stats.PoolsFunded, funding.PoolId)
	}

	return types.Funder{
		Address:     funder.Address,
		Moniker:     funder.Moniker,
		Identity:    funder.Identity,
		Website:     funder.Website,
		Contact:     funder.Contact,
		Description: funder.Description,
		Stats:       &stats,
	}
}
