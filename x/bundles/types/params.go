package types

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/KYVENetwork/chain/util"
)

// DefaultUploadTimeout ...
var DefaultUploadTimeout = uint64(600)

// DefaultStorageCosts ...
func DefaultStorageCosts() map[uint32]LegacyDecValue {
	return map[uint32]LegacyDecValue{
		0: {Value: math.LegacyMustNewDecFromStr("0.025")},
	}
}

// DefaultNetworkFee ...
var DefaultNetworkFee = math.LegacyMustNewDecFromStr("0.01")

// DefaultMaxPoints ...
var DefaultMaxPoints = uint64(24)

// NewParams creates a new Params instance
func NewParams(
	uploadTimeout uint64,
	storageCosts map[uint32]LegacyDecValue,
	networkFee math.LegacyDec,
	maxPoints uint64,
) Params {
	return Params{
		UploadTimeout: uploadTimeout,
		StorageCosts:  storageCosts,
		NetworkFee:    networkFee,
		MaxPoints:     maxPoints,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultUploadTimeout,
		DefaultStorageCosts(),
		DefaultNetworkFee,
		DefaultMaxPoints,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidatePositiveNumber(p.UploadTimeout); err != nil {
		return err
	}

	if len(p.StorageCosts) == 0 {
		return fmt.Errorf("storage costs cannot be empty")
	}
	for _, v := range p.StorageCosts {
		if err := util.ValidateDecimal(v.Value); err != nil {
			return err
		}
	}

	if err := util.ValidatePercentage(p.NetworkFee); err != nil {
		return err
	}

	if err := util.ValidatePositiveNumber(p.MaxPoints); err != nil {
		return err
	}

	return nil
}
