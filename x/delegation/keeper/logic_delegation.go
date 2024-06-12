package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegate performs a safe delegation with all necessary checks
// Warning: does not transfer the amount (only the rewards)
func (k Keeper) performDelegation(ctx sdk.Context, stakerAddress string, delegatorAddress string, amount uint64) {
	// Update in-memory staker index for efficient queries
	k.RemoveStakerIndex(ctx, stakerAddress)
	defer k.SetStakerIndex(ctx, stakerAddress)

	if k.DoesDelegatorExist(ctx, stakerAddress, delegatorAddress) {
		// If the sender is already a delegator, first perform an undelegation, before delegating.
		// "perform a withdrawal"
		if _, err := k.performWithdrawal(ctx, stakerAddress, delegatorAddress); err != nil {
			util.PanicHalt(k.upgradeKeeper, ctx, "no money left in module")
		}

		// Perform delegation by fully undelegating and then delegating the new amount
		unDelegateAmount := k.f1RemoveDelegator(ctx, stakerAddress, delegatorAddress)
		newDelegationAmount := unDelegateAmount + amount
		k.f1CreateDelegator(ctx, stakerAddress, delegatorAddress, newDelegationAmount)
	} else {
		// If the sender isn't a delegator, simply create a new delegation entry.
		k.f1CreateDelegator(ctx, stakerAddress, delegatorAddress, amount)
	}
}

// performUndelegation performs immediately an undelegation of the given amount from the given staker
// If the amount is greater than the available amount, only the available amount will be undelegated.
// This method also transfers the rewards back to the given user.
func (k Keeper) performUndelegation(ctx sdk.Context, stakerAddress string, delegatorAddress string, amount uint64) uint64 {
	// Update in-memory staker index for efficient queries
	k.RemoveStakerIndex(ctx, stakerAddress)
	defer k.SetStakerIndex(ctx, stakerAddress)

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

	return undelegatedAmount - redelegation
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
