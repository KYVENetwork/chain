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

TEST CASES - msg_server_update_commission.go

* Check if initial commission is correct
* Update commission to 50% from previous commission
* Update commission to 0% from previous commission
* Update commission to 100% from previous commission
* Update commission with a negative number from previous commission
* Update commission with a too high number from previous commission
* Update commission multiple times during the commission change time
* Update commission multiple times during the commission change time with the same value
* Update commission with multiple pools

*/

var _ = Describe("msg_server_update_commission.go", Ordered, func() {
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

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			Valaddress:    i.VALADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Get the default commission from a newly joined pool", func() {
		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
	})

	It("Update commission to 50% from previous commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
	})

	It("Update commission to 0% from previous commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyZeroDec(),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyZeroDec()))
	})

	It("Update commission to 100% from previous commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyOneDec(),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyOneDec()))
	})

	It("Update commission with a negative number from previous commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("-0.5"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
	})

	It("Update commission with a too high number from previous commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyNewDec(2),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
	})

	It("Update commission multiple times during the commission change time", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})
		s.PerformValidityChecks()

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.2"),
		})
		s.PerformValidityChecks()

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.3"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.3")))
	})

	It("Update commission multiple times during the commission change time with the same value", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.2"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.1"),
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccount, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))
	})

	It("Update commission with multiple pools", func() {
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
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.PerformValidityChecks()

		// ASSERT
		valaccount0, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount0.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))

		valaccount1, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)
		Expect(valaccount1.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.1")))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccount0, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccount0.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.5")))

		valaccount1, _ = s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)
		Expect(valaccount1.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
	})
})
