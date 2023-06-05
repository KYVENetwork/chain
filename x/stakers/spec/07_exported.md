<!--
order: 7
-->

# Exported

The `x/stakers` module exports the following functions, which can be used
outside the module.

```go
type StakersKeeper interface {

    // LeavePool removes a staker from a pool and emits the corresponding event.
    // The staker is no longer able to participate in the given pool.
    // All points the staker had in that pool are deleted.
	LeavePool(ctx sdk.Context, staker string, poolId uint64)

    // GetAllStakerAddressesOfPool returns a list of all stakers
    // which have currently a valaccount registered for the given pool
    // and are therefore allowed to participate in that pool.
    GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)

    // GetCommission returns the commission of a staker as a parsed sdk.Dec
	GetCommission(ctx sdk.Context, stakerAddress string) sdk.Dec

    // AssertValaccountAuthorized checks if the given `valaddress` is allowed to vote in pool
    // with id `poolId` to vote in favor of `stakerAddress`.
    // If the valaddress is not authorized the appropriate error is returned.
    // Otherwise, it returns `nil`
    AssertValaccountAuthorized(ctx sdk.Context, poolId uint64, stakerAddress string, valaddress string) error

    // GetActiveStakers returns all staker-addresses that are
    // currently participating in at least one pool.
    GetActiveStakers(ctx sdk.Context) []string

    // TotalBondedTokens returns all tokens which are currently bonded by the protocol
    // I.e. the sum of all delegation of all stakers that are currently participating
    // in at least one pool
    TotalBondedTokens(ctx sdk.Context) math.Int

    // GetActiveValidators returns all protocol-node information which
    // are needed by the governance to calculate the voting powers.
    // The interface needs to correspond to github.com/cosmos/cosmos-sdk/x/gov/types/v1.ValidatorGovInfo
    // But as there is no direct dependency in the cosmos-sdk-fork this value is passed as an interface{}
    GetActiveValidators(ctx sdk.Context) (validators []interface{})

    // GetDelegations returns the address and the delegation amount of all active protocol-stakers the
    // delegator as delegated to. This is used to calculate the vote weight each delegator has.
    GetDelegations(ctx sdk.Context, delegator string) (validators []string, amounts []sdk.Dec)

    // IncrementPoints increments to Points for a staker in a given pool.
    // Returns the amount of the current points (including the current incrementation)
    IncrementPoints(ctx sdk.Context, poolId uint64, stakerAddress string) uint64

    // ResetPoints sets the point count for the staker in the given pool back to zero.
    // Returns the amount of points the staker had before the reset.
    ResetPoints(ctx sdk.Context, poolId uint64, stakerAddress string) (previousPoints uint64)

    // DoesValaccountExist only checks if the key is present in the KV-Store
    // without loading and unmarshalling to full entry
    DoesValaccountExist(ctx sdk.Context, poolId uint64, stakerAddress string) bool

    DoesStakerExist(ctx sdk.Context, staker string) bool

    // IncreaseStakerCommissionRewards increases the commission rewards of a
    // staker by a specific amount. It can not be decreased, only the
    // MsgClaimCommissionRewards message can decrease this value.
    IncreaseStakerCommissionRewards(ctx sdk.Context, address string, amount uint64)
}
```
