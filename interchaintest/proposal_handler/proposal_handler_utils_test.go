package proposal_handler_test

import (
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	"github.com/KYVENetwork/chain/app"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/icza/dyno"
	. "github.com/onsi/gomega"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"strconv"
	"strings"
	"time"
)

const (
	uidGid         = "1025:1025"
	consensusSpeed = 2 * time.Second
	maxTxBytes     = 5_000
)

func encodingConfig() *sdktestutil.TestEncodingConfig {
	cfg := sdktestutil.TestEncodingConfig{}
	a := app.Setup()

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
			Type:                "cosmos",
			Name:                "kyve",
			ChainID:             "kyve-1",
			Bin:                 "kyved",
			Bech32Prefix:        "kyve",
			Denom:               "ukyve",
			GasPrices:           "0.02ukyve",
			GasAdjustment:       5,
			TrustingPeriod:      "112h",
			NoHostMount:         false,
			EncodingConfig:      encodingConfig(),
			ModifyGenesis:       ModifyGenesis,
			ConfigFileOverrides: configFileOverrides(),
			Images: []ibc.DockerImage{{
				Repository: "kyve",
				Version:    "local",
				UidGid:     uidGid,
			}},
		},
	}
}

func configFileOverrides() testutil.Toml {
	override := make(testutil.Toml)
	override["config/config.toml"] = testutil.Toml{
		"consensus": testutil.Toml{
			"timeout_propose":   consensusSpeed.String(),
			"timeout_prevote":   consensusSpeed.String(),
			"timeout_precommit": consensusSpeed.String(),
			"timeout_commit":    consensusSpeed.String(),
		},
	}
	return override
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

	_ = dyno.Set(genesis, math.LegacyMustNewDecFromStr("0.5"),
		"app_state", "global", "params", "min_initial_deposit_ratio",
	)

	_ = dyno.Set(genesis, "20s",
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

	// set the max tx bytes
	_ = dyno.Set(genesis, strconv.Itoa(maxTxBytes),
		"consensus", "params", "block", "max_bytes",
	)
	_ = dyno.Set(genesis, strconv.Itoa(maxTxBytes),
		"consensus", "params", "evidence", "max_bytes",
	)

	newGenesis, _ := json.Marshal(genesis)
	return newGenesis, nil
}

type TxData struct {
	Body struct {
		Messages []struct {
			Type    string `json:"@type"`
			Creator string `json:"creator"`
		} `json:"messages"`
	} `json:"body"`
}

func broadcastMsgs(ctx context.Context, broadcaster *cosmos.Broadcaster, broadcastingUser cosmos.User, msgs ...sdk.Msg) {
	f, err := broadcaster.GetFactory(ctx, broadcastingUser)
	ExpectWithOffset(1, err).To(BeNil())

	cc, err := broadcaster.GetClientContext(ctx, broadcastingUser)
	ExpectWithOffset(1, err).To(BeNil())

	err = clienttx.BroadcastTx(cc, f, msgs...)
	ExpectWithOffset(1, err).To(BeNil())
}

func duplicateMsg(msg sdk.Msg, size int) []sdk.Msg {
	var msgs []sdk.Msg
	for i := 0; i < size; i++ {
		msgs = append(msgs, msg)
	}
	return msgs
}

func checkTxsOrder(ctx context.Context, chain *cosmos.CosmosChain, height int64, expectedOrder []string) {
	txs, err := chain.FindTxs(ctx, height)
	ExpectWithOffset(1, err).To(BeNil())

	var order []string
	for _, tx := range txs {
		var result TxData
		err := json.Unmarshal(tx.Data, &result)
		ExpectWithOffset(1, err).To(BeNil())

		for _, msg := range result.Body.Messages {
			order = append(order, msg.Type[strings.LastIndex(msg.Type, ".")+1:])
		}
	}
	ExpectWithOffset(1, order).To(Equal(expectedOrder))
}
