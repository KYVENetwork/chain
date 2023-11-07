package keeper

import (
	"fmt"

	"cosmossdk.io/errors"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CreateFundingState(ctx sdk.Context, poolId uint64) {
	fundingState := types.FundingState{
		PoolId:                poolId,
		ActiveFunderAddresses: []string{},
	}
	k.SetFundingState(ctx, &fundingState)
}

func (k Keeper) GetTotalActiveFunding(ctx sdk.Context, poolId uint64) (amount uint64) {
	state, found := k.GetFundingState(ctx, poolId)
	if !found {
		return 0
	}
	for _, address := range state.ActiveFunderAddresses {
		funding, _ := k.GetFunding(ctx, address, poolId)
		amount += funding.Amount
	}
	return amount
}

// ChargeFundersOfPool charges all funders of a pool with their amount_per_bundle
// If the amount is lower than the amount_per_bundle,
// the max amount is charged and the funder is removed from the active funders list.
// The amount is transferred from the funders to the pool module account where it can be paid out.
// If there are no more active funders, an event is emitted.
func (k Keeper) ChargeFundersOfPool(ctx sdk.Context, poolId uint64) (payout uint64, err error) {
	// Get funding state for pool
	fundingState, found := k.GetFundingState(ctx, poolId)
	if !found {
		return 0, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrFundingStateDoesNotExist.Error(), poolId)
	}

	// If there are no active fundings we immediately return
	activeFundings := k.GetActiveFundings(ctx, fundingState)
	if len(activeFundings) == 0 {
		return 0, nil
	}

	// This is the amount every funding will be charged
	for _, funding := range activeFundings {
		payout += funding.ChargeOneBundle()
		if funding.Amount == 0 {
			fundingState.SetInactive(&funding)
		}
		k.SetFunding(ctx, &funding)
	}

	// Save funding state
	k.SetFundingState(ctx, &fundingState)

	// Emit a pool out of funds event if there are no more active funders
	if len(fundingState.ActiveFunderAddresses) == 0 {
		_ = ctx.EventManager().EmitTypedEvent(&types.EventPoolOutOfFunds{
			PoolId: poolId,
		})
	}

	// Move funds to pool module account
	if payout > 0 {
		err = util.TransferFromModuleToModule(k.bankKeeper, ctx, types.ModuleName, pooltypes.ModuleName, payout)
		if err != nil {
			return 0, err
		}
	}

	return payout, nil
}

// GetLowestFunding returns the funding with the lowest amount
// Precondition: len(fundings) > 0
func (k Keeper) GetLowestFunding(fundings []types.Funding) (lowestFunding *types.Funding, err error) {
	if len(fundings) == 0 {
		return nil, fmt.Errorf("no active fundings")
	}

	lowestFundingIndex := 0
	for i := range fundings {
		if fundings[i].Amount < fundings[lowestFundingIndex].Amount {
			lowestFundingIndex = i
		}
	}
	return &fundings[lowestFundingIndex], nil
}

// ensureParamsCompatibility checks compatibility of the provided funding with the pool params.
// i.e.
// - minimum funding per bundle
// - minimum funding amount
// - minimum funding multiple
func (k Keeper) ensureParamsCompatibility(ctx sdk.Context, funding types.Funding) error {
	params := k.GetParams(ctx)
	if funding.AmountPerBundle < params.MinFundingAmountPerBundle {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrAmountPerBundleTooLow.Error(), params.MinFundingAmountPerBundle)
	}
	if funding.Amount < params.MinFundingAmount {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrMinFundingAmount.Error(), params.MinFundingAmount)
	}
	if funding.AmountPerBundle*params.MinFundingMultiple > funding.Amount {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrMinFundingMultiple.Error(), funding.AmountPerBundle, params.MinFundingAmount, funding.Amount)
	}
	return nil
}

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
