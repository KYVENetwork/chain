package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/util"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_funders.go

* Add funders; check total sum
* Add multiple funders; check total sum
* Remove funder
* Remove funder by defunding everything
* Charge Funders with equal amounts
* Charge Funders test remainder
* Charge exactly the lowest funder amount
* Kick out multiple lowest funders
* Charge more than pool has funds
* Charge pool which has no funds at all

*/

func chargeFunders(s *i.KeeperTestSuite, amount uint64) (payout uint64, err error) {
	payout, err = s.App().PoolKeeper.ChargeFundersOfPool(s.Ctx(), 0, amount)
	if err != nil {
		return 0, err
	}

	if err := util.TransferFromModuleToAddress(s.App().BankKeeper, s.Ctx(), pooltypes.ModuleName, i.BURNER, payout); err != nil {
		return 0, err
	}

	return payout, err
}

func fundersCheck(pool *pooltypes.Pool) {
	poolFunds := uint64(0)
	funders := make(map[string]bool)
	for _, funder := range pool.Funders {
		Expect(funders[funder.Address]).To(BeFalse())
		funders[funder.Address] = true
		poolFunds += funder.Amount
	}
	Expect(pool.TotalFunds).To(Equal(poolFunds))
}

var _ = Describe("logic_funders.go", Ordered, func() {
	s := i.NewCleanChain()
	var pool *pooltypes.Pool

	BeforeEach(func() {
		s = i.NewCleanChain()
		pool = &pooltypes.Pool{
			Name:           "PoolTest",
			MaxBundleSize:  100,
			StartKey:       "0",
			MinDelegation:  100 * i.KYVE,
			UploadInterval: 60,
			OperatingCost:  10_000,
			UpgradePlan:    &pooltypes.UpgradePlan{},
		}

		s.App().PoolKeeper.AppendPool(s.Ctx(), *pool)
	})

	AfterEach(func() {
		fundersCheck(pool)
		s.PerformValidityChecks()
	})

	It("Add funders; check total sum", func() {
		// ACT
		pool.AddAmountToFunder(i.ALICE, 1000)
		pool.AddAmountToFunder(i.ALICE, 2000)
		pool.AddAmountToFunder(i.ALICE, 0)
		pool.AddAmountToFunder(i.BOB, 0)
		pool.AddAmountToFunder(i.ALICE, 10)

		// ASSERT
		Expect(pool.TotalFunds).To(Equal(uint64(3010)))
		Expect(pool.Funders).To(HaveLen(1))
	})

	It("Add multiple funders; check total sum", func() {
		// ACT
		pool.AddAmountToFunder(i.ALICE, 1000)
		pool.AddAmountToFunder(i.ALICE, 2000)
		pool.AddAmountToFunder(i.ALICE, 0)
		pool.AddAmountToFunder(i.BOB, 1000)
		pool.AddAmountToFunder(i.ALICE, 10)

		// ASSERT
		Expect(pool.TotalFunds).To(Equal(uint64(4010)))
		Expect(pool.Funders).To(HaveLen(2))
	})

	It("Remove funder", func() {
		// ARRANGE
		pool.AddAmountToFunder(i.ALICE, 1000)
		pool.AddAmountToFunder(i.ALICE, 2000)
		pool.AddAmountToFunder(i.ALICE, 0)
		pool.AddAmountToFunder(i.BOB, 0)
		pool.AddAmountToFunder(i.ALICE, 10)
		pool.AddAmountToFunder(i.CHARLIE, 500)

		Expect(pool.TotalFunds).To(Equal(uint64(3510)))

		// ACT
		// Alice: 3010
		// Charlie: 500
		pool.RemoveFunder(i.CHARLIE)

		// ASSERT
		Expect(pool.TotalFunds).To(Equal(uint64(3010)))
		Expect(pool.Funders).To(HaveLen(1))
	})

	It("Remove funder by defunding everything", func() {
		// ARRANGE
		pool.AddAmountToFunder(i.ALICE, 1000)
		pool.AddAmountToFunder(i.ALICE, 2000)
		pool.AddAmountToFunder(i.ALICE, 0)
		pool.AddAmountToFunder(i.BOB, 0)
		pool.AddAmountToFunder(i.ALICE, 10)
		pool.AddAmountToFunder(i.CHARLIE, 500)

		// ACT
		// Alice: 3010
		// Charlie: 500
		pool.SubtractAmountFromFunder(i.ALICE, 3010)

		// ASSERT
		Expect(pool.TotalFunds).To(Equal(uint64(500)))
		Expect(pool.Funders).To(HaveLen(1))
	})

	It("Charge Funders with equal amounts", func() {
		// ARRANGE
		for k := 0; k < 50; k++ {
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[k],
				Id:      0,
				Amount:  100 * i.KYVE,
			})
		}
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())
		Expect(pool.TotalFunds).To(Equal(50 * 100 * i.KYVE))

		// ACT
		payout, err := chargeFunders(s, 50*10*i.KYVE)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(50 * 10 * i.KYVE))

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.TotalFunds).To(Equal(50 * 90 * i.KYVE))

		for _, funder := range pool.Funders {
			Expect(funder.Amount).To(Equal(90 * i.KYVE))
		}
	})

	It("Charge Funders test remainder", func() {
		// ARRANGE
		for k := 0; k < 50; k++ {
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[k],
				Id:      0,
				Amount:  100 * i.KYVE,
			})
		}
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())
		Expect(pool.TotalFunds).To(Equal(50 * 100 * i.KYVE))

		// ACT
		// Charge 10 $KYVE + 49tkyve
		// the 49 tkyve will be charged to the lowest funder
		payout, err := chargeFunders(s, 50*10*i.KYVE+49)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(50*10*i.KYVE + 49))

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		for _, funder := range pool.Funders {
			if pool.GetLowestFunder().Address == funder.Address {
				Expect(funder.Amount).To(Equal(90*i.KYVE - 49))
			} else {
				Expect(funder.Amount).To(Equal(90 * i.KYVE))
			}
		}
	})

	It("Charge exactly lowest funder amount", func() {
		// ARRANGE
		for k := 0; k < 40; k++ {
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[k],
				Id:      0,
				Amount:  100 * i.KYVE,
			})
		}
		for k := 0; k < 10; k++ {
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[40+k],
				Id:      0,
				Amount:  200 * i.KYVE,
			})
		}
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())
		Expect(pool.TotalFunds).To(Equal((100*40 + 200*10) * i.KYVE))

		// ACT
		payout, err := chargeFunders(s, 50*100*i.KYVE)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(50 * 100 * i.KYVE))

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.Funders).To(HaveLen(10))
	})

	It("Kick out multiple lowest funders", func() {
		// Arrange
		for k := 0; k < 40; k++ {
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[k],
				Id:      0,
				Amount:  50 * i.KYVE,
			})
		}
		for k := 0; k < 10; k++ {
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[40+k],
				Id:      0,
				Amount:  1000 * i.KYVE,
			})
		}
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())
		Expect(pool.TotalFunds).To(Equal((50*40 + 1000*10) * i.KYVE))

		// 40 * 50 = 2000
		// 10 * 1000 = 10000
		// Charge 5000

		// Act
		payout, err := chargeFunders(s, 5000*i.KYVE)

		// Assert
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(3000 * i.KYVE))

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.Funders).To(HaveLen(10))

		for _, funder := range pool.Funders {
			Expect(funder.Amount).To(Equal(900 * i.KYVE))
		}
	})

	It("Charge more than pool has funds", func() {
		// ARRANGE
		for k := 0; k < 50; k++ {
			s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
				Creator: i.DUMMY[k],
				Id:      0,
				Amount:  50 * i.KYVE,
			})
		}
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())
		Expect(pool.TotalFunds).To(Equal((50 * 50) * i.KYVE))

		// ACT
		payout, err := chargeFunders(s, 5000*i.KYVE)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(Equal(2500 * i.KYVE))

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.Funders).To(HaveLen(0))
	})

	It("Charge pool which has no funds at all", func() {
		// ARRANGE
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())
		Expect(pool.TotalFunds).To(BeZero())

		// ACT
		payout, err := chargeFunders(s, 5000*i.KYVE)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
		Expect(payout).To(BeZero())

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.Funders).To(HaveLen(0))
		Expect(pool.TotalFunds).To(BeZero())
	})
})
