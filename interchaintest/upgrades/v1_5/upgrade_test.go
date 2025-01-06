package v2_0_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"cosmossdk.io/math"

	"github.com/KYVENetwork/chain/app/upgrades/v2_0"
	"github.com/strangelove-ventures/interchaintest/v8"

	"github.com/docker/docker/client"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"go.uber.org/zap/zaptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var UpgradeContainerVersion = "local"

func TestV1P2Upgrade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("%s Upgrade Test Suite", v2_0.UpgradeName))
}

var _ = Describe(fmt.Sprintf("%s Upgrade Tests", v2_0.UpgradeName), Ordered, func() {
	var kaon *cosmos.CosmosChain
	var kyve *cosmos.CosmosChain

	var ctx context.Context
	var client *client.Client
	var network string
	var interchain *interchaintest.Interchain

	var kaonWallet *cosmos.CosmosWallet
	var kyveWallet *cosmos.CosmosWallet

	BeforeAll(func() {
		numFullNodes := 0
		numValidators := 2
		factory := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(GinkgoT()), []*interchaintest.ChainSpec{
			{
				Name:          "kaon",
				ChainConfig:   testnetConfig,
				NumValidators: &numValidators,
				NumFullNodes:  &numFullNodes,
			},
			{
				Name:          "kyve",
				ChainConfig:   mainnetConfig,
				NumValidators: &numValidators,
				NumFullNodes:  &numFullNodes,
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
			GinkgoT(), ctx, GinkgoT().Name(), math.NewInt(10_000_000_000), kaon, kyve,
		)
		kaonWallet, kyveWallet = wallets[0].(*cosmos.CosmosWallet), wallets[1].(*cosmos.CosmosWallet)
	})

	AfterAll(func() {
		_ = kaon.StopAllNodes(ctx)
		_ = kyve.StopAllNodes(ctx)
		_ = interchain.Close()
	})

	It("Kaon upgrade test", func() {
		PerformUpgrade(ctx, client, kaon, kaonWallet, 10, "kaon")
	})

	It("Kyve upgrade test", func() {
		PerformUpgrade(ctx, client, kyve, kyveWallet, 10, "kyve")
	})
})

type Plan struct {
	Name   string `json:"name"`
	Height string `json:"height"`
	Info   string `json:"info"`
}

type SoftwareUpgradeProposal struct {
	Typedef   string `json:"@type"`
	Authority string `json:"authority"`
	Plan      Plan   `json:"plan"`
}

func generateUpgradeProposal(chain *cosmos.CosmosChain, height int64) cosmos.TxProposalv1 {
	prop := SoftwareUpgradeProposal{
		Typedef:   "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
		Authority: "kyve10d07y265gmmuvt4z0w9aw880jnsr700jdv7nah",
		Plan: Plan{
			Name:   v2_0.UpgradeName,
			Height: strconv.FormatInt(height, 10),
			Info:   "",
		},
	}
	msg, err := json.Marshal(prop)
	Expect(err).To(BeNil())
	return cosmos.TxProposalv1{
		Messages: []json.RawMessage{msg},
		Metadata: "",
		Deposit:  fmt.Sprintf("1000000000%s", chain.Config().Denom),
		Title:    v2_0.UpgradeName,
		Summary:  v2_0.UpgradeName,
	}
}

func PerformUpgrade(
	ctx context.Context,
	client *client.Client,
	chain *cosmos.CosmosChain,
	wallet *cosmos.CosmosWallet,
	delta int64,
	container string,
) {
	height, _ := chain.Height(ctx)
	haltHeight := height + delta

	upgrade, proposalErr := chain.SubmitProposal(ctx, wallet.KeyName(), generateUpgradeProposal(chain, haltHeight))
	Expect(proposalErr).To(BeNil())
	proposalId, _ := strconv.ParseUint(upgrade.ProposalID, 10, 64)
	voteErr := chain.VoteOnProposalAllValidators(ctx, proposalId, cosmos.ProposalVoteYes)
	Expect(voteErr).To(BeNil())

	proposalId, err := strconv.ParseUint(upgrade.ProposalID, 10, 64)
	Expect(err).To(BeNil())
	_, statusErr := cosmos.PollForProposalStatus(ctx, chain, height, haltHeight, proposalId, govv1beta1.StatusPassed)
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
