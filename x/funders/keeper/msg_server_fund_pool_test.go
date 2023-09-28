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
* Try to fund less $KYVE than the lowest funder with full funding slots
* Try to fund more $KYVE than the lowest funder with full funding slots

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
		s.RunTxFundersSuccess(msg)

		// create funder
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "moniker",
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
		Expect(fundingState.TotalAmount).To(Equal(100 * i.KYVE))
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
		Expect(fundingState.TotalAmount).To(Equal(150 * i.KYVE))
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
		Expect(fundingState.TotalAmount).To(Equal(uint64(0)))
	})
	// TODO: fix this test
	//It("Fund with a new funder less $KYVE than the existing one", func() {
	//	// ARRANGE
	//	s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//		Creator: i.ALICE,
	//		PoolId:  0,
	//		Amount:  100 * i.KYVE,
	//	})
	//
	//	// ACT
	//	s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//		Creator: i.BOB,
	//		PoolId:  0,
	//		Amount:  50 * i.KYVE,
	//	})
	//
	//	// ASSERT
	//	balanceAfter := s.GetBalanceFromAddress(i.BOB)
	//
	//	pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
	//
	//	Expect(initialBalance - balanceAfter).To(Equal(50 * i.KYVE))
	//
	//	Expect(pool.Funders).To(HaveLen(2))
	//	Expect(pool.TotalFunds).To(Equal(150 * i.KYVE))
	//
	//	funderAmount := pool.GetFunderAmount(i.BOB)
	//
	//	Expect(funderAmount).To(Equal(50 * i.KYVE))
	//	Expect(pool.GetLowestFunder().Address).To(Equal(i.BOB))
	//	Expect(pool.GetLowestFunder().Amount).To(Equal(50 * i.KYVE))
	//})
	//
	//It("Fund with a new funder more $KYVE than the existing one", func() {
	//	// ARRANGE
	//	s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//		Creator: i.ALICE,
	//		PoolId:  0,
	//		Amount:  100 * i.KYVE,
	//	})
	//
	//	// ACT
	//	s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//		Creator: i.BOB,
	//		PoolId:  0,
	//		Amount:  200 * i.KYVE,
	//	})
	//
	//	// ASSERT
	//	balanceAfter := s.GetBalanceFromAddress(i.BOB)
	//
	//	pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
	//
	//	Expect(initialBalance - balanceAfter).To(Equal(200 * i.KYVE))
	//
	//	Expect(pool.Funders).To(HaveLen(2))
	//	Expect(pool.TotalFunds).To(Equal(300 * i.KYVE))
	//
	//	funderAmount := pool.GetFunderAmount(i.BOB)
	//	Expect(funderAmount).To(Equal(200 * i.KYVE))
	//
	//	Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
	//	Expect(pool.GetLowestFunder().Amount).To(Equal(100 * i.KYVE))
	//})
	//
	//It("Try to fund less $KYVE than the lowest funder with full funding slots", func() {
	//	// ARRANGE
	//	s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//		Creator: i.ALICE,
	//		PoolId:  0,
	//		Amount:  100 * i.KYVE,
	//	})
	//
	//	for a := 0; a < 49; a++ {
	//		// fill remaining funding slots
	//		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//			Creator: i.DUMMY[a],
	//			PoolId:  0,
	//			Amount:  1000 * i.KYVE,
	//		})
	//	}
	//
	//	pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
	//
	//	Expect(pool.Funders).To(HaveLen(50))
	//	Expect(pool.TotalFunds).To(Equal(49_100 * i.KYVE))
	//
	//	Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
	//	Expect(pool.GetLowestFunder().Amount).To(Equal(100 * i.KYVE))
	//
	//	balanceAfter := s.GetBalanceFromAddress(i.ALICE)
	//
	//	Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))
	//
	//	// ACT
	//	s.RunTxFundersError(&funderstypes.MsgFundPool{
	//		Creator: i.DUMMY[49],
	//		PoolId:  0,
	//		Amount:  50 * i.KYVE,
	//	})
	//
	//	// ASSERT
	//	Expect(pool.Funders).To(HaveLen(50))
	//	Expect(pool.TotalFunds).To(Equal(49_100 * i.KYVE))
	//
	//	Expect(pool.GetFunderAmount(i.DUMMY[49])).To(BeZero())
	//	Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
	//	Expect(pool.GetLowestFunder().Amount).To(Equal(100 * i.KYVE))
	//})
	//
	//It("Fund more $KYVE than the lowest funder with full funding slots", func() {
	//	// ARRANGE
	//	s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//		Creator: i.ALICE,
	//		PoolId:  0,
	//		Amount:  100 * i.KYVE,
	//	})
	//
	//	for a := 0; a < 49; a++ {
	//		// fill remaining funding slots
	//		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//			Creator: i.DUMMY[a],
	//			PoolId:  0,
	//			Amount:  1000 * i.KYVE,
	//		})
	//	}
	//
	//	// ACT
	//	s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
	//		Creator: i.DUMMY[49],
	//		PoolId:  0,
	//		Amount:  200 * i.KYVE,
	//	})
	//
	//	// ASSERT
	//	pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
	//
	//	Expect(pool.Funders).To(HaveLen(50))
	//	Expect(pool.TotalFunds).To(Equal(49_200 * i.KYVE))
	//
	//	Expect(pool.GetFunderAmount(i.DUMMY[49])).To(Equal(200 * i.KYVE))
	//	Expect(pool.GetLowestFunder().Address).To(Equal(i.DUMMY[49]))
	//	Expect(pool.GetLowestFunder().Amount).To(Equal(200 * i.KYVE))
	//
	//	balanceAfter := s.GetBalanceFromAddress(i.ALICE)
	//	Expect(initialBalance - balanceAfter).To(BeZero())
	//})
})
