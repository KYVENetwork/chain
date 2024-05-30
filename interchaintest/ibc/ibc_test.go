package relayer_test

import (
	"context"
	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"go.uber.org/zap/zaptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - IBC

* Transfer Kyve tokens to Osmosis

*/

func TestProposalHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "interchaintest/IBC Test Suite")
}

var _ = Describe("ibc", Ordered, func() {
	var kyve *cosmos.CosmosChain
	var osmosis *cosmos.CosmosChain

	var ctx context.Context
	var interchain *interchaintest.Interchain

	var ibcPath string
	var eRep *testreporter.RelayerExecReporter
	var rel ibc.Relayer

	BeforeAll(func() {
		ctx = context.Background()
		ibcPath = "kyve-osmosis"

		fw, err := interchaintest.CreateLogFile(GinkgoT().Name() + ".log")
		Expect(err).To(BeNil())

		rep := testreporter.NewReporter(fw)
		eRep = rep.RelayerExecReporter(GinkgoT())

		numFullNodes := 1
		numValidators := 2
		logger := zaptest.NewLogger(GinkgoT())
		factory := interchaintest.NewBuiltinChainFactory(
			logger,
			[]*interchaintest.ChainSpec{
				mainnetChainSpec(numValidators, numFullNodes),
				{Name: "osmosis", Version: "v25.0.0", NumValidators: &numValidators, NumFullNodes: &numFullNodes},
			},
		)

		chains, err := factory.Chains(GinkgoT().Name())
		Expect(err).To(BeNil())
		kyve, osmosis = chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

		client, network := interchaintest.DockerSetup(GinkgoT())

		rel = interchaintest.NewBuiltinRelayerFactory(
			// TODO: make the relayer work
			ibc.Hermes,
			logger,
			//relayer.CustomDockerImage("ghcr.io/informalsystems/hermes", "v1.8.0", rly.RlyDefaultUidGid),
			//relayer.CustomDockerImage("ghcr.io/cosmos/relayer", "v2.4.2", rly.RlyDefaultUidGid),
		).Build(GinkgoT(), client, network)

		interchain = interchaintest.NewInterchain().
			AddChain(kyve).
			AddChain(osmosis).
			AddRelayer(rel, "relayer").
			AddLink(interchaintest.InterchainLink{
				Chain1:  kyve,
				Chain2:  osmosis,
				Path:    ibcPath,
				Relayer: rel,
			})

		err = interchain.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
			TestName:          GinkgoT().Name(),
			Client:            client,
			NetworkID:         network,
			SkipPathCreation:  true,
			BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		})
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		_ = rel.StopRelayer(ctx, eRep)
		_ = kyve.StopAllNodes(ctx)
		_ = osmosis.StopAllNodes(ctx)
		_ = interchain.Close()
	})

	It("Transfer Kyve tokens to Osmosis", func() {
		// ARRANGE
		var wallet = interchaintest.GetAndFundTestUsers(
			GinkgoT(), ctx, GinkgoT().Name(), math.NewInt(10_000_000_000), kyve,
		)[0].(*cosmos.CosmosWallet)

		osmosisReceiver := wallet.FormattedAddressWithPrefix(osmosis.Config().Bech32Prefix)

		// Wait a few blocks
		err := testutil.WaitForBlocks(ctx, 3, kyve, osmosis)
		Expect(err).To(BeNil())

		err = rel.StartRelayer(ctx, eRep)
		Expect(err).To(BeNil())

		kyveChans, err := rel.GetChannels(ctx, eRep, kyve.Config().ChainID)
		Expect(err).To(BeNil())
		Expect(kyveChans).To(HaveLen(1))
		kyveChan := kyveChans[0]

		// ACT
		transfer := ibc.WalletAmount{
			Address: osmosisReceiver,
			Denom:   kyve.Config().Denom,
			Amount:  math.NewInt(1_000_000_000),
		}
		tx, err := kyve.SendIBCTransfer(ctx, ibcPath, wallet.KeyName(), transfer, ibc.TransferOptions{})
		Expect(err).To(BeNil())

		height, err := kyve.Height(ctx)
		Expect(err).To(BeNil())

		_, err = testutil.PollForAck(ctx, kyve, height, height+10, tx.Packet)
		Expect(err).To(BeNil())

		// ASSERT
		userBalance, err := kyve.GetBalance(ctx, wallet.FormattedAddress(), kyve.Config().Denom)
		Expect(err).To(BeNil())
		Expect(userBalance).To(Equal(math.NewInt(9_000_000_000)))

		prefixedDenom := transfertypes.GetPrefixedDenom(kyveChan.Counterparty.PortID, kyveChan.Counterparty.ChannelID, kyve.Config().Denom)
		denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
		ibcDenom := denomTrace.IBCDenom()

		receiverBalance, err := osmosis.GetBalance(ctx, osmosisReceiver, ibcDenom)
		Expect(err).To(BeNil())
		Expect(receiverBalance).To(Equal(math.NewInt(1_000_000_000)))
	})
})
