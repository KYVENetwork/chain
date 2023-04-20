package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/oracle/sender/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Request(c context.Context, req *types.QueryRequest) (*types.QueryRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	request, _ := k.GetRequest(ctx, req.Sequence)
	return &types.QueryRequestResponse{Request: &request}, nil
}

func (k Keeper) Response(c context.Context, req *types.QueryResponse) (*types.QueryResponseResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	response, _ := k.GetResponse(ctx, req.Sequence)
	return &types.QueryResponseResponse{Response: &response}, nil
}
