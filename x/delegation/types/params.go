package types

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
)

// DefaultUnbondingDelegationTime ...
var DefaultUnbondingDelegationTime = uint64(60 * 60 * 24 * 5)

// DefaultRedelegationCooldown ...
var DefaultRedelegationCooldown = uint64(60 * 60 * 24 * 5)

// DefaultRedelegationMaxAmount ...
var DefaultRedelegationMaxAmount = uint64(5)

// DefaultVoteSlash ...
var DefaultVoteSlash = math.LegacyMustNewDecFromStr("0.1")

// DefaultUploadSlash ...
var DefaultUploadSlash = math.LegacyMustNewDecFromStr("0.2")

// DefaultTimeoutSlash ...
var DefaultTimeoutSlash = math.LegacyMustNewDecFromStr("0.02")

// NewParams creates a new Params instance
func NewParams(
	unbondingDelegationTime uint64,
	redelegationCooldown uint64,
	redelegationMaxAmount uint64,
	voteSlash math.LegacyDec,
	uploadSlash math.LegacyDec,
	timeoutSlash math.LegacyDec,
) Params {
	return Params{
		UnbondingDelegationTime: unbondingDelegationTime,
		RedelegationCooldown:    redelegationCooldown,
		RedelegationMaxAmount:   redelegationMaxAmount,
		VoteSlash:               voteSlash,
		UploadSlash:             uploadSlash,
		TimeoutSlash:            timeoutSlash,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingDelegationTime,
		DefaultRedelegationCooldown,
		DefaultRedelegationMaxAmount,
		DefaultVoteSlash,
		DefaultUploadSlash,
		DefaultTimeoutSlash,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidateNumber(p.UnbondingDelegationTime); err != nil {
		return err
	}

	if err := util.ValidateNumber(p.RedelegationCooldown); err != nil {
		return err
	}

	if err := util.ValidateNumber(p.RedelegationMaxAmount); err != nil {
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
