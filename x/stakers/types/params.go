package types

import (
	"github.com/KYVENetwork/chain/util"
)

// DefaultCommissionChangeTime ...
var DefaultCommissionChangeTime = uint64(60 * 60 * 24 * 5)

// DefaultLeavePoolTime ...
var DefaultLeavePoolTime = uint64(60 * 60 * 24 * 5)

// NewParams creates a new Params instance
func NewParams(
	commissionChangeTime uint64,
	leavePoolTime uint64,
) Params {
	return Params{
		CommissionChangeTime: commissionChangeTime,
		LeavePoolTime:        leavePoolTime,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultCommissionChangeTime,
		DefaultLeavePoolTime,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidateNumber(p.CommissionChangeTime); err != nil {
		return err
	}

	if err := util.ValidateNumber(p.LeavePoolTime); err != nil {
		return err
	}

	return nil
}
