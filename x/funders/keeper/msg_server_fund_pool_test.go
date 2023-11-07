package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_fund_pool.go

* Fund a pool with 100 $KYVE
* Fund additional 50 $KYVE to an existing funding with 100 $KYVE
* Try to fund more $KYVE than available in balance
* Fund with a new funder less $KYVE than the existing one
* Fund with a new funder more $KYVE than the existing one
* Try to fund with a non-existent funder
* Try to fund less $KYVE than the lowest funder with full funding slots
* Fund more $KYVE than the lowest funder with full funding slots
* Refund a funding as the lowest funder
* Try to fund a non-existent pool
* Try to fund below the minimum amount
* Try to fund below the minimum amount per bundle
* Try to fund without fulfilling min_funding_multiple

*/

var _ = Describe("msg_server_fund_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalance := s.GetBalanceFromAddress(i.ALICE)

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

		// create funder
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.BOB,
			Moniker: "Bob",
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Fund a pool with 100 $KYVE", func() {
		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.FunderAddress).To(Equal(i.ALICE))
		Expect(funding.PoolId).To(Equal(uint64(0)))
		Expect(funding.Amount).To(Equal(100 * i.KYVE))
		Expect(funding.AmountPerBundle).To(Equal(1 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(0 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
	})

	It("Fund additional 50 $KYVE to an existing funding with 100 $KYVE", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		Expect(initialBalance - balanceAfter).To(Equal(150 * i.KYVE))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.FunderAddress).To(Equal(i.ALICE))
		Expect(funding.PoolId).To(Equal(uint64(0)))
		Expect(funding.Amount).To(Equal(150 * i.KYVE))
		Expect(funding.AmountPerBundle).To(Equal(1 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(0 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Try to fund more $KYVE than available in balance", func() {
		// ACT
		currentBalance := s.GetBalanceFromAddress(i.ALICE)

		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          currentBalance + 1,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceAfter).To(BeZero())

		_, found := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(found).To(BeFalse())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))
	})

	It("Fund with a new funder less $KYVE than the existing one", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          50 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.BOB)
		Expect(initialBalance - balanceAfter).To(Equal(50 * i.KYVE))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(funding.Amount).To(Equal(50 * i.KYVE))
		Expect(funding.AmountPerBundle).To(Equal(1 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(0 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(2))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.BOB))
	})

	It("Fund with a new funder more $KYVE than the existing one", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.BOB)
		Expect(initialBalance - balanceAfter).To(Equal(200 * i.KYVE))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(funding.Amount).To(Equal(200 * i.KYVE))
		Expect(funding.AmountPerBundle).To(Equal(1 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(0 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(2))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Try to fund with a non-existent funder", func() {
		// ASSERT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:         i.CHARLIE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})
	})

	It("Try to fund less $KYVE than the lowest funder with full funding slots", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		for a := 0; a < funderstypes.MaxFunders-1; a++ {
			s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
				Creator: i.DUMMY[a],
				Moniker: i.DUMMY[a],
			})
			// fill remaining funding slots
			s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
				Creator:         i.DUMMY[a],
				PoolId:          0,
				Amount:          1000 * i.KYVE,
				AmountPerBundle: 1 * i.KYVE,
			})
		}

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		balanceAfter := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		// ACT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          50 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ASSERT
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		Expect(activeFundings).To(HaveLen(50))
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Fund more $KYVE than the lowest funder with full funding slots", func() {
		// ARRANGE
		initialBalanceBob := s.GetBalanceFromAddress(i.BOB)
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		for a := 0; a < funderstypes.MaxFunders-1; a++ {
			s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
				Creator: i.DUMMY[a],
				Moniker: i.DUMMY[a],
			})
			// fill remaining funding slots
			s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
				Creator:         i.DUMMY[a],
				PoolId:          0,
				Amount:          1000 * i.KYVE,
				AmountPerBundle: 1 * i.KYVE,
			})
		}
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          200 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})
		x := s.GetBalanceFromAddress(i.BOB)
		Expect(initialBalanceBob - x).To(Equal(200 * i.KYVE))
		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		Expect(activeFundings).To(HaveLen(50))
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.BOB))

		balanceEnd := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceEnd).To(BeZero())

		balanceAfterBob := s.GetBalanceFromAddress(i.BOB)
		Expect(initialBalanceBob - balanceAfterBob).To(Equal(200 * i.KYVE))
	})

	It("Refund a funding as the lowest funder", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		for a := 0; a < funderstypes.MaxFunders-1; a++ {
			s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
				Creator: i.DUMMY[a],
				Moniker: i.DUMMY[a],
			})
			// fill remaining funding slots
			s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
				Creator:         i.DUMMY[a],
				PoolId:          0,
				Amount:          1000 * i.KYVE,
				AmountPerBundle: 1 * i.KYVE,
			})
		}

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		Expect(activeFundings).To(HaveLen(50))
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          50 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ASSERT
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		balanceEnd := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceEnd).To(Equal(150 * i.KYVE))
	})

	It("Try to fund a non-existent pool", func() {
		// ASSERT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          1,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})
	})

	It("Try to fund below the minimum amount", func() {
		// ASSERT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          1,
			AmountPerBundle: 1 * i.KYVE,
		})
	})

	It("Try to fund below the minimum amount per bundle", func() {
		// ASSERT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1,
		})
	})

	It("Try to fund without fulfilling min_funding_multiple", func() {
		// ASSERT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          2 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})
		Expect(err.Error()).To(Equal("per_bundle_amount (1000000000kyve) times min_funding_multiple (1000000000) is smaller than funded_amount (2000000000kyve): invalid request"))
	})
})
