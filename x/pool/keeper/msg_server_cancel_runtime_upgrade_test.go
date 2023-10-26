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

TEST CASES - msg_server_cancel_runtime_upgrade.go

* Invalid authority (transaction).
* Invalid authority (proposal).
* Cancel scheduled runtime upgrade
* Try to cancel upgrade which is already upgrading

*/

var _ = Describe("msg_server_cancel_runtime_upgrade.go", Ordered, func() {
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
		msg := &types.MsgCancelRuntimeUpgrade{
			Authority: i.DUMMY[0],
			Runtime:   "@kyve/test",
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Invalid authority (proposal).", func() {
		// ARRANGE
		msg := &types.MsgCancelRuntimeUpgrade{
			Authority: i.DUMMY[0],
			Runtime:   "@kyve/test",
		}

		proposal, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, err := s.RunTx(&proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Cancel scheduled runtime upgrade", func() {
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

		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		// ACT
		cancel := &types.MsgCancelRuntimeUpgrade{
			Authority: gov,
			Runtime:   "@kyve/test",
		}

		p, v = BuildGovernanceTxs(s, []sdk.Msg{cancel})

		_, submitErr = s.RunTx(&p)
		_, voteErr = s.RunTx(&v)
		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.UpgradePlan).To(Equal(&types.UpgradePlan{
			Version:     "",
			Binaries:    "",
			ScheduledAt: 0,
			Duration:    0,
		}))
	})

	It("Try to cancel upgrade which is already upgrading", func() {
		// ARRANGE
		msg := &types.MsgScheduleRuntimeUpgrade{
			Authority:   gov,
			Runtime:     "@kyve/test",
			Version:     "1.0.0",
			Binaries:    "{}",
			Duration:    60,
			ScheduledAt: currentTime,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		// ACT
		cancel := &types.MsgCancelRuntimeUpgrade{
			Authority: gov,
			Runtime:   "@kyve/test",
		}

		p, v = BuildGovernanceTxs(s, []sdk.Msg{cancel})

		_, submitErr = s.RunTx(&p)
		_, voteErr = s.RunTx(&v)
		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool.UpgradePlan).To(Equal(&types.UpgradePlan{
			Version:     "",
			Binaries:    "",
			ScheduledAt: 0,
			Duration:    0,
		}))
	})
})
