package keeper_test

import (
	"fmt"

	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Delegation
	"github.com/KYVENetwork/chain/x/delegation/types"
	// Gov
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

/*

TEST CASES - msg_server_update_params.go

* Check default params
* Invalid authority (transaction)
* Invalid authority (proposal)
* Update every param at once
* Update no param
* Update with invalid formatted payload

* Update unbonding delegation time
* Update unbonding delegation time with invalid value

* Update redelegation cooldown
* Update redelegation cooldown with invalid value

* Update redelegation max amount
* Update redelegation max amount with invalid value

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

	minDeposit := s.App().GovKeeper.GetDepositParams(s.Ctx()).MinDeposit
	votingPeriod := s.App().GovKeeper.GetVotingParams(s.Ctx()).VotingPeriod

	delegations := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
	voter := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	BeforeEach(func() {
		s = i.NewCleanChain()

		delegations := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
		voter = sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Check default params", func() {
		// ASSERT
		params := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(params.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(params.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(params.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(params.VoteSlash).To(Equal(types.DefaultVoteSlash))
		Expect(params.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(params.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, err := s.RunTx(proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Update every param at once", func() {
		// ARRANGE
		payload := `{
			"unbonding_delegation_time": 3600,
			"redelegation_cooldown": 3600,
			"redelegation_max_amount": 1,
			"vote_slash": "0.05",
			"upload_slash": "0.05",
			"timeout_slash": "0.05"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(uint64(3600)))
		Expect(updatedParams.RedelegationCooldown).To(Equal(uint64(3600)))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(uint64(1)))
		Expect(updatedParams.VoteSlash).To(Equal("0.05"))
		Expect(updatedParams.UploadSlash).To(Equal("0.05"))
		Expect(updatedParams.TimeoutSlash).To(Equal("0.05"))
	})

	It("Update no param", func() {
		// ARRANGE
		payload := `{}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			"unbonding_delegation_time": 3600,
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
	})

	It("Update unbonding delegation time", func() {
		// ARRANGE
		payload := `{
			"unbonding_delegation_time": 3600
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(uint64(3600)))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update unbonding delegation time with invalid value", func() {
		// ARRANGE
		payload := `{
			"unbonding_delegation_time": "invalid"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update redelegation cooldown", func() {
		// ARRANGE
		payload := `{
			"redelegation_cooldown": 3600
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationCooldown).To(Equal(uint64(3600)))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update redelegation cooldown with invalid value", func() {
		// ARRANGE
		payload := `{
			"redelegation_cooldown": "invalid"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update Update redelegation max amount", func() {
		// ARRANGE
		payload := `{
			"redelegation_max_amount": 1
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(uint64(1)))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})

	It("Update Update redelegation max amount with invalid value", func() {
		// ARRANGE
		payload := `{
			"redelegation_max_amount": -2
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal("0.05"))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		fmt.Println(msg.ValidateBasic())
		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal("0.05"))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
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
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal("0.05"))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().DelegationKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UnbondingDelegationTime).To(Equal(types.DefaultUnbondingDelegationTime))
		Expect(updatedParams.RedelegationCooldown).To(Equal(types.DefaultRedelegationCooldown))
		Expect(updatedParams.RedelegationMaxAmount).To(Equal(types.DefaultRedelegationMaxAmount))
		Expect(updatedParams.UploadSlash).To(Equal(types.DefaultUploadSlash))
		Expect(updatedParams.TimeoutSlash).To(Equal(types.DefaultTimeoutSlash))
		Expect(updatedParams.VoteSlash).To(Equal(types.DefaultVoteSlash))
	})
})
