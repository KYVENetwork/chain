package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
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
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10_000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)

		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		params := s.App().FundersKeeper.GetParams(s.Ctx())
		params.MinFundingAmountPerBundle = 10_000
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

		// create a valid bundle so that uploader earns commission rewards
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 10_000,
		})

		s.CommitAfterSeconds(60)

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

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
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
		treasuryReward := uint64(sdk.NewDec(int64(pool.InflationShareWeight)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.InflationShareWeight - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())
		totalDelegationReward := totalUploaderReward - uploaderPayoutReward

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalDelegationReward))

		// assert commission rewards
		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(Equal(100*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
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
		treasuryReward := uint64(sdk.NewDec(int64(pool.InflationShareWeight)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.InflationShareWeight - treasuryReward - storageReward

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
		treasuryReward := uint64(sdk.NewDec(int64(pool.InflationShareWeight)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.InflationShareWeight - treasuryReward - storageReward

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
		treasuryReward := uint64(sdk.NewDec(int64(pool.InflationShareWeight)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.InflationShareWeight - treasuryReward - storageReward

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
		treasuryReward := uint64(sdk.NewDec(int64(pool.InflationShareWeight)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.InflationShareWeight - treasuryReward - storageReward

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
		treasuryReward := uint64(sdk.NewDec(int64(pool.InflationShareWeight)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward := pool.InflationShareWeight - treasuryReward - storageReward

		uploaderPayoutReward := uint64(sdk.NewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())

		Expect(uploader.CommissionRewards).To(Equal(uploaderPayoutReward + storageReward - 100))

		Expect(s.GetBalanceFromAddress(i.STAKER_0)).To(Equal(initialBalanceStaker0 + 100))

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_0_A,
			Staker:    i.STAKER_0,
			PoolId:    0,
			StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "SsdTPx9adtpwAGIjiHilqVPEfoTiq7eRw6khbVxKetQ",
			DataSize:      100,
			DataHash:      "test_hash3",
			FromIndex:     200,
			BundleSize:    100,
			FromKey:       "200",
			ToKey:         "299",
			BundleSummary: "test_value3",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "SsdTPx9adtpwAGIjiHilqVPEfoTiq7eRw6khbVxKetQ",
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "iW1jN99yH_gdQtRhf5J_lVwOIu8p_i7FyxEgoQAkWxU",
			DataSize:      100,
			DataHash:      "test_hash4",
			FromIndex:     300,
			BundleSize:    100,
			FromKey:       "300",
			ToKey:         "399",
			BundleSummary: "test_value4",
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
		treasuryReward = uint64(sdk.NewDec(int64(pool.InflationShareWeight)).Mul(networkFee).TruncateInt64())
		storageReward = uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx()).MulInt64(100).TruncateInt64())
		totalUploaderReward = pool.InflationShareWeight - treasuryReward - storageReward

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
