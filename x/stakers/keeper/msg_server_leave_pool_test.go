package keeper_test

import (
	"cosmossdk.io/math"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_leave_pool.go

* Leave a pool a staker has just joined as the first one
* Leave a pool multiple other stakers have joined previously
* Leave one of multiple pools a staker has previously joined
* Try to leave a pool again
* Leave a pool a staker has never joined

*/

var _ = Describe("msg_server_leave_pool.go", Ordered, func() {
	s := i.NewCleanChain()
	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create pool
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		// create staker
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		// join pool
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Leave a pool a staker has just joined as the first one", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, active := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(active).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0_A))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeTrue())
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))

		s.PerformValidityChecks()

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(BeEmpty())

		valaccount, active = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(active).To(BeFalse())

		valaccountsOfPool = s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(BeEmpty())

		totalStakeOfPool = s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(BeZero())

		// check if commission and stake fraction is still available
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))
	})

	It("Leave a pool multiple other stakers have joined previously", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, active := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(active).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0_A))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeTrue())
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(2))

		totalStakeOfPool := s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(200 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))

		s.PerformValidityChecks()

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(BeEmpty())

		valaccount, active = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(active).To(BeFalse())

		valaccountsOfPool = s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool = s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		// check if commission and stake fraction is still available
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))
	})

	It("Try to leave a pool again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})
		s.PerformValidityChecks()

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(HaveLen(1))

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Leave one of multiple pools a staker has previously joined", func() {
		// ARRANGE
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        1,
			Valaddress:    i.VALADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})
		s.PerformValidityChecks()

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  1,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(2))

		valaccount, active := s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)

		Expect(active).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(Equal(uint64(1)))
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_1_A))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeTrue())
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 1)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 1)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, active = s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)

		Expect(active).To(BeFalse())

		valaccountsOfPool = s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 1)

		Expect(valaccountsOfPool).To(BeEmpty())

		totalStakeOfPool = s.App().StakersKeeper.GetDelegationOfPool(s.Ctx(), 1)
		Expect(totalStakeOfPool).To(BeZero())

		// check if commission and stake fraction is still available
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))
	})

	It("Leave a pool a staker has never joined", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_1,
			PoolId:  0,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_1)
		Expect(valaccountsOfStaker).To(BeEmpty())

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_1)
		Expect(valaccountsOfStaker).To(BeEmpty())

		valaccount, active := s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)

		Expect(active).To(BeFalse())

		// check if commission and stake fraction is still available
		Expect(valaccount.Commission).To(Equal(math.LegacyZeroDec()))
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyZeroDec()))
	})
})
