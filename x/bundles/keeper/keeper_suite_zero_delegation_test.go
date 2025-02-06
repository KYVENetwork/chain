package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - zero delegation

* Staker votes with zero delegation
* Staker receives vote slash with zero delegation
* Staker submit bundle proposal with zero delegation
* Staker receives upload slash with zero delegation
* Staker receives timeout slash because votes were missed
* Stakers try to produce valid bundle but all stakers have zero delegation

*/

var _ = Describe("zero delegation", Ordered, func() {
	var s *i.KeeperTestSuite
	var initialBalanceStaker0, initialBalancePoolAddress0, initialBalanceStaker1, initialBalancePoolAddress1 uint64

	amountPerBundle := uint64(10_000)

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalancePoolAddress0 = s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetBalanceFromAddress(i.POOL_ADDRESS_1_A)

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
			MinDelegation:        0 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		params := funderstypes.DefaultParams()
		params.CoinWhitelist[0].MinFundingAmountPerBundle = math.NewIntFromUint64(amountPerBundle)
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

		// create funders
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(int64(amountPerBundle)),
		})

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Staker votes with zero delegation", func() {
		// ARRANGE
		// create normal validator
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
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

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		// create zero delegation validator
		s.CreateZeroDelegationValidator(i.STAKER_2, "Staker-2")

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_2))
		Expect(bundleProposal.VotersInvalid).NotTo(ContainElement(i.STAKER_2))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_2))

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
	})

	It("Staker receives vote slash with zero delegation", func() {
		// ARRANGE
		// create normal validator
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
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

		// create zero delegation validator
		s.CreateZeroDelegationValidator(i.STAKER_2, "Staker-2")

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).NotTo(ContainElement(i.STAKER_2))
		Expect(bundleProposal.VotersInvalid).To(ContainElement(i.STAKER_2))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_2))

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))

		// calculate voter slashes
		fraction := s.App().StakersKeeper.GetVoteSlash(s.Ctx())
		slashAmountVoter := uint64(math.LegacyNewDec(int64(0 * i.KYVE)).Mul(fraction).TruncateInt64())
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.STAKER_2)).To(Equal(0*i.KYVE - slashAmountVoter))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200*i.KYVE - slashAmountVoter))
	})

	It("Staker submit bundle proposal with zero delegation", func() {
		// ARRANGE

		// create zero delegation validator
		staker0 := s.CreateNewValidator("Staker-0", 100*i.KYVE)

		validators, _ := s.App().StakingKeeper.GetValidators(s.Ctx(), 20)
		_ = validators

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       staker0.Address,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// create normal validator
		staker1 := s.CreateNewValidator("Staker-1", 100*i.KYVE)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       staker1.Address,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// create normal validator
		staker2 := s.CreateNewValidator("Staker-2", 100*i.KYVE)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       staker2.Address,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// manually set next uploader
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = staker0.Address
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        staker0.Address,
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

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    staker1.Address,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    staker2.Address,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)

		initialBalanceStaker0 = s.GetBalanceFromAddress(staker0.Address)
		initialBalancePoolAddress0 = s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)

		initialBalanceStaker1 = s.GetBalanceFromAddress(staker1.Address)
		initialBalancePoolAddress1 = s.GetBalanceFromAddress(i.POOL_ADDRESS_1_A)

		s.SetDelegationToZero(staker0.Address)
		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        staker1.Address,
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
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(staker0.Address))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(staker1.Address))
		Expect(bundleProposal.NextUploader).To(Equal(staker1.Address))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(staker1.Address))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), staker0.Address, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetBalanceFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		balanceUploader := s.GetBalanceFromAddress(poolAccountUploader.Staker)

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), staker1.Address, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetBalanceFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetBalanceFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// calculate uploader rewards
		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := pool.InflationShareWeight.Mul(networkFee)
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.CurrentStorageProviderId).MulInt64(100)
		totalUploaderReward := pool.InflationShareWeight.Sub(treasuryReward).Sub(storageReward)

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		// assert commission rewards
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), staker0.Address).AmountOf(globaltypes.Denom).Int64()).To(Equal(totalUploaderReward.Add(storageReward).TruncateInt64()))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), staker0.Address, staker0.Address)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100*i.KYVE - 1*amountPerBundle))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Staker receives upload slash with zero delegation", func() {
		// ARRANGE
		// create zero delegation validator
		staker0 := s.CreateNewValidator("Staker-0", 1*i.KYVE)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       staker0.Address,
			PoolId:        0,
			PoolAddress:   staker0.PoolAccount[0],
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// create normal validator
		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// create normal validator
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// manually set next uploader
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = staker0.Address
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       staker0.PoolAccount[0],
			Staker:        staker0.Address,
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

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.CommitAfterSeconds(60)

		initialBalanceStaker0 = s.GetBalanceFromAddress(staker0.Address)
		initialBalancePoolAddress0 = s.GetBalanceFromAddress(staker0.PoolAccount[0])

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetBalanceFromAddress(i.POOL_ADDRESS_1_A)

		s.SetDelegationToZero(staker0.Address)
		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
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
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
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
		_, uploaderActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), staker0.Address, 0)
		Expect(uploaderActive).To(BeFalse())

		balancePoolAddress := s.GetBalanceFromAddress(staker0.PoolAccount[0])
		Expect(balancePoolAddress).To(Equal(initialBalancePoolAddress0))

		balanceUploader := s.GetBalanceFromAddress(staker0.Address)
		_, uploaderFound := s.App().StakersKeeper.GetValidator(s.Ctx(), staker0.Address)
		Expect(uploaderFound).To(BeTrue())

		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), staker0.Address, staker0.Address)).To(BeEmpty())

		// calculate uploader slashes
		fraction := s.App().StakersKeeper.GetUploadSlash(s.Ctx())
		slashAmount := uint64(math.LegacyNewDec(int64(0 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), staker0.Address, staker0.Address)).To(Equal(0*i.KYVE - slashAmount))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200*i.KYVE - slashAmount))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetBalanceFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetBalanceFromAddress(poolAccountVoter.Staker)

		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100 * i.KYVE))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Staker receives timeout slash because votes were missed", func() {
		// ARRANGE
		// create normal validator
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// create zero delegation validator
		staker1 := s.CreateNewValidator("Staker-1", 100*i.KYVE)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       staker1.Address,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// create normal validator
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// manually set next uploader
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_0
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

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
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)
		s.SetDelegationToZero(staker1.Address)

		// ACT
		maxPoints := int(s.App().BundlesKeeper.GetMaxPoints(s.Ctx()))

		for r := 1; r <= maxPoints; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.POOL_ADDRESS_0_A,
				Staker:        i.STAKER_0,
				PoolId:        0,
				StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				DataSize:      100,
				DataHash:      "test_hash",
				FromIndex:     uint64(r * 100),
				BundleSize:    100,
				FromKey:       "test_key",
				ToKey:         "test_key",
				BundleSummary: "test_value",
			})

			s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
				Creator:   i.POOL_ADDRESS_2_A,
				Staker:    i.STAKER_2,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// ASSERT
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, stakerFound := s.App().StakersKeeper.GetValidator(s.Ctx(), staker1.Address)
		Expect(stakerFound).To(BeTrue())

		_, poolAccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), staker1.Address, 0)
		Expect(poolAccountActive).To(BeFalse())

		// check if voter got slashed
		slashAmountRatio := s.App().StakersKeeper.GetTimeoutSlash(s.Ctx())
		expectedBalance := 0*i.KYVE - uint64(math.LegacyNewDec(int64(0*i.KYVE)).Mul(slashAmountRatio).TruncateInt64())

		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), staker1.Address, staker1.Address)))
	})

	It("Stakers try to produce valid bundle but all stakers have zero delegation", func() {
		// ARRANGE
		s.CreateZeroDelegationValidator(i.STAKER_0, "Staker-0")

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateZeroDelegationValidator(i.STAKER_1, "Staker-1")

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// manually set next uploader
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_0
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		s.CommitAfterSeconds(60)

		s.RunTxBundlesError(&bundletypes.MsgSubmitBundleProposal{
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
	})
})
