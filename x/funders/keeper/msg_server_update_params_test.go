package keeper_test

import (
	"time"

	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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

* Update existing coin whitelist entry
* Update existing coin whitelist entry with invalid value
* Update multiple coin whitelist entries
* Update coin whitelist entry without the native kyve coin

* Update min-funding-multiple
* Update min-funding-multiple with invalid value

*/

var _ = Describe("msg_server_update_params.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var gov string
	var minDeposit sdk.Coins
	var votingPeriod *time.Duration
	var voter sdk.AccAddress
	var delegations []stakingtypes.Delegation

	BeforeEach(func() {
		s = i.NewCleanChain()

		// set whitelist
		s.App().FundersKeeper.SetParams(s.Ctx(), types.DefaultParams())

		gov = s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		params, err := s.App().GovKeeper.Params.Get(s.Ctx())
		Expect(err).NotTo(HaveOccurred())

		minDeposit = params.MinDeposit
		votingPeriod = params.VotingPeriod

		delegations, err = s.App().StakingKeeper.GetAllDelegations(s.Ctx())
		Expect(err).NotTo(HaveOccurred())

		voter = sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Check default params", func() {
		// ASSERT
		params := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(params.CoinWhitelist).To(Equal(types.DefaultParams().CoinWhitelist))
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
			"coin_whitelist": [{"coin_denom":"tkyve","coin_decimals":6,"min_funding_amount":20000000000,"min_funding_amount_per_bundle":2000000000,"coin_weight":"5"}],
			"min_funding_multiple": 25
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(HaveLen(1))
		Expect(updatedParams.CoinWhitelist[0].CoinDenom).To(Equal("tkyve"))
		Expect(updatedParams.CoinWhitelist[0].CoinDecimals).To(Equal(uint32(6)))
		Expect(updatedParams.CoinWhitelist[0].MinFundingAmount).To(Equal(uint64(20000000000)))
		Expect(updatedParams.CoinWhitelist[0].MinFundingAmountPerBundle).To(Equal(uint64(2000000000)))
		Expect(updatedParams.CoinWhitelist[0].CoinWeight.TruncateInt64()).To(Equal(int64(5)))

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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(Equal(types.DefaultParams().CoinWhitelist))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			"min_funding_amount_multiple": abc,
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(Equal(types.DefaultParams().CoinWhitelist))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update existing coin whitelist entry", func() {
		// ARRANGE
		payload := `{
			"coin_whitelist": [{"coin_denom":"tkyve","coin_decimals":9,"min_funding_amount":20000000000,"min_funding_amount_per_bundle":200000,"coin_weight":"7"}]
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(HaveLen(1))
		Expect(updatedParams.CoinWhitelist[0].CoinDenom).To(Equal("tkyve"))
		Expect(updatedParams.CoinWhitelist[0].CoinDecimals).To(Equal(uint32(9)))
		Expect(updatedParams.CoinWhitelist[0].MinFundingAmount).To(Equal(uint64(20000000000)))
		Expect(updatedParams.CoinWhitelist[0].MinFundingAmountPerBundle).To(Equal(uint64(200000)))
		Expect(updatedParams.CoinWhitelist[0].CoinWeight.TruncateInt64()).To(Equal(int64(7)))

		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update existing coin whitelist entry with invalid value", func() {
		// ARRANGE
		payload := `{
			"coin_whitelist": [{"coin_denom":"tkyve","coin_decimals":6,"min_funding_amount":invalid,"min_funding_amount_per_bundle":100000,"coin_weight":"1"}]
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(Equal(types.DefaultParams().CoinWhitelist))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update multiple coin whitelist entries", func() {
		// ARRANGE
		payload := `{
			"coin_whitelist": [{"coin_denom":"tkyve","coin_decimals":9,"min_funding_amount":20000000000,"min_funding_amount_per_bundle":200000,"coin_weight":"5"},{"coin_denom":"acoin","coin_decimals":12,"min_funding_amount":10000000000,"min_funding_amount_per_bundle":100000,"coin_weight":"2"}]
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(HaveLen(2))

		Expect(updatedParams.CoinWhitelist[0].CoinDenom).To(Equal("tkyve"))
		Expect(updatedParams.CoinWhitelist[0].CoinDecimals).To(Equal(uint32(9)))
		Expect(updatedParams.CoinWhitelist[0].MinFundingAmount).To(Equal(uint64(20000000000)))
		Expect(updatedParams.CoinWhitelist[0].MinFundingAmountPerBundle).To(Equal(uint64(200000)))
		Expect(updatedParams.CoinWhitelist[0].CoinWeight.TruncateInt64()).To(Equal(int64(5)))

		Expect(updatedParams.CoinWhitelist[1].CoinDenom).To(Equal("acoin"))
		Expect(updatedParams.CoinWhitelist[1].CoinDecimals).To(Equal(uint32(12)))
		Expect(updatedParams.CoinWhitelist[1].MinFundingAmount).To(Equal(uint64(10000000000)))
		Expect(updatedParams.CoinWhitelist[1].MinFundingAmountPerBundle).To(Equal(uint64(100000)))
		Expect(updatedParams.CoinWhitelist[1].CoinWeight.TruncateInt64()).To(Equal(int64(2)))

		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})

	It("Update coin whitelist entry without the native kyve coin", func() {
		// ARRANGE
		payload := `{
			"coin_whitelist": [{"coin_denom":"acoin","coin_decimals":6,"min_funding_amount":10000000000,"min_funding_amount_per_bundle":100000,"coin_weight":"2"}]
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())
		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(Equal(types.DefaultParams().CoinWhitelist))
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
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(Equal(types.DefaultParams().CoinWhitelist))
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
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.CoinWhitelist).To(Equal(types.DefaultParams().CoinWhitelist))
		Expect(updatedParams.MinFundingMultiple).To(Equal(types.DefaultMinFundingMultiple))
	})
})
