package proposal_handler_test

import (
	"context"
	"cosmossdk.io/math"
	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"reflect"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"go.uber.org/zap/zaptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - proposal_handler.go

* Execute multiple transactions and check their order
* Execute transactions that exceed max tx bytes

*/

func TestProposalHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "interchaintest/ProposalHandler Test Suite")
}

var _ = Describe("proposal_handler.go", Ordered, func() {
	var chain *cosmos.CosmosChain

	var ctx context.Context
	var interchain *interchaintest.Interchain

	var broadcaster *cosmos.Broadcaster
	var wallets []*cosmos.CosmosWallet

	BeforeAll(func() {
		numFullNodes := 0
		numValidators := 2
		factory := interchaintest.NewBuiltinChainFactory(
			zaptest.NewLogger(GinkgoT()),
			[]*interchaintest.ChainSpec{mainnetChainSpec(numValidators, numFullNodes)},
		)

		chains, err := factory.Chains(GinkgoT().Name())
		Expect(err).To(BeNil())
		chain = chains[0].(*cosmos.CosmosChain)

		interchain = interchaintest.NewInterchain().
			AddChain(chain)

		broadcaster = cosmos.NewBroadcaster(GinkgoT(), chain)
		broadcaster.ConfigureClientContextOptions(func(clientContext sdkclient.Context) sdkclient.Context {
			return clientContext.
				WithBroadcastMode(flags.BroadcastAsync)
		})
		broadcaster.ConfigureFactoryOptions(func(factory tx.Factory) tx.Factory {
			return factory.
				WithGas(flags.DefaultGasLimit * 10)
		})

		ctx = context.Background()
		client, network := interchaintest.DockerSetup(GinkgoT())

		err = interchain.Build(ctx, nil, interchaintest.InterchainBuildOptions{
			TestName:         GinkgoT().Name(),
			Client:           client,
			NetworkID:        network,
			SkipPathCreation: true,
		})
		Expect(err).To(BeNil())

		for i := 0; i < 10; i++ {
			wallets = append(wallets, interchaintest.GetAndFundTestUsers(
				GinkgoT(), ctx, GinkgoT().Name(), math.NewInt(10_000_000_000), chain,
			)[0].(*cosmos.CosmosWallet))
		}
	})

	AfterAll(func() {
		_ = chain.StopAllNodes(ctx)
		_ = interchain.Close()
	})

	It("Execute multiple transactions and check their order", func() {
		// ARRANGE
		err := testutil.WaitForBlocks(ctx, 1, chain)
		Expect(err).To(BeNil())

		height, err := chain.Height(ctx)
		Expect(err).To(BeNil())

		// ACT

		// Execute different transactions
		// We don't care about the results, they only have to be included in a block
		broadcastMsg(ctx, broadcaster, wallets[0], &banktypes.MsgSend{FromAddress: wallets[0].FormattedAddress()})
		broadcastMsg(ctx, broadcaster, wallets[1], &stakerstypes.MsgCreateStaker{Creator: wallets[1].FormattedAddress()})
		broadcastMsg(ctx, broadcaster, wallets[2], &bundlestypes.MsgClaimUploaderRole{Creator: wallets[2].FormattedAddress()}) // priority msg
		broadcastMsg(ctx, broadcaster, wallets[3], &stakerstypes.MsgJoinPool{Creator: wallets[3].FormattedAddress(), Valaddress: wallets[0].FormattedAddress(), PoolId: 0})
		broadcastMsg(ctx, broadcaster, wallets[4], &banktypes.MsgSend{FromAddress: wallets[4].FormattedAddress()})
		broadcastMsg(ctx, broadcaster, wallets[5], &bundlestypes.MsgVoteBundleProposal{Creator: wallets[5].FormattedAddress()})   // priority msg
		broadcastMsg(ctx, broadcaster, wallets[6], &bundlestypes.MsgSkipUploaderRole{Creator: wallets[6].FormattedAddress()})     // priority msg
		broadcastMsg(ctx, broadcaster, wallets[7], &bundlestypes.MsgSubmitBundleProposal{Creator: wallets[7].FormattedAddress()}) // priority msg

		expectedOrder := []string{
			// priority msgs
			reflect.TypeOf(bundlestypes.MsgClaimUploaderRole{}).Name(),
			reflect.TypeOf(bundlestypes.MsgVoteBundleProposal{}).Name(),
			reflect.TypeOf(bundlestypes.MsgSkipUploaderRole{}).Name(),
			reflect.TypeOf(bundlestypes.MsgSubmitBundleProposal{}).Name(),
			// default msgs
			reflect.TypeOf(banktypes.MsgSend{}).Name(),
			reflect.TypeOf(stakerstypes.MsgCreateStaker{}).Name(),
			reflect.TypeOf(stakerstypes.MsgJoinPool{}).Name(),
			reflect.TypeOf(banktypes.MsgSend{}).Name(),
		}

		afterHeight, err := chain.Height(ctx)
		Expect(err).To(BeNil())
		Expect(afterHeight).To(Equal(height))

		// Wait for the transactions to be included in a block
		err = testutil.WaitForBlocks(ctx, 2, chain)
		Expect(err).To(BeNil())

		// ASSERT

		// Check the order of the transactions
		checkTxsOrder(ctx, chain, height+1, expectedOrder)
	})

	It("Execute transactions that exceed max tx bytes", func() {
		// ARRANGE
		err := testutil.WaitForBlocks(ctx, 1, chain)
		Expect(err).To(BeNil())

		height, err := chain.Height(ctx)
		Expect(err).To(BeNil())

		// ACT
		const duplications = 40
		broadcastMsg(ctx, broadcaster, wallets[0], &stakerstypes.MsgCreateStaker{Creator: wallets[0].FormattedAddress()})
		broadcastMsgs(ctx, broadcaster, wallets[1], duplicateMsg(&banktypes.MsgSend{FromAddress: wallets[1].FormattedAddress()}, duplications)...)
		broadcastMsg(ctx, broadcaster, wallets[2], &bundlestypes.MsgSkipUploaderRole{Creator: wallets[2].FormattedAddress()}) // priority msg

		// this will not make it into the actual block, so it goes into the next one with all following msgs
		broadcastMsgs(ctx, broadcaster, wallets[3], duplicateMsg(&banktypes.MsgSend{FromAddress: wallets[4].FormattedAddress()}, duplications)...)
		broadcastMsg(ctx, broadcaster, wallets[4], &stakerstypes.MsgJoinPool{Creator: wallets[5].FormattedAddress(), Valaddress: wallets[0].FormattedAddress(), PoolId: 0})
		broadcastMsg(ctx, broadcaster, wallets[5], &bundlestypes.MsgVoteBundleProposal{Creator: wallets[6].FormattedAddress()}) // priority msg

		afterHeight, err := chain.Height(ctx)
		Expect(err).To(BeNil())
		Expect(afterHeight).To(Equal(height))

		// Wait for the transactions to be included in a block
		err = testutil.WaitForBlocks(ctx, 2, chain)
		Expect(err).To(BeNil())

		// ASSERT
		var msgTypes []string
		for i := 0; i < duplications; i++ {
			msgTypes = append(msgTypes, reflect.TypeOf(banktypes.MsgSend{}).Name())
		}

		// Check that only the first block contains the first transactions
		checkTxsOrder(ctx, chain, height+1, append(
			[]string{
				reflect.TypeOf(bundlestypes.MsgSkipUploaderRole{}).Name(), // priority msg
				reflect.TypeOf(stakerstypes.MsgCreateStaker{}).Name(),
			},
			msgTypes...,
		))
		// The second block should contain the rest of the transactions
		checkTxsOrder(ctx, chain, height+2, append(
			msgTypes,
			[]string{
				reflect.TypeOf(bundlestypes.MsgVoteBundleProposal{}).Name(), // priority msg
				reflect.TypeOf(stakerstypes.MsgJoinPool{}).Name(),
			}...,
		))
	})
})
