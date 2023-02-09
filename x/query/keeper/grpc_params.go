package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	// Query
	"github.com/KYVENetwork/chain/x/query/types"
)

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	bp := k.bundleKeeper.GetParams(ctx)
	dp := k.delegationKeeper.GetParams(ctx)
	globalParams := k.globalKeeper.GetParams(ctx)
	govParams := govTypes.QueryParamsResponse{}
	sp := k.stakerKeeper.GetParams(ctx)

	govVotingParams := k.govKeeper.GetVotingParams(ctx)
	govParams.VotingParams = &govVotingParams
	govDepositParams := k.govKeeper.GetDepositParams(ctx)
	govParams.DepositParams = &govDepositParams
	govTallyParams := k.govKeeper.GetTallyParams(ctx)
	govParams.TallyParams = &govTallyParams

	return &types.QueryParamsResponse{BundlesParams: &bp, DelegationParams: &dp, GlobalParams: &globalParams, GovParams: &govParams, StakersParams: &sp}, nil
}
