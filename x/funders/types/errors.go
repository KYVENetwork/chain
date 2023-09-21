package types

import (
	"cosmossdk.io/errors"
)

// x/funders module sentinel errors
var (
	ErrFunderAlreadyExists   = errors.Register(ModuleName, 1100, "funder with address %v already exists")
	ErrFunderDoesNotExist    = errors.Register(ModuleName, 1101, "funder with address %v does not exist")
	ErrFundsTooLow           = errors.Register(ModuleName, 1102, "minimum funding amount of %vkyve not reached")
	ErrAmountPerBundleTooLow = errors.Register(ModuleName, 1103, "minimum amount per bundle of %vkyve not reached")
	ErrMinFundingAmount      = errors.Register(ModuleName, 1104, "minimum funding amount of %vkyve not reached")
	//ErrDefundTooHigh = errors.Register(ModuleName, 1102, "maximum defunding amount of %vkyve surpassed")
)
