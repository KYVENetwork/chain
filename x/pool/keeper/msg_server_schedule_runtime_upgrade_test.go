package keeper_test

import (
	"time"

	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_schedule_runtime_upgrade.go

* Invalid authority (transaction).
* Invalid authority (proposal).
* Schedule runtime upgrade with no version
* Schedule runtime upgrade with no binaries
* Schedule runtime upgrade in the past
* Schedule runtime upgrade in the future
* Schedule runtime upgrade while another one is ongoing

*/

var _ = Describe("msg_server_schedule_runtime_upgrade.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	votingPeriod := s.App().GovKeeper.GetParams(s.Ctx()).VotingPeriod

	var currentTime uint64

	BeforeEach(func() {
		s = i.NewCleanChain()

		createPoolWithEmptyValues(s)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.UpgradePlan = &types.UpgradePlan{}
		pool.Protocol = &types.Protocol{
			Version:     "0.0.0",
			Binaries:    "{\"linux\":\"test\"}",
			LastUpgrade: 0,
		}
		pool.Runtime = "@kyve/test"
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		currentTime = uint64(time.Now().Unix())
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Invalid authority (transaction).", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   i.DUMMY[0],
			Runtime:     "@kyve/test",
			Version:     "1.0.0",
			ScheduledAt: currentTime,
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Invalid authority (proposal).", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   i.DUMMY[0],
			Runtime:     "@kyve/test",
			Version:     "1.0.0",
			Binaries:    "{}",
			Duration:    60,
			ScheduledAt: currentTime,
		}

		proposal, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, err := s.RunTx(&proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Schedule runtime upgrade with no version", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   gov,
			Runtime:     "@kyve/test",
			Version:     "",
			Binaries:    "{}",
			Duration:    60,
			ScheduledAt: currentTime,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusFailed))
	})

	It("Schedule runtime upgrade with no binaries", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   gov,
			Runtime:     "@kyve/test",
			Version:     "1.0.0",
			Binaries:    "",
			Duration:    60,
			ScheduledAt: currentTime,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusFailed))
	})

	It("Schedule runtime upgrade in the past", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   gov,
			Runtime:     "@kyve/test",
			Version:     "1.0.0",
			Binaries:    "{}",
			Duration:    60,
			ScheduledAt: currentTime - 7*24*3600,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.UpgradePlan).To(Equal(&types.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix()),
			Duration:    60,
		}))
	})

	It("Schedule runtime upgrade in the future", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   gov,
			Runtime:     "@kyve/test",
			Version:     "1.0.0",
			Binaries:    "{}",
			Duration:    60,
			ScheduledAt: currentTime + 7*24*3600,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.UpgradePlan).To(Equal(&types.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: currentTime + 7*24*3600,
			Duration:    60,
		}))
	})

	It("Schedule runtime upgrade while another one is ongoing", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   gov,
			Runtime:     "@kyve/test",
			Version:     "1.0.0",
			Binaries:    "{}",
			Duration:    60,
			ScheduledAt: currentTime + 7*24*3600,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.UpgradePlan).To(Equal(&types.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: currentTime + 7*24*3600,
			Duration:    60,
		}))

		// ACT
		msg = &types.MsgScheduleRuntimeUpgrade{
			Authority:   gov,
			Runtime:     "@kyve/test",
			Version:     "2.0.0",
			Binaries:    "{}",
			Duration:    60,
			ScheduledAt: currentTime + 7*24*3600,
		}

		p, v = BuildGovernanceTxs(s, []sdk.Msg{msg})

		_, submitErr = s.RunTx(&p)
		_, voteErr = s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ = s.App().GovKeeper.GetProposal(s.Ctx(), 2)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.UpgradePlan).To(Equal(&types.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: currentTime + 7*24*3600,
			Duration:    60,
		}))
	})
})
