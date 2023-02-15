package types

import (
	"cosmossdk.io/errors"
)

var ErrPoolNotFound = errors.Register(ModuleName, 1100, "pool with id %v does not exist")

// funding errors
var (
	ErrFundsTooLow   = errors.Register(ModuleName, 1101, "minimum funding amount of %vkyve not reached")
	ErrDefundTooHigh = errors.Register(ModuleName, 1102, "maximum defunding amount of %vkyve surpassed")
	ErrInvalidJson   = errors.Register(ModuleName, 1103, "invalid json object: %v")
	ErrInvalidArgs   = errors.Register(ModuleName, 1104, "invalid args")
)
