package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StakeInfo returns the current staking information for a given user and pool.
// It is used by the protocol node to determine whether a node is able to participate
// and to adjust the initial stake.
func (k Keeper) StakeInfo(goCtx context.Context, req *types.QueryStakeInfoRequest) (*types.QueryStakeInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Create empty response, as the query shall return the default values if a specific entry does not exist
	response := types.QueryStakeInfoResponse{
		Balance:      "0",
		CurrentStake: "0",
		MinimumStake: "0",
	}

	// Query balance
	account, addressError := sdk.AccAddressFromBech32(req.Staker)
	if addressError != nil {
		return nil, sdkErrors.ErrNotFound
	}

	coin := k.bankKeeper.GetBalance(ctx, account, "tkyve")
	response.Balance = strings.TrimSuffix(coin.String(), coin.Denom)

	// Current Stake amount
	staker, exists := k.GetStaker(ctx, req.Staker, req.PoolId)
	if exists {
		response.CurrentStake = strconv.FormatUint(staker.Amount, 10)
	}

	response.Status = staker.Status

	// Fetch pool
	pool, exists := k.GetPool(ctx, req.PoolId)
	if !exists {
		return nil, sdkErrors.ErrNotFound
	}

	// Fetch current lowest staker only if all stacker slots are occupied
	if len(pool.Stakers) >= types.MaxStakers {
		lowestStaker, _ := k.GetStaker(ctx, pool.LowestStaker, req.PoolId)
		response.MinimumStake = strconv.FormatUint(lowestStaker.Amount, 10)
	}

	return &response, nil
}
