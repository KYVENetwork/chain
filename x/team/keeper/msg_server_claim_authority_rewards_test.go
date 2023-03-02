package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/team/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_claim_authority_rewards.go

* invalid_authority
* invalid_second_authority
* claim_more_rewards_than_available
* partially_claim_rewards_once
* claim_rewards_with_3_months_interval

*/

var _ = Describe("msg_server_claim_authority_rewards.go", Ordered, func() {
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
		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())

		// ACT
		_, err := s.RunTx(&types.MsgClaimAuthorityRewards{
			Authority: i.ALICE,
			Amount:    info.AvailableAuthorityRewards,
			Recipient: i.BOB,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("invalid_second_authority", func() {
		// ARRANGE
		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())

		// ACT
		_, err := s.RunTx(&types.MsgClaimAuthorityRewards{
			Authority: types.BCP_ADDRESS,
			Amount:    info.AvailableAuthorityRewards,
			Recipient: i.BOB,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("claim_more_rewards_than_available", func() {
		// ARRANGE
		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())

		// ACT
		_, err := s.RunTx(&types.MsgClaimAuthorityRewards{
			Authority: types.FOUNDATION_ADDRESS,
			Amount:    info.AvailableAuthorityRewards + 1,
			Recipient: i.ALICE,
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("partially_claim_rewards_once", func() {
		// ACT
		s.RunTxTeamSuccess(&types.MsgClaimAuthorityRewards{
			Authority: types.FOUNDATION_ADDRESS,
			Amount:    100,
			Recipient: i.ALICE,
		})

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_000*i.KYVE + 100))

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.ClaimedAuthorityRewards).To(Equal(uint64(100)))
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards - uint64(100)))
		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards - uint64(100)))
		Expect(info.TeamModuleBalance).To(Equal(info.RequiredModuleBalance))
	})

	It("claim_rewards_with_3_months_interval", func() {
		// ARRANGE
		totalClaimed := uint64(0)

		// ACT
		for m := 1; m <= 12; m++ {
			s.CommitAfterSeconds(3 * MONTH)

			info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())

			s.RunTxTeamSuccess(&types.MsgClaimAuthorityRewards{
				Authority: types.FOUNDATION_ADDRESS,
				Amount:    info.AvailableAuthorityRewards,
				Recipient: i.ALICE,
			})

			totalClaimed += info.AvailableAuthorityRewards
		}

		// ASSERT
		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())

		Expect(info.ClaimedAuthorityRewards).To(Equal(totalClaimed))
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards - totalClaimed))

		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(1_000*i.KYVE + totalClaimed))
	})
})
