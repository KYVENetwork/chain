package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CanVote(c context.Context, req *types.QueryCanVoteRequest) (*types.QueryCanVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := k.bundleKeeper.AssertCanVote(ctx, req.PoolId, req.Staker, req.Voter, req.StorageId); err != nil {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   err.Error(),
		}, nil
	}

	bundleProposal, _ := k.bundleKeeper.GetBundleProposal(ctx, req.PoolId)
	hasVotedAbstain := util.ContainsString(bundleProposal.VotersAbstain, req.Staker)

	if hasVotedAbstain {
		return &types.QueryCanVoteResponse{
			Possible: true,
			Reason:   "KYVE_VOTE_NO_ABSTAIN_ALLOWED",
		}, nil
	}

	return &types.QueryCanVoteResponse{
		Possible: true,
		Reason:   "",
	}, nil
}
