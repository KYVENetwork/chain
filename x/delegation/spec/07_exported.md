<!--
order: 7
-->

# Exported

The `x/delegation` module exports the following functions, which can be used
outside the module. These functions will not return an error, as everything is
handled internally.

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

    // PayoutRewards transfers `amount` $nKYVE from the `payerModuleName`-module to the delegation module.
    // It then awards these tokens internally to all delegators of staker `staker`.
    // Delegators can then receive these rewards if they call the `withdraw`-transaction.
    // This method returns false if the payout fails. This happens usually if there are no
    // delegators for that staker. If this happens one should do something else with the rewards.
    PayoutRewards(ctx sdk.Context, staker string, amount uint64, payerModuleName string) (success bool)

    // SlashDelegators reduces the delegation of all delegators of `staker` by fraction
    // and transfers the amount to the Treasury.
    SlashDelegators(ctx sdk.Context, poolId uint64, staker string, slashType stakertypes.SlashType)

    // GetOutstandingRewards calculates the current rewards a delegator has collected for
    // the given staker.
    GetOutstandingRewards(ctx sdk.Context, staker string, delegator string) uint64

}
```
