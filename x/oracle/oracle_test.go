package oracle_test

import (
	"context"
	"encoding/json"
	"fmt"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/strangelove-ventures/interchaintest/v6/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
	"github.com/strangelove-ventures/interchaintest/v6/testreporter"
	"github.com/strangelove-ventures/interchaintest/v6/testutil"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v6"
	"go.uber.org/zap/zaptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOracle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oracle Test Suite")
}

var config = ibc.ChainConfig{
	Type:    "cosmos",
	Name:    "kyve",
	ChainID: "kyve-1",
	Images: []ibc.DockerImage{{
		Repository: "kyve",
		Version:    "local",
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

func ModifyGenesis(_ ibc.ChainConfig, genbz []byte) ([]byte, error) {
	genesis := map[string]interface{}{}
	_ = json.Unmarshal(genbz, &genesis)
	appState := genesis["app_state"].(map[string]interface{})

	bankModuleBytes, _ := json.Marshal(appState[bankTypes.ModuleName])
	newBankModuleBytes := InjectTeamBalance(bankModuleBytes)

	bankModuleState := map[string]interface{}{}
	_ = json.Unmarshal(newBankModuleBytes, &bankModuleState)

	appState[bankTypes.ModuleName] = bankModuleState
	genesis["app_state"] = appState

	newGenesis, _ := json.MarshalIndent(genesis, "", "  ")
	return newGenesis, nil
}

var _ = Describe("Oracle Tests", Ordered, func() {
	var ctx context.Context
	var interchain *interchaintest.Interchain
	var relayerReporter *testreporter.RelayerExecReporter
	var relayer ibc.Relayer

	BeforeEach(func() {
		chainFactory := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(GinkgoT()), []*interchaintest.ChainSpec{
			{
				Name:    "juno",
				Version: "v14.0.0",
			},
			{
				Name:        "kyve",
				ChainConfig: config,
			},
		})

		chains, err := chainFactory.Chains(GinkgoT().Name())
		Expect(err).To(BeNil())
		juno, kyve := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

		client, network := interchaintest.DockerSetup(GinkgoT())
		relayerFactory := interchaintest.NewBuiltinRelayerFactory(ibc.Hermes, zaptest.NewLogger(GinkgoT()))
		relayer = relayerFactory.Build(GinkgoT(), client, network)

		interchain = interchaintest.NewInterchain().
			AddChain(juno).
			AddChain(kyve).
			AddRelayer(relayer, "rly").
			AddLink(interchaintest.InterchainLink{
				Chain1:  kyve,
				Chain2:  juno,
				Relayer: relayer,
				Path:    "kyve-juno",
			})

		ctx = context.Background()
		reporter := testreporter.NewNopReporter()
		relayerReporter = reporter.RelayerExecReporter(GinkgoT())
		err = interchain.Build(
			ctx,
			relayerReporter,
			interchaintest.InterchainBuildOptions{
				TestName:          GinkgoT().Name(),
				Client:            client,
				NetworkID:         network,
				SkipPathCreation:  false,
				BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
			},
		)
		Expect(err).To(BeNil())

		_ = relayer.StartRelayer(ctx, relayerReporter, "kyve-juno")

		users := interchaintest.GetAndFundTestUsers(GinkgoT(), ctx, GinkgoT().Name(), 100_000_000, juno, kyve)

		junoUser, kyveUser := users[0].(*cosmos.CosmosWallet), users[1].(*cosmos.CosmosWallet)
		fmt.Println(junoUser, kyveUser)
		_ = testutil.WaitForBlocks(ctx, 5, juno, kyve)

		fmt.Println(juno.GetBalance(ctx, junoUser.FormattedAddress(), "ujuno"))
		fmt.Println(kyve.GetBalance(ctx, kyveUser.FormattedAddress(), "ukyve"))
	})

	AfterEach(func() {
		_ = relayer.StopRelayer(ctx, relayerReporter)
		_ = interchain.Close()
	})

	It("", func() {
		Expect(true).To(BeTrue())
	})
})
