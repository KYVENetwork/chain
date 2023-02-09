package keeper

import (
	sdkErrors "cosmossdk.io/errors"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// consumeRedelegationSpell checks if the user has still redelegation spells
// available. If so, one spell is used and set on a cooldown.
// If all slots are currently on cooldown the function returns an error
func (k Keeper) consumeRedelegationSpell(ctx sdk.Context, address string) error {
	// Check if cooldowns are over,
	// Remove all expired entries
	for _, creationDate := range k.GetRedelegationCooldownEntries(ctx, address) {
		if ctx.BlockTime().Unix()-int64(creationDate) > int64(k.GetRedelegationCooldown(ctx)) {
			k.RemoveRedelegationCooldown(ctx, address, creationDate)
		} else {
			break
		}
	}

	// Get list of active cooldowns
	creationDates := k.GetRedelegationCooldownEntries(ctx, address)

	// Check if there are still free slots
	if len(creationDates) >= int(k.GetRedelegationMaxAmount(ctx)) {
		return sdkErrors.Wrapf(errorsTypes.ErrLogic, types.ErrRedelegationOnCooldown.Error())
	}

	// Check that no Redelegation occurred in this block, as it will lead to errors, as
	// the block-time is used for an index key.
	if len(creationDates) > 0 && creationDates[len(creationDates)-1] == uint64(ctx.BlockTime().Unix()) {
		return sdkErrors.Wrapf(errorsTypes.ErrLogic, types.ErrMultipleRedelegationInSameBlock.Error())
	}

	// All checks passed, create cooldown entry
	k.SetRedelegationCooldown(ctx, types.RedelegationCooldown{
		Address:      address,
		CreationDate: uint64(ctx.BlockTime().Unix()),
	})

	return nil
}
