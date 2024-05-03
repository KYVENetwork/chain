<!--
order: 7
-->

# Exported

The `x/delegation` module exports the following functions, which can be used
outside the module.

```go
type DelegationKeeper interface {

    // GetDelegationAmount returns the sum of all delegations for a specific staker.
    // If the staker does not exist, it returns zero as the staker has zero delegations
    GetDelegationAmount(ctx sdk.Context, staker string) uint64

    // GetDelegationAmountOfDelegator returns the amount of how many $KYVE `delegatorAddress`
    // has delegated to `stakerAddress`. If one of the addresses does not exist, it returns zero.
    GetDelegationAmountOfDelegator(ctx sdk.Context, stakerAddress string, delegatorAddress string) uint64

    // GetDelegationOfPool returns the amount of how many $KYVE users have delegated
    // to stakers that are participating in the given pool
    GetDelegationOfPool(ctx sdk.Context, poolId uint64) uint64

    // GetTotalAndHighestDelegationOfPool returns the total delegation amount of all validators in the given pool and
    // the highest total delegation amount of a single validator in a pool
    GetTotalAndHighestDelegationOfPool(ctx sdk.Context, poolId uint64) (totalDelegation, highestDelegation uint64)

    // PayoutRewards transfers `amount` from the `payerModuleName`-module to the delegation module.
    // It then awards these tokens internally to all delegators of staker `staker`.
    // Delegators can then receive these rewards if they call the `withdraw`-transaction.
    // If the staker has no delegators or the module to module transfer fails, the method fails and
    // returns an error.
	PayoutRewards(ctx sdk.Context, staker string, amount sdk.Coins, payerModuleName string) error

    // SlashDelegators reduces the delegation of all delegators of `staker` by fraction
    // and transfers the amount to the Treasury.
    SlashDelegators(ctx sdk.Context, poolId uint64, staker string, slashType stakertypes.SlashType)

    // GetOutstandingRewards calculates the current rewards a delegator has collected for
    // the given staker.
	GetOutstandingRewards(ctx sdk.Context, staker string, delegator string) sdk.Coins

}
```
