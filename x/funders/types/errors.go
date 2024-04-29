package types

import (
	"cosmossdk.io/errors"
)

// x/funders module sentinel errors
var (
	ErrFunderAlreadyExists      = errors.Register(ModuleName, 1100, "funder with address %v already exists")
	ErrFunderDoesNotExist       = errors.Register(ModuleName, 1101, "funder with address %v does not exist")
	ErrFundsTooLow              = errors.Register(ModuleName, 1102, "can not defund zero coins")
	ErrMinAmountPerBundle       = errors.Register(ModuleName, 1103, "minimum amount per bundle of coin not reached")
	ErrMinFundingAmount         = errors.Register(ModuleName, 1104, "minimum funding amount of coin not reached")
	ErrFundingDoesNotExist      = errors.Register(ModuleName, 1105, "funding for pool %v and funder %v does not exist")
	ErrFundingStateDoesNotExist = errors.Register(ModuleName, 1106, "funding state for pool %v does not exist")
	ErrMinFundingMultiple       = errors.Register(ModuleName, 1107, "per_bundle_amount times min_funding_multiple is smaller than funded_amount")
	ErrCoinNotWhitelisted       = errors.Register(ModuleName, 1108, "coin not in whitelist")
)
