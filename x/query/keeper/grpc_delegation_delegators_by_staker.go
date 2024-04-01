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

func (k Keeper) DelegatorsByStaker(goCtx context.Context, req *types.QueryDelegatorsByStakerRequest) (*types.QueryDelegatorsByStakerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var delegators []types.StakerDelegatorResponse

	storeAdapter := runtime.KVStoreAdapter(k.delegationStoreService.OpenKVStore(ctx))
	delegatorStore := prefix.NewStore(storeAdapter, util.GetByteKey(delegationtypes.DelegatorKeyPrefix, req.Staker))

	pageRes, err := query.FilteredPaginate(delegatorStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var delegator delegationtypes.Delegator

			if err := k.cdc.Unmarshal(value, &delegator); err != nil {
				return false, nil
			}

			delegators = append(delegators, types.StakerDelegatorResponse{
				Delegator:        delegator.Delegator,
				CurrentReward:    k.delegationKeeper.GetOutstandingRewards(ctx, req.Staker, delegator.Delegator),
				DelegationAmount: k.delegationKeeper.GetDelegationAmountOfDelegator(ctx, req.Staker, delegator.Delegator),
				Staker:           req.Staker,
			})
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationData, _ := k.delegationKeeper.GetDelegationData(ctx, req.Staker)

	return &types.QueryDelegatorsByStakerResponse{
		Delegators:          delegators,
		TotalDelegation:     delegationData.TotalDelegation,
		TotalDelegatorCount: delegationData.DelegatorCount,
		Pagination:          pageRes,
	}, nil
}
