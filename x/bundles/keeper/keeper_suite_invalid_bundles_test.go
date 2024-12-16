package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - invalid bundles

* Produce an invalid bundle with multiple validators and no foreign delegations
* Produce an invalid bundle with multiple validators and foreign delegations
* Produce an invalid bundle with multiple validators although some voted valid

*/

var _ = Describe("invalid bundles", Ordered, func() {
	var s *i.KeeperTestSuite
	var initialBalanceStaker0, initialBalanceValaddress0, initialBalanceStaker1, initialBalanceValaddress1, initialBalanceStaker2, initialBalanceValaddress2 uint64

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0_A)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceValaddress1 = s.GetBalanceFromAddress(i.VALADDRESS_1_A)

		initialBalanceStaker2 = s.GetBalanceFromAddress(i.STAKER_2)
		initialBalanceValaddress2 = s.GetBalanceFromAddress(i.VALADDRESS_2_A)

		// create clean pool for every test case
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
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

		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(1 * i.T_KYVE),
		})

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0_A)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceValaddress1 = s.GetBalanceFromAddress(i.VALADDRESS_1_A)

		initialBalanceStaker2 = s.GetBalanceFromAddress(i.STAKER_2)
		initialBalanceValaddress2 = s.GetBalanceFromAddress(i.VALADDRESS_2_A)

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce an invalid bundle with multiple validators and no foreign delegations", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
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

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_2)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got not finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal(""))
		Expect(pool.CurrentSummary).To(BeEmpty())
		Expect(pool.CurrentIndex).To(BeZero())
		Expect(pool.TotalBundles).To(BeZero())

		// check if finalized bundle exists
		_, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeFalse())

		// check if bundle proposal got dropped
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
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

		// check uploader status
		_, uploaderActive := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(uploaderActive).To(BeFalse())

		balanceValaddress := s.GetBalanceFromAddress(i.VALADDRESS_0_A)
		Expect(balanceValaddress).To(Equal(initialBalanceValaddress0))

		balanceUploader := s.GetBalanceFromAddress(i.STAKER_0)

		_, uploaderFound := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(uploaderFound).To(BeTrue())

		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeEmpty())

		// calculate uploader slashes
		fraction := s.App().StakersKeeper.GetUploadSlash(s.Ctx())
		slashAmount := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100*i.KYVE - slashAmount))
		Expect(s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check voter status
		valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress := s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter := s.GetBalanceFromAddress(valaccountVoter.Staker)

		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())

		// check voter 2 status
		valaccountVoter, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_2)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress = s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter = s.GetBalanceFromAddress(valaccountVoter.Staker)

		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_2, i.STAKER_2)).To(BeEmpty())

		// check pool funds
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100 * i.KYVE))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce an invalid bundle with multiple validators and foreign delegations", func() {
		// ARRANGE
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(100*i.KYVE)),
		))

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.BOB,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(100*i.KYVE)),
		))

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.CHARLIE,
			util.MustValaddressFromOperatorAddress(i.STAKER_2),
			sdk.NewInt64Coin(globalTypes.Denom, int64(100*i.KYVE)),
		))

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
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

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got not finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal(""))
		Expect(pool.CurrentSummary).To(BeEmpty())
		Expect(pool.CurrentIndex).To(BeZero())
		Expect(pool.TotalBundles).To(BeZero())

		// check if finalized bundle exists
		_, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeFalse())

		// check if bundle proposal got dropped
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
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

		// check uploader status
		_, uploaderActive := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(uploaderActive).To(BeFalse())

		balanceValaddress := s.GetBalanceFromAddress(i.VALADDRESS_0_A)
		Expect(balanceValaddress).To(Equal(initialBalanceValaddress0))

		balanceUploader := s.GetBalanceFromAddress(i.STAKER_0)
		_, uploaderFound := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(uploaderFound).To(BeTrue())

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeEmpty())

		// calculate uploader slashes
		fraction := s.App().StakersKeeper.GetUploadSlash(s.Ctx())
		slashAmountUploader := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())
		slashAmountDelegator := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100*i.KYVE - slashAmountUploader))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.ALICE)).To(Equal(100*i.KYVE - slashAmountDelegator))

		Expect(s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal(400 * i.KYVE))

		// check voter status
		valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress := s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter := s.GetBalanceFromAddress(valaccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.BOB)).To(BeEmpty())

		// check voter 2 status
		valaccountVoter, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress = s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter = s.GetBalanceFromAddress(valaccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.BOB)).To(BeEmpty())

		// check pool funds
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100 * i.KYVE))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce an invalid bundle with multiple validators although some voted valid", func() {
		// ARRANGE
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(100*i.KYVE)),
		))

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.BOB,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(100*i.KYVE)),
		))

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.CHARLIE,
			util.MustValaddressFromOperatorAddress(i.STAKER_2),
			sdk.NewInt64Coin(globalTypes.Denom, int64(100*i.KYVE)),
		))

		s.CreateValidator(i.STAKER_3, "Staker-3", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_3,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_3_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// delegate a bit more so invalid voters have more than 50%
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.DAVID,
			util.MustValaddressFromOperatorAddress(i.STAKER_3),
			sdk.NewInt64Coin(globalTypes.Denom, int64(150*i.KYVE)),
		))

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
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

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceStaker2 = s.GetBalanceFromAddress(i.STAKER_2)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_3_A,
			Staker:    i.STAKER_3,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_3_A,
			Staker:        i.STAKER_3,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got not finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal(""))
		Expect(pool.CurrentSummary).To(BeEmpty())
		Expect(pool.CurrentIndex).To(BeZero())
		Expect(pool.TotalBundles).To(BeZero())

		// check if finalized bundle exists
		_, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeFalse())

		// check if bundle proposal got dropped
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
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

		// check uploader status
		_, uploaderActive := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(uploaderActive).To(BeFalse())

		balanceValaddress := s.GetBalanceFromAddress(i.VALADDRESS_0_A)
		Expect(balanceValaddress).To(Equal(initialBalanceValaddress0))

		balanceUploader := s.GetBalanceFromAddress(i.STAKER_0)
		_, uploaderFound := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(uploaderFound).To(BeTrue())

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeEmpty())

		// calculate uploader slashes
		fraction := s.App().StakersKeeper.GetUploadSlash(s.Ctx())
		slashAmountUploader := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())
		slashAmountDelegator1 := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100*i.KYVE - slashAmountUploader))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.ALICE)).To(Equal(100*i.KYVE - slashAmountDelegator1))

		// calculate voter slashes
		fraction = s.App().StakersKeeper.GetVoteSlash(s.Ctx())
		slashAmountVoter := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())
		slashAmountDelegator2 := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(Equal(100*i.KYVE - slashAmountVoter))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.BOB)).To(Equal(100*i.KYVE - slashAmountDelegator2))

		Expect(s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal(450 * i.KYVE))

		// check voter status
		_, voterActive := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(voterActive).To(BeFalse())

		balanceVoterValaddress := s.GetBalanceFromAddress(i.VALADDRESS_1_A)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress2))

		balanceVoter := s.GetBalanceFromAddress(i.STAKER_1)
		Expect(balanceVoter).To(Equal(initialBalanceStaker2))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.BOB)).To(BeEmpty())

		// check voter2 status
		valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_2)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress = s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter = s.GetBalanceFromAddress(valaccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_2, i.STAKER_2)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_2, i.CHARLIE)).To(BeEmpty())

		// check voter3 status
		valaccountVoter, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_3)
		Expect(valaccountVoter.Points).To(BeZero())

		balanceVoterValaddress = s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter = s.GetBalanceFromAddress(valaccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_3, i.STAKER_3)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_3, i.DAVID)).To(BeEmpty())

		// check pool funds
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100 * i.KYVE))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})
})
