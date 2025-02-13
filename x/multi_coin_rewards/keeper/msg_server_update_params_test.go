package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Gov
	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
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

* Update distribution pending time
* Update distribution pending time with invalid value

* Update policy admin address
* Update policy admin address with invalid value

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
		params := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(params.MultiCoinDistributionPendingTime).To(Equal(types.DefaultMultiCoinDistributionPendingTime))
		Expect(params.MultiCoinDistributionPolicyAdminAddress).To(BeEmpty())
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
			"multi_coin_distribution_pending_time": 10000,
			"multi_coin_distribution_policy_admin_address": "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"
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
		updatedParams := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MultiCoinDistributionPendingTime).To(Equal(uint64(10000)))
		Expect(updatedParams.MultiCoinDistributionPolicyAdminAddress).To(Equal("kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"))
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
		updatedParams := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MultiCoinDistributionPendingTime).To(Equal(types.DefaultMultiCoinDistributionPendingTime))
		Expect(updatedParams.MultiCoinDistributionPolicyAdminAddress).To(BeEmpty())
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			"multi_coin_distribution_pending_time": "0.5",
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
		updatedParams := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MultiCoinDistributionPendingTime).To(Equal(types.DefaultMultiCoinDistributionPendingTime))
		Expect(updatedParams.MultiCoinDistributionPolicyAdminAddress).To(BeEmpty())
	})

	It("Update distribution pending time", func() {
		// ARRANGE
		payload := `{
			"multi_coin_distribution_pending_time": 5
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
		updatedParams := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MultiCoinDistributionPendingTime).To(Equal(uint64(5)))
		Expect(updatedParams.MultiCoinDistributionPolicyAdminAddress).To(BeEmpty())
	})

	It("Update distribution pending time with invalid value", func() {
		// ARRANGE
		payload := `{
			"multi_coin_distribution_pending_time": "5.0"
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
		updatedParams := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MultiCoinDistributionPendingTime).To(Equal(types.DefaultMultiCoinDistributionPendingTime))
		Expect(updatedParams.MultiCoinDistributionPolicyAdminAddress).To(BeEmpty())
	})

	It("Update policy admin address", func() {
		// ARRANGE
		payload := `{
			"multi_coin_distribution_policy_admin_address": "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"
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
		updatedParams := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MultiCoinDistributionPendingTime).To(Equal(types.DefaultMultiCoinDistributionPendingTime))
		Expect(updatedParams.MultiCoinDistributionPolicyAdminAddress).To(Equal("kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"))
	})

	It("Update policy admin address with invalid value", func() {
		// ARRANGE
		payload := `{
			"multi_coin_distribution_policy_admin_address": "kyveabc"
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
		updatedParams := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MultiCoinDistributionPendingTime).To(Equal(types.DefaultMultiCoinDistributionPendingTime))
		Expect(updatedParams.MultiCoinDistributionPolicyAdminAddress).To(BeEmpty())
	})
})
