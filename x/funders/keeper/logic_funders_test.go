package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_funders.go

* Charge funders once with one coin
* Charge funders once with multiple coins
* Charge funders until one funder runs out of funds
* Charge funders with multiple coins until he is completely out of funds
* Charge funders until all funders run out of funds
* Charge funder with less funds than amount_per_bundle
* Charge funder that has coins which are not in the whitelist
* Charge without fundings
* Check if the lowest funding is returned correctly with one coin
* Check if the lowest funding is returned correctly with multiple coins
* Check if the lowest funding is returned correctly with coins which are not whitelisted

*/

var _ = Describe("logic_funders.go", Ordered, func() {
	s := i.NewCleanChain()

	var whitelist []*funderstypes.WhitelistCoinEntry

	BeforeEach(func() {
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

		// set whitelist
		whitelist = []*funderstypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.A_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(3),
			},
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams(whitelist, 20))

		params := s.App().FundersKeeper.GetParams(s.Ctx())
		params.MinFundingMultiple = 5
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

		// create funder
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.BOB,
			Moniker: "Bob",
		})

		// fund pool
		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(50 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(10 * i.T_KYVE),
		})

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		Expect(fundersBalance.String()).To(Equal(i.ACoins(150 * i.T_KYVE).String()))
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Charge funders once with one coin", func() {
		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
		Expect(err).NotTo(HaveOccurred())

		// ASSERT
		Expect(payout.String()).To(Equal(i.ACoins(11 * i.T_KYVE).String()))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.String()).To(Equal(i.ACoins(99 * i.T_KYVE).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.String()).To(Equal(i.ACoins(40 * i.T_KYVE).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
		Expect(fundingState.ActiveFunderAddresses[1]).To(Equal(i.BOB))

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.String()).To(Equal(i.ACoins(139 * i.T_KYVE).String()))
		Expect(poolBalance.String()).To(Equal(i.ACoins(11 * i.T_KYVE).String()))
	})

	It("Charge funders once with multiple coins", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.BCoins(1000 * i.T_KYVE),
			AmountsPerBundle: i.BCoins(20 * i.T_KYVE),
		})
		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.CCoins(100 * i.T_KYVE),
			AmountsPerBundle: i.CCoins(2 * i.T_KYVE),
		})

		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
		Expect(err).NotTo(HaveOccurred())

		// ASSERT
		Expect(payout.String()).To(Equal(sdk.NewCoins(i.ACoin(11*i.T_KYVE), i.BCoin(20*i.T_KYVE), i.CCoin(2*i.T_KYVE)).String()))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(99*i.T_KYVE), i.BCoin(980*i.T_KYVE)).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(20*i.T_KYVE)).String()))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(40*i.T_KYVE), i.CCoin(98*i.T_KYVE)).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.CCoin(2*i.T_KYVE)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
		Expect(fundingState.ActiveFunderAddresses[1]).To(Equal(i.BOB))

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.String()).To(Equal(sdk.NewCoins(i.ACoin(139*i.T_KYVE), i.BCoin(980*i.T_KYVE), i.CCoin(98*i.T_KYVE)).String()))
		Expect(poolBalance.String()).To(Equal(sdk.NewCoins(i.ACoin(11*i.T_KYVE), i.BCoin(20*i.T_KYVE), i.CCoin(2*i.T_KYVE)).String()))
	})

	It("Charge funders until one funder runs out of funds", func() {
		// ACT
		for range [5]struct{}{} {
			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout.String()).To(Equal(i.ACoins(11 * i.T_KYVE).String()))
		}

		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.String()).To(Equal(i.ACoins(95 * i.T_KYVE).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(5 * i.T_KYVE).String()))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.IsZero()).To(BeTrue())
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.String()).To(Equal(i.ACoins(95 * i.T_KYVE).String()))
		Expect(poolBalance.String()).To(Equal(i.ACoins(55 * i.T_KYVE).String()))
	})

	It("Charge funders with multiple coins until he is completely out of funds", func() {
		// ACT
		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.BCoins(1000 * i.T_KYVE),
			AmountsPerBundle: i.BCoins(20 * i.T_KYVE),
		})
		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.CCoins(100 * i.T_KYVE),
			AmountsPerBundle: i.CCoins(10 * i.T_KYVE),
		})

		for range [5]struct{}{} {
			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout.String()).To(Equal(sdk.NewCoins(i.ACoin(11*i.T_KYVE), i.BCoin(20*i.T_KYVE), i.CCoin(10*i.T_KYVE)).String()))
		}

		for range [5]struct{}{} {
			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout.String()).To(Equal(sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(20*i.T_KYVE), i.CCoin(10*i.T_KYVE)).String()))
		}

		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(90*i.T_KYVE), i.BCoin(800*i.T_KYVE)).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(10*i.T_KYVE), i.BCoin(200*i.T_KYVE)).String()))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.IsZero()).To(BeTrue())
		Expect(fundingBob.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(50*i.T_KYVE), i.CCoin(100*i.T_KYVE)).String()))

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.String()).To(Equal(sdk.NewCoins(i.ACoin(90*i.T_KYVE), i.BCoin(800*i.T_KYVE)).String()))
		Expect(poolBalance.String()).To(Equal(sdk.NewCoins(i.ACoin(60*i.T_KYVE), i.BCoin(200*i.T_KYVE), i.CCoin(100*i.T_KYVE)).String()))
	})

	It("Charge funders until all funders run out of funds", func() {
		// ARRANGE
		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		funding.AmountsPerBundle = i.ACoins(10 * i.T_KYVE)
		s.App().FundersKeeper.SetFunding(s.Ctx(), &funding)

		// ACT / ASSERT
		for range [5]struct{}{} {
			fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
			Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))

			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout.String()).To(Equal(i.ACoins(20 * i.T_KYVE).String()))
		}
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.IsZero()).To(BeTrue())
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))

		for range [5]struct{}{} {
			fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
			Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))

			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))
		}
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundingAlice, foundAlice = s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.IsZero()).To(BeTrue())
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(100 * i.T_KYVE).String()))

		fundingBob, foundBob = s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.IsZero()).To(BeTrue())
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))

		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
		Expect(err).NotTo(HaveOccurred())
		Expect(payout.IsZero()).To(BeTrue())

		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.IsZero()).To(BeTrue())
		Expect(poolBalance.String()).To(Equal(i.ACoins(150 * i.T_KYVE).String()))
	})

	It("Charge funder with less funds than amount_per_bundle", func() {
		// ARRANGE
		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		funding.AmountsPerBundle = i.ACoins(105 * i.T_KYVE)
		s.App().FundersKeeper.SetFunding(s.Ctx(), &funding)

		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
		Expect(err).NotTo(HaveOccurred())
		Expect(payout.String()).To(Equal(i.ACoins(110 * i.T_KYVE).String()))

		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.BOB))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.IsZero()).To(BeTrue())
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(100 * i.T_KYVE).String()))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.String()).To(Equal(i.ACoins(40 * i.T_KYVE).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.String()).To(Equal(i.ACoins(40 * i.T_KYVE).String()))
		Expect(poolBalance.String()).To(Equal(i.ACoins(110 * i.T_KYVE).String()))
	})

	It("Charge funder that has coins which are not in the whitelist", func() {
		// ARRANGE
		whitelist = []*funderstypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(3),
			},
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams(whitelist, 20))

		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
		Expect(err).NotTo(HaveOccurred())

		// ASSERT
		Expect(payout.String()).To(Equal(i.ACoins(11 * i.T_KYVE).String()))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.String()).To(Equal(i.ACoins(99 * i.T_KYVE).String()))
		Expect(fundingAlice.TotalFunded.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.String()).To(Equal(i.ACoins(40 * i.T_KYVE).String()))
		Expect(fundingBob.TotalFunded.String()).To(Equal(i.ACoins(10 * i.T_KYVE).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
		Expect(fundingState.ActiveFunderAddresses[1]).To(Equal(i.BOB))

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.String()).To(Equal(i.ACoins(139 * i.T_KYVE).String()))
		Expect(poolBalance.String()).To(Equal(i.ACoins(11 * i.T_KYVE).String()))
	})

	It("Charge without fundings", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})
		s.RunTxFundersSuccess(&funderstypes.MsgDefundPool{
			Creator: i.BOB,
			PoolId:  0,
			Amounts: i.ACoins(50 * i.T_KYVE),
		})

		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
		Expect(payout.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amounts.IsZero()).To(BeTrue())
		Expect(fundingAlice.TotalFunded.IsZero()).To(BeTrue())

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amounts.IsZero()).To(BeTrue())
		Expect(fundingBob.TotalFunded.IsZero()).To(BeTrue())

		fundersBalance := s.GetCoinsFromModule(funderstypes.ModuleName)
		poolBalance := s.GetCoinsFromModule(pooltypes.ModuleName)
		Expect(fundersBalance.IsZero()).To(BeTrue())
		Expect(poolBalance.IsZero()).To(BeTrue())
	})

	It("Check if the lowest funding is returned correctly with one coin", func() {
		whitelist := []*funderstypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:  i.A_DENOM,
				CoinWeight: math.LegacyNewDec(1),
			},
			{
				CoinDenom:  i.B_DENOM,
				CoinWeight: math.LegacyNewDec(2),
			},
			{
				CoinDenom:  i.C_DENOM,
				CoinWeight: math.LegacyNewDec(3),
			},
		}

		fundings := []funderstypes.Funding{
			{
				FunderAddress:    i.DUMMY[0],
				PoolId:           0,
				Amounts:          i.ACoins(1000 * i.T_KYVE),
				AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
			},
			{
				FunderAddress:    i.DUMMY[1],
				PoolId:           0,
				Amounts:          i.ACoins(1100 * i.T_KYVE),
				AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
			},
			{
				FunderAddress:    i.DUMMY[2],
				PoolId:           0,
				Amounts:          i.ACoins(900 * i.T_KYVE),
				AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
			},
		}

		getLowestFunding, err := s.App().FundersKeeper.GetLowestFunding(fundings, whitelist)
		Expect(err).NotTo(HaveOccurred())
		Expect(getLowestFunding.FunderAddress).To(Equal(i.DUMMY[2]))
	})

	It("Check if the lowest funding is returned correctly with multiple coins", func() {
		whitelist := []*funderstypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:  i.A_DENOM,
				CoinWeight: math.LegacyNewDec(1),
			},
			{
				CoinDenom:  i.B_DENOM,
				CoinWeight: math.LegacyNewDec(2),
			},
			{
				CoinDenom:  i.C_DENOM,
				CoinWeight: math.LegacyNewDec(3),
			},
		}

		fundings := []funderstypes.Funding{
			{
				FunderAddress:    i.DUMMY[0],
				PoolId:           0,
				Amounts:          sdk.NewCoins(i.ACoin(1000*i.T_KYVE), i.BCoin(500*i.T_KYVE), i.CCoin(200*i.T_KYVE)),
				AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(1*i.T_KYVE), i.CCoin(1)),
			},
			{
				FunderAddress:    i.DUMMY[1],
				PoolId:           0,
				Amounts:          sdk.NewCoins(i.ACoin(1100*i.T_KYVE), i.BCoin(600*i.T_KYVE), i.CCoin(5*i.T_KYVE)),
				AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(1*i.T_KYVE), i.CCoin(1)),
			},
			{
				FunderAddress:    i.DUMMY[2],
				PoolId:           0,
				Amounts:          sdk.NewCoins(i.ACoin(500*i.T_KYVE), i.CCoin(700*i.T_KYVE)),
				AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.CCoin(1)),
			},
		}

		getLowestFunding, err := s.App().FundersKeeper.GetLowestFunding(fundings, whitelist)
		Expect(err).NotTo(HaveOccurred())
		Expect(getLowestFunding.FunderAddress).To(Equal(i.DUMMY[1]))
	})

	It("Check if the lowest funding is returned correctly with coins which are not whitelisted", func() {
		whitelist := []*funderstypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:  i.A_DENOM,
				CoinWeight: math.LegacyNewDec(1),
			},
			{
				CoinDenom:  i.B_DENOM,
				CoinWeight: math.LegacyNewDec(2),
			},
		}

		fundings := []funderstypes.Funding{
			{
				FunderAddress:    i.DUMMY[0],
				PoolId:           0,
				Amounts:          sdk.NewCoins(i.ACoin(1000*i.T_KYVE), i.BCoin(500*i.T_KYVE), i.CCoin(200*i.T_KYVE)),
				AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(1*i.T_KYVE), i.CCoin(1)),
			},
			{
				FunderAddress:    i.DUMMY[1],
				PoolId:           0,
				Amounts:          sdk.NewCoins(i.ACoin(1100*i.T_KYVE), i.BCoin(600*i.T_KYVE), i.CCoin(5*i.T_KYVE)),
				AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(1*i.T_KYVE), i.CCoin(1)),
			},
			{
				FunderAddress:    i.DUMMY[2],
				PoolId:           0,
				Amounts:          sdk.NewCoins(i.ACoin(500*i.T_KYVE), i.CCoin(700*i.T_KYVE)),
				AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.CCoin(1)),
			},
		}

		getLowestFunding, err := s.App().FundersKeeper.GetLowestFunding(fundings, whitelist)
		Expect(err).NotTo(HaveOccurred())
		Expect(getLowestFunding.FunderAddress).To(Equal(i.DUMMY[2]))
	})
})
