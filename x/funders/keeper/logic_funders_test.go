package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_funders.go

* Charge funders once
* Charge funders until one funder runs out of funds
* Charge funders until all funders run out of funds
* Charge funder with less funds than amount_per_bundle
* Charge without fundings
* Check if the lowest funding is returned correctly

*/

var _ = Describe("logic_funders.go", Ordered, func() {
	s := i.NewCleanChain()
	fundersModuleAcc := s.App().AccountKeeper.GetModuleAccount(s.Ctx(), funderstypes.ModuleName).GetAddress()
	poolModuleAcc := s.App().AccountKeeper.GetModuleAccount(s.Ctx(), pooltypes.ModuleName).GetAddress()

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
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})
		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          50 * i.KYVE,
			AmountPerBundle: 10 * i.KYVE,
		})

		fundersBalance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAcc, globaltypes.Denom).Amount.Uint64()
		Expect(fundersBalance).To(Equal(150 * i.KYVE))
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Charge funders once", func() {
		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0)
		Expect(err).NotTo(HaveOccurred())

		// ASSERT
		Expect(payout).To(Equal(11 * i.KYVE))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amount).To(Equal(99 * i.KYVE))
		Expect(fundingAlice.TotalFunded).To(Equal(1 * i.KYVE))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amount).To(Equal(40 * i.KYVE))
		Expect(fundingBob.TotalFunded).To(Equal(10 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
		Expect(fundingState.ActiveFunderAddresses[1]).To(Equal(i.BOB))

		fundersBalance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAcc, globaltypes.Denom).Amount.Uint64()
		poolBalance := s.App().BankKeeper.GetBalance(s.Ctx(), poolModuleAcc, globaltypes.Denom).Amount.Uint64()
		Expect(fundersBalance).To(Equal(139 * i.KYVE))
		Expect(poolBalance).To(Equal(11 * i.KYVE))
	})

	It("Charge funders until one funder runs out of funds", func() {
		// ACT
		for range [5]struct{}{} {
			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout).To(Equal(11 * i.KYVE))
		}

		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amount).To(Equal(95 * i.KYVE))
		Expect(fundingAlice.TotalFunded).To(Equal(5 * i.KYVE))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amount).To(Equal(0 * i.KYVE))
		Expect(fundingBob.TotalFunded).To(Equal(50 * i.KYVE))

		fundersBalance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAcc, globaltypes.Denom).Amount.Uint64()
		poolBalance := s.App().BankKeeper.GetBalance(s.Ctx(), poolModuleAcc, globaltypes.Denom).Amount.Uint64()
		Expect(fundersBalance).To(Equal(95 * i.KYVE))
		Expect(poolBalance).To(Equal(55 * i.KYVE))
	})

	It("Charge funders until all funders run out of funds", func() {
		// ARRANGE
		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		funding.AmountPerBundle = 10 * i.KYVE
		s.App().FundersKeeper.SetFunding(s.Ctx(), &funding)

		// ACT / ASSERT
		for range [5]struct{}{} {
			fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
			Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))

			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout).To(Equal(20 * i.KYVE))
		}
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amount).To(Equal(50 * i.KYVE))
		Expect(fundingAlice.TotalFunded).To(Equal(50 * i.KYVE))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amount).To(Equal(0 * i.KYVE))
		Expect(fundingBob.TotalFunded).To(Equal(50 * i.KYVE))

		for range [5]struct{}{} {
			fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
			Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))

			payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(payout).To(Equal(10 * i.KYVE))
		}
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundingAlice, foundAlice = s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amount).To(Equal(0 * i.KYVE))
		Expect(fundingAlice.TotalFunded).To(Equal(100 * i.KYVE))

		fundingBob, foundBob = s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amount).To(Equal(0 * i.KYVE))
		Expect(fundingBob.TotalFunded).To(Equal(50 * i.KYVE))

		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(0 * i.KYVE))

		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundersBalance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAcc, globaltypes.Denom).Amount.Uint64()
		poolBalance := s.App().BankKeeper.GetBalance(s.Ctx(), poolModuleAcc, globaltypes.Denom).Amount.Uint64()
		Expect(fundersBalance).To(Equal(0 * i.KYVE))
		Expect(poolBalance).To(Equal(150 * i.KYVE))
	})

	It("Charge funder with less funds than amount_per_bundle", func() {
		// ARRANGE
		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		funding.AmountPerBundle = 105 * i.KYVE
		s.App().FundersKeeper.SetFunding(s.Ctx(), &funding)

		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(110 * i.KYVE))

		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.BOB))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amount).To(Equal(0 * i.KYVE))
		Expect(fundingAlice.TotalFunded).To(Equal(100 * i.KYVE))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amount).To(Equal(40 * i.KYVE))
		Expect(fundingBob.TotalFunded).To(Equal(10 * i.KYVE))

		fundersBalance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAcc, globaltypes.Denom).Amount.Uint64()
		poolBalance := s.App().BankKeeper.GetBalance(s.Ctx(), poolModuleAcc, globaltypes.Denom).Amount.Uint64()
		Expect(fundersBalance).To(Equal(40 * i.KYVE))
		Expect(poolBalance).To(Equal(110 * i.KYVE))
	})

	It("Charge without fundings", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  100 * i.KYVE,
		})
		s.RunTxFundersSuccess(&funderstypes.MsgDefundPool{
			Creator: i.BOB,
			PoolId:  0,
			Amount:  50 * i.KYVE,
		})

		// ACT
		payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0)

		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(0))

		fundingAlice, foundAlice := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(foundAlice).To(BeTrue())
		Expect(fundingAlice.Amount).To(Equal(0 * i.KYVE))
		Expect(fundingAlice.TotalFunded).To(Equal(0 * i.KYVE))

		fundingBob, foundBob := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(foundBob).To(BeTrue())
		Expect(fundingBob.Amount).To(Equal(0 * i.KYVE))
		Expect(fundingBob.TotalFunded).To(Equal(0 * i.KYVE))

		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(0 * i.KYVE))
		fundersBalance := s.App().BankKeeper.GetBalance(s.Ctx(), fundersModuleAcc, globaltypes.Denom).Amount.Uint64()
		poolBalance := s.App().BankKeeper.GetBalance(s.Ctx(), poolModuleAcc, globaltypes.Denom).Amount.Uint64()
		Expect(fundersBalance).To(Equal(0 * i.KYVE))
		Expect(poolBalance).To(Equal(0 * i.KYVE))
	})

	It("Check if the lowest funding is returned correctly", func() {
		fundings := []funderstypes.Funding{
			{
				FunderAddress:   i.DUMMY[0],
				PoolId:          0,
				Amount:          1000 * i.KYVE,
				AmountPerBundle: 1 * i.KYVE,
			},
			{
				FunderAddress:   i.DUMMY[1],
				PoolId:          0,
				Amount:          900 * i.KYVE,
				AmountPerBundle: 1 * i.KYVE,
			},
			{
				FunderAddress:   i.DUMMY[2],
				PoolId:          0,
				Amount:          1100 * i.KYVE,
				AmountPerBundle: 1 * i.KYVE,
			},
		}

		getLowestFunding, err := s.App().FundersKeeper.GetLowestFunding(fundings)
		Expect(err).NotTo(HaveOccurred())
		Expect(getLowestFunding.Amount).To(Equal(900 * i.KYVE))
	})
})
