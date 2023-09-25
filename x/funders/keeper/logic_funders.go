package keeper

import (
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ChargeFundersOfPool equally splits the amount between all funders and removes
// the appropriate amount from each funder.
// All funders who can't afford the amount, are kicked out.
// Their remaining amount is transferred to the Treasury.
// This method does not transfer any funds. The bundles-module
// is responsible for transferring the rewards out of the module.
func (k Keeper) ChargeFundersOfPool(ctx sdk.Context, poolId uint64, amount uint64) (payout uint64, err error) {
	pool, poolErr := k.poolKeeper.GetPoolWithError(ctx, poolId)
	if poolErr != nil {
		return 0, poolErr
	}

	// if pool has no funders we immediately return
	if len(pool.Funders) == 0 {
		return payout, err
	}

	// This is the amount every funder will be charged
	amountPerFunder := amount / uint64(len(pool.Funders))

	// Due to discrete division there will be a remainder which can not be split
	// equally among all funders. This amount is charged to the lowest funder
	amountRemainder := amount - amountPerFunder*uint64(len(pool.Funders))

	funders := pool.Funders

	for _, funder := range funders {
		if funder.Amount < amountPerFunder {
			pool.RemoveFunder(funder.Address)
			payout += funder.Amount
		} else {
			pool.SubtractAmountFromFunder(funder.Address, amountPerFunder)
			payout += amountPerFunder
		}
	}

	lowestFunder := pool.GetLowestFunder()

	if lowestFunder.Address != "" {
		if lowestFunder.Amount < amountRemainder {
			pool.RemoveFunder(lowestFunder.Address)
			payout += lowestFunder.Amount
		} else {
			pool.SubtractAmountFromFunder(lowestFunder.Address, amountRemainder)
			payout += amountRemainder
		}
	}

	if len(pool.Funders) == 0 {
		_ = ctx.EventManager().EmitTypedEvent(&poolTypes.EventPoolOutOfFunds{
			PoolId: pool.Id,
		})
	}

	k.SetPool(ctx, pool)
	return payout, nil
}
