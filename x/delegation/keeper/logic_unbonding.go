package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FullyProcessDelegatorUnbondingQueue ...
func (k Keeper) FullyProcessDelegatorUnbondingQueue(ctx sdk.Context) {
	allEntries := k.GetAllUnbondingDelegationQueueEntries(ctx)
	for _, entry := range allEntries {
		// Perform undelegation and save undelegated amount to then transfer back to the user
		undelegatedAmount := k.PerformUndelegation(ctx, entry.Staker, entry.Delegator, entry.Amount)

		// Transfer the money
		if err := util.TransferFromModuleToAddress(
			k.bankKeeper,
			ctx,
			types.ModuleName,
			entry.Delegator,
			undelegatedAmount,
		); err != nil {
			util.PanicHalt(k.upgradeKeeper, ctx, "Not enough money in delegation module - logic_unbonding")
		}

		// Emit a delegation event.
		_ = ctx.EventManager().EmitTypedEvent(&types.EventUndelegate{
			Address: entry.Delegator,
			Staker:  entry.Staker,
			Amount:  undelegatedAmount,
		})

	}
}
