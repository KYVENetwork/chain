package types

import (
	"cosmossdk.io/errors"
)

// x/funders module sentinel errors
var (
	ErrFunderAlreadyExists      = errors.Register(ModuleName, 1100, "funder with address %v already exists")
	ErrFunderDoesNotExist       = errors.Register(ModuleName, 1101, "funder with address %v does not exist")
	ErrFundsTooLow              = errors.Register(ModuleName, 1102, "minimum funding amount of %vkyve not reached")
	ErrMinAmountPerBundle       = errors.Register(ModuleName, 1103, "minimum amount per bundle of %d %s not reached")
	ErrMinFundingAmount         = errors.Register(ModuleName, 1104, "minimum funding amount  %d %s not reached")
	ErrFundingDoesNotExist      = errors.Register(ModuleName, 1105, "funding for pool %v and funder %v does not exist")
	ErrFundingIsUsedUp          = errors.Register(ModuleName, 1106, "funding for pool %v and funder %v is used up")
	ErrFundingStateDoesNotExist = errors.Register(ModuleName, 1107, "funding state for pool %v does not exist")
	ErrMinFundingMultiple       = errors.Register(ModuleName, 1108, "per_bundle_amount (%dkyve) times min_funding_multiple (%d) is smaller than funded_amount (%vkyve)")
	ErrInvalidAmountLength      = errors.Register(ModuleName, 1109, "funding amounts has length %d while amounts_per_bundle has length %d")
	ErrAmountsPerBundleNoSubset = errors.Register(ModuleName, 1110, "amounts_per_bundle is no subset of amounts")
	ErrCoinNotWhitelisted       = errors.Register(ModuleName, 1111, "coin of denom %s not in whitelist")
	ErrDifferentDenom           = errors.Register(ModuleName, 1111, "found denom %s in amount and denom %s in amount_per_bundle")
)
