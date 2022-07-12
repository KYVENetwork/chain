package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Staker returns all staker info
func (k Keeper) Staker(goCtx context.Context, req *types.QueryStakerRequest) (*types.QueryStakerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryStakerResponse{}

	// Load pool
	_, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	// Load staker
	staker, found := k.GetStaker(ctx, req.Staker, req.PoolId)

	// Load unbondingStaker
	unbondingStaker, _ := k.GetUnbondingStaker(ctx, req.PoolId, req.Staker)

	if !found {
		return nil, sdkErrors.ErrNotFound
	}

	stakerResponse := types.StakerResponse{
		Staker:                  staker.Account,
		PoolId:                  staker.PoolId,
		Account:                 staker.Account,
		Amount:                  staker.Amount,
		TotalDelegation:         0,
		Commission:              staker.Commission,
		Moniker:                 staker.Moniker,
		Website:                 staker.Website,
		Logo:                    staker.Logo,
		Points:                  staker.Points,
		UnbondingAmount:         unbondingStaker.UnbondingAmount,
		UploadProbability:       "0",
		Status:                  staker.Status,
		PendingCommissionChange: nil,
	}

	if staker.Status == types.STAKER_STATUS_ACTIVE {
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

	// Fetch total delegation for staker, as it is stored in DelegationPoolData
	poolDelegationData, _ := k.GetDelegationPoolData(ctx, staker.PoolId, staker.Account)
	stakerResponse.TotalDelegation = poolDelegationData.TotalDelegation

	response.Staker = &stakerResponse

	return &response, nil
}
