package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Gov
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_disabled_pool.go

* Invalid authority (transaction)
* Invalid authority (proposal)
* Disable a non-existing pool
* Disable pool which is active
* Disable pool which is active and has a balance
* Disable pool which is already disabled
* Disable multiple pools
* Kick out all stakers from pool
* Kick out all stakers from pool which are still members of another pool
* Drop current bundle proposal when pool gets disabled

*/

var _ = Describe("msg_server_disable_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	params, _ := s.App().GovKeeper.Params.Get(s.Ctx())
	votingPeriod := params.VotingPeriod
	fundingAmount := 100 * i.KYVE

	BeforeEach(func() {
		s = i.NewCleanChain()

		msg := &types.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: math.LegacyNewDec(10_000),
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxPoolSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		params := s.App().FundersKeeper.GetParams(s.Ctx())
		params.CoinWhitelist[0].MinFundingAmount = math.NewInt(100_000_000)
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(int64(fundingAmount)),
			AmountsPerBundle: i.KYVECoins(1 * i.T_KYVE),
		})

		msg = &types.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest2",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: math.LegacyNewDec(10_000),
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           1,
			Amounts:          i.KYVECoins(int64(fundingAmount)),
			AmountsPerBundle: i.KYVECoins(1 * i.T_KYVE),
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Invalid authority (transaction).", func() {
		// ARRANGE
		msg := &types.MsgDisablePool{
			Authority: i.DUMMY[0],
			Id:        0,
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.Disabled).To(BeFalse())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())
	})

	It("Invalid authority (proposal).", func() {
		// ARRANGE
		msg := &types.MsgDisablePool{
			Authority: i.DUMMY[0],
			Id:        0,
		}

		proposal, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, err := s.RunTx(&proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.Disabled).To(BeFalse())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())
	})

	It("Disable a non-existing pool", func() {
		// ARRANGE
		msg := &types.MsgDisablePool{
			Authority: gov,
			Id:        42,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusFailed))
	})

	It("Disable pool which is active", func() {
		// ARRANGE
		msg := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		fundersModuleAddr := s.App().AccountKeeper.GetModuleAddress(funderstypes.ModuleName)
		balanceBefore := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAddr, globalTypes.Denom).Amount.Uint64()
		Expect(balanceBefore).To(BeNumerically(">", uint64(0)))

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(pool.Disabled).To(BeTrue())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())

		// assert same funding balance
		balance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAddr, globalTypes.Denom).Amount.Uint64()
		Expect(balance).To(Equal(balanceBefore))
	})

	It("Disable pool which is active and has a balance", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		s.App().PoolKeeper.SetParams(s.Ctx(), types.Params{
			ProtocolInflationShare:  math.LegacyMustNewDecFromStr("0.1"),
			PoolInflationPayoutRate: math.LegacyMustNewDecFromStr("0.05"),
		})

		for i := 0; i < 100; i++ {
			s.Commit()
		}

		fundersModuleAddr := s.App().AccountKeeper.GetModuleAddress(funderstypes.ModuleName)
		balanceBefore := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAddr, globalTypes.Denom).Amount.Uint64()
		Expect(balanceBefore).To(BeNumerically(">", uint64(0)))

		msg := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(pool.Disabled).To(BeTrue())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())

		// assert same funding balance
		balance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAddr, globalTypes.Denom).Amount.Uint64()
		Expect(balance).To(Equal(balanceBefore))
	})

	It("Disable pool which is active", func() {
		// ARRANGE
		msg := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(pool.Disabled).To(BeTrue())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())
	})

	It("Disable pool which is already disabled", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		msg := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusFailed))
		Expect(pool.Disabled).To(BeTrue())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())
	})

	It("Disable multiple pools", func() {
		// ARRANGE
		msgFirstPool := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}
		msgSecondPool := &types.MsgDisablePool{
			Authority: gov,
			Id:        1,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msgFirstPool, msgSecondPool})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		firstPool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		secondPool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(firstPool.Disabled).To(BeTrue())
		Expect(secondPool.Disabled).To(BeTrue())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())

		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 1)
		Expect(bundleProposal.StorageId).To(BeEmpty())
	})

	It("Kick out all stakers from pool", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		msgFirstPool := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}

		Expect(s.App().StakersKeeper.GetAllPoolAccounts(s.Ctx())).To(HaveLen(2))

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msgFirstPool})

		msgVoteStaker0 := govV1Types.NewMsgVote(sdk.MustAccAddressFromBech32(i.STAKER_0), 1, govV1Types.VoteOption_VOTE_OPTION_YES, "")
		msgVoteStaker1 := govV1Types.NewMsgVote(sdk.MustAccAddressFromBech32(i.STAKER_1), 1, govV1Types.VoteOption_VOTE_OPTION_YES, "")

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)
		_, voteErr0 := s.RunTx(msgVoteStaker0)
		_, voteErr1 := s.RunTx(msgVoteStaker1)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)

		Expect(s.App().StakersKeeper.GetAllPoolAccounts(s.Ctx())).To(HaveLen(0))

		firstPool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))
		Expect(voteErr0).To(Not(HaveOccurred()))
		Expect(voteErr1).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(firstPool.Disabled).To(BeTrue())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())
	})

	It("Kick out all stakers from pool which are still members of another pool", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        1,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		msgFirstPool := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}

		Expect(s.App().StakersKeeper.GetAllPoolAccounts(s.Ctx())).To(HaveLen(3))

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msgFirstPool})

		msgVoteStaker0 := govV1Types.NewMsgVote(sdk.MustAccAddressFromBech32(i.STAKER_0), 1, govV1Types.VoteOption_VOTE_OPTION_YES, "")
		msgVoteStaker1 := govV1Types.NewMsgVote(sdk.MustAccAddressFromBech32(i.STAKER_1), 1, govV1Types.VoteOption_VOTE_OPTION_YES, "")

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)
		_, voteErr0 := s.RunTx(msgVoteStaker0)
		_, voteErr1 := s.RunTx(msgVoteStaker1)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)

		Expect(s.App().StakersKeeper.GetAllPoolAccounts(s.Ctx())).To(HaveLen(1))

		firstPool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))
		Expect(voteErr0).To(Not(HaveOccurred()))
		Expect(voteErr1).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(firstPool.Disabled).To(BeTrue())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(BeEmpty())
	})

	It("Drop current bundle proposal when pool gets disabled", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.POOL_ADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		msgFirstPool := &types.MsgDisablePool{
			Authority: gov,
			Id:        0,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msgFirstPool})

		msgVoteStaker0 := govV1Types.NewMsgVote(sdk.MustAccAddressFromBech32(i.STAKER_0), 1, govV1Types.VoteOption_VOTE_OPTION_YES, "")
		msgVoteStaker1 := govV1Types.NewMsgVote(sdk.MustAccAddressFromBech32(i.STAKER_1), 1, govV1Types.VoteOption_VOTE_OPTION_YES, "")

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)
		_, voteErr0 := s.RunTx(msgVoteStaker0)
		_, voteErr1 := s.RunTx(msgVoteStaker1)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))
		Expect(voteErr0).To(Not(HaveOccurred()))
		Expect(voteErr1).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(pool.Disabled).To(BeTrue())

		// check if bundle proposal got dropped
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.DataSize).To(BeZero())
		Expect(bundleProposal.DataHash).To(BeEmpty())
		Expect(bundleProposal.BundleSize).To(BeZero())
		Expect(bundleProposal.FromKey).To(BeEmpty())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.BundleSummary).To(BeEmpty())
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())
	})
})
