package keeper_test

import (
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/delegation/types"
)

/*

TEST CASES - msg_server_delegate.go

* Delegate 10 $KYVE to ALICE
* Delegate additional 50 $KYVE to ALICE
* Try delegating to non-existent staker
* Delegate more than available
* Payout delegators
* Don't pay out rewards twice
* Delegate to validator with 0 $KYVE
* Delegate to multiple validators

*/

var _ = Describe("msg_server_delegate.go", Ordered, func() {
	s := i.NewCleanChain()

	const aliceSelfDelegation = 100 * i.KYVE
	const bobSelfDelegation = 200 * i.KYVE

	BeforeEach(func() {
		s = i.NewCleanChain()

		CreatePool(s)

		s.CreateValidator(i.ALICE, "Alice", int64(aliceSelfDelegation))
		// Shadow delegation inside delegation module (temporary work-around)
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.ALICE,
			Staker:  i.ALICE,
			Amount:  aliceSelfDelegation,
		})

		s.CreateValidator(i.BOB, "Bob", int64(bobSelfDelegation))
		// Shadow delegation inside delegation module (temporary work-around)
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.BOB,
			Amount:  bobSelfDelegation,
		})

		_, stakerFound := s.App().StakersKeeper.GetValidator(s.Ctx(), i.ALICE)
		Expect(stakerFound).To(BeTrue())

		_, stakerFound = s.App().StakersKeeper.GetValidator(s.Ctx(), i.BOB)
		Expect(stakerFound).To(BeTrue())

		s.CommitAfterSeconds(7)
	})

	AfterEach(func() {
		CheckAndContinueChainForOneMonth(s)
	})

	It("Delegate 10 $KYVE to ALICE", func() {
		// ARRANGE
		bobBalance := s.GetBalanceFromAddress(i.BOB)

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		// ASSERT
		CheckAndContinueChainForOneMonth(s)
		bobBalanceAfter := s.GetBalanceFromAddress(i.BOB)
		Expect(bobBalanceAfter).To(Equal(bobBalance - 10*i.KYVE))

		aliceDelegation := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegation).To(Equal(10*i.KYVE + aliceSelfDelegation))
	})

	It("Delegate 10 $KYVE to ALICE and then another 50 $KYVE", func() {
		// ARRANGE
		bobBalance := s.GetBalanceFromAddress(i.BOB)
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})
		CheckAndContinueChainForOneMonth(s)

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.ALICE,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		CheckAndContinueChainForOneMonth(s)
		bobBalanceAfter := s.GetBalanceFromAddress(i.BOB)
		Expect(bobBalanceAfter).To(Equal(bobBalance - 60*i.KYVE))

		aliceDelegation := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegation).To(Equal(60*i.KYVE + aliceSelfDelegation))
	})

	It("Try delegating to non-existent staker", func() {
		// ARRANGE
		bobBalance := s.GetBalanceFromAddress(i.BOB)
		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorError(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.CHARLIE,
			Amount:  10 * i.KYVE,
		})

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.BOB)).To(Equal(bobBalance))

		aliceDelegation := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegation).To(Equal(aliceSelfDelegation))
	})

	It("Delegate more than available", func() {
		// ARRANGE
		bobBalance := s.GetBalanceFromAddress(i.BOB)
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorError(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.ALICE,
			Amount:  bobBalance + 1,
		})

		// ASSERT
		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationBefore).To(Equal(aliceDelegationAfter))

		bobBalanceAfter := s.GetBalanceFromAddress(i.BOB)
		Expect(bobBalanceAfter).To(Equal(bobBalance))
	})

	It("Payout delegators", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  100 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  209 * i.KYVE,
		})

		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeEmpty())
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(BeEmpty())

		s.PerformValidityChecks()

		// ACT
		PayoutRewards(s, i.ALICE, sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(10*i.KYVE))))

		// ASSERT

		// Name    amount   shares
		// Alice:   100		100/(409) * 10 * 1e9 = 2.444.987.775
		// Dummy0:  100		100/(409) * 10 * 1e9 = 2.444.987.775
		// Dummy1:  209		209/(409) * 10 * 1e9 = 5.110.024.449
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.ALICE).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(2_444_987_775)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(2_444_987_775)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(5_110_024_449)))

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})

		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.ALICE).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(2_444_987_775)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(0)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(5_110_024_449)))

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(900*i.KYVE + 2_444_987_775))
		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal((200+409)*i.KYVE + uint64(2_444_987_775+5_110_024_449+1)))
	})

	It("Don't pay out rewards twice", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  100 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  200 * i.KYVE,
		})

		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeEmpty())
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(BeEmpty())

		// ACT
		PayoutRewards(s, i.ALICE, sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(10*i.KYVE))))

		// ASSERT

		// Alice: 100
		// Dummy0: 100
		// Dummy1: 200
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(2_500_000_000)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(5_000_000_000)))

		s.PerformValidityChecks()

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})

		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeEmpty())
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(5_000_000_000)))

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(900*i.KYVE + 2_500_000_000))
		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal(600*i.KYVE + 7_500_000_000))
	})
})
