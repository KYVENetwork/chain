package keeper

import (
	"context"
	"encoding/binary"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountDelegationList returns all staker with their pools the given user has delegated to.
// It calculates the current rewards.
// Pagination is supported
func (k Keeper) AccountDelegationList(goCtx context.Context, req *types.QueryAccountDelegationListRequest) (*types.QueryAccountDelegationListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var delegated []types.DelegatorResponse
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegatorPrefix := types.KeyPrefixBuilder{Key: types.DelegatorKeyPrefixIndex2}.AString(req.Address).Key
	delegatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), delegatorPrefix)

	pageRes, err := query.FilteredPaginate(delegatorStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {

		staker := string(key[9:52])
		poolId := binary.BigEndian.Uint64(key[0:8])

		if accumulate {
			var delegator, found = k.GetDelegator(ctx, poolId, staker, req.Address)
			if !found {
				k.Logger(ctx).Error("Delegator entry does not exist: {delegator: %s, staker: %s, poolId: %d}",
					req.Address, staker, poolId)
				return false, nil
			}

			pool, _ := k.GetPool(ctx, delegator.Id)

			f1 := F1Distribution{
				k:                k,
				ctx:              ctx,
				poolId:           pool.Id,
				stakerAddress:    delegator.Staker,
				delegatorAddress: delegator.Delegator,
			}

			delegationPoolData, _ := k.GetDelegationPoolData(ctx, pool.Id, delegator.Staker)

			commissionChangeEntry, foundCommissionChange := k.GetCommissionChangeQueueEntryByIndex2(ctx, staker, pool.Id)

			var pendingCommissionChange *types.PendingCommissionChange = nil
			if foundCommissionChange {
				pendingCommissionChange = &types.PendingCommissionChange{
					NewCommission: commissionChangeEntry.Commission,
					CreationDate:  commissionChangeEntry.CreationDate,
					FinishDate:    commissionChangeEntry.CreationDate + int64(k.CommissionChangeTime(ctx)),
				}
			}

			delegated = append(delegated, types.DelegatorResponse{
				Account:                 req.Address,
				Pool:                    &pool,
				CurrentReward:           f1.getCurrentReward(),
				DelegationAmount:        delegator.DelegationAmount,
				Staker:                  delegator.Staker,
				PendingCommissionChange: pendingCommissionChange,
				DelegationPoolData:      &delegationPoolData,
			})
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAccountDelegationListResponse{
		Delegations: delegated,
		Pagination:  pageRes,
	}, nil
}
