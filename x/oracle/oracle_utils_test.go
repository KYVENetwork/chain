package oracle_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Bank
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func InjectTeamBalance(bankModuleBytes []byte) []byte {
	var bankModuleState bankTypes.GenesisState
	bankTypes.ModuleCdc.MustUnmarshalJSON(bankModuleBytes, &bankModuleState)

	teamBalance := bankTypes.Balance{
		// Address: authTypes.NewModuleAddress(teamTypes.ModuleName).String(),
		Address: "kyve1e29j95xmsw3zmvtrk4st8e89z5n72v7nf70ma4",
		Coins: sdk.NewCoins(sdk.NewCoin(
			"ukyve", sdk.NewIntFromUint64(165_000_000_000_000),
		)),
	}
	bankModuleState.Balances = append(bankModuleState.Balances, teamBalance)

	newBankModuleBytes := bankTypes.ModuleCdc.MustMarshalJSON(&bankModuleState)
	return newBankModuleBytes
}
