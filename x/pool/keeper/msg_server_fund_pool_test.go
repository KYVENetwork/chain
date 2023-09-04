package keeper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_fund_pool.go

* Create funder by funding a pool with 100 $KYVE
* Fund additional 50 $KYVE to an existing funder with 100 $KYVE
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
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name: "PoolTest",
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Create funder by funding a pool with 100 $KYVE", func() {
		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.TotalFunds).To(Equal(100 * i.KYVE))

		funderAmount := pool.GetFunderAmount(i.ALICE)

		Expect(funderAmount).To(Equal(100 * i.KYVE))
		Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
		Expect(pool.GetLowestFunder().Amount).To(Equal(100 * i.KYVE))
	})

	It("Fund additional 50 $KYVE to an existing funder with 100 $KYVE", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(Equal(150 * i.KYVE))

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.TotalFunds).To(Equal(150 * i.KYVE))

		funderAmount := pool.GetFunderAmount(i.ALICE)

		Expect(funderAmount).To(Equal(150 * i.KYVE))
		Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
		Expect(pool.GetLowestFunder().Amount).To(Equal(150 * i.KYVE))
	})

	It("Try to fund more $KYVE than available in balance", func() {
		// ACT
		currentBalance := s.GetBalanceFromAddress(i.ALICE)

		s.RunTxPoolError(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  currentBalance + 1,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(BeZero())

		Expect(pool.Funders).To(BeEmpty())
		Expect(pool.TotalFunds).To(BeZero())

		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(0 * i.KYVE))
		Expect(pool.GetLowestFunder().Address).To(Equal(""))
		Expect(pool.GetLowestFunder().Amount).To(Equal(0 * i.KYVE))
	})

	It("Fund with a new funder less $KYVE than the existing one", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.BOB)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(Equal(50 * i.KYVE))

		Expect(pool.Funders).To(HaveLen(2))
		Expect(pool.TotalFunds).To(Equal(150 * i.KYVE))

		funderAmount := pool.GetFunderAmount(i.BOB)

		Expect(funderAmount).To(Equal(50 * i.KYVE))
		Expect(pool.GetLowestFunder().Address).To(Equal(i.BOB))
		Expect(pool.GetLowestFunder().Amount).To(Equal(50 * i.KYVE))
	})

	It("Fund with a new funder more $KYVE than the existing one", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  200 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.BOB)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(Equal(200 * i.KYVE))

		Expect(pool.Funders).To(HaveLen(2))
		Expect(pool.TotalFunds).To(Equal(300 * i.KYVE))

		funderAmount := pool.GetFunderAmount(i.BOB)
		Expect(funderAmount).To(Equal(200 * i.KYVE))

		Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
		Expect(pool.GetLowestFunder().Amount).To(Equal(100 * i.KYVE))
	})

	It("Try to fund less $KYVE than the lowest funder with full funding slots", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		for a := 0; a < 49; a++ {
			// fill remaining funding slots
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[a],
				Id:      0,
				Amount:  1000 * i.KYVE,
			})
		}

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Funders).To(HaveLen(50))
		Expect(pool.TotalFunds).To(Equal(49_100 * i.KYVE))

		Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
		Expect(pool.GetLowestFunder().Amount).To(Equal(100 * i.KYVE))

		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		// ACT
		s.RunTxPoolError(&pooltypes.MsgFundPool{
			Creator: i.DUMMY[49],
			Id:      0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		Expect(pool.Funders).To(HaveLen(50))
		Expect(pool.TotalFunds).To(Equal(49_100 * i.KYVE))

		Expect(pool.GetFunderAmount(i.DUMMY[49])).To(BeZero())
		Expect(pool.GetLowestFunder().Address).To(Equal(i.ALICE))
		Expect(pool.GetLowestFunder().Amount).To(Equal(100 * i.KYVE))
	})

	It("Fund more $KYVE than the lowest funder with full funding slots", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		for a := 0; a < 49; a++ {
			// fill remaining funding slots
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[a],
				Id:      0,
				Amount:  1000 * i.KYVE,
			})
		}

		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.DUMMY[49],
			Id:      0,
			Amount:  200 * i.KYVE,
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Funders).To(HaveLen(50))
		Expect(pool.TotalFunds).To(Equal(49_200 * i.KYVE))

		Expect(pool.GetFunderAmount(i.DUMMY[49])).To(Equal(200 * i.KYVE))
		Expect(pool.GetLowestFunder().Address).To(Equal(i.DUMMY[49]))
		Expect(pool.GetLowestFunder().Amount).To(Equal(200 * i.KYVE))

		balanceAfter := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceAfter).To(BeZero())
	})
})
