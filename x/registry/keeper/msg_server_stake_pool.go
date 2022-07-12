package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// StakePool handles the logic of an SDK message that allows protocol nodes to stake in a specified pool.
func (k msgServer) StakePool(goCtx context.Context, msg *types.MsgStakePool) (*types.MsgStakePoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is already a staker.
	staker, stakerExists := k.GetStaker(ctx, msg.Creator, msg.Id)

	if stakerExists {
		staker.Amount += msg.Amount
		k.SetStaker(ctx, staker)
	} else {
		// Check if we have reached the maximum number of stakers.
		// If we are staking more than the lowest staker, remove them.
		if len(pool.Stakers) == types.MaxStakers {
			lowestStaker, _ := k.GetStaker(ctx, pool.LowestStaker, msg.Id)

			if msg.Amount > lowestStaker.Amount {

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

		pool.Stakers = append(pool.Stakers, msg.Creator)

		k.SetStaker(ctx, types.Staker{
			Account:    msg.Creator,
			PoolId:     msg.Id,
			Amount:     msg.Amount,
			Commission: types.DefaultCommission,
			Status:     types.STAKER_STATUS_ACTIVE,
		})
	}

	// Transfer tokens from sender to this module.
	err := k.transferToRegistry(ctx, msg.Creator, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Event a stake event.
	errEmit := ctx.EventManager().EmitTypedEvent(&types.EventStakePool{
		PoolId:  msg.Id,
		Address: msg.Creator,
		Amount:  msg.Amount,
	})
	if errEmit != nil {
		return nil, errEmit
	}

	staker, _ = k.GetStaker(ctx, msg.Creator, msg.Id)
	if staker.Status == types.STAKER_STATUS_ACTIVE {
		pool.TotalStake += msg.Amount
	} else if staker.Status == types.STAKER_STATUS_INACTIVE {
		pool.TotalInactiveStake += msg.Amount
	}

	k.updateLowestStaker(ctx, &pool)
	k.SetPool(ctx, pool)

	return &types.MsgStakePoolResponse{}, nil
}
