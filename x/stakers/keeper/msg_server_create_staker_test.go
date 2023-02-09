package keeper_test

import (
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/stakers/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_create_staker.go

* Create a first new staker and delegate 100 $KYVE
* Do an additional 50 $KYVE self delegation after staker has already delegated 100 $KYVE
* Try to create staker with more $KYVE than available in balance
* Create a second staker by staking 150 $KYVE
* Try to create a staker again

*/

var _ = Describe("msg_server_create_staker.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalance := s.GetBalanceFromAddress(i.STAKER_0)

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Create a first new staker and delegate 100 $KYVE", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.STAKER_0)

		staker, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		valaccounts := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		Expect(staker.Address).To(Equal(i.STAKER_0))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))
		Expect(staker.Commission).To(Equal(types.DefaultCommission))

		Expect(staker.Moniker).To(BeEmpty())
		Expect(staker.Logo).To(BeEmpty())
		Expect(staker.Website).To(BeEmpty())

		Expect(valaccounts).To(BeEmpty())
	})

	It("Do an additional 50 $KYVE self delegation after staker has already delegated 100 $KYVE", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.STAKER_0,
			Staker:  i.STAKER_0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.STAKER_0)

		staker, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		valaccounts := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(initialBalance - balanceAfter).To(Equal(150 * i.KYVE))

		Expect(staker.Address).To(Equal(i.STAKER_0))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(150 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(150 * i.KYVE))

		Expect(staker.Commission).To(Equal(types.DefaultCommission))

		Expect(staker.Moniker).To(BeEmpty())
		Expect(staker.Logo).To(BeEmpty())
		Expect(staker.Website).To(BeEmpty())

		Expect(valaccounts).To(HaveLen(0))
	})

	It("Try to create staker with more $KYVE than available in balance", func() {
		// ACT
		currentBalance := s.GetBalanceFromAddress(i.STAKER_0)

		s.RunTxStakersError(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  currentBalance + 1,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.STAKER_0)
		Expect(initialBalance - balanceAfter).To(BeZero())

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeFalse())
	})

	It("Create a second staker by staking 150 $KYVE", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  150 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.BOB)

		staker, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.BOB)
		valaccounts := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.BOB)

		Expect(found).To(BeTrue())

		Expect(initialBalance - balanceAfter).To(Equal(150 * i.KYVE))

		Expect(staker.Address).To(Equal(i.BOB))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)).To(Equal(150 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.BOB, i.BOB)).To(Equal(150 * i.KYVE))

		Expect(staker.Commission).To(Equal(types.DefaultCommission))

		Expect(staker.Moniker).To(BeEmpty())
		Expect(staker.Logo).To(BeEmpty())
		Expect(staker.Website).To(BeEmpty())

		Expect(valaccounts).To(BeEmpty())
	})

	It("Try to create a staker again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.STAKER_0)

		staker, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		valaccounts := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		Expect(staker.Address).To(Equal(i.STAKER_0))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))
		Expect(staker.Commission).To(Equal(types.DefaultCommission))

		Expect(staker.Moniker).To(BeEmpty())
		Expect(staker.Logo).To(BeEmpty())
		Expect(staker.Website).To(BeEmpty())

		Expect(valaccounts).To(BeEmpty())
	})
})
