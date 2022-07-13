package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyVoteSlash            = []byte("VoteSlash")
	DefaultVoteSlash string = "0.1"
)

var (
	KeyUploadSlash            = []byte("UploadSlash")
	DefaultUploadSlash string = "0.2"
)

var (
	KeyTimeoutSlash            = []byte("TimeoutSlash")
	DefaultTimeoutSlash string = "0.02"
)

var (
	KeyUploadTimeout            = []byte("UploadTimeout")
	DefaultUploadTimeout uint64 = 600
)

var (
	KeyStorageCost            = []byte("StorageCost")
	DefaultStorageCost uint64 = 100000
)

var (
	KeyNetworkFee            = []byte("NetworkFee")
	DefaultNetworkFee string = "0.01"
)

var (
	KeyMaxPoints            = []byte("MaxPoints")
	DefaultMaxPoints uint64 = 5
)

var (
	KeyUnbondingStakingTime            = []byte("UnbondingStakingTime")
	DefaultUnbondingStakingTime uint64 = 60 * 60 * 24 * 5
)

var (
	KeyUnbondingDelegationTime            = []byte("UnbondingDelegationTime")
	DefaultUnbondingDelegationTime uint64 = 60 * 60 * 24 * 5
)

var (
	KeyRedelegationCooldown            = []byte("RedelegationCooldown")
	DefaultRedelegationCooldown uint64 = 60 * 60 * 24 * 5
)

var (
	KeyRedelegationMaxAmount            = []byte("KeyRedelegationMaxAmount")
	DefaultRedelegationMaxAmount uint64 = 5
)

var (
	KeyCommissionChangeTime            = []byte("KeyCommissionChangeTime")
	DefaultCommissionChangeTime uint64 = 60 * 60 * 24 * 5
)

// ParamKeyTable the param Key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	voteSlash string,
	uploadSlash string,
	timeoutSlash string,
	uploadTimeout uint64,
	storageCost uint64,
	networkFee string,
	maxPoints uint64,
	unbondingStakingTime uint64,
	unbondingDelegationTime uint64,
	redelegationCooldown uint64,
	redelegationMaxAmount uint64,
	commissionChangeTime uint64,
) Params {
	return Params{
		VoteSlash:               voteSlash,
		UploadSlash:             uploadSlash,
		TimeoutSlash:            timeoutSlash,
		UploadTimeout:           uploadTimeout,
		StorageCost:             storageCost,
		NetworkFee:              networkFee,
		MaxPoints:               maxPoints,
		UnbondingStakingTime:    unbondingStakingTime,
		UnbondingDelegationTime: unbondingDelegationTime,
		RedelegationCooldown:    redelegationCooldown,
		RedelegationMaxAmount:   redelegationMaxAmount,
		CommissionChangeTime:    commissionChangeTime,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultVoteSlash,
		DefaultUploadSlash,
		DefaultTimeoutSlash,
		DefaultUploadTimeout,
		DefaultStorageCost,
		DefaultNetworkFee,
		DefaultMaxPoints,
		DefaultUnbondingStakingTime,
		DefaultUnbondingDelegationTime,
		DefaultRedelegationCooldown,
		DefaultRedelegationMaxAmount,
		DefaultCommissionChangeTime,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVoteSlash, &p.VoteSlash, validateVoteSlash),
		paramtypes.NewParamSetPair(KeyUploadSlash, &p.UploadSlash, validateUploadSlash),
		paramtypes.NewParamSetPair(KeyTimeoutSlash, &p.TimeoutSlash, validateTimeoutSlash),
		paramtypes.NewParamSetPair(KeyUploadTimeout, &p.UploadTimeout, validateUploadTimeout),
		paramtypes.NewParamSetPair(KeyStorageCost, &p.StorageCost, validateStorageCost),
		paramtypes.NewParamSetPair(KeyNetworkFee, &p.NetworkFee, validateNetworkFee),
		paramtypes.NewParamSetPair(KeyMaxPoints, &p.MaxPoints, validateMaxPoints),
		paramtypes.NewParamSetPair(KeyUnbondingStakingTime, &p.UnbondingStakingTime, validateUnbondingStakingTime),
		paramtypes.NewParamSetPair(KeyUnbondingDelegationTime, &p.UnbondingDelegationTime, validateUnbondingDelegationTime),
		paramtypes.NewParamSetPair(KeyRedelegationCooldown, &p.RedelegationCooldown, validateTrue),
		paramtypes.NewParamSetPair(KeyRedelegationMaxAmount, &p.RedelegationMaxAmount, validateTrue),
		paramtypes.NewParamSetPair(KeyCommissionChangeTime, &p.CommissionChangeTime, validateTrue),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateVoteSlash(p.VoteSlash); err != nil {
		return err
	}

	if err := validateUploadSlash(p.UploadSlash); err != nil {
		return err
	}

	if err := validateTimeoutSlash(p.TimeoutSlash); err != nil {
		return err
	}

	if err := validateUploadTimeout(p.UploadTimeout); err != nil {
		return err
	}

	if err := validateStorageCost(p.StorageCost); err != nil {
		return err
	}

	if err := validateNetworkFee(p.NetworkFee); err != nil {
		return err
	}

	if err := validateMaxPoints(p.MaxPoints); err != nil {
		return err
	}

	if err := validateUnbondingStakingTime(p.UnbondingStakingTime); err != nil {
		return err
	}

	if err := validateUnbondingDelegationTime(p.UnbondingDelegationTime); err != nil {
		return err
	}

	return nil
}

func validateTrue(v interface{}) error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateVoteSlash validates the VoteSlash param
func validateVoteSlash(v interface{}) error {
	return validatePercentage(v)
}

// validateUploadSlash validates the UploadSlash param
func validateUploadSlash(v interface{}) error {
	return validatePercentage(v)
}

// validateTimeoutSlash validates the TimeoutSlash param
func validateTimeoutSlash(v interface{}) error {
	return validatePercentage(v)
}

// validateUploadTimeout validates the uploadTimeout param
func validateUploadTimeout(v interface{}) error {
	uploadTimeout, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = uploadTimeout

	return nil
}

// validateStorageCost validates the StorageCost param
func validateStorageCost(v interface{}) error {
	storageCost, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = storageCost

	return nil
}

// validateNetworkFee validates the NetworkFee param
func validateNetworkFee(v interface{}) error {
	return validatePercentage(v)
}

// validateMaxPoints validates the MaxPoints param
func validateMaxPoints(v interface{}) error {
	maxPoints, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = maxPoints

	return nil
}

// validateMaxPoints validates the MaxPoints param
func validateUnbondingStakingTime(v interface{}) error {
	_, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	return nil
}

// validateMaxPoints validates the MaxPoints param
func validateUnbondingDelegationTime(v interface{}) error {
	_, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	return nil
}

// validatePercentage ...
func validatePercentage(v interface{}) error {
	val, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	parsedVal, err := sdk.NewDecFromStr(val)
	if err != nil {
		return fmt.Errorf("invalid decimal representation: %T", v)
	}

	if parsedVal.LT(sdk.NewDec(0)) {
		return fmt.Errorf("percentage should be greater than or equal to 0")
	}
	if parsedVal.GT(sdk.NewDec(1)) {
		return fmt.Errorf("percentage should be less than or equal to 1")
	}

	return nil
}
