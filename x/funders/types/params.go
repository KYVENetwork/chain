package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// DefaultMinFundingAmount 1000 Kyve
	DefaultMinFundingAmount = uint64(1_000_000_000)
	// DefaultMinFundingAmountPerBundle 1 Kyve
	DefaultMinFundingAmountPerBundle = uint64(1_000_000)
)

// NewParams creates a new Params instance
func NewParams(minFundingAmount uint64, minFundingAmountPerBundle uint64) Params {
	return Params{
		MinFundingAmount:          minFundingAmount,
		MinFundingAmountPerBundle: minFundingAmountPerBundle,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultMinFundingAmount,
		DefaultMinFundingAmountPerBundle,
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

	return nil
}
