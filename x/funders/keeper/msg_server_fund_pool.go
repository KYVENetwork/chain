package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// FundPool handles the logic to fund a pool.
// A funder is added to the active funders list with the specified amount
// If the funders list is full, it checks if the funder wants to fund
// more than the current lowest funder. If so, the current lowest funder
// will get their tokens back and removed form the active funders list.
func (k msgServer) FundPool(goCtx context.Context, msg *types.MsgFundPool) (*types.MsgFundPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Funder has to exist
	if !k.DoesFunderExist(ctx, msg.Creator) {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrFunderDoesNotExist.Error(), msg.Creator)
	}

	// Pool has to exist
	if err := k.poolKeeper.AssertPoolExists(ctx, msg.PoolId); err != nil {
		return nil, err
	}

	// Get funding state for pool
	fundingState, found := k.GetFundingState(ctx, msg.PoolId)
	if !found {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrFundingStateDoesNotExist.Error(), msg.PoolId)
	}

	amountsPerBundle := msg.AmountsPerBundle

	// Check if funding already exists
	funding, found := k.GetFunding(ctx, msg.Creator, msg.PoolId)
	if found {
		// If so, update funding amounts
		funding.Amounts = funding.Amounts.Add(msg.Amounts...)

		// Replace all coins in funding.AmountsPerBundle with the values of msg.AmountsPerBundle
		for _, coin := range funding.AmountsPerBundle {
			if f, _ := amountsPerBundle.Find(coin.Denom); !f {
				amountsPerBundle = amountsPerBundle.Add(coin)
			}
		}
		funding.AmountsPerBundle = amountsPerBundle
	} else {
		// If not, create new funding
		funding = types.Funding{
			FunderAddress:    msg.Creator,
			PoolId:           msg.PoolId,
			Amounts:          msg.Amounts,
			AmountsPerBundle: amountsPerBundle,
			TotalFunded:      sdk.NewCoins(),
		}
	}

	// Check if updated (or new) funding is compatible with module params
	if err := k.ensureParamsCompatibility(ctx, &funding); err != nil {
		return nil, err
	}

	// Kicks out lowest funder if all slots are taken and new funder is about to fund more.
	// Otherwise, an error is thrown
	// funding and fundingState are not written to the KV-Store. Everything else is handled safely.
	if err := k.ensureFreeSlot(ctx, &funding, &fundingState); err != nil {
		return nil, err
	}

	// All checks passed, transfer funds from funder to module
	sender := sdk.MustAccAddressFromBech32(msg.Creator)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.Amounts); err != nil {
		return nil, err
	}

	// Funding must be active
	fundingState.SetActive(&funding)

	// Save funding and funding state
	k.SetFunding(ctx, &funding)
	k.SetFundingState(ctx, &fundingState)

	// Emit a fund event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventFundPool{
		PoolId:           msg.PoolId,
		Address:          msg.Creator,
		Amounts:          msg.Amounts,
		AmountsPerBundle: msg.AmountsPerBundle,
	})

	return &types.MsgFundPoolResponse{}, nil
}
