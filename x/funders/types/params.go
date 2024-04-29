package types

import (
	"cosmossdk.io/math"
	"errors"
	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
)

const (

	// DefaultMinFundingMultiple 20
	DefaultMinFundingMultiple = uint64(20)
)

var (
	// DefaultCoinWhitelist ukyve
	DefaultCoinWhitelist = []*WhitelistCoinEntry{
		{
			CoinDenom:                 globalTypes.Denom,
			MinFundingAmount:          uint64(1_000_000_000), // 1,000 $KYVE
			MinFundingAmountPerBundle: uint64(100_000),       // 0.1 $KYVE
			CoinWeight:                math.LegacyNewDec(1),
		},
	}
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
		DefaultCoinWhitelist,
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
