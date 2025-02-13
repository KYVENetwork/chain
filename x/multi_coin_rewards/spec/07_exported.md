<!--
order: 7
-->

# Exported

The `x/multi_coin_rewards` module exports the following functions, which can be used
outside the module.

```go
type MultiCoinRewardsKeeper interface {

    // HandleMultiCoinRewards checks if the user has opted in to receive multi-coin rewards
	// and returns the amount which can get paid out.
	HandleMultiCoinRewards(goCtx context.Context, withdrawAddress sdk.AccAddress, coins sdk.Coins) sdk.Coins
}
```
