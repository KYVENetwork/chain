package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StakersList returns a list of all validators for a given pool with their current stake,
// total delegation and additional information (moniker, website, etc.)
func (k Keeper) StakersList(goCtx context.Context, req *types.QueryStakersListRequest) (*types.QueryStakersListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryStakersListResponse{}

	// Load pool
	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	var stakers []string

	// Filter by Staker status
	switch req.Status {
	case types.STAKER_STATUS_UNSPECIFIED:
		return &types.QueryStakersListResponse{}, nil
	case types.STAKER_STATUS_ACTIVE:
		stakers = pool.Stakers
	case types.STAKER_STATUS_INACTIVE:
		stakers = pool.InactiveStakers
	}

	// Pagination
	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: 50,
		}
	}

	start := req.Pagination.Offset
	end := req.Pagination.Offset + req.Pagination.Limit

	if start >= uint64(len(stakers)) {
		return &types.QueryStakersListResponse{}, nil
	} else if end > uint64(len(stakers)) {
		end = uint64(len(stakers))
	}

	// Iterate paginated stakers and collect stats
	for _, account := range stakers[start:end] {
		staker, _ := k.GetStaker(ctx, account, req.PoolId)
		// Load unbondingStaker
		unbondingStaker, _ := k.GetUnbondingStaker(ctx, req.PoolId, account)

		// Fetch total delegation for staker, as it is stored in DelegationPoolData
		poolDelegationData, _ := k.GetDelegationPoolData(ctx, staker.PoolId, staker.Account)

		stakerResponse := types.StakerResponse{
			Staker:            staker.Account,
			PoolId:            staker.PoolId,
			Account:           staker.Account,
			Amount:            staker.Amount,
			TotalDelegation:   poolDelegationData.TotalDelegation,
			Commission:        staker.Commission,
			Moniker:           staker.Moniker,
			Website:           staker.Website,
			Logo:              staker.Logo,
			Points:            staker.Points,
			UnbondingAmount:   unbondingStaker.UnbondingAmount,
			UploadProbability: "0",
			Status:            staker.Status,
		}

		// Only active stakers have an upload probability
		if req.Status == types.STAKER_STATUS_ACTIVE {
			stakerResponse.UploadProbability = k.GetUploadProbability(ctx, staker.Account, staker.PoolId).String()
		}

		commissionChangeEntry, foundCommissionChange := k.GetCommissionChangeQueueEntryByIndex2(ctx, staker.Account, staker.PoolId)
		if foundCommissionChange {
			stakerResponse.PendingCommissionChange = &types.PendingCommissionChange{
				NewCommission: commissionChangeEntry.Commission,
				CreationDate:  commissionChangeEntry.CreationDate,
				FinishDate:    commissionChangeEntry.CreationDate + int64(k.CommissionChangeTime(ctx)),
			}
		}

		response.Stakers = append(response.Stakers, &stakerResponse)
	}

	// Pagination results
	response.Pagination = &query.PageResponse{
		NextKey: nil,
		Total:   end - start,
	}

	return &response, nil
}
