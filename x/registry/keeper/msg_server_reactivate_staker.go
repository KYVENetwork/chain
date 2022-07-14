package keeper

import (
	"context"

	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReactivateStaker ...
func (k msgServer) ReactivateStaker(
	goCtx context.Context, msg *types.MsgReactivateStaker,
) (*types.MsgReactivateStakerResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetPool(ctx, msg.PoolId)
	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.PoolId)
	}

	staker, stakerFound := k.GetStaker(ctx, msg.Creator, msg.PoolId)
	if !stakerFound {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrNoStaker.Error())
	}

	if staker.Status != types.STAKER_STATUS_INACTIVE {
		// TODO custom error
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrNoStaker.Error())
	}

	if len(pool.Stakers) >= types.MaxStakers {
		lowestStaker, _ := k.GetStaker(ctx, pool.LowestStaker, msg.PoolId)

		if staker.Amount > lowestStaker.Amount {

			if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventStakerStatusChanged{
				PoolId:  pool.Id,
				Address: lowestStaker.Account,
				Status:  types.STAKER_STATUS_INACTIVE,
			}); errEmit != nil {
				return nil, errEmit
			}

			// Move the lowest staker to inactive staker set
			deactivateStaker(&pool, &lowestStaker)
			k.SetStaker(ctx, lowestStaker)

		} else {
			return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrStakeTooLow.Error(), lowestStaker.Amount)
		}
	}

	pool.InactiveStakers = removeStringFromList(pool.InactiveStakers, staker.Account)
	pool.Stakers = append(pool.Stakers, staker.Account)
	pool.TotalStake += staker.Amount
	pool.TotalInactiveStake -= staker.Amount
	staker.Status = types.STAKER_STATUS_ACTIVE

	k.SetStaker(ctx, staker)
	k.updateLowestStaker(ctx, &pool)
	k.SetPool(ctx, pool)

	// Emit a delegation event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventStakerStatusChanged{
		PoolId:  msg.PoolId,
		Address: msg.Creator,
		Status:  types.STAKER_STATUS_ACTIVE,
	}); errEmit != nil {
		return nil, errEmit
	}

	return &types.MsgReactivateStakerResponse{}, nil
}
