package util

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateDecimal(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("invalid decimal: %s", v)
	}

	return nil
}

func ValidateNumber(i uint64) error {
	v := math.NewIntFromUint64(i)

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("invalid number: %s", v)
	}

	return nil
}

func ValidatePercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid type: %T", i)
	}

	if v.IsNil() || v.IsNegative() || v.GT(sdk.OneDec()) {
		return fmt.Errorf("invalid percentage: %s", v)
	}

	return nil
}
