package keeper_test

import (
	"cosmossdk.io/math"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_update_stake_fraction.go

* Get the default stake fraction from a newly joined pool
* Increase stake fraction to 50% from previous stake fraction
* Decrease stake fraction to 0% from previous stake fraction
* Decrease stake fraction to 1% from previous stake fraction
* Increase stake fraction to 100% from previous stake fraction
* Update stake fraction to same value from previous stake fraction
* Update stake fraction with a negative number from previous stake fraction
* Update stake fraction with a too high number from previous stake fraction
* Increase stake fraction after stake fraction has been decreased before during change time
* Decrease stake fraction after stake fraction has been decreased before during change time
* Decrease stake fraction after stake fraction has been increased before
* Update stake fraction with multiple pools
* Validator stake increases while stake fraction stays the same
* Validator stake decreases while stake fraction stays the same

*/

var _ = Describe("msg_server_update_stake_fraction.go", Ordered, func() {
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

		s.SetMaxVotingPower("1")

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.1"),
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Get the default stake fraction from a newly joined pool", func() {
		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))
	})

	It("Increase stake fraction to 50% from previous stake fraction", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.5"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(50 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(50 * i.KYVE))
	})

	It("Decrease stake fraction to 0% from previous stake fraction", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetStakeFractionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(BeZero())
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(BeZero())
	})

	It("Decrease stake fraction to 1% from previous stake fraction", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.01"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetStakeFractionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.01")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(1 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(1 * i.KYVE))
	})

	It("Increase stake fraction to 100% from previous stake fraction", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(100 * i.KYVE))
	})

	It("Update stake fraction to same value from previous stake fraction", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.1"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))
	})

	It("Update stake fraction with a negative number from previous stake fraction", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("-0.5"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))
	})

	It("Update stake fraction with a too high number from previous stake fraction during change time", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("2"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))
	})

	It("Increase stake fraction after stake fraction has been decreased before during change time", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.05"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.2"),
		})

		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.2")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(20 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(20 * i.KYVE))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.2")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(20 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(20 * i.KYVE))
	})

	It("Decrease stake fraction after stake fraction has been decreased before during change time", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.05"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.01"),
		})

		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.01")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(1 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(1 * i.KYVE))
	})

	It("Decrease stake fraction after stake fraction has been increased before", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.5"),
		})

		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(50 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(50 * i.KYVE))

		s.PerformValidityChecks()

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.02"),
		})

		// ASSERT
		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(50 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(50 * i.KYVE))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.02")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(2 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(2 * i.KYVE))
	})

	It("Update stake fraction with multiple pools", func() {
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
			Valaddress:    i.VALADDRESS_0_B,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.1"),
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        0,
			StakeFraction: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateStakeFraction{
			Creator:       i.STAKER_0,
			PoolId:        1,
			StakeFraction: math.LegacyMustNewDecFromStr("0.03"),
		})

		s.PerformValidityChecks()

		// ASSERT
		valaccount0, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount0.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(50 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(50 * i.KYVE))

		valaccount1, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)
		Expect(valaccount1.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 1)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 1)).To(Equal(10 * i.KYVE))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount0, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount0.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(50 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(50 * i.KYVE))

		valaccount1, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)
		Expect(valaccount1.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.03")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 1)).To(Equal(3 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 1)).To(Equal(3 * i.KYVE))
	})

	It("Validator stake increases while stake fraction stays the same", func() {
		// ARRANGE
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))

		// ACT
		s.SelfDelegateValidator(i.STAKER_0, 50*i.KYVE)

		// ASSERT
		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(15 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(15 * i.KYVE))
	})

	It("Validator stake decreases while stake fraction stays the same", func() {
		// ARRANGE
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(10 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(10 * i.KYVE))

		// ACT
		s.SelfUndelegateValidator(i.STAKER_0, 50*i.KYVE)

		// wait for update
		unbondingTime, _ := s.App().StakingKeeper.UnbondingTime(s.Ctx())
		s.CommitAfterSeconds(uint64(unbondingTime.Seconds()))
		s.CommitAfterSeconds(1)

		// ASSERT
		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(5 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(5 * i.KYVE))
	})
})
