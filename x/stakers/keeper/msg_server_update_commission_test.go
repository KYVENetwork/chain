package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_update_commission.go

* Get the default commission from a newly created staker
* Update commission to 50% from previously default commission
* Update commission to 0% from previously default commission
* Update commission to 100% from previously default commission
* Update commission with an invalid number from previously default commission
* Update commission with a negative number from previously default commission
* Update commission with a too high number from previously default commission
* Update commission multiple times during the commission change time
* Update commission multiple times during the commission change time with the same value
* Update commission with multiple stakers

*/

var _ = Describe("msg_server_update_commission.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create staker
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Get the default commission from a newly created staker", func() {
		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission to 50% from previously default commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
	})

	It("Update commission to 0% from previously default commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyZeroDec(),
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(math.LegacyZeroDec()))
	})

	It("Update commission to 100% from previously default commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyOneDec(),
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(sdk.OneDec()))
	})

	It("Update commission with a negative number from previously default commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("-0.5"),
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission with a too high number from previously default commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyNewDec(2),
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission multiple times during the commission change time", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})
		s.PerformValidityChecks()

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("0.2"),
		})
		s.PerformValidityChecks()

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("0.3"),
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.3")))
	})

	It("Update commission multiple times during the commission change time with the same value", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("0.2"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: stakerstypes.DefaultCommission,
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission with multiple stakers", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_1,
			Commission: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.PerformValidityChecks()

		// ASSERT
		staker0, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker0.Commission).To(Equal(stakerstypes.DefaultCommission))

		staker1, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_1)
		Expect(staker1.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.GetCommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker0, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker0.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.5")))

		staker1, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_1)
		Expect(staker1.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.5")))
	})
})
