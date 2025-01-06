package keeper

import (
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) migration_RemoveBranch(ctx sdk.Context, keyPrefix []byte) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, keyPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	keys := make([][]byte, 0)
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
	}

	for _, key := range keys {
		store.Delete(key)
	}
}

func (k Keeper) Migration_ResetOldState(ctx sdk.Context) {
	k.migration_RemoveBranch(ctx, types.StakerKeyPrefix)

	k.migration_RemoveBranch(ctx, types.PoolAccountPrefix)
	k.migration_RemoveBranch(ctx, types.PoolAccountPrefixIndex2)

	k.migration_RemoveBranch(ctx, types.CommissionChangeEntryKeyPrefix)
	k.migration_RemoveBranch(ctx, types.CommissionChangeEntryKeyPrefixIndex2)

	k.migration_RemoveBranch(ctx, types.LeavePoolEntryKeyPrefix)
	k.migration_RemoveBranch(ctx, types.LeavePoolEntryKeyPrefixIndex2)

	k.migration_RemoveBranch(ctx, types.ActiveStakerIndex)

	k.SetQueueState(ctx, types.QUEUE_IDENTIFIER_COMMISSION, types.QueueState{
		LowIndex:  0,
		HighIndex: 0,
	})

	k.SetQueueState(ctx, types.QUEUE_IDENTIFIER_LEAVE, types.QueueState{
		LowIndex:  0,
		HighIndex: 0,
	})

	k.SetQueueState(ctx, types.QUEUE_IDENTIFIER_STAKE_FRACTION, types.QueueState{
		LowIndex:  0,
		HighIndex: 0,
	})
}
