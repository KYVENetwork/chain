package keeper_test

import (
	"fmt"
	"sort"

	"cosmossdk.io/math"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_bundles.go

* Correctly load round-robin state
* Correctly save and load round-robin state
* Empty round-robin set
* Partially filled round-robin set (one staker with 0 delegation)
* Frequency analysis
* Frequency analysis with stake fractions
* Frequency analysis with maximum voting power cap
* Frequency analysis (rounding)
* Frequency analysis (excluded)
* Exclude everybody
* Exclude all but one
* Leave set
* Join set

*/

func joinDummy(s *i.KeeperTestSuite, index, kyveAmount uint64) {
	joinDummyWithStakeFraction(s, index, kyveAmount, math.LegacyOneDec())
}

func joinDummyWithStakeFraction(s *i.KeeperTestSuite, index, kyveAmount uint64, stakeFraction math.LegacyDec) {
	s.CreateValidator(i.DUMMY[index], fmt.Sprintf("dummy-%d", index), int64(kyveAmount*i.KYVE))

	s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
		Creator:       i.DUMMY[index],
		PoolId:        0,
		PoolAddress:   i.VALDUMMY[index],
		Amount:        0,
		Commission:    math.LegacyMustNewDecFromStr("0.1"),
		StakeFraction: stakeFraction,
	})
}

func leaveDummy(s *i.KeeperTestSuite, index uint64) {
	s.RunTxStakersSuccess(&stakertypes.MsgLeavePool{
		Creator: i.DUMMY[index],
		PoolId:  0,
	})
	s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
	s.CommitAfterSeconds(1)
}

var _ = Describe("logic_round_robin.go", Ordered, func() {
	var s *i.KeeperTestSuite

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// sort dummy accounts alphabetically
		sort.Slice(i.DUMMY, func(k, j int) bool {
			return i.DUMMY[k] < i.DUMMY[j]
		})

		// create clean pool for every test case
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: math.LegacyNewDec(int64(2 * i.KYVE)),
			MinDelegation:        1_000_000 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Protocol.LastUpgrade = uint64(s.Ctx().BlockTime().Unix())
		pool.UpgradePlan = &pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix()),
			Duration:    60,
		}
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		s.SetMaxVotingPower("1")
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Correctly load round-robin state", func() {
		// ARRANGE
		joinDummy(s, 0, 100)
		joinDummy(s, 1, 200)
		joinDummy(s, 2, 300)

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)

		// ASSERT
		Expect(rrvs.Validators).To(HaveLen(3))
		Expect(rrvs.Validators[0].Address).To(Equal(i.DUMMY[0]))
		Expect(rrvs.Validators[0].Power).To(Equal(100 * int64(i.KYVE)))
		Expect(rrvs.Validators[1].Address).To(Equal(i.DUMMY[1]))
		Expect(rrvs.Validators[1].Power).To(Equal(200 * int64(i.KYVE)))
		Expect(rrvs.Validators[2].Address).To(Equal(i.DUMMY[2]))
		Expect(rrvs.Validators[2].Power).To(Equal(300 * int64(i.KYVE)))

		Expect(rrvs.Progress).To(HaveLen(3))
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(int64(0)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(int64(0)))
		Expect(rrvs.Progress[i.DUMMY[2]]).To(Equal(int64(0)))
	})

	It("Correctly save and load round-robin state", func() {
		// ARRANGE
		joinDummy(s, 0, 100)
		joinDummy(s, 1, 200)
		joinDummy(s, 2, 300)

		// ACT
		rrvs1 := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		rrvs1.Progress[i.DUMMY[0]] = 1
		rrvs1.Progress[i.DUMMY[1]] = 2
		rrvs1.Progress[i.DUMMY[2]] = 3
		s.App().BundlesKeeper.SaveRoundRobinValidatorSet(s.Ctx(), rrvs1)

		state := rrvs1.GetRoundRobinProgress()
		// loading round-robin performs normalising, which shifts values from 1,2,3 to -1,0,1
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)

		// ASSERT
		Expect(state[0].Address).To(Equal(i.DUMMY[0]))
		Expect(state[0].Progress).To(Equal(int64(1)))
		Expect(state[1].Address).To(Equal(i.DUMMY[1]))
		Expect(state[1].Progress).To(Equal(int64(2)))
		Expect(state[2].Address).To(Equal(i.DUMMY[2]))
		Expect(state[2].Progress).To(Equal(int64(3)))

		Expect(rrvs.Validators).To(HaveLen(3))
		Expect(rrvs.Validators[0].Address).To(Equal(i.DUMMY[0]))
		Expect(rrvs.Validators[1].Address).To(Equal(i.DUMMY[1]))
		Expect(rrvs.Validators[2].Address).To(Equal(i.DUMMY[2]))

		Expect(rrvs.Progress).To(HaveLen(3))
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(int64(-1)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(int64(0)))
		Expect(rrvs.Progress[i.DUMMY[2]]).To(Equal(int64(1)))
	})

	It("Partially filled round-robin set (one staker with 0 delegation)", func() {
		// ARRANGE
		joinDummy(s, 1, 10)
		joinDummy(s, 2, 5)

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		state := rrvs.GetRoundRobinProgress()

		nextProposer := rrvs.NextProposer()

		// ASSERT
		Expect(rrvs.Validators).To(HaveLen(2))
		Expect(rrvs.Progress).To(HaveLen(2))

		Expect(state).To(HaveLen(2))

		Expect(nextProposer).To(Equal(i.DUMMY[1]))
	})

	It("Frequency analysis", func() {
		// ARRANGE
		joinDummy(s, 0, 2)
		joinDummy(s, 1, 31)
		joinDummy(s, 2, 67)

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)

		frequency1 := make(map[string]int, 0)
		for i := 0; i < 10; i++ {
			frequency1[rrvs.NextProposer()] += 1
		}

		frequency2 := make(map[string]int, 0)
		for i := 0; i < 100; i++ {
			frequency2[rrvs.NextProposer()] += 1
		}

		frequency3 := make(map[string]int, 0)
		for i := 0; i < 100000; i++ {
			frequency3[rrvs.NextProposer()] += 1
		}

		// ASSERT

		Expect(frequency1[i.DUMMY[0]]).To(Equal(0))
		Expect(frequency1[i.DUMMY[1]]).To(Equal(3))
		Expect(frequency1[i.DUMMY[2]]).To(Equal(7))

		Expect(frequency2[i.DUMMY[0]]).To(Equal(2))
		Expect(frequency2[i.DUMMY[1]]).To(Equal(31))
		Expect(frequency2[i.DUMMY[2]]).To(Equal(67))

		Expect(frequency3[i.DUMMY[0]]).To(Equal(2000))
		Expect(frequency3[i.DUMMY[1]]).To(Equal(31000))
		Expect(frequency3[i.DUMMY[2]]).To(Equal(67000))
	})

	It("Frequency analysis with stake fractions", func() {
		// ARRANGE
		joinDummyWithStakeFraction(s, 0, 100, math.LegacyMustNewDecFromStr("0"))
		joinDummyWithStakeFraction(s, 1, 100, math.LegacyMustNewDecFromStr("0.5"))
		joinDummyWithStakeFraction(s, 2, 100, math.LegacyMustNewDecFromStr("1"))

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)

		frequency1 := make(map[string]int, 0)
		for i := 0; i < 10; i++ {
			frequency1[rrvs.NextProposer()] += 1
		}

		frequency2 := make(map[string]int, 0)
		for i := 0; i < 100; i++ {
			frequency2[rrvs.NextProposer()] += 1
		}

		frequency3 := make(map[string]int, 0)
		for i := 0; i < 100000; i++ {
			frequency3[rrvs.NextProposer()] += 1
		}

		// ASSERT
		Expect(frequency1[i.DUMMY[0]]).To(Equal(0))
		Expect(frequency1[i.DUMMY[1]]).To(Equal(3))
		Expect(frequency1[i.DUMMY[2]]).To(Equal(7))

		Expect(frequency2[i.DUMMY[0]]).To(Equal(0))
		Expect(frequency2[i.DUMMY[1]]).To(Equal(34))
		Expect(frequency2[i.DUMMY[2]]).To(Equal(66))

		Expect(frequency3[i.DUMMY[0]]).To(Equal(0))
		Expect(frequency3[i.DUMMY[1]]).To(Equal(33333))
		Expect(frequency3[i.DUMMY[2]]).To(Equal(66667))
	})

	It("Frequency analysis with maximum voting power cap", func() {
		// ARRANGE
		s.SetMaxVotingPower("0.5")

		// NOTE that dummy with index 2 has more than 50% voting power, so his effective stake
		// will be lower
		joinDummy(s, 0, 2)
		joinDummy(s, 1, 31)
		joinDummy(s, 2, 67)

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)

		frequency1 := make(map[string]int, 0)
		for i := 0; i < 10; i++ {
			frequency1[rrvs.NextProposer()] += 1
		}

		frequency2 := make(map[string]int, 0)
		for i := 0; i < 100; i++ {
			frequency2[rrvs.NextProposer()] += 1
		}

		frequency3 := make(map[string]int, 0)
		for i := 0; i < 100000; i++ {
			frequency3[rrvs.NextProposer()] += 1
		}

		// ASSERT
		Expect(frequency1[i.DUMMY[0]]).To(Equal(0))
		Expect(frequency1[i.DUMMY[1]]).To(Equal(5))
		Expect(frequency1[i.DUMMY[2]]).To(Equal(5))

		Expect(frequency2[i.DUMMY[0]]).To(Equal(3))
		Expect(frequency2[i.DUMMY[1]]).To(Equal(47))
		Expect(frequency2[i.DUMMY[2]]).To(Equal(50))

		Expect(frequency3[i.DUMMY[0]]).To(Equal(3031))  // 2/66
		Expect(frequency3[i.DUMMY[1]]).To(Equal(46969)) // 31/66
		Expect(frequency3[i.DUMMY[2]]).To(Equal(50000)) // 33/66
	})

	It("Frequency analysis (rounding)", func() {
		// ARRANGE
		joinDummy(s, 0, 1)
		joinDummy(s, 1, 1)
		joinDummy(s, 2, 1)

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)

		frequency1 := make(map[string]int, 0)
		for i := 0; i < 10; i++ {
			frequency1[rrvs.NextProposer()] += 1
		}

		frequency2 := make(map[string]int, 0)
		for i := 0; i < 100; i++ {
			frequency2[rrvs.NextProposer()] += 1
		}

		frequency3 := make(map[string]int, 0)
		for i := 0; i < 100000; i++ {
			frequency3[rrvs.NextProposer()] += 1
		}

		// ASSERT

		// First one is selected one more time, because of alphabetical order
		Expect(frequency1[i.DUMMY[0]]).To(Equal(4))
		Expect(frequency1[i.DUMMY[1]]).To(Equal(3))
		Expect(frequency1[i.DUMMY[2]]).To(Equal(3))

		Expect(frequency2[i.DUMMY[0]]).To(Equal(33))
		// The state is not reset between rounds, the first one already got selected one more time in
		// the previous round, hence its progress is already reset. Therefore, the second one is slected
		// one more than the others.
		Expect(frequency2[i.DUMMY[1]]).To(Equal(34))
		Expect(frequency2[i.DUMMY[2]]).To(Equal(33))

		Expect(frequency3[i.DUMMY[0]]).To(Equal(33333))
		Expect(frequency3[i.DUMMY[1]]).To(Equal(33333))
		// Same argument as above
		Expect(frequency3[i.DUMMY[2]]).To(Equal(33334))
	})

	It("Frequency analysis (excluded)", func() {
		// ARRANGE
		joinDummy(s, 0, 5)
		joinDummy(s, 1, 10)
		joinDummy(s, 2, 15)
		// total stake = 30
		// Do 1000 rounds, in the first 500 exclude Dummy0, in the second 500 exclude Dummy1
		// Frequencies for all three validators:
		// P(0) = 1/1000 * (0 + 500 * 5/(5+15)) = 0.125
		// P(1) = 1/1000 * (500 * 10/(10+15) + 0) = 0.2
		// P(2) = 1/1000 * (500 * 15/(10+15) + 500 * 15/(5+15)) = 0.675
		// P(1) + P(2) + P(3) = 1

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)

		frequency := make(map[string]int, 0)
		for j := 0; j < 500; j++ {
			frequency[rrvs.NextProposer(i.DUMMY[0])] += 1
		}
		for j := 0; j < 500; j++ {
			frequency[rrvs.NextProposer(i.DUMMY[1])] += 1
		}

		// ASSERT
		Expect(frequency[i.DUMMY[0]]).To(Equal(125))
		Expect(frequency[i.DUMMY[1]]).To(Equal(200))
		Expect(frequency[i.DUMMY[2]]).To(Equal(675))
	})

	It("Exclude everybody", func() {
		// ARRANGE
		joinDummy(s, 0, 5)
		joinDummy(s, 1, 10)
		joinDummy(s, 2, 15)

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		nextProposer := rrvs.NextProposer(i.DUMMY[0], i.DUMMY[1], i.DUMMY[2])

		// ASSERT
		Expect(nextProposer).To(Equal(i.DUMMY[2]))
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(5 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(10 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[2]]).To(Equal(-15 * int64(i.KYVE)))
	})

	It("Exclude all but one", func() {
		// ARRANGE
		joinDummy(s, 0, 5)
		joinDummy(s, 1, 10)
		joinDummy(s, 2, 15)

		// ACT
		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		nextProposer := rrvs.NextProposer(i.DUMMY[1], i.DUMMY[2])

		// ASSERT
		Expect(nextProposer).To(Equal(i.DUMMY[0]))
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(0 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(0 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[2]]).To(Equal(0 * int64(i.KYVE)))
	})

	It("Leave set", func() {
		// ARRANGE
		joinDummy(s, 0, 1000)
		joinDummy(s, 1, 1000)
		joinDummy(s, 2, 1000)
		joinDummy(s, 3, 1000)
		joinDummy(s, 4, 1000)

		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		nextProposer := rrvs.NextProposer()
		s.App().BundlesKeeper.SaveRoundRobinValidatorSet(s.Ctx(), rrvs)
		Expect(nextProposer).To(Equal(i.DUMMY[0]))
		leaveDummy(s, 2)
		leaveDummy(s, 3)
		leaveDummy(s, 4)

		// ACT, ASSERT
		rrvs = s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(-2000 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(2000 * int64(i.KYVE)))

		nextProposer = rrvs.NextProposer()
		Expect(nextProposer).To(Equal(i.DUMMY[1]))
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(-1000 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(1000 * int64(i.KYVE)))

		nextProposer = rrvs.NextProposer()
		Expect(nextProposer).To(Equal(i.DUMMY[1]))
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(0 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(0 * int64(i.KYVE)))
	})

	It("Join set", func() {
		// ARRANGE
		joinDummy(s, 0, 100)
		joinDummy(s, 1, 200)
		joinDummy(s, 2, 300)

		rrvs := s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		nextProposer := rrvs.NextProposer()
		Expect(nextProposer).To(Equal(i.DUMMY[2]))
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(100 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(200 * int64(i.KYVE)))
		Expect(rrvs.Progress[i.DUMMY[2]]).To(Equal(-300 * int64(i.KYVE)))
		s.App().BundlesKeeper.SaveRoundRobinValidatorSet(s.Ctx(), rrvs)

		// ACT
		joinDummy(s, 3, 400)
		rrvs = s.App().BundlesKeeper.LoadRoundRobinValidatorSet(s.Ctx(), 0)
		shift := int64(1125_000_000 / 4)

		// Assert
		Expect(rrvs.Progress[i.DUMMY[0]]).To(Equal(100*int64(i.KYVE) + shift))
		Expect(rrvs.Progress[i.DUMMY[1]]).To(Equal(200*int64(i.KYVE) + shift))
		Expect(rrvs.Progress[i.DUMMY[2]]).To(Equal(-300*int64(i.KYVE) + shift))
		Expect(rrvs.Progress[i.DUMMY[3]]).To(Equal(-1125*int64(i.KYVE) + shift))
	})
})
