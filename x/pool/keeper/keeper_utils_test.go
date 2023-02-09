package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func BuildGovernanceTxs(s *i.KeeperTestSuite, msgs []sdk.Msg) (govV1Types.MsgSubmitProposal, govV1Types.MsgVote) {
	minDeposit := s.App().GovKeeper.GetDepositParams(s.Ctx()).MinDeposit
	delegations := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
	voter := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	proposal, _ := govV1Types.NewMsgSubmitProposal(
		msgs, minDeposit, i.DUMMY[0], "",
	)

	proposalId, _ := s.App().GovKeeper.GetProposalID(s.Ctx())

	vote := govV1Types.NewMsgVote(
		voter, proposalId, govV1Types.VoteOption_VOTE_OPTION_YES, "",
	)

	return *proposal, *vote
}
