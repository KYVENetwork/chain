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

TEST CASES - inflation splitting

* Produce a valid bundle with no funders and 0% inflation splitting
* Produce a valid bundle with no funders and 10% inflation splitting
* Produce a valid bundle with no funders and 100% inflation splitting

* Produce a valid bundle with sufficient funders and 0% inflation splitting
* Produce a valid bundle with sufficient funders and 10% inflation splitting
* Produce a valid bundle with sufficient funders and 100% inflation splitting

* Produce a valid bundle with insufficient funders and 0% inflation splitting
* Produce a valid bundle with insufficient funders and 10% inflation splitting
* Produce a valid bundle with insufficient funders and 100% inflation splitting

* Produce a valid bundle with some insufficient funders and 0% inflation splitting
* Produce a valid bundle with some insufficient funders and 10% inflation splitting
* Produce a valid bundle with some insufficient funders and 100% inflation splitting

*/

var _ = Describe("inflation splitting", Ordered, func() {
	s := i.NewCleanChain()

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

		params := s.App().FundersKeeper.GetParams(s.Ctx())
		params.MinFundingAmountPerBundle = 1_000
		params.MinFundingAmount = 100
		params.MinFundingMultiple = 0
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

		// create funders
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.BOB,
			Moniker: "Bob",
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
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

		b1 := s.GetBalanceFromPool(0)

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
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

		b1 := s.GetBalanceFromPool(0)

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with sufficient funders and 0% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the inflation share weight because the funding is sufficient
		// and there is no additional inflation
		totalPayout := pool.InflationShareWeight

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with sufficient funders and 10% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0.1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.3"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := pool.InflationShareWeight + payout

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with sufficient funders and 100% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := pool.InflationShareWeight + payout

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with insufficient funders and 0% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the total funds
		totalPayout := uint64(300)

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with insufficient funders and 30% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0.1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.3"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := 300 + payout

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with insufficient funders and 10% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := 300 + payout

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with some insufficient funders and 0% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the total funds
		totalPayout := (pool.InflationShareWeight / 2) + 200

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with some insufficient funders and 30% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0.1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.3"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := (pool.InflationShareWeight / 2) + 200 + payout

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with some insufficient funders and 10% inflation splitting", func() {
		// ARRANGE
		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 5_000,
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200,
			AmountPerBundle: 5_000,
		})

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

		b1 := s.GetBalanceFromPool(0)

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

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := (pool.InflationShareWeight / 2) + 200 + payout

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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})
})
