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
    // Their remaining amount is transferred to the Treasury.
    // The function throws an error if pool ran out of funds.
    // This method does not transfer any funds. The bundles-module
    // is responsible for transferring the rewards out of the module.
    ChargeFundersOfPool(ctx sdk.Context, poolId uint64, amount uint64) error
}
```
