<!--
order: 6
-->

# Exported

The `x/funders` module exports the following functions, which can be used
outside the module.

```go
type FundersKeeper interface {
    // ChargeFundersOfPool equally splits the amount between all funders and removes
    // the appropriate amount from each funder.
    // All funders who can't afford the amount, are kicked out.
    // The method returns the payout amount the pool was able to charge from the funders.
    ChargeFundersOfPool(ctx sdk.Context, poolId uint64) (payout sdk.Coins, err error)
}
```
