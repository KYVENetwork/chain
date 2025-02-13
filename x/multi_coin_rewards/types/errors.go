package types

import (
	"cosmossdk.io/errors"
)

// staking errors
var (
	ErrMultiCoinDistributionPolicyInvalidAdminAddress = errors.Register(ModuleName, 1122, "multi coin distribution policy admin address invalid")
	ErrMultiCoinDistributionPolicyInvalid             = errors.Register(ModuleName, 1123, "multi coin distribution policy invalid")
	ErrMultiCoinRewardsAlreadyEnabled                 = errors.Register(ModuleName, 1124, "multi coin rewards already enabled")
	ErrMultiCoinRewardsAlreadyDisabled                = errors.Register(ModuleName, 1125, "multi coin rewards already disabled")
)
