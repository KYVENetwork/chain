package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
	"github.com/cosmos/cosmos-sdk/types"
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

* Produce a valid bundle with no funders, 0% inflation splitting and 0 inflation-share-weight
* Produce a valid bundle with no funders, 10% inflation splitting and pool-0 = 0.1 weight and pool-1 = 1.0 weight
* Produce a valid bundle with no funders, 10% inflation splitting and pool-0 = 1.0 weight and pool-1 = 1.0 weight

*/

var _ = Describe("inflation splitting", Ordered, func() {
	var s *i.KeeperTestSuite

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
			InflationShareWeight: math.LegacyNewDec(10_000),
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		params := funderstypes.DefaultParams()
		params.CoinWhitelist[0].MinFundingAmount = 100
		params.CoinWhitelist[0].MinFundingAmountPerBundle = 1_000
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
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

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
		Expect(uploader.CommissionRewards).To(BeEmpty())
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with no funders and 10% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

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

		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is here just the inflation payout
		totalPayout := payout

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(math.LegacyNewDec(int64(totalPayout)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateInt64())
		totalUploaderReward := totalPayout - treasuryReward - storageReward

		uploaderPayoutReward := uint64(math.LegacyNewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())
		uploaderDelegationReward := totalUploaderReward - uploaderPayoutReward

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).Uint64()).To(Equal(uploaderPayoutReward + storageReward))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with no funders and 100% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.2")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

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

		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is here just the inflation payout
		totalPayout := payout

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := uint64(math.LegacyNewDec(int64(totalPayout)).Mul(networkFee).TruncateInt64())
		storageReward := uint64(s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateInt64())
		totalUploaderReward := totalPayout - treasuryReward - storageReward

		uploaderPayoutReward := uint64(math.LegacyNewDec(int64(totalUploaderReward)).Mul(uploader.Commission).TruncateInt64())
		uploaderDelegationReward := totalUploaderReward - uploaderPayoutReward

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).Uint64()).To(Equal(uploaderPayoutReward + storageReward))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with sufficient funders and 0% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
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
		treasuryReward := totalPayout.Mul(networkFee)
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100)
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission)
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with sufficient funders and 10% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.3")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
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

		b1 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))

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
		b2 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))
		Expect(b1.TruncateInt64()).To(BeNumerically(">", b2.TruncateInt64()))

		payout := b1.Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateDec()
		Expect(b1.Sub(b2)).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := pool.InflationShareWeight.Add(payout)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with sufficient funders and 100% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
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

		b1 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))

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
		b2 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))
		Expect(b1.TruncateInt64()).To(BeNumerically(">", b2.TruncateInt64()))

		payout := b1.Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateDec()
		Expect(b1.Sub(b2)).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := pool.InflationShareWeight.Add(payout)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with insufficient funders and 0% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(5_000),
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
		totalPayout := math.LegacyNewDec(300)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with insufficient funders and 30% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.3")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(5_000),
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

		b1 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))

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
		b2 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))
		Expect(b1.TruncateInt64()).To(BeNumerically(">", b2.TruncateInt64()))

		payout := b1.Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateDec()
		Expect(b1.Sub(b2)).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := payout.Add(math.LegacyNewDec(300))

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with insufficient funders and 10% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(5_000),
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

		b1 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))

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
		b2 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))
		Expect(b1.TruncateInt64()).To(BeNumerically(">", b2.TruncateInt64()))

		payout := b1.Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateDec()
		Expect(b1.Sub(b2)).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := payout.Add(math.LegacyNewDec(300))

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with some insufficient funders and 0% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(5_000),
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
		totalPayout := pool.InflationShareWeight.Quo(math.LegacyNewDec(2)).Add(math.LegacyNewDec(200))

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with some insufficient funders and 30% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.3")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(5_000),
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

		b1 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))

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
		b2 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))
		Expect(b1.TruncateInt64()).To(BeNumerically(">", b2.TruncateInt64()))

		payout := b1.Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateDec()
		Expect(b1.Sub(b2)).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := pool.InflationShareWeight.Quo(math.LegacyNewDec(2)).Add(math.LegacyNewDec(200)).Add(payout)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with some insufficient funders and 10% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(5_000),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(5_000),
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

		b1 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))

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
		b2 := math.LegacyNewDec(int64(s.GetBalanceFromPool(0)))
		Expect(b1.TruncateInt64()).To(BeNumerically(">", b2.TruncateInt64()))

		payout := b1.Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateDec()
		Expect(b1.Sub(b2)).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is the inflation share weight plus the inflation payout
		totalPayout := pool.InflationShareWeight.Quo(math.LegacyNewDec(2)).Add(math.LegacyNewDec(200)).Add(payout)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with no funders, 0% inflation splitting and 0 inflation-share-weight", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.InflationShareWeight = math.LegacyNewDec(0)
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

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
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

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
		Expect(uploader.CommissionRewards).To(BeEmpty())
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with no funders, 10% inflation splitting and pool-0 = 0.1 weight and pool-1 = 1.0 weight", func() {
		// ARRANGE

		// Enable inflation share for pools
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// set team-share to zero to not interfere with inflation splitting
		teamTypes.TEAM_ALLOCATION = 0
		_ = s.App().BankKeeper.SendCoinsFromModuleToAccount(s.Ctx(), "team", types.MustAccAddressFromBech32(i.CHARLIE), s.GetCoinsFromModule("team"))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.InflationShareWeight = math.LegacyMustNewDecFromStr("0.1")
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// Add a second pool
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest 2",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: math.LegacyMustNewDecFromStr("1"),
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_0_B,
		})
		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     1,
			Valaddress: i.VALADDRESS_1_B,
		})

		preMineBalance := s.App().BankKeeper.GetSupply(s.Ctx(), globalTypes.Denom)
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

		postMineBalance := s.App().BankKeeper.GetSupply(s.Ctx(), globalTypes.Denom)

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

		inflationAmount := postMineBalance.Sub(preMineBalance)
		// Reward calculation:
		// (inflationAmount - teamRewards) * protocolInflationShare -> Both pools equally
		// (340112344399tkyve - 847940tkyve) * 0.1 -> rewards for both pools, but it is split according to the different weights
		// teamAuthority rewards are hard to set to zero from this test-suite without using reflection.
		// therefore we ignore the small amount.
		Expect(inflationAmount.String()).To(Equal("340112344399tkyve"))

		// assert if bundle go finalized
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert pool balance
		finalBalancePool0 := s.GetBalanceFromPool(0)
		finalBalancePool1 := s.GetBalanceFromPool(1)
		// First pool has weight: 0.1, second pool has weight 1
		// additionally, pool-0 produced a bundle -> subtract PoolInflationPayoutRate (1 - 0.1 = 0.9)
		// formula: (inflation - teamRewards) * inflationShare * inflationShareWeighOfPool * (1-PoolInflationPayoutRate)
		// (340112344399 - 847940) * 0.1 * 1 / 11 * 0.9
		// Evaluates to 2782730425, however due to multiple roundings to actual amount is 2782730381
		// second pool
		// (340112344399 - 847940) * 0.1 * 10 / 11
		// Evaluates to 30919226950
		Expect(finalBalancePool0).To(Equal(uint64(2782730381)))
		Expect(finalBalancePool1).To(Equal(uint64(30919226950)))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is here just the inflation payout
		// (inflation - teamRewards)*inflationShare - balancePool0 - balancePool1
		// (340112344399 - 847940) * 0.1 * 1 / 11 * 0.1
		// evaluates to 309192269, due to multiple rounding: 309192264
		totalPayout := math.LegacyNewDec(309192264)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())

		// Important: Reset changes of global variables as they will not be reverted by the s.NewCleanChain()
		teamTypes.TEAM_ALLOCATION = 165000000000000000
	})

	It("Produce a valid bundle with no funders, 10% inflation splitting and pool-0 = 1.0 weight and pool-1 = 1.0 weight", func() {
		// ARRANGE

		// Enable inflation share for pools
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("0.1")
		params.PoolInflationPayoutRate = math.LegacyMustNewDecFromStr("0.2")
		s.App().PoolKeeper.SetParams(s.Ctx(), params)

		// set team-share to zero to not interfere with inflation splitting
		teamTypes.TEAM_ALLOCATION = 0
		_ = s.App().BankKeeper.SendCoinsFromModuleToAccount(s.Ctx(), "team", types.MustAccAddressFromBech32(i.CHARLIE), s.GetCoinsFromModule("team"))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.InflationShareWeight = math.LegacyMustNewDecFromStr("1")
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// Add a second pool
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest 2",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: math.LegacyMustNewDecFromStr("1"),
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_0_B,
		})
		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     1,
			Valaddress: i.VALADDRESS_1_B,
		})

		// Both pools now have inflation_share=1 and zero balance

		preMineBalance := s.App().BankKeeper.GetSupply(s.Ctx(), globalTypes.Denom)
		// mine some blocks
		for i := 1; i < 100; i++ {
			s.Commit()
		}

		// Prepare bundle proposal
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

		postMineBalance := s.App().BankKeeper.GetSupply(s.Ctx(), globalTypes.Denom)

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

		inflationAmount := postMineBalance.Sub(preMineBalance)
		// Reward calculation:
		// (inflationAmount - teamRewards) * protocolInflationShare -> Both pools equally
		// (340112344399tkyve - 847940tkyve) * 0.1  -> (//2) -> 17005574822 for both pools
		// teamAuthority rewards are hard to set to zero from this test-suite without using reflection.
		// therefore we ignore the small amount.
		Expect(inflationAmount.String()).To(Equal("340112344399tkyve"))

		// assert if bundle go finalized
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert pool balance
		finalBalancePool0 := s.GetBalanceFromPool(0)
		finalBalancePool1 := s.GetBalanceFromPool(1)
		// Both pools have inflation-weight 1
		// however, pool-0 produced a bundle -> subtract PoolInflationPayoutRate (1 - 0.2 = 0.8)
		// 17005574822 * 0.8
		Expect(finalBalancePool0).To(Equal(uint64(13604459858)))
		Expect(finalBalancePool1).To(Equal(uint64(17005574822)))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// the total payout is here just the inflation payout
		totalPayout := math.LegacyNewDec(17005574822 - 13604459858)

		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := totalPayout.Mul(networkFee).TruncateDec()
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.GetCurrentStorageProviderId()).MulInt64(100).TruncateDec()
		totalUploaderReward := totalPayout.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(uploader.Commission).TruncateDec()
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		// assert commission rewards
		Expect(uploader.CommissionRewards.AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderPayoutReward.Add(storageReward)))
		// assert uploader self delegation rewards
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).ToLegacyDec()).To(Equal(uploaderDelegationReward))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())

		// Important: Reset changes of global variables as they will not be reverted by the s.NewCleanChain()
		teamTypes.TEAM_ALLOCATION = 165000000000000000
	})
})
