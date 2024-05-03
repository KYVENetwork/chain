package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakersByDelegator(goCtx context.Context, req *types.QueryStakersByDelegatorRequest) (*types.QueryStakersByDelegatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var stakers []types.DelegationForStakerResponse

	storeAdapter := runtime.KVStoreAdapter(k.delegationStoreService.OpenKVStore(ctx))
	delegatorStore := prefix.NewStore(storeAdapter, util.GetByteKey(delegationtypes.DelegatorKeyPrefixIndex2, req.Delegator))

	pageRes, err := query.FilteredPaginate(delegatorStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			staker := string(key[0:43])

			stakers = append(stakers, types.DelegationForStakerResponse{
				Staker:           k.GetFullStaker(ctx, staker),
				CurrentRewards:   k.delegationKeeper.GetOutstandingRewards(ctx, staker, req.Delegator),
				DelegationAmount: k.delegationKeeper.GetDelegationAmountOfDelegator(ctx, staker, req.Delegator),
			})
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryStakersByDelegatorResponse{
		Delegator:  req.Delegator,
		Stakers:    stakers,
		Pagination: pageRes,
	}, nil
}
