package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	"github.com/KYVENetwork/chain/x/team/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_claim_account_rewards.go

* invalid_authority
* claim_more_rewards_than_available
* partially_claim_rewards_once
* partially_claim_rewards_once_with_other_authority
* claim_rewards_with_3_months_interval

*/

var _ = Describe("msg_server_claim_account_rewards.go", Ordered, func() {
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
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE,
		})

		s.CommitAfterSeconds(2 * YEAR)

		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)

		Expect(tva.TotalRewards).To(BeNumerically(">", uint64(0)))
		Expect(tva.RewardsClaimed).To(BeZero())

		// ACT
		_, err := s.RunTx(&types.MsgClaimAccountRewards{
			Authority: i.ALICE,
			Id:        0,
			Amount:    tva.TotalRewards,
			Recipient: i.BOB,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("claim_more_rewards_than_available", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE,
		})

		s.CommitAfterSeconds(2 * YEAR)

		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)

		Expect(tva.TotalRewards).To(BeNumerically(">", uint64(0)))
		Expect(tva.RewardsClaimed).To(BeZero())

		// ACT
		_, err := s.RunTx(&types.MsgClaimAccountRewards{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Amount:    tva.TotalRewards + 1,
			Recipient: i.ALICE,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("partially_claim_rewards_once", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE,
		})

		s.CommitAfterSeconds(2 * YEAR)

		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)

		Expect(tva.TotalRewards).To(BeNumerically(">", uint64(0)))
		Expect(tva.RewardsClaimed).To(BeZero())
		s.PerformValidityChecks()

		// ACT
		s.RunTxTeamSuccess(&types.MsgClaimAccountRewards{
			Authority: types.FOUNDATION_ADDRESS,
			Id:        0,
			Amount:    100,
			Recipient: i.ALICE,
		})

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_000*i.KYVE + 100))

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.TotalRewards).To(BeNumerically(">", uint64(0)))
		Expect(tva.RewardsClaimed).To(Equal(uint64(100)))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.ClaimedAccountRewards).To(Equal(uint64(100)))
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards - uint64(100)))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - uint64(100)))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
	})

	It("partially_claim_rewards_once_with_other_authority", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.BCP_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE,
		})

		s.CommitAfterSeconds(2 * YEAR)

		tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)

		Expect(tva.TotalRewards).To(BeNumerically(">", uint64(0)))
		Expect(tva.RewardsClaimed).To(BeZero())
		s.PerformValidityChecks()

		// ACT
		s.RunTxTeamSuccess(&types.MsgClaimAccountRewards{
			Authority: types.BCP_ADDRESS,
			Id:        0,
			Amount:    100,
			Recipient: i.ALICE,
		})

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_000*i.KYVE + 100))

		tva, _ = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(tva.TotalRewards).To(BeNumerically(">", uint64(0)))
		Expect(tva.RewardsClaimed).To(Equal(uint64(100)))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.ClaimedAccountRewards).To(Equal(uint64(100)))
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards - uint64(100)))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - uint64(100)))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
	})

	It("claim_rewards_with_3_months_interval", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE,
		})

		totalClaimed := uint64(0)
		s.PerformValidityChecks()

		// ACT
		for m := 1; m <= 16; m++ {
			s.CommitAfterSeconds(3 * MONTH)

			tva, _ := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
			status := teamKeeper.GetVestingStatus(tva, uint64(s.Ctx().BlockTime().Unix()))
			rewards := tva.TotalRewards - tva.RewardsClaimed

			// account should only receive inflation rewards if it has vested $KYVE
			if m < 4 {
				Expect(rewards).To(BeZero())
			} else if m <= 12 {
				Expect(rewards).To(BeNumerically(">", uint64(0)))
			} else {
				Expect(rewards).To(BeZero())
			}

			s.RunTxTeamSuccess(&types.MsgClaimUnlocked{
				Authority: types.FOUNDATION_ADDRESS,
				Id:        0,
				Amount:    status.CurrentClaimableAmount,
				Recipient: i.BOB,
			})

			s.RunTxTeamSuccess(&types.MsgClaimAccountRewards{
				Authority: types.FOUNDATION_ADDRESS,
				Id:        0,
				Amount:    rewards,
				Recipient: i.ALICE,
			})

			totalClaimed += rewards
		}

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_000*i.KYVE + totalClaimed))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.ClaimedAccountRewards).To(Equal(totalClaimed))
		Expect(info.AvailableAccountRewards).To(Equal(info.TotalAccountRewards - totalClaimed))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards + info.TotalAccountRewards - totalClaimed - 1_000_000*i.KYVE))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
	})
})
