package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"

	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
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
		[]*WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				CoinDecimals:              uint32(6),
				MinFundingAmount:          math.NewInt(100_000_000), // 100 $KYVE
				MinFundingAmountPerBundle: math.NewInt(100_000),     // 0.1 $KYVE
				CoinWeight:                math.LegacyNewDec(1),
			},
		},
		DefaultMinFundingMultiple,
	)
}

// Validate validates the set of params
func (p *Params) Validate() error {
	if err := util.ValidateNumber(p.MinFundingMultiple); err != nil {
		return err
	}

	// the native $KYVE coin has to always be whitelisted
	kyveWhitelisted := false

	for _, entry := range p.CoinWhitelist {
		if entry.CoinDenom == "" {
			return errors.New("coin denom is empty")
		}

		if err := util.ValidateInt(entry.MinFundingAmount); err != nil {
			return err
		}

		if err := util.ValidateInt(entry.MinFundingAmountPerBundle); err != nil {
			return err
		}

		if err := util.ValidateDecimal(entry.CoinWeight); err != nil {
			return err
		}

		if entry.CoinDenom == globalTypes.Denom {
			kyveWhitelisted = true
		}
	}

	if !kyveWhitelisted {
		return fmt.Errorf("native KYVE coin \"%s\" not whitelisted", globalTypes.Denom)
	}

	return nil
}
