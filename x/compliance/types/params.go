package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// DefaultMultiCoinRefundPendingTime ...
var DefaultMultiCoinRefundPendingTime = uint64(60 * 60 * 24 * 14)

// NewParams creates a new Params instance
func NewParams(
	multiCoinRefundPendingTime uint64,
	multiCoinRefundPolicyAdminAddress string,
) Params {
	return Params{
		MultiCoinRefundPendingTime:        multiCoinRefundPendingTime,
		MultiCoinRefundPolicyAdminAddress: multiCoinRefundPolicyAdminAddress,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultMultiCoinRefundPendingTime,
		"",
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.MultiCoinRefundPolicyAdminAddress != "" {
		_, err := sdk.AccAddressFromBech32(p.MultiCoinRefundPolicyAdminAddress)
		if err != nil {
			return err
		}
	}

	return nil
}
