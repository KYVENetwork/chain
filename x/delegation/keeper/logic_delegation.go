package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) payoutUndelegation(ctx sdk.Context, stakerAddress, delegatorAddress string, amount uint64) {
	// Transfer the money
	if err := util.TransferFromModuleToAddress(
		k.bankKeeper,
		ctx,
		types.ModuleName,
		delegatorAddress,
		amount,
	); err != nil {
		util.PanicHalt(k.upgradeKeeper, ctx, "Not enough money in delegation module - logic_unbonding")
	}

	// Emit a delegation event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventUndelegate{
		Address: delegatorAddress,
		Staker:  stakerAddress,
		Amount:  amount,
	})
}

// PerformFullUndelegation performs a full undelegation of the staker delegator pair
func (k Keeper) PerformFullUndelegation(ctx sdk.Context, stakerAddress string, delegatorAddress string) uint64 {
	// Withdraw all outstanding rewards
	if _, err := k.performWithdrawal(ctx, stakerAddress, delegatorAddress); err != nil {
		util.PanicHalt(k.upgradeKeeper, ctx, "no money left in module")
	}

	amount := k.f1RemoveDelegator(ctx, stakerAddress, delegatorAddress)

	k.payoutUndelegation(ctx, stakerAddress, delegatorAddress, amount)

	return amount
}

// PerformUndelegation performs immediately an undelegation of the given amount from the given staker
// If the amount is greater than the available amount, only the available amount will be undelegated.
// This method also transfers the rewards back to the given user.
func (k Keeper) PerformUndelegation(ctx sdk.Context, stakerAddress string, delegatorAddress string, amount uint64) uint64 {
	// Withdraw all outstanding rewards
	if _, err := k.performWithdrawal(ctx, stakerAddress, delegatorAddress); err != nil {
		util.PanicHalt(k.upgradeKeeper, ctx, "no money left in module")
	}

	// Perform an internal re-delegation.
	undelegatedAmount := k.f1RemoveDelegator(ctx, stakerAddress, delegatorAddress)

	redelegation := uint64(0)
	if undelegatedAmount > amount {
		// if user didnt undelegate everything ...
		redelegation = undelegatedAmount - amount
		// ... create a new delegator entry with the remaining amount
		k.f1CreateDelegator(ctx, stakerAddress, delegatorAddress, redelegation)
	}

	remainingAmount := undelegatedAmount - redelegation

	k.payoutUndelegation(ctx, stakerAddress, delegatorAddress, remainingAmount)

	return remainingAmount
}

// performWithdrawal withdraws all pending rewards from a user and transfers it.
// The amount is returned by the function.
func (k Keeper) performWithdrawal(ctx sdk.Context, stakerAddress, delegatorAddress string) (sdk.Coins, error) {
	reward := k.f1WithdrawRewards(ctx, stakerAddress, delegatorAddress)
	recipient, errAddress := sdk.AccAddressFromBech32(delegatorAddress)
	if errAddress != nil {
		return nil, errAddress
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, reward); err != nil {
		return nil, err
	}

	// Emit withdraw event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventWithdrawRewards{
		Address: delegatorAddress,
		Staker:  stakerAddress,
		Amounts: reward.String(),
	})

	return reward, nil
}
