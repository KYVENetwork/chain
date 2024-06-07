package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func BuildGovernanceTxs(s *i.KeeperTestSuite, msgs []sdk.Msg) (govV1Types.MsgSubmitProposal, govV1Types.MsgVote) {
	params, _ := s.App().GovKeeper.Params.Get(s.Ctx())
	minDeposit := params.MinDeposit
	delegations, _ := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
	voter := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	proposal, _ := govV1Types.NewMsgSubmitProposal(
		msgs, minDeposit, i.DUMMY[0], "", "title", "summary", false,
	)

	proposalId, _ := s.App().GovKeeper.ProposalID.Peek(s.Ctx())

	vote := govV1Types.NewMsgVote(
		voter, proposalId, govV1Types.VoteOption_VOTE_OPTION_YES, "",
	)

	return *proposal, *vote
}

func createPoolWithEmptyValues(s *i.KeeperTestSuite) {
	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	msg := &types.MsgCreatePool{
		Authority:            gov,
		UploadInterval:       60,
		MaxBundleSize:        100,
		InflationShareWeight: math.LegacyZeroDec(),
		Binaries:             "{}",
	}
	s.RunTxPoolSuccess(msg)

	poolId := s.App().PoolKeeper.GetPoolCount(s.Ctx()) - 1
	pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), poolId)
	pool.UploadInterval = 0
	pool.MaxBundleSize = 0
	pool.Protocol = &types.Protocol{}
	pool.UpgradePlan = &types.UpgradePlan{}
	s.App().PoolKeeper.SetPool(s.Ctx(), pool)
}
