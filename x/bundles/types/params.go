package types

import (
	"fmt"

	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultUploadTimeout ...
var DefaultUploadTimeout = uint64(600)

// DefaultStorageCost ...
var DefaultStorageCost = sdk.MustNewDecFromStr("0.025")

// DefaultNetworkFee ...
var DefaultNetworkFee = "0.01"

// DefaultMaxPoints ...
var DefaultMaxPoints = uint64(5)

// NewParams creates a new Params instance
func NewParams(
	uploadTimeout uint64,
	storageCost sdk.Dec,
	networkFee string,
	maxPoints uint64,
) Params {
	return Params{
		UploadTimeout: uploadTimeout,
		StorageCost:   storageCost,
		NetworkFee:    networkFee,
		MaxPoints:     maxPoints,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultUploadTimeout,
		DefaultStorageCost,
		DefaultNetworkFee,
		DefaultMaxPoints,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidateUint64(p.UploadTimeout); err != nil {
		return err
	}

	if err := validateStorageCost(p.StorageCost); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.NetworkFee); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.MaxPoints); err != nil {
		return err
	}

	return nil
}

// validateStorageCost ...
func validateStorageCost(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", i)
	}

	return nil
}
