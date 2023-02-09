package keeper

import (
	"context"
	"time"

	"github.com/KYVENetwork/chain/x/team/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TeamVestingStatus(c context.Context, req *types.QueryTeamVestingStatusRequest) (*types.QueryTeamVestingStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	teamVesting, err := k.TeamVestingStatusByTime(ctx, &types.QueryTeamVestingStatusByTimeRequest{
		Id:   req.Id,
		Time: uint64(ctx.BlockTime().Unix()),
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryTeamVestingStatusResponse{
		RequestDate: time.Unix(ctx.BlockTime().Unix(), 0).String(),
		Plan:        teamVesting.Plan,
		Status:      teamVesting.Status,
	}, nil
}
