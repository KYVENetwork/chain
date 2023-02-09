package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalidAuthority     = sdkerrors.Register(ModuleName, 1100, "invalid authority; expected %v, got %v")
	ErrClaimAmountTooHigh   = sdkerrors.Register(ModuleName, 1101, "tried to claim %v tkyve, unlocked amount is only %v tkyve")
	ErrAvailableFundsTooLow = sdkerrors.Register(ModuleName, 1102, "team has %v tkyve available, asking for %v tkyve")
	ErrInvalidClawbackDate  = sdkerrors.Register(ModuleName, 1103, "The clawback can not be set earlier than the last claimed amount")
)
