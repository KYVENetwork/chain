package util

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateUint64(v interface{}) error {
	_, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	return nil
}

func ValidatePercentage(v interface{}) error {
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
