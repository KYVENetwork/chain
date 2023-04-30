package v1_2_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KYVENetwork/chain/app/upgrades/v1_2"
	"github.com/docker/docker/client"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"go.uber.org/zap/zaptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var UpgradeContainerVersion = "local"

func TestV1P2Upgrade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("%s Upgrade Test Suite", v1_2.UpgradeName))
}

var _ = Describe(fmt.Sprintf("%s Upgrade Tests", v1_2.UpgradeName), Ordered, func() {
	var kaon *cosmos.CosmosChain
	var kyve *cosmos.CosmosChain

	var ctx context.Context
	var client *client.Client
	var network string
	var interchain *interchaintest.Interchain

	var kaonWallet *cosmos.CosmosWallet
	var kyveWallet *cosmos.CosmosWallet

	BeforeAll(func() {
		factory := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(GinkgoT()), []*interchaintest.ChainSpec{
			{
				Name:        "kaon",
				ChainConfig: testnetConfig,
			},
			{
				Name:        "kyve",
				ChainConfig: mainnetConfig,
			},
		})

		chains, err := factory.Chains(GinkgoT().Name())
		Expect(err).To(BeNil())
		kaon, kyve = chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

		interchain = interchaintest.NewInterchain().
			AddChain(kaon).
			AddChain(kyve)

		ctx = context.Background()
		client, network = interchaintest.DockerSetup(GinkgoT())

		err = interchain.Build(ctx, nil, interchaintest.InterchainBuildOptions{
			TestName:         GinkgoT().Name(),
			Client:           client,
			NetworkID:        network,
			SkipPathCreation: true,
		})
		Expect(err).To(BeNil())

		wallets := interchaintest.GetAndFundTestUsers(
			GinkgoT(), ctx, GinkgoT().Name(), 10_000_000_000, kaon, kyve,
		)
		kaonWallet, kyveWallet = wallets[0].(*cosmos.CosmosWallet), wallets[1].(*cosmos.CosmosWallet)
	})

	AfterAll(func() {
		_ = interchain.Close()
	})

	It("", func() {
		PerformUpgrade(ctx, client, kaon, kaonWallet, 10, "kaon")
	})

	It("", func() {
		PerformUpgrade(ctx, client, kyve, kyveWallet, 10, "kyve")
	})

	// TODO: Figure out why concurrency doesn't work.
	PIt("", func() {
		// KAON
		go PerformUpgrade(ctx, client, kaon, kaonWallet, 10, "kaon")

		// KYVE
		go PerformUpgrade(ctx, client, kyve, kyveWallet, 10, "kyve")
	})
})

func PerformUpgrade(
	ctx context.Context,
	client *client.Client,
	chain *cosmos.CosmosChain,
	wallet *cosmos.CosmosWallet,
	delta uint64,
	container string,
) {
	height, _ := chain.Height(ctx)
	haltHeight := height + delta

	proposal := cosmos.SoftwareUpgradeProposal{
		Deposit:     "1" + chain.Config().Denom,
		Title:       fmt.Sprintf("%s Software Upgrade", v1_2.UpgradeName),
		Name:        v1_2.UpgradeName,
		Description: "description",
		Height:      haltHeight,
	}

	upgrade, proposalErr := chain.UpgradeProposal(ctx, wallet.KeyName(), proposal)
	Expect(proposalErr).To(BeNil())
	voteErr := chain.VoteOnProposalAllValidators(ctx, upgrade.ProposalID, cosmos.ProposalVoteYes)
	Expect(voteErr).To(BeNil())

	_, statusErr := cosmos.PollForProposalStatus(ctx, chain, height, haltHeight, upgrade.ProposalID, cosmos.ProposalStatusPassed)
	Expect(statusErr).To(BeNil())

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	height, _ = chain.Height(ctx)
	_ = testutil.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, chain)

	height, _ = chain.Height(ctx)
	Expect(height).To(Equal(haltHeight))

	stopErr := chain.StopAllNodes(ctx)
	Expect(stopErr).To(BeNil())
	chain.UpgradeVersion(ctx, client, container, UpgradeContainerVersion)
	startErr := chain.StartAllNodes(ctx)
	Expect(startErr).To(BeNil())

	waitErr := testutil.WaitForBlocks(ctx, int(delta), chain)
	Expect(waitErr).To(BeNil())
}
