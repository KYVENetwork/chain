package global_test

import (
	"context"
	"encoding/json"
	"strconv"

	"cosmossdk.io/math"

	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icza/dyno"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"go.uber.org/zap/zaptest"
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

func mainnetChainSpec(numValidators int, numFullNodes int, burnRatio math.LegacyDec) *interchaintest.ChainSpec {
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
			GasPrices:      "1ukyve",
			GasAdjustment:  10,
			TrustingPeriod: "112h",
			NoHostMount:    false,
			EncodingConfig: encodingConfig(),
			ModifyGenesis:  modifyGenesis(burnRatio),
			Images: []ibc.DockerImage{{
				Repository: "kyve",
				Version:    "local",
				UidGid:     uidGid,
			}},
		},
	}
}

func modifyGenesis(burnRatio math.LegacyDec) func(config ibc.ChainConfig, genbz []byte) ([]byte, error) {
	return func(config ibc.ChainConfig, genbz []byte) ([]byte, error) {
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

		// Update total supply
		coin := totalSupply[0].(map[string]interface{})
		amountStr := coin["amount"].(string)
		amount, _ := strconv.Atoi(amountStr)
		totalSupply[0] = sdk.NewCoin(config.Denom, math.NewInt(int64(amount)+teamSupply.Int64()))
		_ = dyno.Set(genesis, totalSupply, "app_state", "bank", "supply")

		// Set burn ratio to given value
		_ = dyno.Set(genesis, burnRatio,
			"app_state", "global", "params", "burn_ratio",
		)

		// Set inflation to zero
		_ = dyno.Set(genesis, math.LegacyZeroDec(),
			"app_state", "mint", "params", "inflation_min",
		)
		_ = dyno.Set(genesis, math.LegacyZeroDec(),
			"app_state", "mint", "params", "inflation_max",
		)

		newGenesis, _ := json.Marshal(genesis)
		return newGenesis, nil
	}
}

// startNewChainWithCustomBurnRatio starts a new chain with the given burn ratio and an inflation of zero.
func startNewChainWithCustomBurnRatio(ctx context.Context, burnRatio math.LegacyDec) (*cosmos.CosmosChain, *interchaintest.Interchain, *cosmos.Broadcaster, *cosmos.CosmosWallet) {
	numFullNodes := 0
	numValidators := 2
	factory := interchaintest.NewBuiltinChainFactory(
		zaptest.NewLogger(GinkgoT()),
		[]*interchaintest.ChainSpec{mainnetChainSpec(numValidators, numFullNodes, burnRatio)},
	)

	chains, err := factory.Chains(GinkgoT().Name())
	Expect(err).To(BeNil())
	chain := chains[0].(*cosmos.CosmosChain)

	interchain := interchaintest.NewInterchain().AddChain(chain)

	broadcaster := cosmos.NewBroadcaster(GinkgoT(), chain)

	dockerClient, network := interchaintest.DockerSetup(GinkgoT())

	err = interchain.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		TestName:         GinkgoT().Name(),
		Client:           dockerClient,
		NetworkID:        network,
		SkipPathCreation: true,
	})
	Expect(err).To(BeNil())

	wallet := interchaintest.GetAndFundTestUsers(GinkgoT(), ctx, GinkgoT().Name(), math.NewInt(10*i.T_KYVE), chain)[0].(*cosmos.CosmosWallet)
	return chain, interchain, broadcaster, wallet
}
