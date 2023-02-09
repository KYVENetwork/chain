package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// FundPool handles the logic to fund a pool.
// A funder is added to the funders list with the specified amount
// If the funders list is full, it checks if the funder wants to fund
// more than the current lowest funder. If so, the current lowest funder
// will get their tokens back and removed form the funders list.
func (k msgServer) FundPool(goCtx context.Context, msg *types.MsgFundPool) (*types.MsgFundPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, poolFound := k.GetPool(ctx, msg.Id)

	if !poolFound {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if funder already exists
	// If sender is not a funder, check if a free funding slot is still available
	if pool.GetFunderAmount(msg.Creator) == 0 {
		// If funder does not exist, check if limit is already exceeded.
		if len(pool.Funders) >= types.MaxFunders {
			// If so, check if funder wants to fund more than current lowest funder.
			lowestFunder := pool.GetLowestFunder()
			if msg.Amount > lowestFunder.Amount {
				// Unstake lowest Funder
				err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, lowestFunder.Address, lowestFunder.Amount)
				if err != nil {
					return nil, err
				}

				// Emit a defund event.
				_ = ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
					PoolId:  msg.Id,
					Address: lowestFunder.Address,
					Amount:  lowestFunder.Amount,
				})

				// Remove from pool
				pool.RemoveFunder(lowestFunder.Address)
			} else {
				return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrFundsTooLow.Error(), lowestFunder.Amount)
			}
		}
	}

	// User is allowed to fund
	pool.AddAmountToFunder(msg.Creator, msg.Amount)

	if err := util.TransferFromAddressToModule(k.bankKeeper, ctx, msg.Creator, types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventFundPool{
		PoolId:  msg.Id,
		Address: msg.Creator,
		Amount:  msg.Amount,
	})

	k.SetPool(ctx, pool)

	return &types.MsgFundPoolResponse{}, nil
}
