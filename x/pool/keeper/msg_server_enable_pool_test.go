package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Gov
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_enable_pool.go

* Invalid authority (transaction)
* Invalid authority (proposal)
* Enable a non-existing pool
* Enable pool which is active
* Enable pool which is disabled
* Enable multiple pools

*/

var _ = Describe("msg_server_enable_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	votingPeriod := s.App().GovKeeper.GetParams(s.Ctx()).VotingPeriod

	BeforeEach(func() {
		s = i.NewCleanChain()

		createPoolWithEmptyValues(s)
		createPoolWithEmptyValues(s)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Invalid authority (transaction)", func() {
		// ARRANGE
		msg := &types.MsgEnablePool{
			Authority: i.DUMMY[0],
			Id:        0,
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Invalid authority (proposal)", func() {
		// ARRANGE
		msg := &types.MsgEnablePool{
			Authority: i.DUMMY[0],
			Id:        0,
		}

		proposal, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, err := s.RunTx(&proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Enable a non-existing pool", func() {
		// ARRANGE
		msg := &types.MsgEnablePool{
			Authority: gov,
			Id:        42,
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

	It("Enable pool which is active.", func() {
		// ARRANGE
		msg := &types.MsgEnablePool{
			Authority: gov,
			Id:        0,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusFailed))
		Expect(pool.Disabled).To(BeFalse())
	})

	It("Enable pool which is disabled", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		msg := &types.MsgEnablePool{
			Authority: gov,
			Id:        0,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(pool.Disabled).To(BeFalse())
	})

	It("Enable multiple pools", func() {
		// ARRANGE
		firstPool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		firstPool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), firstPool)

		secondPool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		secondPool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), secondPool)

		msgFirstPool := &types.MsgEnablePool{
			Authority: gov,
			Id:        0,
		}
		msgSecondPool := &types.MsgEnablePool{
			Authority: gov,
			Id:        1,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msgFirstPool, msgSecondPool})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)
		firstPool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		secondPool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))
		Expect(firstPool.Disabled).To(BeFalse())
		Expect(secondPool.Disabled).To(BeFalse())
	})
})
