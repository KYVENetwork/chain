package keeper

import (
	goerrors "errors"
	"fmt"

	"cosmossdk.io/errors"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: should this be here or when we call the getter of the funding state?
func (k Keeper) CreateFundingState(ctx sdk.Context, poolId uint64) {
	fundingState := types.FundingState{
		PoolId:                poolId,
		ActiveFunderAddresses: []string{},
		TotalAmount:           0,
	}
	k.SetFundingState(ctx, &fundingState)
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
	payout = 0
	for _, funding := range activeFundings {
		payout += funding.ChargeOneBundle()
		if funding.Amount == 0 {
			fundingState.SetInactive(&funding)
		}
		k.SetFunding(ctx, &funding)
	}
	fundingState.SubtractAmount(payout)

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
		return nil, goerrors.New(fmt.Sprintf("no active fundings"))
	}

	lowestFunding = &fundings[0]
	for _, funding := range fundings {
		if funding.Amount < lowestFunding.Amount {
			lowestFunding = &funding
		}
	}
	return lowestFunding, nil
}
