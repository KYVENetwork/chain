package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StakersList returns a list of all validators for a given pool with their current stake,
// total delegation and additional information (moniker, website, etc.)
// This query is not paginated as it contains a maximum of types.MAX_STAKERS entries
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

	for _, account := range pool.Stakers {
		staker, _ := k.GetStaker(ctx, account, req.PoolId)

		stakerResponse := types.StakerResponse{
			Staker:          staker.Account,
			PoolId:          staker.PoolId,
			Account:         staker.Account,
			Amount:          staker.Amount,
			TotalDelegation: 0,
			Commission:      staker.Commission,
			Moniker:         staker.Moniker,
			Website:         staker.Website,
			Logo:            staker.Logo,
			Points:          staker.Points,
		}

		// Fetch total delegation for staker, as it is stored in DelegationPoolData
		poolDelegationData, _ := k.GetDelegationPoolData(ctx, staker.PoolId, staker.Account)
		stakerResponse.TotalDelegation = poolDelegationData.TotalDelegation

		response.Stakers = append(response.Stakers, &stakerResponse)
	}

	return &response, nil
}
