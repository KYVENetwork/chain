package types

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
)

// DefaultCommissionChangeTime ...
var DefaultCommissionChangeTime = uint64(60 * 60 * 24 * 5)

// DefaultLeavePoolTime ...
var DefaultLeavePoolTime = uint64(60 * 60 * 24 * 5)

// DefaultStakeFractionChangeTime ...
var DefaultStakeFractionChangeTime = uint64(60 * 60 * 24 * 5)

// DefaultMultiCoinRefundPendingTime ...
var DefaultMultiCoinRefundPendingTime = uint64(60 * 60 * 24 * 14)

// DefaultVoteSlash ...
var DefaultVoteSlash = math.LegacyMustNewDecFromStr("0.01")

// DefaultUploadSlash ...
var DefaultUploadSlash = math.LegacyMustNewDecFromStr("0.02")

// DefaultTimeoutSlash ...
var DefaultTimeoutSlash = math.LegacyMustNewDecFromStr("0.002")

// NewParams creates a new Params instance
func NewParams(
	commissionChangeTime uint64,
	leavePoolTime uint64,
	stakeFractionChangeTime uint64,
	multiCoinRefundPendingTime uint64,
	voteSlash math.LegacyDec,
	uploadSlash math.LegacyDec,
	timeoutSlash math.LegacyDec,
) Params {
	return Params{
		CommissionChangeTime:       commissionChangeTime,
		LeavePoolTime:              leavePoolTime,
		StakeFractionChangeTime:    stakeFractionChangeTime,
		MultiCoinRefundPendingTime: multiCoinRefundPendingTime,
		VoteSlash:                  voteSlash,
		UploadSlash:                uploadSlash,
		TimeoutSlash:               timeoutSlash,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultCommissionChangeTime,
		DefaultLeavePoolTime,
		DefaultStakeFractionChangeTime,
		DefaultMultiCoinRefundPendingTime,
		DefaultVoteSlash,
		DefaultUploadSlash,
		DefaultTimeoutSlash,
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

	if err := util.ValidateNumber(p.StakeFractionChangeTime); err != nil {
		return err
	}

	if err := util.ValidateNumber(p.MultiCoinRefundPendingTime); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.VoteSlash); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.UploadSlash); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.TimeoutSlash); err != nil {
		return err
	}

	return nil
}
