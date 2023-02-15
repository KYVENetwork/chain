package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidAuthority     = errors.Register(ModuleName, 1100, "invalid authority; expected %v, got %v")
	ErrClaimAmountTooHigh   = errors.Register(ModuleName, 1101, "tried to claim %v tkyve, unlocked amount is only %v tkyve")
	ErrAvailableFundsTooLow = errors.Register(ModuleName, 1102, "team has %v tkyve available, asking for %v tkyve")
	ErrInvalidClawbackDate  = errors.Register(ModuleName, 1103, "The clawback can not be set earlier than the last claimed amount")
)
