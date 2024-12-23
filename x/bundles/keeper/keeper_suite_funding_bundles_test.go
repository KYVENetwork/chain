package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - funding bundles

* Produce a valid bundle with only one funder
* Produce a valid bundle with multiple funders and same funding amounts
* Produce a valid bundle with multiple funders and different funding amounts
* Produce a valid bundle with multiple funders and different funding amounts where not everyone can afford the funds
* Produce a valid bundle although the only funder can not pay for the full bundle reward
* Produce a valid bundle although multiple funders with same amount can not pay for the bundle reward
* Produce a valid bundle although multiple funders with different amount can not pay for the bundle reward
* Produce a valid bundle although there are no funders at all
* Produce a valid bundle with only one funder and multiple coins
* Produce a valid bundle with multiple funders and multiple coins
* Produce a valid bundle although the only funder can not pay for the full bundle reward

*/

var _ = Describe("funding bundles", Ordered, func() {
	var s *i.KeeperTestSuite
	var initialBalanceAlice sdk.Coins
	var initialBalanceBob sdk.Coins

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
		initialBalanceAlice = s.GetCoinsFromAddress(i.ALICE)
		initialBalanceBob = s.GetCoinsFromAddress(i.BOB)

		params := s.App().FundersKeeper.GetParams(s.Ctx())
		params.MinFundingMultiple = 0
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

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

		// set whitelist
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams([]*funderstypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewIntFromUint64(1 * i.KYVE),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.A_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewIntFromUint64(1 * i.KYVE),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewIntFromUint64(1 * i.KYVE),
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewIntFromUint64(1 * i.KYVE),
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

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce a valid bundle with only one funder", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(10 * i.T_KYVE),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(90 * i.T_KYVE).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)

		// assert individual funds
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(90 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(100 * i.T_KYVE)).String()))
	})

	It("Produce a valid bundle with multiple funders and same funding amounts", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(10 * i.T_KYVE),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(10 * i.T_KYVE),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(180 * i.T_KYVE).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))

		fundingAlice, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		fundingBob, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)

		// assert individual funds
		Expect(fundingAlice.Amounts.String()).To(Equal(i.ACoins(90 * i.T_KYVE).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))
		Expect(fundingBob.Amounts.String()).To(Equal(i.ACoins(90 * i.T_KYVE).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(100 * i.T_KYVE)).String()))

		balanceBob := s.GetCoinsFromAddress(i.BOB)
		Expect(balanceBob.String()).To(Equal(initialBalanceBob.Sub(i.ACoin(100 * i.T_KYVE)).String()))
	})

	It("Produce a valid bundle with multiple funders and different funding amounts", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(150 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(15 * i.T_KYVE),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(50 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(5 * i.T_KYVE),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(180 * i.T_KYVE).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))

		fundingAlice, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		fundingBob, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)

		// assert individual funds
		Expect(fundingAlice.Amounts.String()).To(Equal(i.ACoins(135 * i.T_KYVE).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(15 * i.T_KYVE).String()))
		Expect(fundingBob.Amounts.String()).To(Equal(i.ACoins(45 * i.T_KYVE).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(5 * i.T_KYVE).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(150 * i.T_KYVE)).String()))

		balanceBob := s.GetCoinsFromAddress(i.BOB)
		Expect(balanceBob.String()).To(Equal(initialBalanceBob.Sub(i.ACoin(50 * i.T_KYVE)).String()))
	})

	It("Produce a valid bundle with multiple funders and different funding amounts where not everyone can afford the full funds", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(10 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(50 * i.T_KYVE),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(50 * i.T_KYVE),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.BOB))

		fundingAlice, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		fundingBob, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)

		// assert individual funds
		Expect(fundingAlice.Amounts.IsZero()).To(BeTrue())
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))
		Expect(fundingBob.Amounts.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(10 * i.T_KYVE)).String()))

		balanceBob := s.GetCoinsFromAddress(i.BOB)
		Expect(balanceBob.String()).To(Equal(initialBalanceBob.Sub(i.ACoin(100 * i.T_KYVE)).String()))
	})

	It("Produce a valid bundle although the only funder can not pay for the full bundle reward", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(10 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(20 * i.T_KYVE),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundingAlice, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)

		// assert individual funds
		Expect(fundingAlice.Amounts.IsZero()).To(BeTrue())
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(10 * i.T_KYVE)).String()))
	})

	It("Produce a valid bundle although multiple funders with same amount can not pay for the bundle reward", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(10 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(20 * i.T_KYVE),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(10 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(20 * i.T_KYVE),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert total pool funds
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(10 * i.T_KYVE)).String()))

		balanceBob := s.GetCoinsFromAddress(i.BOB)
		Expect(balanceBob.String()).To(Equal(initialBalanceBob.Sub(i.ACoin(10 * i.T_KYVE)).String()))
	})

	It("Produce a dropped bundle because multiple funders with different amount can not pay for the bundle reward", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(10 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(10 * i.T_KYVE),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(20 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(20 * i.T_KYVE),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundingAlice, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)

		// assert individual funds
		Expect(fundingAlice.Amounts.IsZero()).To(BeTrue())
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(10 * i.T_KYVE)).String()))
	})

	It("Produce a valid bundle although there are no funders at all", func() {
		// ARRANGE
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

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).IsZero()).To(BeTrue())
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))
	})

	It("Produce a valid bundle with only one funder and multiple coins", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.BCoin(2*i.T_KYVE)),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.ACoin(90*i.T_KYVE), i.BCoin(48*i.T_KYVE)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)

		// assert individual funds
		Expect(funding.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(90*i.T_KYVE), i.BCoin(48*i.T_KYVE)).String()))
		Expect(funding.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.BCoin(2*i.T_KYVE)).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))
	})

	It("Produce a valid bundle with multiple funders and multiple coins", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.BCoin(2*i.T_KYVE)),
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.BCoin(100*i.T_KYVE), i.CCoin(200*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.BCoin(10*i.T_KYVE), i.CCoin(20*i.T_KYVE)),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.ACoin(90*i.T_KYVE), i.BCoin(138*i.T_KYVE), i.CCoin(180*i.T_KYVE)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))

		fundingAlice, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		fundingBob, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)

		// assert individual funds
		Expect(fundingAlice.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(90*i.T_KYVE), i.BCoin(48*i.T_KYVE)).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.BCoin(2*i.T_KYVE)).String()))
		Expect(fundingBob.Amounts.String()).To(Equal(sdk.NewCoins(i.BCoin(90*i.T_KYVE), i.CCoin(180*i.T_KYVE)).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(sdk.NewCoins(i.BCoin(10*i.T_KYVE), i.CCoin(20*i.T_KYVE)).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))

		balanceBob := s.GetCoinsFromAddress(i.BOB)
		Expect(balanceBob.String()).To(Equal(initialBalanceBob.Sub(i.BCoin(100*i.T_KYVE), i.CCoin(200*i.T_KYVE)).String()))
	})

	It("Produce a valid bundle although the only funder can not pay for the full bundle reward", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
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

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.ACoin(90 * i.T_KYVE)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))

		fundingAlice, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)

		// assert individual funds
		Expect(fundingAlice.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(90 * i.T_KYVE)).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))

		// assert individual balances
		balanceAlice := s.GetCoinsFromAddress(i.ALICE)
		Expect(balanceAlice.String()).To(Equal(initialBalanceAlice.Sub(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))
	})
})
