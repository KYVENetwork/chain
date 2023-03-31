package types

import (
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultUploadTimeout ...
var DefaultUploadTimeout = uint64(600)

// DefaultStorageCost ...
var DefaultStorageCost = sdk.MustNewDecFromStr("0.025")

// DefaultNetworkFee ...
var DefaultNetworkFee = sdk.MustNewDecFromStr("0.01")

// DefaultMaxPoints ...
var DefaultMaxPoints = uint64(24)

// NewParams creates a new Params instance
func NewParams(
	uploadTimeout uint64,
	storageCost sdk.Dec,
	networkFee sdk.Dec,
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
	if err := util.ValidatePositiveNumber(p.UploadTimeout); err != nil {
		return err
	}

	if err := util.ValidateDecimal(p.StorageCost); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.NetworkFee); err != nil {
		return err
	}

	if err := util.ValidatePositiveNumber(p.MaxPoints); err != nil {
		return err
	}

	return nil
}
