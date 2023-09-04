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

TEST CASES - msg_server_claim_commission_rewards.go

* Produce a valid bundle and check commission rewards
* Claim with non-staker account
* Claim more rewards than available
* Claim zero rewards
* Claim partial rewards
* Claim partial rewards twice
* Claim all rewards

*/

var _ = Describe("msg_server_claim_commission_rewards.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceStaker0 := s.GetBalanceFromAddress(i.STAKER_0)

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

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)

		// create a valid bundle so that uploader earns commission rewards
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

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
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce a valid bundle and check commission rewards", func() {
		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// calculate uploader rewards
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		balanceUploader := s.GetBalanceFromAddress(i.STAKER_0)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(pool.OperatingCost)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.OperatingCost - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())
		totalDelegationReward := totalUploaderReward - uploaderPayoutReward

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalDelegationReward))

		// assert commission rewards
		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward))

		// check pool funds
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(100*i.KYVE - pool.OperatingCost))
	})

	It("Claim with non-staker account", func() {
		// ACT
		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_1,
			Amount:  1,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())

		// assert commission rewards
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(pool.OperatingCost)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.OperatingCost - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())

		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward))

		Expect(s.GetBalanceFromAddress(i.STAKER_0)).To(Equal(initialBalanceStaker0))
	})

	It("Claim more rewards than available", func() {
		// ACT
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_1,
			Amount:  uploader.CommissionRewards + 1,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())

		// assert commission rewards
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(pool.OperatingCost)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.OperatingCost - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())

		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward))
	})

	It("Claim zero rewards", func() {
		// ACT
		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  0,
		})

		// ASSERT
		Expect(err).NotTo(HaveOccurred())

		// assert commission rewards
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(pool.OperatingCost)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.OperatingCost - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())

		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward))

		Expect(s.GetBalanceFromAddress(i.STAKER_0)).To(Equal(initialBalanceStaker0))
	})

	It("Claim partial rewards", func() {
		// ACT
		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  100,
		})

		// ASSERT
		Expect(err).NotTo(HaveOccurred())

		// assert commission rewards
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(pool.OperatingCost)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.OperatingCost - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())

		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward - 100))

		Expect(s.GetBalanceFromAddress(i.STAKER_0)).To(Equal(initialBalanceStaker0 + 100))
	})

	It("Claim partial rewards twice", func() {
		// ACT
		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  100,
		})

		// ASSERT
		Expect(err).NotTo(HaveOccurred())

		// assert commission rewards
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(sdk.NewDec(int64(pool.OperatingCost)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.OperatingCost - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())

		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward - 100))

		Expect(s.GetBalanceFromAddress(i.STAKER_0)).To(Equal(initialBalanceStaker0 + 100))

		// ACT
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash3",
			FromIndex:     200,
			BundleSize:    100,
			FromKey:       "200",
			ToKey:         "299",
			BundleSummary: "test_value3",
		})

		_, err = s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  200,
		})

		// ASSERT
		Expect(err).NotTo(HaveOccurred())

		// assert commission rewards
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		networkFee = s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward = uint64(sdk.NewDec(int64(pool.OperatingCost)).Mul(networkFee).TruncateInt64())
		storageReward = uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward = pool.OperatingCost - treasuryReward - storageReward

		uploaderPayoutReward = uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())

		Expect(uploader.CommissionRewards).To(Equal(2*(uploaderPayoutReward+storageReward) - 300))

		Expect(s.GetBalanceFromAddress(i.STAKER_0)).To(Equal(initialBalanceStaker0 + 300))
	})

	It("Claim all rewards", func() {
		// ACT
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		rewards := uploader.CommissionRewards

		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  rewards,
		})

		// ASSERT
		Expect(err).NotTo(HaveOccurred())

		// assert commission rewards
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(uploader.CommissionRewards).To(BeZero())

		Expect(s.GetBalanceFromAddress(i.STAKER_0)).To(Equal(initialBalanceStaker0 + rewards))
	})
})
