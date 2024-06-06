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

* Produce a valid bundle with multiple coins funded and 0% inflation splitting
* Produce a valid bundle with multiple coins funded and 10% inflation splitting
* Produce a valid bundle with multiple coins funded and 100% inflation splitting

* Produce a valid bundle with no funders, 0% inflation splitting and 0 inflation-share-weight
* Produce a valid bundle with no funders, 10% inflation splitting and pool-0 = 0.1 weight and pool-1 = 1.0 weight
* Produce a valid bundle with no funders, 10% inflation splitting and pool-0 = 1.0 weight and pool-1 = 1.0 weight

*/

var _ = Describe("inflation splitting", Ordered, func() {
	var s *i.KeeperTestSuite

	amountPerBundle := int64(5_000)

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
			StorageProviderId:    1,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		// set storage cost to 0.5
		bundleParams := s.App().BundlesKeeper.GetParams(s.Ctx())
		bundleParams.StorageCosts = append(bundleParams.StorageCosts, bundletypes.StorageCost{StorageProviderId: 1, Cost: math.LegacyMustNewDecFromStr("0.5")})
		s.App().BundlesKeeper.SetParams(s.Ctx(), bundleParams)

		// set funders params
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams([]*funderstypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				MinFundingAmount:          math.NewInt(100),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.A_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(3),
			},
		}, 0))

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

		// Important: Reset changes of global variables as they will not be reverted by the s.NewCleanChain()
		originalTeamAllocation := teamTypes.TEAM_ALLOCATION
		DeferCleanup(func() {
			teamTypes.TEAM_ALLOCATION = originalTeamAllocation
		})
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
		c1 := s.GetCoinsFromCommunityPool()

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

		// assert treasury payout since inflation is zero here
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...)).To(BeEmpty())

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

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (total_bundle_payout - treasury_reward - storage_cost) * (1 - commission)
		// storage_cost = byte_size * usd_per_byte / len(coins) * coin_weight
		// (2471669 - (2471669 * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(244_740).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (total_bundle_payout - treasury_reward - storage_cost) * commission + storage_cost
		// storage_cost = byte_size * usd_per_byte / len(coins) * coin_weight
		// (2471669 - (2471669 * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(2_202_213).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)).To(BeZero())
		Expect(fundingState.ActiveFunderAddresses).To(BeEmpty())
	})

	It("Produce a valid bundle with no funders and 100% inflation splitting", func() {
		// ARRANGE
		params := pooltypes.DefaultParams()
		params.ProtocolInflationShare = math.LegacyMustNewDecFromStr("1")
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

		// inflation payout is 49433483tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (49433483 - (49433483 * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(4_893_959).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (49433483 - (49433483 * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(44_045_190).String()))

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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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
		c1 := s.GetCoinsFromCommunityPool()

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

		// check uploader rewards
		// we assert no kyve coins here since inflation is zero
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(1035).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(8865).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout since inflation is zero here
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(BeZero())
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(BeZero())
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(BeZero())

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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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

		// inflation payout is 7415009tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (7_415_009 + 10_000 - ((7_415_009 + 10_000) * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(735_120).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (7_415_009 + 10_000 - ((7_415_009 + 10_000) * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(6_615_639).String()))

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

		// inflation payout is 24716741tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (2_471_6741 + 10_000 - ((2_471_6741 + 10_000) * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(2_447_992).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (2_471_6741 + 10_000 - ((2_471_6741 + 10_000) * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(22_031_482).String()))

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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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
		c1 := s.GetCoinsFromCommunityPool()

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

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (300 - (300 * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(74).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (300 - (300 * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(223).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout since inflation is zero here
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(3)))
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(BeZero())
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(BeZero())
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(BeZero())

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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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

		// inflation payout is 7_415_009tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (7_415_009 + 300 - ((7_415_009 + 300) * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(734_160).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (7_415_009 + 300 - ((7_415_009 + 300) * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(6_606_996).String()))

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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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

		// inflation payout is 24716741tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (24_716_741 + 300 - ((24_716_741 + 300) * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(2_447_032).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (24_716_741 + 300 - ((24_716_741 + 300) * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(22_022_839).String()))

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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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
		c1 := s.GetCoinsFromCommunityPool()

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

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (5_000 + 200 - ((5_000 + 200) * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(559).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (5_000 + 200 - ((5_000 + 200) * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(4_589).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout since inflation is zero here
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(52)))
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(BeZero())
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(BeZero())
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(BeZero())

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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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

		// inflation payout is 7415009tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (7_415_009 + 5_000 + 200 - ((7_415_009 + 5_000 + 200) * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(734_645).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (7_415_009 + 5_000 + 200 - ((7_415_009 + 5_000 + 200) * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(6_611_362).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with some insufficient funders and 100% inflation splitting", func() {
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
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.KYVECoins(200),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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

		// inflation payout is 24716741tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (24_716_741 + 5_000 + 200 - ((24_716_741 + 5_000 + 200) * 0.01) - _((100 * 0.5) / (1 * 1))_) * 0.1 + _((100 * 0.5) / (1 * 1))_
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(2_447_517).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (24_716_741 + 5_000 + 200 - ((24_716_741 + 5_000 + 200) * 0.01) - _((100 * 0.5) / (1 * 1))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(22_027_205).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100*i.KYVE - 5_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple coins funded and 0% inflation splitting", func() {
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
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(amountPerBundle), i.BCoin(2*amountPerBundle)),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(amountPerBundle), i.BCoin(2*amountPerBundle)),
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
		c1 := s.GetCoinsFromCommunityPool()

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

		// check uploader rewards
		// we assert no kyve coins here since inflation is zero
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (2 * coin_weight))_) * 0.1 + _((100 * 0.5) / (2 * coin_weight))_
		Expect(uploader.CommissionRewards.String()).To(Equal(sdk.NewCoins(i.ACoin(1012), i.BCoin(1990)).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (20_000 - (20_000 * 0.01) - _((100 * 0.5) / (2 * coin_weight))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(sdk.NewCoins(i.ACoin(8888), i.BCoin(17810)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(Equal(uint64(200)))
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(BeZero())

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.ACoin(200*i.T_KYVE-2*amountPerBundle), i.BCoin(200*i.T_KYVE-4*amountPerBundle)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with multiple coins funded and 10% inflation splitting", func() {
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
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(amountPerBundle), i.BCoin(2*amountPerBundle)),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(amountPerBundle), i.BCoin(2*amountPerBundle)),
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
		c1 := s.GetCoinsFromCommunityPool()

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

		// inflation payout is 7415009tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// for kyve coin (7_415_009 - (7_415_009 * 0.01) - _((100 * 0.5) / (3 * 1))_) * 0.1 + _((100 * 0.5) / (3 * 1))_
		// for acoin (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (3 * 1))_) * 0.1 + _((100 * 0.5) / (3 * 1))_
		// for bcoin coins (20_000 - (20_000 * 0.01) - _((100 * 0.5) / (3 * 2))_) * 0.1 + _((100 * 0.5) / (3 * 2))_
		Expect(uploader.CommissionRewards.String()).To(Equal(sdk.NewCoins(i.KYVECoin(734_100), i.ACoin(1004), i.BCoin(1987)).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// for kyve coin (7_415_009 - (7_415_009 * 0.01) - _((100 * 0.5) / (3 * 1))_) * (1 - 0.1)
		// for acoin (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (3 * 1))_) * (1 - 0.1)
		// for bcoin (20_000 - (20_000 * 0.01) - _((100 * 0.5) / (3 * 2))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(sdk.NewCoins(i.KYVECoin(6_606_759), i.ACoin(8896), i.BCoin(17813)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(Equal(uint64(200)))
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(BeZero())

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
	})

	It("Produce a valid bundle with multiple coins funded and 100% inflation splitting", func() {
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
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(amountPerBundle), i.BCoin(2*amountPerBundle)),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(amountPerBundle), i.BCoin(2*amountPerBundle)),
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
		c1 := s.GetCoinsFromCommunityPool()

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

		// inflation payout is 24716741tkyve
		payout := uint64(math.LegacyNewDec(int64(b1)).Mul(s.App().PoolKeeper.GetPoolInflationPayoutRate(s.Ctx())).TruncateInt64())
		Expect(b1 - b2).To(Equal(payout))

		// assert bundle reward
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// for kyve coin (24_716_741 - (24_716_741 * 0.01) - _((100 * 0.5) / (3 * 1))_) * 0.1 + _((100 * 0.5) / (3 * 1))_
		// for acoin (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (3 * 1))_) * 0.1 + _((100 * 0.5) / (3 * 1))_
		// for bcoin coins (20_000 - (20_000 * 0.01) - _((100 * 0.5) / (3 * 2))_) * 0.1 + _((100 * 0.5) / (3 * 2))_
		Expect(uploader.CommissionRewards.String()).To(Equal(sdk.NewCoins(i.KYVECoin(2_446_971), i.ACoin(1004), i.BCoin(1987)).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// for kyve coin (24_716_741 - (24_716_741 * 0.01) - _((100 * 0.5) / (3 * 1))_) * (1 - 0.1)
		// for acoin (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (3 * 1))_) * (1 - 0.1)
		// for bcoin (20_000 - (20_000 * 0.01) - _((100 * 0.5) / (3 * 2))_) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(sdk.NewCoins(i.KYVECoin(22_022_603), i.ACoin(8896), i.BCoin(17813)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(Equal(uint64(200)))
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(BeZero())

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(200*i.KYVE - 10_000))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
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
		_ = s.App().BankKeeper.SendCoinsFromModuleToAccount(s.Ctx(), "team", sdk.MustAccAddressFromBech32(i.CHARLIE), s.GetCoinsFromModule("team"))

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
		_ = s.App().BankKeeper.SendCoinsFromModuleToAccount(s.Ctx(), "team", sdk.MustAccAddressFromBech32(i.CHARLIE), s.GetCoinsFromModule("team"))

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
	})
})
