package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// DefaultMultiCoinDistributionPendingTime ...
var DefaultMultiCoinDistributionPendingTime = uint64(60 * 60 * 24 * 14)

// NewParams creates a new Params instance
func NewParams(
	multiCoinDistributionPendingTime uint64,
	multiCoinDistributionPolicyAdminAddress string,
) Params {
	return Params{
		MultiCoinDistributionPendingTime:        multiCoinDistributionPendingTime,
		MultiCoinDistributionPolicyAdminAddress: multiCoinDistributionPolicyAdminAddress,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultMultiCoinDistributionPendingTime,
		"",
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.MultiCoinDistributionPolicyAdminAddress != "" {
		_, err := sdk.AccAddressFromBech32(p.MultiCoinDistributionPolicyAdminAddress)
		if err != nil {
			return err
		}
	}

	return nil
}
