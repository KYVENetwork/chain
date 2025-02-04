package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	"github.com/KYVENetwork/chain/x/team/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_clawback.go

* try_with_invalid_authority
* try_to_apply_clawback_before_tjoin
* try_to_apply_clawback_before_last_claim_time
* apply_clawback
* clawback_multiple_times
* clawback_multiple_accounts

*/

var _ = Describe("msg_server_clawback.go", Ordered, func() {
	s := i.NewCleanChainAtTime(int64(types.TGE))

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChainAtTime(int64(types.TGE))
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("try_with_invalid_authority", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - YEAR,
		})

		s.CommitAfterSeconds(1 * YEAR)
		s.CommitAfterSeconds(1 * MONTH) // One month of unlock
		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(uint64(694_444_444_444)))
		Expect(status.TotalUnlockedAmount).To(Equal(uint64(28_935_185_185)))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(305_555_555_556)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(665_509_259_259)))
		Expect(status.CurrentClaimableAmount).To(Equal(uint64(28_935_185_185)))

		// ACT
		s.RunTxTeamError(&types.MsgClawback{
			Authority: i.ALICE,
			Id:        0,
			Clawback:  uint64(s.Ctx().BlockTime().Unix()),
		})

		// ASSERT
		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.Clawback).To(Equal(uint64(0)))
	})

	It("try_to_apply_clawback_before_tjoin", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - YEAR,
		})

		s.CommitAfterSeconds(1 * YEAR)
		s.CommitAfterSeconds(1 * MONTH) // One month of unlock
		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(uint64(694_444_444_444)))
		Expect(status.TotalUnlockedAmount).To(Equal(uint64(28_935_185_185)))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(305_555_555_556)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(665_509_259_259)))
		Expect(status.CurrentClaimableAmount).To(Equal(uint64(28_935_185_185)))
		s.PerformValidityChecks()

		// ACT
		s.RunTxTeamError(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Clawback:  types.TGE - YEAR - MONTH, // one month before tjoin
		})

		// ASSERT
		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.Clawback).To(BeZero())
	})

	It("try_to_apply_clawback_before_last_claim_time", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - YEAR,
		})

		s.CommitAfterSeconds(1 * YEAR)
		s.CommitAfterSeconds(1 * MONTH) // One month of unlock
		acc, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(acc.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(acc.LastClaimedTime).To(Equal(uint64(0)))

		status := teamKeeper.GetVestingStatus(acc, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(uint64(694_444_444_444)))
		Expect(status.TotalUnlockedAmount).To(Equal(uint64(28_935_185_185)))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(305_555_555_556)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(665_509_259_259)))
		Expect(status.CurrentClaimableAmount).To(Equal(uint64(28_935_185_185)))
		s.PerformValidityChecks()

		// ACT
		s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Amount:    10_935_185_185,
			Recipient: i.ALICE,
		})

		s.RunTxTeamError(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Clawback:  types.TGE, // before unlock claim
		})

		// ASSERT
		acc, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(acc.Clawback).To(BeZero())
		Expect(acc.UnlockedClaimed).To(Equal(uint64(10_935_185_185)))

		status = teamKeeper.GetVestingStatus(acc, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(uint64(694_444_444_444)))
		Expect(status.TotalUnlockedAmount).To(Equal(uint64(28_935_185_185)))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(305_555_555_556)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(665_509_259_259)))
		Expect(status.CurrentClaimableAmount).To(Equal(uint64(18_000_000_000)))
	})

	It("apply_clawback", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - YEAR,
		})

		s.CommitAfterSeconds(1 * YEAR)
		s.CommitAfterSeconds(1 * MONTH) // One month of unlock
		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(uint64(694_444_444_444)))
		Expect(status.TotalUnlockedAmount).To(Equal(uint64(28_935_185_185)))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(305_555_555_556)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(665_509_259_259)))
		Expect(status.CurrentClaimableAmount).To(Equal(uint64(28_935_185_185)))
		s.PerformValidityChecks()

		s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Amount:    20_000 * i.KYVE,
			Recipient: i.ALICE,
		})
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(21_000 * i.KYVE))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(1_000_000 * i.KYVE))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))

		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards))

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 20_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
		s.PerformValidityChecks()

		// ACT
		s.RunTxTeamSuccess(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Clawback:  uint64(s.Ctx().BlockTime().Unix()),
		})

		// ASSERT
		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(uint64(694_444_444_444)))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - uint64(694_444_444_444)))

		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards))

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 20_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status2 := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status2.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(tva.Clawback).To(Equal(uint64(s.Ctx().BlockTime().Unix())))

		s.CommitAfterSeconds(2 * YEAR)

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status3 := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(uint64(694_444_444_444)))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - uint64(694_444_444_444)))

		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards))

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 20_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))

		Expect(status3.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(status3.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status3.CurrentClaimableAmount).To(Equal(status3.TotalUnlockedAmount - 20_000*i.KYVE))
		Expect(status3.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status3.TotalVestedAmount).To(Equal(status.TotalVestedAmount))
	})

	It("apply_clawback_with_other_authority", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.BCP_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - YEAR,
		})

		s.CommitAfterSeconds(1 * YEAR)
		s.CommitAfterSeconds(1 * MONTH) // One month of unlock
		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(uint64(694_444_444_444)))
		Expect(status.TotalUnlockedAmount).To(Equal(uint64(28_935_185_185)))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(305_555_555_556)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(665_509_259_259)))
		Expect(status.CurrentClaimableAmount).To(Equal(uint64(28_935_185_185)))
		s.PerformValidityChecks()

		s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
			Authority: types.BCP_ADDRESS,
			Id:        0,
			Amount:    20_000 * i.KYVE,
			Recipient: i.ALICE,
		})
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(21_000 * i.KYVE))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(1_000_000 * i.KYVE))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))

		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards))

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 20_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
		s.PerformValidityChecks()

		// ACT
		s.RunTxTeamSuccess(&types.MsgClawback{
			Authority: types.BCP_ADDRESS,
			Id:        0,
			Clawback:  uint64(s.Ctx().BlockTime().Unix()),
		})

		// ASSERT
		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(uint64(694_444_444_444)))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - uint64(694_444_444_444)))

		Expect(info.TotalAuthorityRewards).To(BeNumerically(">", uint64(0)))
		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.TotalAccountRewards).To(BeNumerically(">", uint64(0)))
		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards))

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 20_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status2 := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status2.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(tva.Clawback).To(Equal(uint64(s.Ctx().BlockTime().Unix())))

		s.CommitAfterSeconds(2 * YEAR)

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status3 := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(uint64(694_444_444_444)))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - uint64(694_444_444_444)))

		Expect(info.TotalAuthorityRewards).To(BeNumerically(">", uint64(0)))
		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.TotalAccountRewards).To(BeNumerically(">", uint64(0)))
		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards))

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 20_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))

		Expect(status3.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(status3.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status3.CurrentClaimableAmount).To(Equal(status3.TotalUnlockedAmount - 20_000*i.KYVE))
		Expect(status3.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status3.TotalVestedAmount).To(Equal(status.TotalVestedAmount))
	})

	It("clawback_multiple_times", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE,
		})

		s.CommitAfterSeconds(3 * YEAR) // vesting is done and nothing has claimed yet
		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.Clawback).To(Equal(uint64(0)))
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.TotalUnlockedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status.CurrentClaimableAmount).To(Equal(1_000_000 * i.KYVE))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))
		s.PerformValidityChecks()

		// ACT
		// clawback before cliff
		s.RunTxTeamSuccess(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Clawback:  types.TGE + MONTH,
		})

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.Clawback).To(Equal(types.TGE + MONTH))
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status = teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(BeZero())
		Expect(status.TotalUnlockedAmount).To(BeZero())
		Expect(status.RemainingUnvestedAmount).To(BeZero())
		Expect(status.LockedVestedAmount).To(BeZero())
		Expect(status.CurrentClaimableAmount).To(BeZero())

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION))

		// clawback right in the middle
		s.RunTxTeamSuccess(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Clawback:  types.TGE + YEAR + 6*MONTH,
		})

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.Clawback).To(Equal(types.TGE + YEAR + 6*MONTH))
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status = teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(500_000 * i.KYVE))
		Expect(status.TotalUnlockedAmount).To(Equal(500_000 * i.KYVE))
		Expect(status.RemainingUnvestedAmount).To(BeZero())
		Expect(status.LockedVestedAmount).To(BeZero())
		Expect(status.CurrentClaimableAmount).To(Equal(500_000 * i.KYVE))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 500_000*i.KYVE))

		// clawback after vesting period
		s.RunTxTeamSuccess(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Clawback:  types.TGE + 3*YEAR + 6*MONTH,
		})

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.Clawback).To(Equal(types.TGE + 3*YEAR + 6*MONTH))
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status = teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.TotalUnlockedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status.CurrentClaimableAmount).To(Equal(1_000_000 * i.KYVE))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))

		// reset clawback
		s.RunTxTeamSuccess(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Clawback:  0,
		})

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.Clawback).To(Equal(uint64(0)))
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status = teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.TotalUnlockedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status.CurrentClaimableAmount).To(Equal(1_000_000 * i.KYVE))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))
	})

	It("clawback_multiple_accounts", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE,
		})

		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 2_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE + YEAR,
		})

		s.CommitAfterSeconds(4 * YEAR) // vesting is done and nothing has claimed yet
		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 1)
		Expect(tva.Clawback).To(Equal(uint64(0)))
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(2_000_000 * i.KYVE))
		Expect(status.TotalUnlockedAmount).To(Equal(2_000_000 * i.KYVE))
		Expect(status.RemainingUnvestedAmount).To(Equal(uint64(0)))
		Expect(status.LockedVestedAmount).To(Equal(uint64(0)))
		Expect(status.CurrentClaimableAmount).To(Equal(2_000_000 * i.KYVE))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 3_000_000*i.KYVE))
		s.PerformValidityChecks()

		// ACT
		// clawback right in the middle
		s.RunTxTeamSuccess(&types.MsgClawback{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        1,
			Clawback:  types.TGE + 2*YEAR + 6*MONTH,
		})

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 1)
		Expect(tva.Clawback).To(Equal(types.TGE + 2*YEAR + 6*MONTH))
		Expect(tva.UnlockedClaimed).To(Equal(uint64(0)))
		Expect(tva.LastClaimedTime).To(Equal(uint64(0)))

		status = teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.TotalVestedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.TotalUnlockedAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(status.RemainingUnvestedAmount).To(BeZero())
		Expect(status.LockedVestedAmount).To(BeZero())
		Expect(status.CurrentClaimableAmount).To(Equal(1_000_000 * i.KYVE))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE - 1_000_000*i.KYVE))
	})
})
