package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/team/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TeamInfo(c context.Context, req *types.QueryTeamInfoRequest) (*types.QueryTeamInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	return k.GetTeamInfo(ctx), nil
}
