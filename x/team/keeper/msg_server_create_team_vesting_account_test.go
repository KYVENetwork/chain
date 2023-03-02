package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/team/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_create_team_vesting_account.go

* Create a first TVA with invalid authority
* Create a first TVA with zero allocation
* Create a first TVA with zero allocation and other authority
* Create a first TVA with zero commencement
* Create a first TVA with commencement 3 years before TGE
* Create a first TVA with commencement 3 years before TGE with other authority
* Create TVA with more Allocation than available
* Create multiple TVAs

*/

var _ = Describe("msg_server_create_team_vesting_account.go", Ordered, func() {
	// init new clean chain at TGE time
	s := i.NewCleanChainAtTime(int64(types.TGE))

	BeforeEach(func() {
		// init new clean chain at TGE time
		s = i.NewCleanChainAtTime(int64(types.TGE))
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Create a first TVA with invalid authority", func() {
		// ACT
		s.RunTxTeamError(&types.MsgCreateTeamVestingAccount{
			Authority:       i.ALICE,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - types.VESTING_DURATION,
		})

		// ASSERT
		tvas := s.App().TeamKeeper.GetTeamVestingAccounts(s.Ctx())
		Expect(tvas).To(HaveLen(0))
	})

	It("Create a first TVA with zero allocation", func() {
		// ACT
		s.RunTxTeamError(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 0, // 1m
			Commencement:    types.TGE - types.VESTING_DURATION,
		})

		// ASSERT
		tvas := s.App().TeamKeeper.GetTeamVestingAccounts(s.Ctx())
		Expect(tvas).To(HaveLen(0))
	})

	It("Create a first TVA with zero commencement", func() {
		// ACT
		s.RunTxTeamError(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    0,
		})

		// ASSERT
		tvas := s.App().TeamKeeper.GetTeamVestingAccounts(s.Ctx())
		Expect(tvas).To(HaveLen(0))
	})

	It("Create a first TVA with commencement 3 years before TGE", func() {
		// ACT
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - types.VESTING_DURATION,
		})

		// ASSERT
		tva, found := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(found).To(BeTrue())
		Expect(tva.Commencement).To(Equal(types.TGE - types.VESTING_DURATION))
		Expect(tva.TotalAllocation).To(Equal(1_000_000 * i.KYVE))
		Expect(tva.Clawback).To(BeZero())
		Expect(tva.UnlockedClaimed).To(BeZero())
		Expect(tva.LastClaimedTime).To(BeZero())
		Expect(tva.TotalRewards).To(BeZero())
		Expect(tva.RewardsClaimed).To(BeZero())

		_, found = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 1)
		Expect(found).To(BeFalse())

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(1_000_000 * i.KYVE))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))

		// NOTE: Disable because there is no inflation rewards here
		// Expect(info.TotalAuthorityRewards).To(BeNumerically(">", uint64(0)))
		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.TotalAccountRewards).To(BeZero())
		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(BeZero())

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards))
		Expect(info.TeamModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards))
	})

	It("Create a first TVA with commencement 3 years before TGE with other authority", func() {
		// ACT
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.BCP_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - types.VESTING_DURATION,
		})

		// ASSERT
		tva, found := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(found).To(BeTrue())
		Expect(tva.Commencement).To(Equal(types.TGE - types.VESTING_DURATION))
		Expect(tva.TotalAllocation).To(Equal(1_000_000 * i.KYVE))
		Expect(tva.Clawback).To(BeZero())
		Expect(tva.UnlockedClaimed).To(BeZero())
		Expect(tva.LastClaimedTime).To(BeZero())
		Expect(tva.TotalRewards).To(BeZero())
		Expect(tva.RewardsClaimed).To(BeZero())

		_, found = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 1)
		Expect(found).To(BeFalse())

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(1_000_000 * i.KYVE))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 1_000_000*i.KYVE))

		// NOTE: Disable because there is no inflation rewards here
		// Expect(info.TotalAuthorityRewards).To(BeNumerically(">", uint64(0)))
		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.TotalAccountRewards).To(BeZero())
		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(BeZero())

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards))
		Expect(info.TeamModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards))
	})

	It("Create TVA with more Allocation than available", func() {
		// ACT
		s.RunTxTeamError(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: types.TEAM_ALLOCATION + 1, // 1m
			Commencement:    types.TGE - types.VESTING_DURATION,
		})

		// ASSERT
		_, found := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(found).To(BeFalse())

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(BeZero())
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
	})

	It("Create multiple TVAs", func() {
		// ARRANGE
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 1_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE - types.VESTING_DURATION,
		})

		// ACT
		s.RunTxTeamSuccess(&types.MsgCreateTeamVestingAccount{
			Authority:       types.FOUNDATION_ADDRESS,
			TotalAllocation: 2_000_000 * i.KYVE, // 1m
			Commencement:    types.TGE + types.VESTING_DURATION,
		})

		// ASSERT
		tvas := s.App().TeamKeeper.GetTeamVestingAccounts(s.Ctx())
		Expect(tvas).To(HaveLen(2))

		tva, found := s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 0)
		Expect(found).To(BeTrue())
		Expect(tva.Commencement).To(Equal(types.TGE - types.VESTING_DURATION))
		Expect(tva.TotalAllocation).To(Equal(1_000_000 * i.KYVE))
		Expect(tva.Clawback).To(BeZero())
		Expect(tva.UnlockedClaimed).To(BeZero())
		Expect(tva.LastClaimedTime).To(BeZero())
		Expect(tva.TotalRewards).To(BeZero())
		Expect(tva.RewardsClaimed).To(BeZero())

		tva, found = s.App().TeamKeeper.GetTeamVestingAccount(s.Ctx(), 1)
		Expect(found).To(BeTrue())
		Expect(tva.Commencement).To(Equal(types.TGE + types.VESTING_DURATION))
		Expect(tva.TotalAllocation).To(Equal(2_000_000 * i.KYVE))
		Expect(tva.Clawback).To(BeZero())
		Expect(tva.UnlockedClaimed).To(BeZero())
		Expect(tva.LastClaimedTime).To(BeZero())
		Expect(tva.TotalRewards).To(BeZero())
		Expect(tva.RewardsClaimed).To(BeZero())

		info := s.App().TeamKeeper.GetTeamInfo(s.Ctx())
		Expect(info.FoundationAuthority).To(Equal(types.FOUNDATION_ADDRESS))
		Expect(info.BcpAuthority).To(Equal(types.BCP_ADDRESS))
		Expect(info.TotalTeamAllocation).To(Equal(types.TEAM_ALLOCATION))
		Expect(info.IssuedTeamAllocation).To(Equal(3_000_000 * i.KYVE))
		Expect(info.AvailableTeamAllocation).To(Equal(types.TEAM_ALLOCATION - 3_000_000*i.KYVE))

		// NOTE: Disable because there is no inflation rewards here
		// Expect(info.TotalAuthorityRewards).To(BeNumerically(">", uint64(0)))
		Expect(info.ClaimedAuthorityRewards).To(BeZero())
		Expect(info.AvailableAuthorityRewards).To(Equal(info.TotalAuthorityRewards))

		Expect(info.TotalAccountRewards).To(BeZero())
		Expect(info.ClaimedAccountRewards).To(BeZero())
		Expect(info.AvailableAccountRewards).To(BeZero())

		Expect(info.RequiredModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards))
		Expect(info.TeamModuleBalance).To(Equal(types.TEAM_ALLOCATION + info.TotalAuthorityRewards))
	})
})
