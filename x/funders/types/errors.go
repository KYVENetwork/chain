package types

import (
	"cosmossdk.io/errors"
)

// x/funders module sentinel errors
var (
	ErrFunderAlreadyExists      = errors.Register(ModuleName, 1100, "funder with address %v already exists")
	ErrFunderDoesNotExist       = errors.Register(ModuleName, 1101, "funder with address %v does not exist")
	ErrFundsTooLow              = errors.Register(ModuleName, 1102, "minimum funding amount of %vkyve not reached")
	ErrAmountPerBundleTooLow    = errors.Register(ModuleName, 1103, "minimum amount per bundle of %vkyve not reached")
	ErrMinFundingAmount         = errors.Register(ModuleName, 1104, "minimum funding amount of %vkyve not reached")
	ErrFundingDoesNotExist      = errors.Register(ModuleName, 1105, "funding for pool %v and funder %v does not exist")
	ErrFundingIsUsedUp          = errors.Register(ModuleName, 1106, "funding for pool %v and funder %v is used up")
	ErrFundingStateDoesNotExist = errors.Register(ModuleName, 1107, "funding state for pool %v does not exist")
	ErrMinFundingMultiple       = errors.Register(ModuleName, 1108, "per_bundle_amount (%dkyve) times min_funding_multiple (%d) is smaller than funded_amount (%vkyve)")
)
