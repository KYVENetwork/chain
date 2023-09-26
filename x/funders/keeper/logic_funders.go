package keeper

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// ChargeFundersOfPool equally splits the amount between all funders and removes
// the appropriate amount from each funder.
// All funders who can't afford the amount, are kicked out.
// Their remaining amount is transferred to the Treasury.
// This method does not transfer any funds. The bundles-module
// is responsible for transferring the rewards out of the module.
// TODO: update text
func (k Keeper) ChargeFundersOfPool(ctx sdk.Context, poolId uint64) (payout uint64, err error) {
	// Get funding state for pool
	fundingState, found := k.GetFundingState(ctx, poolId)
	if !found {
		return 0, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrFundingStateDoesNotExist.Error(), poolId)
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
			fundingState.SetInactive(funding)
		}
		k.setFunding(ctx, funding)
	}

	// Save funding state
	k.setFundingState(ctx, fundingState)

	// Emit a pool out of funds event if there are no more active funders
	if len(fundingState.ActiveFunderAddresses) == 0 {
		_ = ctx.EventManager().EmitTypedEvent(&types.EventPoolOutOfFunds{
			PoolId: poolId,
		})
	}
	return payout, nil
}
