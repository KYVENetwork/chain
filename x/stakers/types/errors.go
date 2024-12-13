package types

import (
	"cosmossdk.io/errors"
)

// staking errors
var (
	ErrStakeTooLow             = errors.Register(ModuleName, 1103, "minimum staking amount of %vkyve not reached")
	ErrUnstakeTooHigh          = errors.Register(ModuleName, 1104, "maximum unstaking amount of %vkyve surpassed")
	ErrNoStaker                = errors.Register(ModuleName, 1105, "sender is no staker")
	ErrAlreadyJoinedPool       = errors.Register(ModuleName, 1106, "already joined pool")
	ErrAlreadyLeftPool         = errors.Register(ModuleName, 1107, "already left pool")
	ValaddressAlreadyUsed      = errors.Register(ModuleName, 1108, "valaddress already used")
	ErrStringMaxLengthExceeded = errors.Register(ModuleName, 1109, "String length exceeded: %d vs %d")
	ErrStakerAlreadyCreated    = errors.Register(ModuleName, 1110, "Staker already created")
	ErrValaddressSameAsStaker  = errors.Register(ModuleName, 1111, "Valaddress has same address as Valaddress")
	ErrCanNotJoinDisabledPool  = errors.Register(ModuleName, 1112, "can not join disabled pool")
	ErrInvalidIdentityString   = errors.Register(ModuleName, 1113, "invalid identity: %s")
	ErrNotEnoughRewards        = errors.Register(ModuleName, 1114, "claim amount is larger than current rewards")

	ErrPoolLeaveAlreadyInProgress = errors.Register(ModuleName, 1117, "Pool leave is already in progress")
	ErrValaccountUnauthorized     = errors.Register(ModuleName, 1118, "valaccount unauthorized")
	ErrValidatorJailed            = errors.Register(ModuleName, 1119, "validator jailed")
)
