package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Gov
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	// Stakers
	"github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_update_params.go

* Check default params
* Invalid authority (transaction)
* Invalid authority (proposal)
* Update every param at once
* Update no param
* Update with invalid formatted payload

* Update unbonding staking time
* Update unbonding staking time with invalid value

* Update commission change time
* Update commission change time with invalid value

* Update leave pool time
* Update leave pool time with invalid value

* Update stake fraction change time
* Update stake fraction change time with invalid value

* Update vote slash
* Update vote slash with invalid value

* Update upload slash
* Update upload slash with invalid value

* Update timeout slash
* Update timeout slash with invalid value

*/

var _ = Describe("msg_server_update_params.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	params, _ := s.App().GovKeeper.Params.Get(s.Ctx())
	minDeposit := params.MinDeposit
	votingPeriod := params.VotingPeriod

	delegations, _ := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
	voter := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	BeforeEach(func() {
		s = i.NewCleanChain()

		delegations, _ := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
		voter = sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Check default params", func() {
		// ASSERT
		params := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(params.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(params.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
	})

	It("Invalid authority (transaction)", func() {
		// ARRANGE
		msg := &types.MsgUpdateParams{
			Authority: i.DUMMY[0],
			Payload:   "{}",
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Invalid authority (proposal)", func() {
		// ARRANGE
		msg := &types.MsgUpdateParams{
			Authority: i.DUMMY[0],
			Payload:   "{}",
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, err := s.RunTx(proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Update every param at once", func() {
		// ARRANGE
		payload := `{
			"unbonding_staking_time": 5,
			"commission_change_time": 5,
			"leave_pool_time": 5
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(uint64(5)))
		Expect(updatedParams.LeavePoolTime).To(Equal(uint64(5)))
	})

	It("Update no params", func() {
		// ARRANGE
		payload := `{}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(updatedParams.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			"vote_slash": "0.5",
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(updatedParams.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
	})

	It("Update commission change time", func() {
		// ARRANGE
		payload := `{
			"commission_change_time": 5
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(uint64(5)))
		Expect(updatedParams.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
		Expect(updatedParams.StakeFractionChangeTime).To(Equal(types.DefaultStakeFractionChangeTime))
	})

	It("Update commission change time with invalid value", func() {
		// ARRANGE
		payload := `{
			"commission_change_time": "5"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(updatedParams.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
		Expect(updatedParams.StakeFractionChangeTime).To(Equal(types.DefaultStakeFractionChangeTime))
	})

	It("Update leave pool time", func() {
		// ARRANGE
		payload := `{
			"leave_pool_time": 5
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(updatedParams.LeavePoolTime).To(Equal(uint64(5)))
		Expect(updatedParams.StakeFractionChangeTime).To(Equal(types.DefaultStakeFractionChangeTime))
	})

	It("Update leave pool time with invalid value", func() {
		// ARRANGE
		payload := `{
			"leave_pool_time": -5
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(updatedParams.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
		Expect(updatedParams.StakeFractionChangeTime).To(Equal(types.DefaultStakeFractionChangeTime))
	})

	It("Update stake fraction change time", func() {
		// ARRANGE
		payload := `{
			"stake_fraction_change_time": 5
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(updatedParams.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
		Expect(updatedParams.StakeFractionChangeTime).To(Equal(uint64(5)))
	})

	It("Update stake fraction change time with invalid value", func() {
		// ARRANGE
		payload := `{
			"stake_fraction_change_time": "5"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CommissionChangeTime).To(Equal(types.DefaultCommissionChangeTime))
		Expect(updatedParams.LeavePoolTime).To(Equal(types.DefaultLeavePoolTime))
		Expect(updatedParams.StakeFractionChangeTime).To(Equal(types.DefaultStakeFractionChangeTime))
	})

	It("Update vote slash", func() {
		// ARRANGE
		payload := `{
			"vote_slash": "0.05"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
	})

	It("Update vote slash with invalid value", func() {
		// ARRANGE
		payload := `{
			"vote_slash": "invalid"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update upload slash", func() {
		// ARRANGE
		payload := `{
			"upload_slash": "0.05"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadSlash).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update upload slash with invalid value", func() {
		// ARRANGE
		payload := `{
			"upload_slash": "1.5"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update timeout slash", func() {
		// ARRANGE
		payload := `{
			"timeout_slash": "0.05"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update timeout slash with invalid value", func() {
		// ARRANGE
		payload := `{
			"upload_slash": "-0.5"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().StakersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})
})
