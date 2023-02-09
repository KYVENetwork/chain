package keeper

import (
	"github.com/KYVENetwork/chain/util"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ChargeFundersOfPool equally splits the amount between all funders and removes
// the appropriate amount from each funder.
// All funders who can't afford the amount, are kicked out.
// Their remaining amount is transferred to the Treasury.
// The function throws an error if pool ran out of funds.
// This method does not transfer any funds. The bundles-module
// is responsible for transferring the rewards out of the module.
func (k Keeper) ChargeFundersOfPool(ctx sdk.Context, poolId uint64, amount uint64) error {
	pool, poolErr := k.GetPoolWithError(ctx, poolId)
	if poolErr != nil {
		return poolErr
	}

	// This is the amount every funder will be charged
	var amountPerFunder uint64

	// Due to discrete division there will be a reminder which can not be split
	// equally among all funders. This amount is charged to the lowest funder
	var amountRemainder uint64

	// When a funder is not able to pay, all the remaining funds will be moved
	// to the treasury.
	var slashedFunds uint64

	// Remove all funders who can not afford amountPerFunder
	for len(pool.Funders) > 0 {
		amountPerFunder = amount / uint64(len(pool.Funders))
		amountRemainder = amount - amountPerFunder*uint64(len(pool.Funders))

		lowestFunder := pool.GetLowestFunder()

		if amountRemainder+amountPerFunder > lowestFunder.Amount {
			pool.RemoveFunder(lowestFunder.Address)

			_ = ctx.EventManager().EmitTypedEvent(&pooltypes.EventPoolFundsSlashed{
				PoolId:  poolId,
				Address: lowestFunder.Address,
				Amount:  lowestFunder.Amount,
			})

			slashedFunds += lowestFunder.Amount
		} else {
			break
		}
	}

	if slashedFunds > 0 {
		// send slash to treasury
		if err := util.TransferFromModuleToTreasury(k.accountKeeper, k.distrkeeper, ctx, pooltypes.ModuleName, slashedFunds); err != nil {
			util.PanicHalt(k.upgradeKeeper, ctx, "pool module out of funds")
		}
	}

	if len(pool.Funders) == 0 {
		k.SetPool(ctx, pool)
		return pooltypes.ErrFundsTooLow
	}

	// Remove amount from funders
	for _, funder := range pool.Funders {
		pool.SubtractAmountFromFunder(funder.Address, amountPerFunder)
	}

	lowestFunder := pool.GetLowestFunder()
	pool.SubtractAmountFromFunder(lowestFunder.Address, amountRemainder)

	k.SetPool(ctx, pool)

	return nil
}
