package keeper

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
This file is responsible for implementing the F1-Fee distribution as described in
https://drops.dagstuhl.de/opus/volltexte/2020/11974/pdf/OASIcs-Tokenomics-2019-10.pdf

We recommend reading the paper first before reading the code.
This file covers all relevant methods to fully implement the algorithm.
It also takes fully care of the entire state. The only interaction needed
is covered by the available methods.
*/

// f1StartNewPeriod finishes the current period according to the F1-Paper
// It returns the index of the new period.
// delegationData is passed as a pointer and updated in this method
// it's the responsibility of the caller to save the meta-data state.
// This method only writes to the entries.
func (k Keeper) f1StartNewPeriod(ctx sdk.Context, staker string, delegationData *types.DelegationData) uint64 {
	// Ending the current period is performed by getting the entry
	// of the previous index and adding the current quotient of
	// $T_f / n_f$

	// Get previous entry
	// F1: corresponds to $Entry_{f-1}$
	previousEntry, found := k.GetDelegationEntry(ctx, staker, delegationData.LatestIndexK)
	if !found {
		previousEntry.Value = sdk.NewDecCoins()
	}

	// Calculate quotient of current period
	// If totalDelegation is zero the quotient is also zero
	currentPeriodValue := sdk.NewDecCoins()
	if delegationData.TotalDelegation != 0 {
		decCurrentRewards := sdk.NewDecCoinsFromCoins(delegationData.CurrentRewards...)
		decTotalDelegation := math.LegacyNewDec(int64(delegationData.TotalDelegation))

		// F1: $T_f / n_f$
		currentPeriodValue = decCurrentRewards.QuoDec(decTotalDelegation)
	}

	// Add previous entry to current one
	currentPeriodValue = currentPeriodValue.Add(previousEntry.Value...)

	// Increment index for the next period
	indexF := delegationData.LatestIndexK + 1

	// Add entry for new period to KV-Store
	k.SetDelegationEntry(ctx, types.DelegationEntry{
		Value:  currentPeriodValue,
		Staker: staker,
		KIndex: indexF,
	})

	// Reset the rewards for the next period back to zero
	// and update to the new index
	delegationData.CurrentRewards = sdk.NewCoins()
	delegationData.LatestIndexK = indexF

	if delegationData.LatestIndexWasUndelegation {
		k.RemoveDelegationEntry(ctx, previousEntry.Staker, previousEntry.KIndex)
		delegationData.LatestIndexWasUndelegation = false
	}

	return indexF
}

// f1CreateDelegator creates a new delegator within the f1-logic.
// It is assumed that no delegator exists.
func (k Keeper) f1CreateDelegator(ctx sdk.Context, staker string, delegator string, amount uint64) {
	if amount == 0 {
		return
	}

	// Fetch metadata
	delegationData, found := k.GetDelegationData(ctx, staker)

	// Init default data-set, if this is the first delegator
	if !found {
		delegationData = types.DelegationData{
			Staker:         staker,
			CurrentRewards: sdk.NewCoins(),
		}
	}

	// Finish current period
	k.f1StartNewPeriod(ctx, staker, &delegationData)

	// Update metadata
	delegationData.TotalDelegation += amount
	delegationData.DelegatorCount += 1
	k.SetDelegationData(ctx, delegationData)

	k.SetDelegator(ctx, types.Delegator{
		Staker:        staker,
		Delegator:     delegator,
		InitialAmount: amount,
		KIndex:        delegationData.LatestIndexK,
	})
}

// f1RemoveDelegator performs a full undelegation and removes the delegator from the f1-logic
// This method returns the amount of tokens that got undelegated
// Due to slashing the undelegated amount can be lower than the initial delegated amount
func (k Keeper) f1RemoveDelegator(ctx sdk.Context, stakerAddress string, delegatorAddress string) (amount uint64) {
	// Check if delegator exists
	delegator, found := k.GetDelegator(ctx, stakerAddress, delegatorAddress)
	if !found {
		return 0
	}

	// Fetch metadata
	delegationData, found := k.GetDelegationData(ctx, stakerAddress)
	if !found {
		// Should never happen, if so there is an error in the f1-implementation
		util.PanicHalt(k.upgradeKeeper, ctx, "No delegationData although somebody is delegating")
	}

	balance := k.f1GetCurrentDelegation(ctx, stakerAddress, delegatorAddress)

	// Start new period
	k.f1StartNewPeriod(ctx, stakerAddress, &delegationData)

	delegationData.LatestIndexWasUndelegation = true

	// Update Metadata
	delegationData.TotalDelegation -= balance
	delegationData.DelegatorCount -= 1

	// Remove Delegator
	k.RemoveDelegator(ctx, delegator.Staker, delegator.Delegator)
	// Remove old entry
	k.RemoveDelegationEntry(ctx, stakerAddress, delegator.KIndex)

	// Final cleanup
	if delegationData.DelegatorCount == 0 {
		k.RemoveDelegationEntry(ctx, stakerAddress, delegationData.LatestIndexK)
	}
	k.SetDelegationData(ctx, delegationData)

	return balance
}

// f1WithdrawRewards calculates all outstanding rewards and withdraws them from
// the f1-logic. A new period starts.
func (k Keeper) f1WithdrawRewards(ctx sdk.Context, stakerAddress string, delegatorAddress string) sdk.Coins {
	delegator, found := k.GetDelegator(ctx, stakerAddress, delegatorAddress)
	if !found {
		return sdk.NewCoins()
	}

	// Fetch metadata
	delegationData, found := k.GetDelegationData(ctx, stakerAddress)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, "No delegationData although somebody is delegating")
	}

	// End current period and use it for calculating the reward
	endIndex := k.f1StartNewPeriod(ctx, stakerAddress, &delegationData)
	k.SetDelegationData(ctx, delegationData)

	// According to F1 the reward is calculated as the difference between two entries multiplied by the
	// delegation amount for the period.
	// To incorporate slashing one needs to iterate all slashes and calculate the reward for every period
	// separately and then sum it.
	reward := sdk.NewDecCoins()
	k.f1IterateConstantDelegationPeriods(ctx, stakerAddress, delegatorAddress, delegator.KIndex, endIndex,
		func(startIndex uint64, endIndex uint64, delegation math.LegacyDec) {
			// entry difference
			difference := k.f1GetEntryDifference(ctx, stakerAddress, startIndex, endIndex)

			periodReward := safeMulDec(difference, delegation)

			reward = reward.Add(periodReward...)
		})

	// Delete Delegator entry as he has no outstanding rewards anymore.
	// To account for slashes, also update the initial amount.
	k.RemoveDelegationEntry(ctx, stakerAddress, delegator.KIndex)
	// Delegator now starts at the latest index.
	delegator.KIndex = endIndex
	delegator.InitialAmount = k.f1GetCurrentDelegation(ctx, delegator.Staker, delegator.Delegator)
	k.SetDelegator(ctx, delegator)

	return truncateDecCoins(reward)
}

// f1IterateConstantDelegationPeriods iterates all periods between minIndex and maxIndex (both inclusive)
// and calls handler() for every period with constant delegation amount
// This method iterates all slashes and additionally calls handler at least once if no slashes occurred
func (k Keeper) f1IterateConstantDelegationPeriods(ctx sdk.Context, stakerAddress string, delegatorAddress string,
	minIndex uint64, maxIndex uint64, handler func(startIndex uint64, endIndex uint64, delegation math.LegacyDec),
) {
	slashes := k.GetAllDelegationSlashesBetween(ctx, stakerAddress, minIndex, maxIndex)

	delegator, _ := k.GetDelegator(ctx, stakerAddress, delegatorAddress)
	delegatorBalance := math.LegacyNewDec(int64(delegator.InitialAmount))

	if len(slashes) == 0 {
		handler(minIndex, maxIndex, delegatorBalance)
		return
	}

	prevIndex := minIndex
	for _, slash := range slashes {
		handler(prevIndex, slash.KIndex, delegatorBalance)
		slashedAmount := delegatorBalance.MulTruncate(slash.Fraction)
		delegatorBalance = delegatorBalance.Sub(slashedAmount)
		prevIndex = slash.KIndex
	}
	handler(prevIndex, maxIndex, delegatorBalance)
}

// f1GetCurrentDelegation calculates the current delegation of a delegator.
// I.e. the initial amount minus the slashes
func (k Keeper) f1GetCurrentDelegation(ctx sdk.Context, stakerAddress string, delegatorAddress string) uint64 {
	delegator, found := k.GetDelegator(ctx, stakerAddress, delegatorAddress)
	if !found {
		return 0
	}

	// Fetch metadata
	delegationData, found := k.GetDelegationData(ctx, stakerAddress)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, "No delegationData although somebody is delegating")
	}

	latestBalance := math.LegacyNewDec(int64(delegator.InitialAmount))
	k.f1IterateConstantDelegationPeriods(ctx, stakerAddress, delegatorAddress, delegator.KIndex, delegationData.LatestIndexK,
		func(startIndex uint64, endIndex uint64, delegation math.LegacyDec) {
			latestBalance = delegation
		})

	return latestBalance.TruncateInt().Uint64()
}

// f1GetOutstandingRewards calculates the current outstanding rewards without modifying the f1-state.
// This method can be used for queries.
func (k Keeper) f1GetOutstandingRewards(ctx sdk.Context, stakerAddress string, delegatorAddress string) sdk.Coins {
	delegator, found := k.GetDelegator(ctx, stakerAddress, delegatorAddress)
	if !found {
		return sdk.NewCoins()
	}

	// Fetch metadata
	delegationData, found := k.GetDelegationData(ctx, stakerAddress)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, "No delegationData although somebody is delegating")
	}

	// End current period and use it for calculating the reward
	endIndex := delegationData.LatestIndexK

	// According to F1 the reward is calculated as the difference between two entries multiplied by the
	// delegation amount for the period.
	// To incorporate slashing one needs to iterate all slashes and calculate the reward for every period
	// separately and then sum it.
	reward := sdk.NewDecCoins()
	latestBalance := math.LegacyNewDec(int64(delegator.InitialAmount))
	k.f1IterateConstantDelegationPeriods(ctx, stakerAddress, delegatorAddress, delegator.KIndex, endIndex,
		func(startIndex uint64, endIndex uint64, delegation math.LegacyDec) {
			difference := k.f1GetEntryDifference(ctx, stakerAddress, startIndex, endIndex)
			// Multiply with delegation for period
			periodReward := safeMulDec(difference, delegation)
			// Add to total rewards
			reward = reward.Add(periodReward...)

			// For calculating the last (ongoing) period
			latestBalance = delegation
		})

	// Append missing rewards from last period to ongoing period
	entry, found := k.GetDelegationEntry(ctx, stakerAddress, delegationData.LatestIndexK)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, "Entry does not exist")
	}
	_ = entry

	currentPeriodValue := sdk.NewDecCoins()
	if delegationData.TotalDelegation != 0 {
		decCurrentRewards := sdk.NewDecCoinsFromCoins(delegationData.CurrentRewards...)
		decTotalDelegation := math.LegacyNewDec(int64(delegationData.TotalDelegation))

		// F1: $T_f / n_f$
		currentPeriodValue = decCurrentRewards.QuoDec(decTotalDelegation)
	}

	ongoingPeriodReward := safeMulDec(currentPeriodValue, latestBalance)

	reward = reward.Add(ongoingPeriodReward...)
	return truncateDecCoins(reward)
}

func (k Keeper) f1GetEntryDifference(ctx sdk.Context, stakerAddress string, lowIndex uint64, highIndex uint64) sdk.DecCoins {
	// entry difference
	firstEntry, found := k.GetDelegationEntry(ctx, stakerAddress, lowIndex)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, "Entry 1 does not exist")
	}

	secondEntry, found := k.GetDelegationEntry(ctx, stakerAddress, highIndex)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, "Entry 2 does not exist")
	}

	return secondEntry.Value.Sub(firstEntry.Value)
}

// truncateDecCoins converts sdm.DecCoins to sdk.Coins by truncating all values to integers.
func truncateDecCoins(decCoins sdk.DecCoins) sdk.Coins {
	coins := sdk.NewCoins()
	for _, coin := range decCoins {
		coins = coins.Add(sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return coins
}

func safeMulDec(coins sdk.DecCoins, scalar math.LegacyDec) sdk.DecCoins {
	if scalar.IsZero() {
		return sdk.NewDecCoins()
	}
	return coins.MulDec(scalar)
}
