package util

import (
	"fmt"

	"cosmossdk.io/math"
)

func ValidateDecimal(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid math.LegacyDec type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("invalid decimal: %s", v)
	}

	return nil
}

func ValidateInt(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid math.Int type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("invalid int: %s", v)
	}

	return nil
}

func ValidateNumber(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid uint64 type: %T", i)
	}

	if math.NewIntFromUint64(v).IsNil() || math.NewIntFromUint64(v).IsNegative() {
		return fmt.Errorf("invalid number: %d", v)
	}

	return nil
}

func ValidatePositiveNumber(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid uint64 type: %T", i)
	}

	if math.NewIntFromUint64(v).IsNil() ||
		math.NewIntFromUint64(v).IsNegative() ||
		math.NewIntFromUint64(v).IsZero() {
		return fmt.Errorf("invalid number: %d", v)
	}

	return nil
}

func ValidatePercentage(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid math.LegacyDec type: %T", i)
	}

	if v.IsNil() || v.IsNegative() || v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("invalid percentage: %s", v)
	}

	return nil
}
