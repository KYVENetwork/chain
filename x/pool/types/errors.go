package types

import (
	"cosmossdk.io/errors"
)

// funding errors
var (
	ErrPoolNotFound = errors.Register(ModuleName, 1100, "pool with id %v does not exist")
	ErrInvalidJson  = errors.Register(ModuleName, 1101, "invalid json object: %v")
	ErrInvalidArgs  = errors.Register(ModuleName, 1102, "invalid args")
)
