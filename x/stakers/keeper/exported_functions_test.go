package keeper_test

import (
	"cosmossdk.io/math"
	"strconv"

	"github.com/cosmos/cosmos-sdk/x/gov/keeper"

	pooltypes "github.com/KYVENetwork/chain/x/pool/types"

	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Delegation
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	// Stakers
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - Protocol Governance Voting

* Protocol validator doesn't vote, delegator votes.
* Protocol validator votes, delegator doesn't.
* Protocol validator votes, delegator votes the same.
* Protocol validator votes, delegator votes different.

*/

var _ = Describe("Protocol Governance Voting", Ordered, func() {
	s := i.NewCleanChain()

	parsedAliceAddr := sdk.MustAccAddressFromBech32(i.ALICE)
	parsedBobAddr := sdk.MustAccAddressFromBech32(i.BOB)

	validatorAmount := 500 * i.KYVE
	delegatorAmount := 250 * i.KYVE

	BeforeEach(func() {
		s = i.NewCleanChain()

		// Create a test proposal.
		proposeTx := CreateTestProposal(s.Ctx(), s.App().GovKeeper)
		_ = s.RunTxSuccess(proposeTx)

		// Initialise a protocol validator.
		createTx := &stakersTypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  validatorAmount,
		}
		_ = s.RunTxSuccess(createTx)

		// Create and join a pool.
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		joinTx := &stakersTypes.MsgJoinPool{
			Creator:    i.ALICE,
			PoolId:     0,
			Valaddress: i.DUMMY[0],
			Amount:     validatorAmount,
		}
		_ = s.RunTxSuccess(joinTx)

		// Delegate to protocol validator.
		delegateTx := &delegationTypes.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.ALICE,
			Amount:  delegatorAmount,
		}
		_ = s.RunTxSuccess(delegateTx)

		Expect(s.App().StakersKeeper.TotalBondedTokens(s.Ctx()).Uint64()).To(Equal(delegatorAmount + validatorAmount))
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Protocol validator doesn't vote, delegator votes.", func() {
		// ARRANGE
		delegatorTx := govTypes.NewMsgVote(
			parsedBobAddr, 1, govTypes.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_ = s.RunTxSuccess(delegatorTx)

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		_, _, tally, _ := s.App().GovKeeper.Tally(s.Ctx(), proposal)

		Expect(tally.YesCount).To(Equal(strconv.Itoa(int(delegatorAmount))))
	})

	It("Protocol validator votes, delegator doesn't.", func() {
		// ARRANGE
		validatorTx := govTypes.NewMsgVote(
			parsedAliceAddr, 1, govTypes.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_ = s.RunTxSuccess(validatorTx)

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		_, _, tally, _ := s.App().GovKeeper.Tally(s.Ctx(), proposal)

		Expect(tally.YesCount).To(Equal(strconv.Itoa(int(delegatorAmount + validatorAmount))))
	})

	It("Protocol validator votes, delegator votes the same.", func() {
		// ARRANGE
		validatorTx := govTypes.NewMsgVote(
			parsedAliceAddr, 1, govTypes.VoteOption_VOTE_OPTION_YES, "",
		)

		delegatorTx := govTypes.NewMsgVote(
			parsedBobAddr, 1, govTypes.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_ = s.RunTxSuccess(validatorTx)
		_ = s.RunTxSuccess(delegatorTx)

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		_, _, tally, _ := s.App().GovKeeper.Tally(s.Ctx(), proposal)

		Expect(tally.YesCount).To(Equal(strconv.Itoa(int(validatorAmount + delegatorAmount))))
	})

	It("Protocol validator votes, delegator votes different.", func() {
		// ARRANGE
		validatorTx := govTypes.NewMsgVote(
			parsedAliceAddr, 1, govTypes.VoteOption_VOTE_OPTION_YES, "",
		)

		delegatorTx := govTypes.NewMsgVote(
			parsedBobAddr, 1, govTypes.VoteOption_VOTE_OPTION_NO, "",
		)

		// ACT
		_ = s.RunTxSuccess(validatorTx)
		_ = s.RunTxSuccess(delegatorTx)

		// ASSERT
		proposal, _ := s.App().GovKeeper.Proposals.Get(s.Ctx(), 1)
		_, _, tally, _ := s.App().GovKeeper.Tally(s.Ctx(), proposal)

		Expect(tally.YesCount).To(Equal(strconv.Itoa(int(validatorAmount))))
		Expect(tally.NoCount).To(Equal(strconv.Itoa(int(delegatorAmount))))
	})
})

func CreateTestProposal(ctx sdk.Context, govKeeper *keeper.Keeper) sdk.Msg {
	params, _ := govKeeper.Params.Get(ctx)

	proposal, _ := govTypes.NewMsgSubmitProposal(
		[]sdk.Msg{}, params.MinDeposit, i.DUMMY[0], "metadata", "title", "summary", false,
	)

	return proposal
}
