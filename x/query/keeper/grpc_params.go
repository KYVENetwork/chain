package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	bp := k.bundleKeeper.GetParams(ctx)
	dp := k.delegationKeeper.GetParams(ctx)
	globalParams := k.globalKeeper.GetParams(ctx)
	govParams := k.govKeeper.GetParams(ctx)
	sp := k.stakerKeeper.GetParams(ctx)

	return &types.QueryParamsResponse{BundlesParams: &bp, DelegationParams: &dp, GlobalParams: &globalParams, GovParams: &govParams, StakersParams: &sp}, nil
}
