package types

import (
	"cosmossdk.io/errors"
)

// staking errors
var (
	ErrMultiCoinRefundPolicyInvalidAdminAddress = errors.Register(ModuleName, 1122, "multi coin refund policy admin address invalid")
	ErrMultiCoinRefundPolicyInvalid             = errors.Register(ModuleName, 1123, "multi coin refund policy invalid")
	ErrMultiCoinRewardsAlreadyEnabled           = errors.Register(ModuleName, 1124, "multi coin rewards already enabled")
	ErrMultiCoinRewardsAlreadyDisabled          = errors.Register(ModuleName, 1125, "multi coin rewards already disabled")
)
