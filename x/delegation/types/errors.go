package types

import (
	sdkErrors "cosmossdk.io/errors"
)

var (
	ErrNotADelegator                   = sdkErrors.Register(ModuleName, 1000, "not a delegator")
	ErrNotEnoughDelegation             = sdkErrors.Register(ModuleName, 1001, "undelegate-amount is larger than current delegation")
	ErrRedelegationOnCooldown          = sdkErrors.Register(ModuleName, 1002, "all redelegation slots are on cooldown")
	ErrMultipleRedelegationInSameBlock = sdkErrors.Register(ModuleName, 1003, "only one redelegation per delegator per block")
	ErrStakerDoesNotExist              = sdkErrors.Register(ModuleName, 1004, "staker does not exist")
	ErrRedelegationToInactiveStaker    = sdkErrors.Register(ModuleName, 1005, "redelegation to inactive staker not allowed")
)
