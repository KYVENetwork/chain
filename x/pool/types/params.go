package types

import (
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultProtocolInflationShare ...
var DefaultProtocolInflationShare = sdk.ZeroDec()

// DefaultPoolInflationPayoutRate ...
var DefaultPoolInflationPayoutRate = sdk.MustNewDecFromStr("0.05")

// NewParams creates a new Params instance
func NewParams(
	protocolInflationShare sdk.Dec,
	poolInflationPayoutRate sdk.Dec,
) Params {
	return Params{
		ProtocolInflationShare:  protocolInflationShare,
		PoolInflationPayoutRate: poolInflationPayoutRate,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultProtocolInflationShare,
		DefaultPoolInflationPayoutRate,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidatePercentage(p.ProtocolInflationShare); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.PoolInflationPayoutRate); err != nil {
		return err
	}

	return nil
}
