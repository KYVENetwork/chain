package types

import (
	"errors"
	"github.com/KYVENetwork/chain/util"
)

const (
	// DefaultMinFundingMultiple 20
	DefaultMinFundingMultiple = uint64(20)
)

// NewParams creates a new Params instance
func NewParams(coinWhitelist []*WhitelistCoinEntry, minFundingMultiple uint64) Params {
	return Params{
		CoinWhitelist:      coinWhitelist,
		MinFundingMultiple: minFundingMultiple,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		make([]*WhitelistCoinEntry, 0),
		DefaultMinFundingMultiple,
	)
}

// Validate validates the set of params
func (p *Params) Validate() error {
	if err := util.ValidateNumber(p.MinFundingMultiple); err != nil {
		return err
	}

	for _, entry := range p.CoinWhitelist {
		if entry.CoinDenom == "" {
			return errors.New("coin denom is empty")
		}

		if err := util.ValidateNumber(entry.MinFundingAmount); err != nil {
			return err
		}

		if err := util.ValidateNumber(entry.MinFundingAmountPerBundle); err != nil {
			return err
		}

		if err := util.ValidateDecimal(entry.CoinWeight); err != nil {
			return err
		}
	}

	return nil
}
