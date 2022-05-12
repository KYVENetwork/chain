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

	if !found {
		return nil, sdkErrors.ErrNotFound
	}

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
	}

	// Fetch total delegation for staker, as it is stored in DelegationPoolData
	poolDelegationData, _ := k.GetDelegationPoolData(ctx, staker.PoolId, staker.Account)
	stakerResponse.TotalDelegation = poolDelegationData.TotalDelegation

	response.Staker = &stakerResponse

	return &response, nil
}
