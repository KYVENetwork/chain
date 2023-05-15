package types

import (
	"github.com/KYVENetwork/chain/util"
)

// DefaultGlobalMinDelegation ...
var DefaultGlobalMinDelegation = uint64(0)

// NewParams creates a new Params instance
func NewParams(
	globalMinDelegation uint64,
) Params {
	return Params{
		GlobalMinDelegation: globalMinDelegation,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultGlobalMinDelegation,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidateNumber(p.GlobalMinDelegation); err != nil {
		return err
	}

	return nil
}
