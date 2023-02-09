package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var ErrPoolNotFound = sdkerrors.Register(ModuleName, 1100, "pool with id %v does not exist")

// funding errors
var (
	ErrFundsTooLow   = sdkerrors.Register(ModuleName, 1101, "minimum funding amount of %vkyve not reached")
	ErrDefundTooHigh = sdkerrors.Register(ModuleName, 1102, "maximum defunding amount of %vkyve surpassed")
	ErrInvalidJson   = sdkerrors.Register(ModuleName, 1103, "invalid json object: %v")
	ErrInvalidArgs   = sdkerrors.Register(ModuleName, 1104, "invalid args")
)
