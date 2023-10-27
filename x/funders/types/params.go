package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// DefaultMinFundingAmount 1000 Kyve
	DefaultMinFundingAmount = uint64(1_000_000_000)
	// DefaultMinFundingAmountPerBundle 0.1 Kyve
	DefaultMinFundingAmountPerBundle = uint64(100_000)
	// DefaultMinFundingMultiple 20
	DefaultMinFundingMultiple = uint64(20)
)

// NewParams creates a new Params instance
func NewParams(minFundingAmount uint64, minFundingAmountPerBundle uint64, minFundingMultiple uint64) Params {
	return Params{
		MinFundingAmount:          minFundingAmount,
		MinFundingAmountPerBundle: minFundingAmountPerBundle,
		MinFundingMultiple:        minFundingMultiple,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultMinFundingAmount,
		DefaultMinFundingAmountPerBundle,
		DefaultMinFundingMultiple,
	)
}

// Validate validates the set of params
func (p *Params) Validate() error {
	if err := util.ValidateNumber(p.MinFundingAmount); err != nil {
		return err
	}

	if err := util.ValidateNumber(p.MinFundingAmountPerBundle); err != nil {
		return err
	}

	if err := util.ValidateNumber(p.MinFundingAmountPerBundle); err != nil {
		return err
	}

	return nil
}
