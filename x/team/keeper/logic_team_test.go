package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	"github.com/KYVENetwork/chain/x/team/types"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_team.go

* leave_minus_join_lt_1y
* leave_minus_join_lt_3y_and_tge_lt_join
* leave_minus_join_lt_3y_and_tge_gt_join
* leave_minus_join_gt_3y_and_tge_gt_join
* leave_minus_join_gt_3y_and_tge_lt_join
* no_clawback_tjoin_lt_tge
* no_clawback_tjoin_gt_tge
* leave_minus_join_gt_3y_and_tge_eq_join

*/

func createTeamAccount(allocation, commencement, duration uint64) types.TeamVestingAccount {
	return types.TeamVestingAccount{
		Id:              0,
		TotalAllocation: allocation,
		Commencement:    commencement,
		Clawback:        commencement + duration,
		UnlockedClaimed: 0,
		LastClaimedTime: 0,
	}
}

var _ = Describe("logic_team_test.go", Ordered, func() {
	const YEAR = uint64(60 * 60 * 24 * 365)
	const MONTH = uint64(5 * 60 * 24 * 365)
	const ALLOCATION = 1_000_000 * i.KYVE

	It("leave_minus_join_lt_1y", func() {
		// ARRANGE
		accountB3y := createTeamAccount(1_000_000*i.KYVE, types.TGE-3*YEAR, YEAR-1)
		accountB2y := createTeamAccount(1_000_000*i.KYVE, types.TGE-2*YEAR, YEAR-1)
		accountB1y := createTeamAccount(1_000_000*i.KYVE, types.TGE-1*YEAR, YEAR-1)
		accountB0y := createTeamAccount(1_000_000*i.KYVE, types.TGE-0*YEAR, YEAR-1)
		accountA1y := createTeamAccount(1_000_000*i.KYVE, types.TGE+1*YEAR, YEAR-1)

		// ASSERT
		status := teamKeeper.GetVestingStatus(accountB3y, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.RemainingUnvestedAmount))

		status = teamKeeper.GetVestingStatus(accountB2y, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.RemainingUnvestedAmount))

		status = teamKeeper.GetVestingStatus(accountB1y, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.RemainingUnvestedAmount))

		status = teamKeeper.GetVestingStatus(accountB0y, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.RemainingUnvestedAmount))

		status = teamKeeper.GetVestingStatus(accountA1y, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.RemainingUnvestedAmount))
	})

	It("leave_minus_join_lt_3y_and_tge_lt_join", func() {
		// ARRANGE
		tjoin := types.TGE + 6*MONTH
		account := createTeamAccount(ALLOCATION, tjoin, 30*MONTH)
		// Maximum is ALLOCATION*30/36

		// ASSERT
		// t < tjoin => everything is unvested
		status := teamKeeper.GetVestingStatus(account, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.LockedVestedAmount))
		Expect(ALLOCATION * 30 / 36).To(Equal(status.RemainingUnvestedAmount))

		// t = tjoin + 1 Year => 1/3 is vested
		status1Y := teamKeeper.GetVestingStatus(account, tjoin+YEAR)
		Expect(ALLOCATION / 3).To(Equal(status1Y.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status1Y.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status1Y.CurrentClaimableAmount))
		Expect(ALLOCATION / 3).To(Equal(status1Y.LockedVestedAmount))
		Expect(ALLOCATION*30/36 - ALLOCATION/3).To(Equal(status1Y.RemainingUnvestedAmount))

		// t = t_leave => max amount (5/6) is vested
		statusAL := teamKeeper.GetVestingStatus(account, tjoin+30*MONTH)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION*30/36))
		Expect(statusAL.TotalUnlockedAmount).To(Equal(statusAL.CurrentClaimableAmount))
		Expect(statusAL.TotalVestedAmount - statusAL.TotalUnlockedAmount).To(Equal(statusAL.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusAL.RemainingUnvestedAmount))

		// t > tjoin + Cliff + Unlock => everything is unlocked
		statusJCU := teamKeeper.GetVestingStatus(account, tjoin+3*YEAR)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(ALLOCATION * 30 / 36).To(Equal(statusJCU.TotalUnlockedAmount))
		Expect(ALLOCATION * 30 / 36).To(Equal(statusJCU.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(statusJCU.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusJCU.RemainingUnvestedAmount))
	})

	It("leave_minus_join_lt_3y_and_tge_gt_join", func() {
		// ARRANGE
		tjoin := types.TGE - 6*MONTH
		account := createTeamAccount(ALLOCATION, tjoin, 30*MONTH)
		// Maximum is ALLOCATION*30/36

		// ASSERT
		// t < tjoin + 1 YR => everything is unvested
		status := teamKeeper.GetVestingStatus(account, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.LockedVestedAmount))
		Expect(ALLOCATION * 30 / 36).To(Equal(status.RemainingUnvestedAmount))

		// t = tjoin + 1 Year => 1/3 is vested
		status1Y := teamKeeper.GetVestingStatus(account, tjoin+YEAR)
		Expect(ALLOCATION / 3).To(Equal(status1Y.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status1Y.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status1Y.CurrentClaimableAmount))
		Expect(ALLOCATION / 3).To(Equal(status1Y.LockedVestedAmount))
		Expect(ALLOCATION*30/36 - ALLOCATION/3).To(Equal(status1Y.RemainingUnvestedAmount))

		// t = t_leave => max amount (5/6) is vested
		statusAL := teamKeeper.GetVestingStatus(account, tjoin+30*MONTH)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION*30/36))
		Expect(statusAL.TotalUnlockedAmount).To(Equal(statusAL.CurrentClaimableAmount))
		Expect(statusAL.TotalVestedAmount - statusAL.TotalUnlockedAmount).To(Equal(statusAL.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusAL.RemainingUnvestedAmount))

		// t = tjoin + Cliff + Unlock => everything is vested but unlock still ongoing
		statusJCU := teamKeeper.GetVestingStatus(account, tjoin+3*YEAR)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusJCU.TotalVestedAmount))
		Expect(statusJCU.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusJCU.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION*30/36))
		Expect(statusJCU.TotalUnlockedAmount).To(Equal(statusJCU.CurrentClaimableAmount))
		Expect(statusJCU.TotalVestedAmount - statusJCU.TotalUnlockedAmount).To(Equal(statusJCU.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusJCU.RemainingUnvestedAmount))

		// t = TGE + Cliff + Unlock => everything is unlocked
		statusAT := teamKeeper.GetVestingStatus(account, types.TGE+3*YEAR)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAT.TotalVestedAmount))
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAT.TotalUnlockedAmount))
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAT.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(statusAT.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusAT.RemainingUnvestedAmount))
	})

	It("leave_minus_join_gt_3y_and_tge_gt_join", func() {
		// ARRANGE
		tjoin := types.TGE - 6*MONTH
		account := createTeamAccount(ALLOCATION, tjoin, 36*MONTH)

		// ASSERT
		// t < tjoin + 1 YR => everything is unvested
		status := teamKeeper.GetVestingStatus(account, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.LockedVestedAmount))
		Expect(ALLOCATION).To(Equal(status.RemainingUnvestedAmount))

		// t = tjoin + 1 Year => 1/3 is vested, t < t_unlock
		status1Y := teamKeeper.GetVestingStatus(account, tjoin+YEAR)
		Expect(ALLOCATION / 3).To(Equal(status1Y.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status1Y.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status1Y.CurrentClaimableAmount))
		Expect(ALLOCATION / 3).To(Equal(status1Y.LockedVestedAmount))
		Expect(ALLOCATION*2/3 + 1).To(Equal(status1Y.RemainingUnvestedAmount))

		// t = t_join * 2.5 Years => max amount (5/6) is vested
		statusAL := teamKeeper.GetVestingStatus(account, tjoin+30*MONTH)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION*30/36))
		Expect(statusAL.TotalUnlockedAmount).To(Equal(statusAL.CurrentClaimableAmount))
		Expect(statusAL.TotalVestedAmount - statusAL.TotalUnlockedAmount).To(Equal(statusAL.LockedVestedAmount))
		Expect(ALLOCATION*1/6 + 1).To(Equal(statusAL.RemainingUnvestedAmount))

		// t = tjoin + Cliff + Unlock => everything is vested but unlock still ongoing
		statusJCU := teamKeeper.GetVestingStatus(account, tjoin+3*YEAR)
		Expect(ALLOCATION).To(Equal(statusJCU.TotalVestedAmount))
		Expect(statusJCU.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusJCU.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION*30/36))
		Expect(statusJCU.TotalUnlockedAmount).To(Equal(statusJCU.CurrentClaimableAmount))
		Expect(statusJCU.TotalVestedAmount - statusJCU.TotalUnlockedAmount).To(Equal(statusJCU.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusJCU.RemainingUnvestedAmount))

		// t = TGE + Cliff + Unlock => everything is unlocked
		statusAT := teamKeeper.GetVestingStatus(account, types.TGE+3*YEAR)
		Expect(ALLOCATION).To(Equal(statusAT.TotalVestedAmount))
		Expect(ALLOCATION).To(Equal(statusAT.TotalUnlockedAmount))
		Expect(ALLOCATION).To(Equal(statusAT.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(statusAT.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusAT.RemainingUnvestedAmount))
	})

	It("leave_minus_join_gt_3y_and_tge_lt_join", func() {
		// ARRANGE
		tjoin := types.TGE + 6*MONTH
		account := createTeamAccount(ALLOCATION, tjoin, 36*MONTH)

		// ASSERT
		// t < tjoin + 1 YR => everything is unvested
		status := teamKeeper.GetVestingStatus(account, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.LockedVestedAmount))
		Expect(ALLOCATION).To(Equal(status.RemainingUnvestedAmount))

		// t = tjoin + 1 Year => 1/3 is vested, t < t_unlock
		status1Y := teamKeeper.GetVestingStatus(account, tjoin+YEAR)
		Expect(ALLOCATION / 3).To(Equal(status1Y.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status1Y.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status1Y.CurrentClaimableAmount))
		Expect(ALLOCATION / 3).To(Equal(status1Y.LockedVestedAmount))
		Expect(ALLOCATION*2/3 + 1).To(Equal(status1Y.RemainingUnvestedAmount))

		// t = t_join * 2.5 Years => (5/6) is vested
		statusAL := teamKeeper.GetVestingStatus(account, tjoin+30*MONTH)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION))
		Expect(statusAL.TotalUnlockedAmount).To(Equal(statusAL.CurrentClaimableAmount))
		Expect(statusAL.TotalVestedAmount - statusAL.TotalUnlockedAmount).To(Equal(statusAL.LockedVestedAmount))
		Expect(ALLOCATION*1/6 + 1).To(Equal(statusAL.RemainingUnvestedAmount))

		// t = tjoin + Cliff + Unlock => everything is vested but unlock still ongoing
		statusJCU := teamKeeper.GetVestingStatus(account, tjoin+3*YEAR)
		Expect(ALLOCATION).To(Equal(statusJCU.TotalVestedAmount))
		Expect(statusJCU.TotalUnlockedAmount).To(Equal(ALLOCATION))
		Expect(statusJCU.CurrentClaimableAmount).To(Equal(ALLOCATION))
		Expect(statusJCU.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(uint64(0)).To(Equal(statusJCU.RemainingUnvestedAmount))
	})

	It("no_clawback_t_join_lt_tge", func() {
		// ARRANGE
		tjoin := types.TGE - 6*MONTH
		account := types.TeamVestingAccount{
			TotalAllocation: ALLOCATION,
			Commencement:    tjoin,
		}

		// ASSERT
		// t < tjoin + 1 YR => everything is unvested
		status := teamKeeper.GetVestingStatus(account, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.LockedVestedAmount))
		Expect(ALLOCATION).To(Equal(status.RemainingUnvestedAmount))

		// t = tjoin + 1 Year => 1/3 is vested, t < t_unlock
		status1Y := teamKeeper.GetVestingStatus(account, tjoin+YEAR)
		Expect(ALLOCATION / 3).To(Equal(status1Y.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status1Y.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status1Y.CurrentClaimableAmount))
		Expect(ALLOCATION / 3).To(Equal(status1Y.LockedVestedAmount))
		Expect(ALLOCATION*2/3 + 1).To(Equal(status1Y.RemainingUnvestedAmount))

		// t = t_join * 2.5 Years => max amount (5/6) is vested
		statusAL := teamKeeper.GetVestingStatus(account, tjoin+30*MONTH)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION*30/36))
		Expect(statusAL.TotalUnlockedAmount).To(Equal(statusAL.CurrentClaimableAmount))
		Expect(statusAL.TotalVestedAmount - statusAL.TotalUnlockedAmount).To(Equal(statusAL.LockedVestedAmount))
		Expect(ALLOCATION*1/6 + 1).To(Equal(statusAL.RemainingUnvestedAmount))

		// t = tjoin + Cliff + Unlock => everything is vested but unlock still ongoing
		statusJCU := teamKeeper.GetVestingStatus(account, tjoin+3*YEAR)
		Expect(ALLOCATION).To(Equal(statusJCU.TotalVestedAmount))
		Expect(statusJCU.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusJCU.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION*30/36))
		Expect(statusJCU.TotalUnlockedAmount).To(Equal(statusJCU.CurrentClaimableAmount))
		Expect(statusJCU.TotalVestedAmount - statusJCU.TotalUnlockedAmount).To(Equal(statusJCU.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusJCU.RemainingUnvestedAmount))

		// t = TGE + Cliff + Unlock => everything is unlocked
		statusAT := teamKeeper.GetVestingStatus(account, types.TGE+3*YEAR)
		Expect(ALLOCATION).To(Equal(statusAT.TotalVestedAmount))
		Expect(ALLOCATION).To(Equal(statusAT.TotalUnlockedAmount))
		Expect(ALLOCATION).To(Equal(statusAT.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(statusAT.LockedVestedAmount))
		Expect(uint64(0)).To(Equal(statusAT.RemainingUnvestedAmount))
	})

	It("no_clawback_t_join_gt_tge", func() {
		// ARRANGE
		tjoin := types.TGE + 6*MONTH
		account := types.TeamVestingAccount{
			TotalAllocation: ALLOCATION,
			Commencement:    tjoin,
		}

		// ASSERT
		// t < tjoin + 1 YR => everything is unvested
		status := teamKeeper.GetVestingStatus(account, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.LockedVestedAmount))
		Expect(ALLOCATION).To(Equal(status.RemainingUnvestedAmount))

		// t = tjoin + 1 Year => 1/3 is vested, t < t_unlock
		status1Y := teamKeeper.GetVestingStatus(account, tjoin+YEAR)
		Expect(ALLOCATION / 3).To(Equal(status1Y.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status1Y.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status1Y.CurrentClaimableAmount))
		Expect(ALLOCATION / 3).To(Equal(status1Y.LockedVestedAmount))
		Expect(ALLOCATION*2/3 + 1).To(Equal(status1Y.RemainingUnvestedAmount))

		// t = t_join * 2.5 Years => (5/6) is vested
		statusAL := teamKeeper.GetVestingStatus(account, tjoin+30*MONTH)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION))
		Expect(statusAL.TotalUnlockedAmount).To(Equal(statusAL.CurrentClaimableAmount))
		Expect(statusAL.TotalVestedAmount - statusAL.TotalUnlockedAmount).To(Equal(statusAL.LockedVestedAmount))
		Expect(ALLOCATION*1/6 + 1).To(Equal(statusAL.RemainingUnvestedAmount))

		// t = tjoin + Cliff + Unlock => everything is vested but unlock still ongoing
		statusJCU := teamKeeper.GetVestingStatus(account, tjoin+3*YEAR)
		Expect(ALLOCATION).To(Equal(statusJCU.TotalVestedAmount))
		Expect(statusJCU.TotalUnlockedAmount).To(Equal(ALLOCATION))
		Expect(statusJCU.CurrentClaimableAmount).To(Equal(ALLOCATION))
		Expect(statusJCU.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(uint64(0)).To(Equal(statusJCU.RemainingUnvestedAmount))
	})

	It("leave_minus_join_gt_3y_and_tge_eq_join", func() {
		// ARRANGE
		tjoin := types.TGE
		account := createTeamAccount(ALLOCATION, tjoin, 36*MONTH)

		// ASSERT
		// t == tjoin => everything is unvested
		status := teamKeeper.GetVestingStatus(account, types.TGE)
		Expect(uint64(0)).To(Equal(status.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status.CurrentClaimableAmount))
		Expect(uint64(0)).To(Equal(status.LockedVestedAmount))
		Expect(ALLOCATION).To(Equal(status.RemainingUnvestedAmount))

		// t = tjoin + 1 Year => 1/3 is vested, t < t_unlock
		status1Y := teamKeeper.GetVestingStatus(account, tjoin+YEAR)
		Expect(ALLOCATION / 3).To(Equal(status1Y.TotalVestedAmount))
		Expect(uint64(0)).To(Equal(status1Y.TotalUnlockedAmount))
		Expect(uint64(0)).To(Equal(status1Y.CurrentClaimableAmount))
		Expect(ALLOCATION / 3).To(Equal(status1Y.LockedVestedAmount))
		Expect(ALLOCATION*2/3 + 1).To(Equal(status1Y.RemainingUnvestedAmount))

		// t = t_join * 2.5 Years => (5/6) is vested
		statusAL := teamKeeper.GetVestingStatus(account, tjoin+30*MONTH)
		Expect(ALLOCATION * 30 / 36).To(Equal(statusAL.TotalVestedAmount))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically(">", 0))
		Expect(statusAL.TotalUnlockedAmount).To(BeNumerically("<=", ALLOCATION))
		Expect(statusAL.TotalUnlockedAmount).To(Equal(statusAL.CurrentClaimableAmount))
		Expect(statusAL.TotalVestedAmount - statusAL.TotalUnlockedAmount).To(Equal(statusAL.LockedVestedAmount))
		Expect(ALLOCATION*1/6 + 1).To(Equal(statusAL.RemainingUnvestedAmount))

		// t = tjoin + Cliff + Unlock => everything is vested but unlock still ongoing
		statusJCU := teamKeeper.GetVestingStatus(account, tjoin+3*YEAR)
		Expect(ALLOCATION).To(Equal(statusJCU.TotalVestedAmount))
		Expect(statusJCU.TotalUnlockedAmount).To(Equal(ALLOCATION))
		Expect(statusJCU.CurrentClaimableAmount).To(Equal(ALLOCATION))
		Expect(statusJCU.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(uint64(0)).To(Equal(statusJCU.RemainingUnvestedAmount))
	})
})
