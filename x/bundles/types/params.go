package types

import (
	"github.com/KYVENetwork/chain/util"
)

// DefaultUploadTimeout ...
var DefaultUploadTimeout = uint64(600)

// DefaultStorageCost ...
var DefaultStorageCost = uint64(100000)

// DefaultNetworkFee ...
var DefaultNetworkFee = "0.01"

// DefaultMaxPoints ...
var DefaultMaxPoints = uint64(5)

// NewParams creates a new Params instance
func NewParams(
	uploadTimeout uint64,
	storageCost uint64,
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

	if err := util.ValidateUint64(p.StorageCost); err != nil {
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
