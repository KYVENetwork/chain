package relayer_test

import (
	"encoding/json"
	"strconv"

	i "github.com/KYVENetwork/chain/testutil/integration"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icza/dyno"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

const (
	uidGid = "1025:1025"
)

func encodingConfig() *sdktestutil.TestEncodingConfig {
	cfg := sdktestutil.TestEncodingConfig{}
	a := i.NewCleanChain().App()

	cfg.Codec = a.AppCodec()
	cfg.TxConfig = authtx.NewTxConfig(a.AppCodec(), authtx.DefaultSignModes)
	cfg.InterfaceRegistry = a.InterfaceRegistry()
	cfg.Amino = a.LegacyAmino()

	return &cfg
}

func mainnetChainSpec(numValidators int, numFullNodes int) *interchaintest.ChainSpec {
	return &interchaintest.ChainSpec{
		NumValidators: &numValidators,
		NumFullNodes:  &numFullNodes,
		ChainConfig: ibc.ChainConfig{
			Type:           "cosmos",
			Name:           "kyve",
			ChainID:        "kyve-1",
			Bin:            "kyved",
			Bech32Prefix:   "kyve",
			Denom:          "ukyve",
			GasPrices:      "0.02ukyve",
			GasAdjustment:  5,
			TrustingPeriod: "112h",
			NoHostMount:    false,
			EncodingConfig: encodingConfig(),
			ModifyGenesis:  ModifyGenesis,
			Images: []ibc.DockerImage{{
				Repository: "kyve",
				Version:    "local",
				UidGid:     uidGid,
			}},
		},
	}
}

func ModifyGenesis(config ibc.ChainConfig, genbz []byte) ([]byte, error) {
	genesis := make(map[string]interface{})
	_ = json.Unmarshal(genbz, &genesis)

	teamSupply := math.NewInt(165_000_000_000_000)
	balances, _ := dyno.GetSlice(genesis, "app_state", "bank", "balances")
	balances = append(balances, bankTypes.Balance{
		Address: "kyve1e29j95xmsw3zmvtrk4st8e89z5n72v7nf70ma4",
		Coins:   sdk.NewCoins(sdk.NewCoin(config.Denom, teamSupply)),
	})
	_ = dyno.Set(genesis, balances, "app_state", "bank", "balances")
	totalSupply, _ := dyno.GetSlice(genesis, "app_state", "bank", "supply")

	// update total supply
	coin := totalSupply[0].(map[string]interface{})
	amountStr := coin["amount"].(string)
	amount, _ := strconv.Atoi(amountStr)
	totalSupply[0] = sdk.NewCoin(config.Denom, math.NewInt(int64(amount)+teamSupply.Int64()))
	_ = dyno.Set(genesis, totalSupply, "app_state", "bank", "supply")

	newGenesis, _ := json.Marshal(genesis)
	return newGenesis, nil
}

func BeBetween(min, max interface{}) types.GomegaMatcher {
	return SatisfyAll(
		BeNumerically(">", min),
		BeNumerically("<", max))
}
