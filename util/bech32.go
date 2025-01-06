package util

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MustValaddressFromOperatorAddress(operatorAddress string) string {
	rawAddressBytes, err := sdk.GetFromBech32(operatorAddress, "kyve")
	if err != nil {
		panic(err)
	}

	valAddress, err := sdk.Bech32ifyAddressBytes("kyvevaloper", rawAddressBytes)
	if err != nil {
		panic(err)
	}

	return valAddress
}

func MustAccountAddressFromValAddress(valAddress string) string {
	rawAddressBytes, err := sdk.GetFromBech32(valAddress, "kyvevaloper")
	if err != nil {
		panic(err)
	}

	accAddress, err := sdk.Bech32ifyAddressBytes("kyve", rawAddressBytes)
	if err != nil {
		panic(err)
	}

	return accAddress
}
