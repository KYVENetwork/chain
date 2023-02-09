package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// staking errors
var (
	ErrStakeTooLow             = sdkerrors.Register(ModuleName, 1103, "minimum staking amount of %vkyve not reached")
	ErrUnstakeTooHigh          = sdkerrors.Register(ModuleName, 1104, "maximum unstaking amount of %vkyve surpassed")
	ErrNoStaker                = sdkerrors.Register(ModuleName, 1105, "sender is no staker")
	ErrAlreadyJoinedPool       = sdkerrors.Register(ModuleName, 1106, "already joined pool")
	ErrAlreadyLeftPool         = sdkerrors.Register(ModuleName, 1107, "already left pool")
	ValaddressAlreadyUsed      = sdkerrors.Register(ModuleName, 1108, "valaddress already used")
	ErrStringMaxLengthExceeded = sdkerrors.Register(ModuleName, 1109, "String length exceeded: %d vs %d")
	ErrStakerAlreadyCreated    = sdkerrors.Register(ModuleName, 1110, "Staker already created")
	ErrValaddressSameAsStaker  = sdkerrors.Register(ModuleName, 1111, "Valaddress has same address as Valaddress")
	ErrCanNotJoinDisabledPool  = sdkerrors.Register(ModuleName, 1112, "can not join disabled pool")

	ErrInvalidCommission          = sdkerrors.Register(ModuleName, 1116, "invalid commission %v")
	ErrPoolLeaveAlreadyInProgress = sdkerrors.Register(ModuleName, 1117, "Pool leave is already in progress")
	ErrValaccountUnauthorized     = sdkerrors.Register(ModuleName, 1118, "valaccount unauthorized")
)
