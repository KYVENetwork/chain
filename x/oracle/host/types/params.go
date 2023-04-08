package types

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
)

// DefaultPricePerByte is 0 (i.e. free).
var DefaultPricePerByte = math.LegacyZeroDec()

// NewParams creates a new Params object.
func NewParams(pricePerByte math.LegacyDec) Params {
	return Params{PricePerByte: pricePerByte}
}

// DefaultParams creates the default Params object.
func DefaultParams() Params {
	return NewParams(DefaultPricePerByte)
}

// Validate validates the params to ensure the expected invariants holds.
func (p *Params) Validate() error {
	return util.ValidateLegacyDecimal(p.PricePerByte)
}
