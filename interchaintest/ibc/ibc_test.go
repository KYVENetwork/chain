package relayer_test

import (
	"context"
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/testutil/integration"
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

* Transfer 1 $KYVE to Osmosis and back

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

		rep := testreporter.NewNopReporter()
		eRep = rep.RelayerExecReporter(GinkgoT())

		numFullNodes := 0
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

		rel = interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, logger).Build(GinkgoT(), client, network)

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
			TestName:         GinkgoT().Name(),
			Client:           client,
			NetworkID:        network,
			SkipPathCreation: false,
		})
		Expect(err).To(BeNil())
	})

	AfterAll(func() {
		_ = rel.StopRelayer(ctx, eRep)
		_ = kyve.StopAllNodes(ctx)
		_ = osmosis.StopAllNodes(ctx)
		_ = interchain.Close()
	})

	It("Transfer 1 $KYVE to Osmosis and back", func() {
		// ARRANGE
		startBalance := math.NewInt(10 * integration.T_KYVE)
		var kyveWallet = interchaintest.GetAndFundTestUsers(GinkgoT(), ctx, GinkgoT().Name(), startBalance, kyve)[0].(*cosmos.CosmosWallet)
		var osmosisWallet = interchaintest.GetAndFundTestUsers(GinkgoT(), ctx, GinkgoT().Name(), startBalance, osmosis)[0].(*cosmos.CosmosWallet)

		kyveChans, err := rel.GetChannels(ctx, eRep, kyve.Config().ChainID)
		Expect(err).To(BeNil())
		Expect(kyveChans).To(HaveLen(1))
		kyveChan := kyveChans[0]

		prefixedDenom := transfertypes.GetPrefixedDenom(kyveChan.Counterparty.PortID, kyveChan.Counterparty.ChannelID, kyve.Config().Denom)
		denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
		ibcDenom := denomTrace.IBCDenom()

		err = rel.StartRelayer(ctx, eRep)
		Expect(err).To(BeNil())

		// ACT 1
		// Transfer 1 $KYVE to Osmosis
		transferAmount := math.NewInt(integration.T_KYVE)
		transfer := ibc.WalletAmount{
			Address: osmosisWallet.FormattedAddress(),
			Denom:   kyve.Config().Denom,
			Amount:  transferAmount,
		}
		tx, err := kyve.SendIBCTransfer(ctx, kyveChan.ChannelID, kyveWallet.KeyName(), transfer, ibc.TransferOptions{})
		Expect(err).To(BeNil())

		height, err := kyve.Height(ctx)
		Expect(err).To(BeNil())

		_, err = testutil.PollForAck(ctx, kyve, height, height+10, tx.Packet)
		Expect(err).To(BeNil())

		// ASSERT 1
		// New balance must be startBalance - transferAmount - fees -> so a bit less than 9 $KYVE
		kyveBal1, err := kyve.GetBalance(ctx, kyveWallet.FormattedAddress(), kyve.Config().Denom)
		Expect(err).To(BeNil())
		newBalance := startBalance.Sub(transferAmount)
		Expect(kyveBal1.Int64()).To(BeBetween(newBalance.Sub(math.NewInt(100_000)).Int64(), newBalance.Int64()))

		osmosisBal1, err := osmosis.GetBalance(ctx, osmosisWallet.FormattedAddress(), ibcDenom)
		Expect(err).To(BeNil())
		Expect(osmosisBal1).To(Equal(math.NewInt(integration.T_KYVE)))

		// ACT 2
		// Transfer 1 $KYVE back to Kyve
		transfer = ibc.WalletAmount{
			Address: kyveWallet.FormattedAddress(),
			Denom:   ibcDenom,
			Amount:  transferAmount,
		}
		tx, err = osmosis.SendIBCTransfer(ctx, kyveChan.Counterparty.ChannelID, osmosisWallet.KeyName(), transfer, ibc.TransferOptions{})
		Expect(err).To(BeNil())

		height, err = osmosis.Height(ctx)
		Expect(err).To(BeNil())

		_, err = testutil.PollForAck(ctx, osmosis, height, height+10, tx.Packet)
		Expect(err).To(BeNil())

		// ASSERT 2
		osmosisBal2, err := osmosis.GetBalance(ctx, osmosisWallet.FormattedAddress(), ibcDenom)
		Expect(err).To(BeNil())
		Expect(osmosisBal2).To(Equal(math.NewInt(0)))

		kyveBal2, err := kyve.GetBalance(ctx, kyveWallet.FormattedAddress(), kyve.Config().Denom)
		Expect(err).To(BeNil())
		Expect(kyveBal2).To(Equal(kyveBal1.Add(transferAmount)))
	})
})
