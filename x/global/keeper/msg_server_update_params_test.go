package keeper_test

import (
	"fmt"

	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Global
	"github.com/KYVENetwork/chain/x/global/types"
	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
* Update min gas price
* Update min gas price with invalid value

* Update burn ratio
* Update burn ratio with invalid value

* Update gas adjustments
* Update gas adjustments with invalid value

* Update gas refunds
* Update gas refunds with invalid value

* Update min initial deposit ratio
* Update min initial deposit ratio with invalid value

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
		params := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(params.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(params.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(params.GasAdjustments).To(BeNil())
		Expect(params.GasRefunds).To(BeNil())
		Expect(params.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
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
		Expect(err.Error()).To(Equal(
			fmt.Sprintf(
				"%s: %s",
				fmt.Sprintf("invalid authority; expected %s, got %s", gov, i.DUMMY[0]),
				govTypes.ErrInvalidSigner,
			),
		))
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
		Expect(err.Error()).To(Equal(
			fmt.Sprintf("%s: %s", i.DUMMY[0], govTypes.ErrInvalidSigner),
		))
	})

	It("Update every param at once", func() {
		// ARRANGE
		payload := `{
			"min_gas_price": "1.5",
			"burn_ratio": "0.2",
			"gas_adjustments": [{
				"type": "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				"amount": 20
			}],
			"gas_refunds": [{
				"type": "/kyve.bundles.v1beta1.MsgSubmitBundleProposal",
				"fraction": "0.75"
			}],
			"min_initial_deposit_ratio": "0.2"
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(sdk.MustNewDecFromStr("1.5")))
		Expect(updatedParams.BurnRatio).To(Equal(sdk.MustNewDecFromStr("0.2")))
		Expect(updatedParams.GasAdjustments).To(Equal([]types.GasAdjustment{
			{
				Type:   "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				Amount: 20,
			},
		}))
		Expect(updatedParams.GasRefunds).To(Equal([]types.GasRefund{
			{
				Type:     "/kyve.bundles.v1beta1.MsgSubmitBundleProposal",
				Fraction: sdk.MustNewDecFromStr("0.75"),
			},
		}))
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(sdk.MustNewDecFromStr("0.2")))
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			min_gas_price: "0.5",
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update min gas price", func() {
		// ARRANGE
		payload := `{
			"min_gas_price": "1.5"
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(sdk.MustNewDecFromStr("1.5")))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update min gas price with invalid value", func() {
		// ARRANGE
		payload := `{
			"min_gas_price": "hello"
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update burn ratio", func() {
		// ARRANGE
		payload := `{
			"burn_ratio": "0.5"
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(sdk.MustNewDecFromStr("0.5")))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update burn ratio with invalid value", func() {
		// ARRANGE
		payload := `{
			"burn_ratio": "1.1"
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update gas adjustments", func() {
		// ARRANGE
		payload := `{
			"gas_adjustments": [{
				"type": "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				"amount": 20
			}]
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(Equal([]types.GasAdjustment{
			{
				Type:   "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				Amount: 20,
			},
		}))
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update gas adjustments with invalid value", func() {
		// ARRANGE
		payload := `{
			"gas_adjustments": [{
				"type": "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				"amount": -20
			}],
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update gas refunds", func() {
		// ARRANGE
		payload := `{
			"gas_refunds": [{
				"type": "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				"fraction": "0.5"
			}]
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(Equal([]types.GasRefund{
			{
				Type:     "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				Fraction: sdk.MustNewDecFromStr("0.5"),
			},
		}))
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update gas refunds with invalid value", func() {
		// ARRANGE
		payload := `{
			"gas_refunds": [{
				"type": "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				"fraction": "-1.5"
			}]
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})

	It("Update min gas price", func() {
		// ARRANGE
		payload := `{
			"min_initial_deposit_ratio": "0.5"
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(sdk.MustNewDecFromStr("0.5")))
	})

	It("Update min gas price with invalid value", func() {
		// ARRANGE
		payload := `{
			"min_initial_deposit_ratio": "1.5"
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
		updatedParams := s.App().GlobalKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.MinGasPrice).To(Equal(types.DefaultMinGasPrice))
		Expect(updatedParams.BurnRatio).To(Equal(types.DefaultBurnRatio))
		Expect(updatedParams.GasAdjustments).To(BeNil())
		Expect(updatedParams.GasRefunds).To(BeNil())
		Expect(updatedParams.MinInitialDepositRatio).To(Equal(types.DefaultMinInitialDepositRatio))
	})
})
