package keeper_test

import (
	"cosmossdk.io/math"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - keeper_suite_effective_stake_test.go

* Test effective stake with all validators below the max pool voting power
* Test effective stake with one validator above the max pool voting power
* Test effective stake with multiple validators above the max pool voting power
* Test effective stake with fewer validators than required to undercut the max pool voting power
* Test effective stake with some validators having zero delegation
* Test effective stake with all validators having zero delegation
* Test effective stake with 0% as max pool stake
* Test effective stake with 100% as max pool stake
* Test effective stake with all validators below the max pool voting power due to stake fractions
* Test effective stake with one validator above the max pool voting power due to stake fractions
* Test effective stake with multiple validators above the max pool voting power due to stake fractions
* Test effective stake with some validators having zero delegation due to stake fractions
* Test effective stake with all validators having zero delegation due to stake fractions

*/

var _ = Describe("keeper_suite_effective_stake_test.go", Ordered, func() {
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

		s.SetMaxVotingPower("0.5")
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Test effective stake with all validators below the max pool voting power", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(300 * i.KYVE))
	})

	It("Test effective stake with one validator above the max pool voting power", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(250*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(200*i.KYVE - 1))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(400*i.KYVE - 1))
	})

	It("Test effective stake with multiple validators above the max pool voting power", func() {
		// ARRANGE
		s.SetMaxVotingPower("0.35")

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(600*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(500*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(120*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(140 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(140 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(120 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(400 * i.KYVE))
	})

	It("Test effective stake with fewer validators than required to undercut the max pool voting power", func() {
		// ARRANGE
		s.SetMaxVotingPower("0.2")

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeTrue())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(BeZero())
	})

	It("Test effective stake with some validators having zero delegation", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(200*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateZeroDelegationValidator(i.STAKER_1, "Staker-1")
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100*i.KYVE - 1))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(0 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200*i.KYVE - 1))
	})

	It("Test effective stake with all validators having zero delegation", func() {
		// ARRANGE
		staker0 := s.CreateNewValidator("Staker-0", 100*i.KYVE)
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       staker0.Address,
			PoolId:        0,
			PoolAddress:   staker0.PoolAccount[0],
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		staker1 := s.CreateNewValidator("Staker-1", 100*i.KYVE)
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       staker1.Address,
			PoolId:        0,
			PoolAddress:   staker1.PoolAccount[0],
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		staker2 := s.CreateNewValidator("Staker-2", 100*i.KYVE)
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       staker2.Address,
			PoolId:        0,
			PoolAddress:   staker2.PoolAccount[0],
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})
		s.SetDelegationToZero(staker0.Address)
		s.SetDelegationToZero(staker1.Address)
		s.SetDelegationToZero(staker2.Address)

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(BeZero())
	})

	It("Test effective stake with 0% as max pool stake", func() {
		// ARRANGE
		s.SetMaxVotingPower("0")

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeTrue())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(BeZero())
	})

	It("Test effective stake with 100% as max pool stake", func() {
		// ARRANGE
		s.SetMaxVotingPower("1")

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(100 * i.KYVE))
	})

	It("Test effective stake with all validators below the max pool voting power due to stake fractions", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(200*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(500*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.2"),
		})

		// ACT

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(300 * i.KYVE))
	})

	It("Test effective stake with one validator above the max pool voting power due to stake fractions", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(200*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.5"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(300*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.9"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(500*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.2"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(200*i.KYVE - 1))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(400*i.KYVE - 1))
	})

	It("Test effective stake with multiple validators above the max pool voting power due to stake fractions", func() {
		// ARRANGE
		s.SetMaxVotingPower("0.35")

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.8"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.7"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.2"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(uint64(23333333)))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(uint64(23333333)))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(20 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(uint64(66666666)))
	})

	It("Test effective stake with some validators having zero delegation due to stake fractions", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(200*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0.5"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(0 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_2, 0)).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))
	})

	It("Test effective stake with all validators having zero delegation due to stake fractions", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("0"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.IsVotingPowerTooHigh(s.Ctx(), 0)).To(BeFalse())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(BeZero())
	})
})
