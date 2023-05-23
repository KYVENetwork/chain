package util

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// PanicHalt performs an emergency upgrade which immediately halts the chain
// The Team has to come up with a solution and develop a patch to handle
// the update.
// In a fully bug-free code this function will never be called.
// This function is there to do assertions and in case of a violation
// it will shut down the chain gracefully, to make it easier to recover from
// a fatal error
func PanicHalt(upgradeKeeper UpgradeKeeper, ctx sdk.Context, message string) {
	// Choose next block for the upgrade
	upgradeBlockHeight := ctx.BlockHeader().Height + 1

	// Create emergency plan
	plan := upgradeTypes.Plan{
		Name:   "emergency_" + strconv.FormatInt(upgradeBlockHeight, 10),
		Height: upgradeBlockHeight,
		Info:   "Emergency Halt; panic occurred; Error:" + message,
	}

	// Directly submit emergency plan
	// Errors can't occur with the current sdk-version
	err := upgradeKeeper.ScheduleUpgrade(ctx, plan)
	if err != nil {
		// Can't happen with current sdk
		panic("Emergency Halt failed: " + message)
	}
}
