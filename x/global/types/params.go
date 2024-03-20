package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// DefaultMinGasPrice is 0 (i.e. disabled)
var DefaultMinGasPrice = math.LegacyNewDec(0)

// DefaultBurnRatio is 0% (i.e. disabled)
var DefaultBurnRatio = math.LegacyNewDec(0)

// DefaultMinInitialDepositRatio is 0% (i.e. disabled)
var DefaultMinInitialDepositRatio = math.LegacyNewDec(0)

// NewParams creates a new Params instance
func NewParams(minGasPrice math.LegacyDec, burnRatio math.LegacyDec, gasAdjustments []GasAdjustment, gasRefunds []GasRefund, minInitialDepositRatio math.LegacyDec) Params {
	return Params{
		MinGasPrice:            minGasPrice,
		BurnRatio:              burnRatio,
		GasAdjustments:         gasAdjustments,
		GasRefunds:             gasRefunds,
		MinInitialDepositRatio: minInitialDepositRatio,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultMinGasPrice, DefaultBurnRatio, []GasAdjustment{}, []GasRefund{}, DefaultMinInitialDepositRatio)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMinGasPrice(p.MinGasPrice); err != nil {
		return err
	}

	if err := validateBurnRatio(p.BurnRatio); err != nil {
		return err
	}

	for _, gasAdjustment := range p.GasAdjustments {
		if err := validateGasAdjustment(gasAdjustment); err != nil {
			return err
		}
	}

	for _, gasRefund := range p.GasRefunds {
		if err := validateGasRefund(gasRefund); err != nil {
			return err
		}
	}

	if err := validateMinInitialDepositRatio(p.MinInitialDepositRatio); err != nil {
		return err
	}

	return nil
}

// validateMinGasPrice ...
func validateMinGasPrice(i interface{}) error {
	v, ok := i.(math.LegacyDec)

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

// validateBurnRatio ...
func validateBurnRatio(i interface{}) error {
	v, ok := i.(math.LegacyDec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", i)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("value cannot be greater than 1: %s", v)
	}

	return nil
}

// validateGasAdjustment ...
func validateGasAdjustment(i interface{}) error {
	v, ok := i.(GasAdjustment)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	amount := math.NewInt(int64(v.Amount))

	if amount.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if amount.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", amount)
	}

	return nil
}

// validateGasRefund ...
func validateGasRefund(i interface{}) error {
	v, ok := i.(GasRefund)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Fraction.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.Fraction.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", v.Fraction)
	}

	if v.Fraction.GT(math.LegacyOneDec()) {
		return fmt.Errorf("value cannot be greater than 1: %s", v.Fraction)
	}

	return nil
}

// validateMinInitialDepositRatio ...
func validateMinInitialDepositRatio(i interface{}) error {
	v, ok := i.(math.LegacyDec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", i)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("value cannot be greater than 1: %s", v)
	}

	return nil
}
