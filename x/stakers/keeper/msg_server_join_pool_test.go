package keeper_test

import (
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_join_pool.go

* Test if a newly created staker is participating in no pools yet
* Join the first pool as the first staker to a newly created pool
* Join a pool with zero delegation
* Join disabled pool
* Join a pool where other stakers have already joined
* Self-Delegate more KYVE after joining a pool
* Join a pool with the same valaddress as the staker address
* Try to join the same pool with the same valaddress again
* Try to join the same pool with a different valaddress
* Try to join another pool with the same valaddress again
* Try to join another pool with a valaddress that is already used by another staker
* Try to join another pool with a different valaddress
* Join a pool with a valaddress which does not exist on chain yet
* Join a pool with a valaddress which does not exist on chain yet and send 0 funds
* Join a pool with an invalid valaddress
* Join a pool and fund the valaddress with more KYVE than available in balance
* Kick out lowest staker by joining a full pool
* Fail to kick out lowest staker because not enough stake
* Kick out lowest staker with respect to stake + delegation
* Fail to kick out lowest staker because not enough stake + delegation

*/

var _ = Describe("msg_server_join_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceStaker0 := uint64(0)
	initialBalanceValaddress0 := uint64(0)

	initialBalanceStaker1 := uint64(0)
	initialBalanceValaddress1 := uint64(0)

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create pool
		msg := &pooltypes.MsgCreatePool{
			Authority:      gov,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Binaries:       "{}",
		}
		s.RunTxPoolSuccess(msg)

		// create staker
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0_A)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceValaddress1 = s.GetBalanceFromAddress(i.VALADDRESS_1_A)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Test if a newly created staker is participating in no pools yet", func() {
		// ASSERT
		valaccounts := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccounts).To(HaveLen(0))
	})

	It("Join the first pool as the first staker to a newly created pool", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterValaddress0 := s.GetBalanceFromAddress(i.VALADDRESS_0_A)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(100 * i.KYVE))
		Expect(balanceAfterValaddress0 - initialBalanceValaddress0).To(Equal(100 * i.KYVE))

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0_A))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join a pool with zero delegation", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  0 * i.KYVE,
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     0 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker1 := s.GetBalanceFromAddress(i.STAKER_1)
		balanceAfterValaddress1 := s.GetBalanceFromAddress(i.VALADDRESS_1_A)

		Expect(initialBalanceStaker1).To(Equal(balanceAfterStaker1))
		Expect(initialBalanceValaddress1).To(Equal(balanceAfterValaddress1))

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_1)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_1))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_1_A))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(BeZero())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_1)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(Equal(totalStakeOfPool))
	})

	It("Join disabled pool", func() {
		// ARRANGE
		msg := &pooltypes.MsgCreatePool{
			Authority:      gov,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Binaries:       "{}",
		}
		s.RunTxPoolSuccess(msg)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		pool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		_, err := s.RunTx(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		Expect(err.Error()).To(Equal("can not join disabled pool: internal logic error"))

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterValaddress0 := s.GetBalanceFromAddress(i.VALADDRESS_0_A)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(0 * i.KYVE))
		Expect(balanceAfterValaddress0 - initialBalanceValaddress0).To(Equal(0 * i.KYVE))

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(0))

		_, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)

		Expect(found).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 1)

		Expect(valaccountsOfPool).To(HaveLen(0))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 1)

		Expect(totalStakeOfPool).To(Equal(0 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetActiveStakers(s.Ctx())).To(HaveLen(0))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("join a pool where other stakers have already joined", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     0 * i.KYVE,
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     0 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterValaddress0 := s.GetBalanceFromAddress(i.VALADDRESS_0_A)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(BeZero())
		Expect(balanceAfterValaddress0 - initialBalanceValaddress0).To(BeZero())

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0_A))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(2))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(200 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Self-Delegate more KYVE after joining a pool", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		// ACT
		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.STAKER_0,
			Staker:  i.STAKER_0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0_A))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool = s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(150 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Try to join the same pool with the same valaddress again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))
	})

	It("join a pool with the same valaddress as the staker address", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.STAKER_0,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Try to join the same pool with a different valaddress", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join another pool with the same valaddress again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		msg := &pooltypes.MsgCreatePool{
			Authority:      gov,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Binaries:       "{}",
		}
		s.RunTxPoolSuccess(msg)

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join pool with a valaddress that is already used by another staker", func() {
		// ARRANGE
		msg := &pooltypes.MsgCreatePool{
			Authority:      gov,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Binaries:       "{}",
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     1,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join pool with a valaddress that is already used by another staker", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_1)
		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Try to join another pool with a different valaddress", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     100 * i.KYVE,
		})

		msg := &pooltypes.MsgCreatePool{
			Authority:      gov,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Binaries:       "{}",
		}
		s.RunTxPoolSuccess(msg)

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_1_A,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(HaveLen(2))
	})

	It("Join a pool with a valaddress which does not exist on chain yet", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x",
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterUnknown := s.GetBalanceFromAddress("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x")

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(100 * i.KYVE))
		Expect(balanceAfterUnknown).To(Equal(100 * i.KYVE))

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x"))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join a pool with a valaddress which does not exist on chain yet and send 0 funds", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x",
			Amount:     0 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterUnknown := s.GetBalanceFromAddress("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x")

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(BeZero())
		Expect(balanceAfterUnknown).To(BeZero())

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x"))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join a pool with an invalid valaddress", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "invalid_valaddress",
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Join a pool and fund the valaddress with more KYVE than available in balance", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "invalid_valaddress",
			Amount:     initialBalanceStaker0 + 1,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.ALICE)

		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Kick out lowest staker by joining a full pool", func() {
		// Arrange
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     1,
		})

		for k := 0; k < 49; k++ {
			s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
				Creator: i.DUMMY[k],
				Amount:  150 * i.KYVE,
			})
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:    i.DUMMY[k],
				PoolId:     0,
				Valaddress: i.VALDUMMY[k],
				Amount:     1,
			})
		}

		// STAKER_0 is lowest staker and all stakers are full now.
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  150 * i.KYVE,
		})

		// Act
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     1,
		})

		// Assert
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 150) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).ToNot(ContainElement(i.STAKER_0))
	})

	It("Fail to kick out lowest staker because not enough stake", func() {
		// Arrange
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     1,
		})

		for k := 0; k < 49; k++ {
			s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
				Creator: i.DUMMY[k],
				Amount:  150 * i.KYVE,
			})
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:    i.DUMMY[k],
				PoolId:     0,
				Valaddress: i.VALDUMMY[k],
				Amount:     1,
			})
		}

		// STAKER_0 is lowest staker and all stakers are full now.
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  50 * i.KYVE,
		})

		// Act
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     1,
		})

		// Assert
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).ToNot(ContainElement(i.STAKER_1))
	})

	It("Kick out lowest staker with respect to stake + delegation", func() {
		// ARRANGE
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     1 * i.KYVE,
		})

		for k := 0; k < 49; k++ {
			s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
				Creator: i.DUMMY[k],
				Amount:  150 * i.KYVE,
			})
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:    i.DUMMY[k],
				PoolId:     0,
				Valaddress: i.VALDUMMY[k],
				Amount:     1 * i.KYVE,
			})
		}

		// Alice is lowest staker and all stakers are full now.
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  150 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.ALICE,
			Staker:  i.STAKER_0,
			Amount:  150 * i.KYVE,
		}) // Staker0 has now 250 delegation

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     1,
		})

		// ASSERT
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 250) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).NotTo(ContainElement(i.STAKER_1))
	})

	It("Fail to kick out lowest staker because not enough stake", func() {
		// Arrange
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     1,
		})

		for k := 0; k < 49; k++ {
			s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
				Creator: i.DUMMY[k],
				Amount:  150 * i.KYVE,
			})
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:    i.DUMMY[k],
				PoolId:     0,
				Valaddress: i.VALDUMMY[k],
				Amount:     1,
			})
		}

		// STAKER_0 is lowest staker and all stakers are full now.
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  50 * i.KYVE,
		})

		// Act
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     1,
		})

		// Assert
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).ToNot(ContainElement(i.STAKER_1))
	})

	It("Fail to kick out lowest staker because not enough stake + delegation", func() {
		// ARRANGE
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     1 * i.KYVE,
		})

		for k := 0; k < 49; k++ {
			s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
				Creator: i.DUMMY[k],
				Amount:  150 * i.KYVE,
			})
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:    i.DUMMY[k],
				PoolId:     0,
				Valaddress: i.VALDUMMY[k],
				Amount:     1 * i.KYVE,
			})
		}

		// Alice is lowest staker and all stakers are full now.
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  50 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.ALICE,
			Staker:  i.STAKER_1,
			Amount:  50 * i.KYVE,
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     0,
		})

		// ASSERT
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).NotTo(ContainElement(i.STAKER_1))
	})
})
