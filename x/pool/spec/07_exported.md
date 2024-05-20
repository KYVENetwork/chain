<!--
order: 7
-->

# Exported

The `x/pool` module exports the following functions, which can be used
outside the module.

```go
type PoolKeeper interface {

    // AssertPoolExists returns nil if the pool exists and types.ErrPoolNotFound if it does not.
    AssertPoolExists(ctx sdk.Context, poolId uint64) error

    // GetPoolWithError returns a pool by its poolId, if the pool does not exist,
    // a types.ErrPoolNotFound error is returned 
    GetPoolWithError(ctx sdk.Context, poolId uint64) (pooltypes.Pool, error)

    // TODO(@troy,@max) double check bundles module ( GetPoolWithError and GetPool)
    GetPool(ctx sdk.Context, id uint64) (val pooltypes.Pool, found bool)

    // IncrementBundleInformation updates the latest finalized bundle of a pool
    IncrementBundleInformation(ctx sdk.Context, poolId uint64, currentHeight uint64, currentKey string, currentValue string)

    GetAllPools(ctx sdk.Context) (list []pooltypes.Pool)

    // ChargeFundersOfPool equally splits the amount between all funders and removes
    // the appropriate amount from each funder.
    // All funders who can't afford the amount, are kicked out.
    // The method returns the payout amount the pool was able to charge from the funders.
    // This method only charges coins which are whitelisted.
    ChargeFundersOfPool(ctx sdk.Context, poolId uint64, amount uint64, recipient string) (payout uint64, err error)
}
```
