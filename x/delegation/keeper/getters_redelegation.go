package keeper

import (
	storeTypes "cosmossdk.io/store/types"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRedelegationCooldown ...
func (k Keeper) SetRedelegationCooldown(ctx sdk.Context, redelegationCooldown types.RedelegationCooldown) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.RedelegationCooldownPrefix)
	store.Set(types.RedelegationCooldownKey(
		redelegationCooldown.Address,
		redelegationCooldown.CreationDate,
	), []byte{1})
}

// GetRedelegationCooldownEntries ...
func (k Keeper) GetRedelegationCooldownEntries(ctx sdk.Context, delegatorAddress string) (creationDates []uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, append(types.RedelegationCooldownPrefix, util.GetByteKey(delegatorAddress)...))
	iterator := storeTypes.KVStorePrefixIterator(store, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		creationDates = append(creationDates, binary.BigEndian.Uint64(iterator.Key()[0:8]))
	}
	return
}

// RemoveRedelegationCooldown ...
func (k Keeper) RemoveRedelegationCooldown(ctx sdk.Context, delegatorAddress string, block uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.RedelegationCooldownPrefix)
	store.Delete(types.RedelegationCooldownKey(delegatorAddress, block))
}

// GetAllRedelegationCooldownEntries ...
func (k Keeper) GetAllRedelegationCooldownEntries(ctx sdk.Context) (list []types.RedelegationCooldown) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.RedelegationCooldownPrefix)
	iterator := storeTypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		val := types.RedelegationCooldown{
			Address:      string(iterator.Key()[0:43]),
			CreationDate: binary.BigEndian.Uint64(iterator.Key()[43:51]),
		}
		list = append(list, val)
	}

	return
}
