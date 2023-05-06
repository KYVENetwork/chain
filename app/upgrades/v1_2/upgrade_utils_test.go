package v1_2_test

import (
	"encoding/json"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icza/dyno"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

var testnetConfig = ibc.ChainConfig{
	Type:    "cosmos",
	Name:    "kaon",
	ChainID: "kaon-1",
	Images: []ibc.DockerImage{{
		Repository: "ghcr.io/kyvenetwork/chain/kaon",
		Version:    "v1.1.0",
		UidGid:     "1025:1025",
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

var mainnetConfig = ibc.ChainConfig{
	Type:    "cosmos",
	Name:    "kyve",
	ChainID: "kyve-1",
	Images: []ibc.DockerImage{{
		Repository: "ghcr.io/kyvenetwork/chain/kyve",
		Version:    "v1.1.0",
		UidGid:     "1025:1025",
	}},
	Bin:                 "kyved",
	Bech32Prefix:        "kyve",
	Denom:               "ukyve",
	GasPrices:           "0.02ukyve",
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
		"app_state", "gov", "voting_params", "voting_period",
	)
	_ = dyno.Set(genesis, "0",
		"app_state", "gov", "deposit_params", "min_deposit", 0, "amount",
	)
	_ = dyno.Set(genesis, config.Denom,
		"app_state", "gov", "deposit_params", "min_deposit", 0, "denom",
	)

	newGenesis, _ := json.Marshal(genesis)
	return newGenesis, nil
}
