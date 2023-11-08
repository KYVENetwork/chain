package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Delegation
	"github.com/KYVENetwork/chain/x/funders/types"
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

* Update min-funding-amount
* Update min-funding-amount with invalid value

* Update min-funding-amount-per-bundle
* Update min-funding-amount-per-bundle with invalid value

* Update min-funding-multiple
* Update min-funding-multiple with invalid value

*/

var _ = Describe("msg_server_update_params.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	minDeposit := s.App().GovKeeper.GetParams(s.Ctx()).MinDeposit
	votingPeriod := s.App().GovKeeper.GetParams(s.Ctx()).VotingPeriod

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
		params := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(params.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(params.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(params.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
		)

		// ACT
		_, err := s.RunTx(proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Update every param at once", func() {
		// ARRANGE
		payload := `{
			"min_funding_amount": 2000000000,
			"min_funding_amount_per_bundle": 500000,
			"min_funding_multiple": 25
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(uint64(2_000_000_000)))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(uint64(500_000)))
		Expect(updatedParams.MinFundingMultiple).To(Equal(uint64(25)))
	})

	It("Update no param", func() {
		// ARRANGE
		payload := `{}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			"min_funding_amount": abc,
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update min-funding-amount", func() {
		// ARRANGE
		payload := `{
			"min_funding_amount": 100000000
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(uint64(100_000_000)))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update min-funding-amount with invalid value", func() {
		// ARRANGE
		payload := `{
			"min_funding_amount": "invalid"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("min-funding-amount-per-bundle", func() {
		// ARRANGE
		payload := `{
			"min_funding_amount_per_bundle": 300000
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(uint64(300000)))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update min-funding-amount-per-bundle", func() {
		// ARRANGE
		payload := `{
			"min_funding_amount_per_bundle": "invalid"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update min-funding-multiple", func() {
		// ARRANGE
		payload := `{
			"min_funding_multiple": 9
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(updatedParams.MinFundingMultiple).To(Equal(uint64(9)))
	})

	It("Update min-funding-multiple with invalid value", func() {
		// ARRANGE
		payload := `{
			"min_funding_multiple": -1
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinFundingAmount).To(Equal(types.DefaultMinFundingAmount))
		Expect(updatedParams.MinFundingAmountPerBundle).To(Equal(types.DefaultMinFundingAmountPerBundle))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})
})
