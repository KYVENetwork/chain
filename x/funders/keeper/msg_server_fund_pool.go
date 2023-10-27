package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensureFreeSlot makes sure that a funder can add funding to a given pool.
// If this is not possible an appropriate error is returned.
// A pool has a fixed amount of funding-slots. If there are still free slots
// a funder can just join (even with the smallest funding possible).
// If all slots are taken, it checks if the new funding has more funds
// than the current lowest funding in that pool.
// If so, the lowest funding gets removed from the pool, so that the
// new funding can be added.
// CONTRACT: no KV Writing on newFunding and fundingState
func (k Keeper) ensureFreeSlot(ctx sdk.Context, newFunding *types.Funding, fundingState *types.FundingState) error {

	activeFundings := k.GetActiveFundings(ctx, *fundingState)
	// check if slots are still available
	if len(activeFundings) < types.MaxFunders {
		return nil
	}

	lowestFunding, _ := k.GetLowestFunding(activeFundings)

	if lowestFunding.FunderAddress == newFunding.FunderAddress {
		// Funder already has a funding slot
		return nil
	}

	// Check if lowest funding is lower than new funding based on amount (amount per bundle is ignored)
	if newFunding.Amount < lowestFunding.Amount {
		return errors.Wrapf(errorsTypes.ErrLogic, types.ErrFundsTooLow.Error(), lowestFunding.Amount)
	}

	// Defund lowest funder
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, lowestFunding.FunderAddress, lowestFunding.Amount); err != nil {
		return err
	}

	subtracted := lowestFunding.SubtractAmount(lowestFunding.Amount)
	fundingState.SetInactive(lowestFunding)
	k.SetFunding(ctx, lowestFunding)

	// Emit a defund event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
		PoolId:  fundingState.PoolId,
		Address: lowestFunding.FunderAddress,
		Amount:  subtracted,
	})

	return nil
}

// FundPool handles the logic to fund a pool.
// A funder is added to the active funders list with the specified amount
// If the funders list is full, it checks if the funder wants to fund
// more than the current lowest funder. If so, the current lowest funder
// will get their tokens back and removed form the active funders list.
// TODO: what if amount_per_bundle is higher than the amount? A funder that knows that he is the next uploader could just fund a huge amount which gets payed to only himself.
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

	// Check if funding already exists
	funding, found := k.GetFunding(ctx, msg.Creator, msg.PoolId)
	if found {
		// If so, update funding amount
		funding.AddAmount(msg.Amount)
		// If the amount per bundle is set, update it
		if msg.AmountPerBundle > 0 {
			funding.AmountPerBundle = msg.AmountPerBundle
		}
	} else {
		// If not, create new funding
		funding = types.Funding{
			FunderAddress:   msg.Creator,
			PoolId:          msg.PoolId,
			Amount:          msg.Amount,
			AmountPerBundle: msg.AmountPerBundle,
			TotalFunded:     0,
		}
	}

	// Check compatibility of updated funding with params
	// i.e minimum funding per bundle
	//     and minimum funding amount
	params := k.GetParams(ctx)
	if funding.AmountPerBundle < params.MinFundingAmountPerBundle {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrAmountPerBundleTooLow.Error(), params.MinFundingAmountPerBundle)
	}
	if funding.Amount < params.MinFundingAmount {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrMinFundingAmount.Error(), params.MinFundingAmount)
	}

	// Kicks out lowest funder if all slots are taken and new funder is about to fund more.
	// Otherwise, an error is thrown
	// funding and fundingState are not written to the KV-Store. Everything else is handled safely.
	if err := k.ensureFreeSlot(ctx, &funding, &fundingState); err != nil {
		return nil, err
	}

	// User is allowed to fund
	// Let's see if he has enough funds
	if err := util.TransferFromAddressToModule(k.bankKeeper, ctx, msg.Creator, types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	// Funding must be active
	fundingState.SetActive(&funding)

	// Save funding and funding state
	k.SetFunding(ctx, &funding)
	k.SetFundingState(ctx, &fundingState)

	// Emit a fund event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventFundPool{
		PoolId:          msg.PoolId,
		Address:         msg.Creator,
		Amount:          msg.Amount,
		AmountPerBundle: msg.AmountPerBundle,
	})

	return &types.MsgFundPoolResponse{}, nil
}
