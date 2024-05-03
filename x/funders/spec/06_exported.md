<!--
order: 6
-->

# Exported

The `x/funders` module exports the following functions, which can be used
outside the module.

```go
type FundersKeeper interface {
    // ChargeFundersOfPool charges all funders of a pool with their amount_per_bundle
    // If the amount is lower than the amount_per_bundle,
    // the max amount is charged and the funder is removed from the active funders list.
    // The amount is transferred from the funders to the recipient module account.
    // If there are no more active funders, an event is emitted.
    ChargeFundersOfPool(ctx sdk.Context, poolId uint64, recipient string) (payout sdk.Coins, err error)
}
```
