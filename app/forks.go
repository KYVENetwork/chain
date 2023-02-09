package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlockForks(ctx sdk.Context, app *App) {
	switch ctx.BlockHeight() {
	default:
		// do nothing
		return
	}
}
