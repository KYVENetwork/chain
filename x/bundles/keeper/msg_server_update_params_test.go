package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Bundles
	"github.com/KYVENetwork/chain/x/bundles/types"
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

* Update upload timeout
* Update upload timeout with invalid value

* Update storage cost
* Update storage cost with invalid value

* Update network fee
* Update network fee with invalid value

* Update max points
* Update max points with invalid value

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
		params := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(params.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(params.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(params.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(params.MaxPoints).To(Equal(types.DefaultMaxPoints))
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
			"upload_timeout": 20,
			"storage_cost": "0.050000000000000000",
			"network_fee": "0.05",
			"max_points": 15
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(uint64(20)))
		Expect(updatedParams.StorageCost).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
		Expect(updatedParams.NetworkFee).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
		Expect(updatedParams.MaxPoints).To(Equal(uint64(15)))
	})

	It("Update no params", func() {
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			"upload_timeout": 20,
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update upload timeout", func() {
		// ARRANGE
		payload := `{
			"upload_timeout": 20
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(uint64(20)))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update upload timeout with invalid value", func() {
		// ARRANGE
		payload := `{
			"upload_timeout": "invalid"
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update storage cost", func() {
		// ARRANGE
		payload := `{
			"storage_cost": "0.050000000000000000"
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update storage cost with invalid value", func() {
		// ARRANGE
		payload := `{
			"storage_cost": -100
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update network fee", func() {
		// ARRANGE
		payload := `{
			"network_fee": "0.05"
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update network fee with invalid value", func() {
		// ARRANGE
		payload := `{
			"network_fee": "invalid"
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})

	It("Update max points", func() {
		// ARRANGE
		payload := `{
			"max_points": 15
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(uint64(15)))
	})

	It("Update max points with invalid value", func() {
		// ARRANGE
		payload := `{
			"max_points": "invalid"
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
		updatedParams := s.App().BundlesKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.UploadTimeout).To(Equal(types.DefaultUploadTimeout))
		Expect(updatedParams.StorageCost).To(Equal(types.DefaultStorageCost))
		Expect(updatedParams.NetworkFee).To(Equal(types.DefaultNetworkFee))
		Expect(updatedParams.MaxPoints).To(Equal(types.DefaultMaxPoints))
	})
})
