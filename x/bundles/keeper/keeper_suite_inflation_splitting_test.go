package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - inflation splitting

* Produce a valid bundle with no funders and 0% inflation splitting
* Produce a valid bundle with no funders and 10% inflation splitting
* Produce a valid bundle with no funders and 100% inflation splitting

* TODO: Produce a valid bundle with sufficient funders and 0% inflation splitting
* TODO: Produce a valid bundle with sufficient funders and 10% inflation splitting
* TODO: Produce a valid bundle with sufficient funders and 100% inflation splitting

* TODO: Produce a valid bundle with insufficient funders and 0% inflation splitting
* TODO: Produce a valid bundle with insufficient funders and 10% inflation splitting
* TODO: Produce a valid bundle with insufficient funders and 100% inflation splitting

*/

var _ = Describe("inflation splitting", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			MaxBundleSize:  100,
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  10_000,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce a valid bundle with no funders and 0% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.CommitAfterSeconds(60)

		b1 := s.GetBalanceFromPool(0)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
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
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert pool balance
		b2 := s.GetBalanceFromPool(0)
		Expect(b1).To(BeZero())
		Expect(b2).To(BeZero())

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards
		Expect(uploader.CommissionRewards).To(BeZero())
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeZero())

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())
	})

	It("Produce a valid bundle with no funders and 10% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0.1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.CommitAfterSeconds(60)

		b1 := s.GetBalanceFromPool(0)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
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
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert pool balance
		b2 := s.GetBalanceFromPool(0)
		Expect(b1).To(BeNumerically(">", b2))

		payout := uint64(sdk.NewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is here just the inflation payout
		totalPayout := payout

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(totalPayout)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := totalPayout - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())
		uploaderDelegationReward := totalUploaderReward - uploaderPayoutReward

		// assert commission rewards
		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(uploaderDelegationReward))

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())
	})

	It("Produce a valid bundle with no funders and 100% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.2"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.CommitAfterSeconds(60)

		b1 := s.GetBalanceFromPool(0)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
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
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert pool balance
		b2 := s.GetBalanceFromPool(0)
		Expect(b1).To(BeNumerically(">", b2))

		payout := uint64(sdk.NewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is here just the inflation payout
		totalPayout := payout

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(totalPayout)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := totalPayout - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())
		uploaderDelegationReward := totalUploaderReward - uploaderPayoutReward

		// assert commission rewards
		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(uploaderDelegationReward))

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())
	})
})
