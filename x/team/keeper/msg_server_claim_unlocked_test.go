package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	"github.com/KYVENetwork/chain/x/team/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_claim_unlocked.go

* invalid_authority
* claim_zero_unlocked
* partially_claim_unlocked_once
* claim_entire_allocation_with_3_months_interval
* claim_twice_in_same_block

*/

func appendTeamVestingAccount(s *i.KeeperTestSuite, commencement, clawback uint64) {
	s.App().TeamKeeper.AppendTeamVestingAccount(s.Ctx(), types.TeamVestingAccount{
		TotalAllocation: 1_000_000 * i.KYVE,
		Commencement:    commencement,
		Clawback:        clawback,
		UnlockedClaimed: 0,
		LastClaimedTime: 0,
	})
}

const (
	YEAR       = uint64(60 * 60 * 24 * 365)
	MONTH      = uint64(5 * 60 * 24 * 365)
	ALLOCATION = 1_000_000 * i.KYVE
)

var _ = Describe("msg_server_claim_unlocked.go", Ordered, func() {
	s := i.NewCleanChainAtTime(int64(types.TGE))

	BeforeEach(func() {
		// init new clean chain at TGE time
		s = i.NewCleanChainAtTime(int64(types.TGE))
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("invalid_authority", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.AUTHORITY_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - YEAR,
		})

		s.CommitAfterSeconds(3 * YEAR)

		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
		Expect(status.CurrentClaimableAmount).To(Equal(1_000_000 * i.KYVE))
		s.PerformValidityChecks()

		// ACT
		_, err := s.RunTx(&types.MsgClaimUnlocked{
			Authority: i.BOB,
			Id:        0,
			Amount:    1_000_000 * i.KYVE,
			Recipient: i.ALICE,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("claim_zero_unlocked", func() {
		// ARRANGE
		appendTeamVestingAccount(s, types.TGE-11*MONTH, 0)

		// ASSERT
		s.RunTxTeamError(&types.MsgClaimUnlocked{
			Authority: types.AUTHORITY_ADDRESS,
			Id:        0,
			Amount:    100,
			Recipient: i.ALICE,
		})

		s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
			Authority: types.AUTHORITY_ADDRESS,
			Id:        0,
			Amount:    0,
			Recipient: i.ALICE,
		})
	})

	It("partially_claim_unlocked_once", func() {
		// ARRANGE
		appendTeamVestingAccount(s, types.TGE, 0)

		s.CommitAfterSeconds(3 * YEAR)

		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

		Expect(status.CurrentClaimableAmount).To(Equal(1_000_000 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_000 * i.KYVE))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
		s.PerformValidityChecks()

		// ACT
		s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
			Authority: types.AUTHORITY_ADDRESS,
			Id:        0,
			Amount:    100_000 * i.KYVE,
			Recipient: i.ALICE,
		})

		// ASSERT
		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status = teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

		Expect(status.CurrentClaimableAmount).To(Equal(900_000 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(101_000 * i.KYVE))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 100_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
	})

	It("claim_entire_allocation_with_3_months_interval", func() {
		// ARRANGE
		appendTeamVestingAccount(s, types.TGE, 0)

		s.CommitAfterSeconds(YEAR)

		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

		Expect(status.CurrentClaimableAmount).To(BeZero())
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_000 * i.KYVE))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
		s.PerformValidityChecks()

		// ACT
		for m := 1; m <= 8; m++ {
			s.CommitAfterSeconds(3 * MONTH)

			tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
			status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

			s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
				Authority: types.AUTHORITY_ADDRESS,
				Id:        0,
				Amount:    status.CurrentClaimableAmount,
				Recipient: i.ALICE,
			})
		}

		// ASSERT
		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status = teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

		Expect(status.CurrentClaimableAmount).To(BeZero())
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_001_000 * i.KYVE))

		info = s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 1_000_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
	})

	It("claim_twice_in_same_block", func() {
		// ARRANGE
		appendTeamVestingAccount(s, types.TGE-YEAR, 0)

		s.CommitAfterSeconds(3 * YEAR)

		// ACT
		s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
			Authority: types.AUTHORITY_ADDRESS,
			Id:        0,
			Amount:    ALLOCATION / 2,
			Recipient: i.ALICE,
		})

		s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
			Authority: types.AUTHORITY_ADDRESS,
			Id:        0,
			Amount:    ALLOCATION / 2,
			Recipient: i.ALICE,
		})

		s.RunTxTeamError(&types.MsgClaimUnlocked{
			Authority: types.AUTHORITY_ADDRESS,
			Id:        0,
			Amount:    1,
			Recipient: i.ALICE,
		})

		// ASSERT
		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))

		Expect(status.CurrentClaimableAmount).To(BeZero())
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_001_000 * i.KYVE))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - 1_000_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
	})
})
