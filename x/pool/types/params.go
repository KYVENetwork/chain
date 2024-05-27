package types

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
)

// DefaultProtocolInflationShare ...
var DefaultProtocolInflationShare = math.LegacyZeroDec()

// DefaultPoolInflationPayoutRate ...
var DefaultPoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.05")

// DefaultMaxVotingPowerPerPool ...
var DefaultMaxVotingPowerPerPool = math.LegacyMustNewDecFromStr("0.5")

// NewParams creates a new Params instance
func NewParams(
	protocolInflationShare math.LegacyDec,
	poolInflationPayoutRate math.LegacyDec,
	maxVotingPowerPerPool math.LegacyDec,
) Params {
	return Params{
		ProtocolInflationShare:  protocolInflationShare,
		PoolInflationPayoutRate: poolInflationPayoutRate,
		MaxVotingPowerPerPool:   maxVotingPowerPerPool,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultProtocolInflationShare,
		DefaultPoolInflationPayoutRate,
		DefaultMaxVotingPowerPerPool,
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

	if err := util.ValidatePercentage(p.MaxVotingPowerPerPool); err != nil {
		return err
	}

	return nil
}
