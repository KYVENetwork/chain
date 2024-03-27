package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountDelegationUnbondings(goCtx context.Context, req *types.QueryAccountDelegationUnbondingsRequest) (*types.QueryAccountDelegationUnbondingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var delegationUnbondings []types.DelegationUnbonding

	store := prefix.NewStore(ctx.KVStore(k.delegationStoreKey), util.GetByteKey(delegationtypes.UndelegationQueueKeyPrefixIndex2, req.Address))
	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			index := binary.BigEndian.Uint64(key[0:8])
			unbondingEntry, _ := k.delegationKeeper.GetUndelegationQueueEntry(ctx, index)

			delegationUnbondings = append(delegationUnbondings, types.DelegationUnbonding{
				Amount:       unbondingEntry.Amount,
				CreationTime: unbondingEntry.CreationTime,
				Staker:       k.GetFullStaker(ctx, unbondingEntry.Staker),
			})
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountDelegationUnbondingsResponse{
		Unbondings: delegationUnbondings,
		Pagination: pageRes,
	}, nil
}
