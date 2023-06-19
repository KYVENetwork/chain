package types

import (
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultProtocolInflationShare ...
var DefaultProtocolInflationShare = sdk.ZeroDec()

// NewParams creates a new Params instance
func NewParams(protocolInflationShare sdk.Dec) Params {
	return Params{
		ProtocolInflationShare: protocolInflationShare,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultProtocolInflationShare)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidatePercentage(p.ProtocolInflationShare); err != nil {
		return err
	}

	return nil
}
