package util

import sdk "github.com/cosmos/cosmos-sdk/types"

// TruncateDecCoins converts sdm.DecCoins to sdk.Coins by truncating all values to integers.
func TruncateDecCoins(decCoins sdk.DecCoins) sdk.Coins {
	coins := sdk.NewCoins()
	for _, coin := range decCoins {
		coins = coins.Add(sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()))
	}
	return coins
}
