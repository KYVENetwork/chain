package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ProposalByHeight(goCtx context.Context, req *types.QueryProposalByHeightRequest) (*types.QueryProposalByHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposalPrefixBuilder := types.KeyPrefixBuilder{Key: types.ProposalKeyPrefixIndex2}.AInt(req.PoolId)
	proposalIndexStore := prefix.NewStore(ctx.KVStore(k.storeKey), proposalPrefixBuilder.Key)
	proposalIndexIterator := proposalIndexStore.ReverseIterator(nil, types.KeyPrefixBuilder{}.AInt(req.Height+1).Key)

	defer proposalIndexIterator.Close()

	if proposalIndexIterator.Valid() {

		bundleId := string(proposalIndexIterator.Value())

		proposal, found := k.GetProposal(ctx, bundleId)
		if found {
			if proposal.FromHeight <= req.Height && proposal.ToHeight > req.Height {
				return &types.QueryProposalByHeightResponse{
					Proposal: proposal,
				}, nil
			}
		}
	}

	return nil, status.Error(codes.NotFound, "no bundle found")
}
