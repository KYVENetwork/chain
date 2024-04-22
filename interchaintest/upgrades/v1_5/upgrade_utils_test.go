package v1_5_test

import (
	"encoding/json"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icza/dyno"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

const (
	previousVersion = "v1.4.0"
	uidGid          = "1025:1025"
)

var testnetConfig = ibc.ChainConfig{
	Type:    "cosmos",
	Name:    "kaon",
	ChainID: "kaon-1",
	Images: []ibc.DockerImage{{
		Repository: "ghcr.io/strangelove-ventures/heighliner/kaon",
		Version:    previousVersion,
		UidGid:     uidGid,
	}},
	Bin:                 "kyved",
	Bech32Prefix:        "kyve",
	Denom:               "tkyve",
	GasPrices:           "0.02tkyve",
	GasAdjustment:       5,
	TrustingPeriod:      "112h",
	NoHostMount:         false,
	ModifyGenesis:       ModifyGenesis,
	ConfigFileOverrides: nil,
	EncodingConfig:      nil,
}

func ModifyGenesis(config ibc.ChainConfig, genbz []byte) ([]byte, error) {
	genesis := make(map[string]interface{})
	_ = json.Unmarshal(genbz, &genesis)

	balances, _ := dyno.GetSlice(genesis, "app_state", "bank", "balances")
	balances = append(balances, bankTypes.Balance{
		Address: "kyve1e29j95xmsw3zmvtrk4st8e89z5n72v7nf70ma4",
		Coins:   sdk.NewCoins(sdk.NewCoin(config.Denom, math.NewInt(165_000_000_000_000))),
	})
	_ = dyno.Set(genesis, balances, "app_state", "bank", "balances")

	switch config.ChainID {
	case "kaon-1":
		_ = dyno.Set(genesis, math.LegacyMustNewDecFromStr("0.25"),
			"app_state", "global", "params", "min_initial_deposit_ratio",
		)
	case "kyve-1":
		_ = dyno.Set(genesis, math.LegacyMustNewDecFromStr("0.5"),
			"app_state", "global", "params", "min_initial_deposit_ratio",
		)
	}

	_ = dyno.Set(genesis, "10s",
		"app_state", "gov", "params", "voting_period",
	)
	_ = dyno.Set(genesis, "0",
		"app_state", "gov", "params", "min_deposit", 0, "amount",
	)
	_ = dyno.Set(genesis, config.Denom,
		"app_state", "gov", "params", "min_deposit", 0, "denom",
	)

	_ = dyno.Set(genesis, "0.169600000000000000",
		"app_state", "pool", "params", "protocol_inflation_share",
	)
	_ = dyno.Set(genesis, "0.050000000000000000",
		"app_state", "pool", "params", "pool_inflation_payout_rate",
	)

	newGenesis, _ := json.Marshal(genesis)
	return newGenesis, nil
}
