package keeper

import (
	"encoding/binary"
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRedelegationCooldown ...
func (k Keeper) SetRedelegationCooldown(ctx sdk.Context, redelegationCooldown types.RedelegationCooldown) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationCooldownPrefix)
	store.Set(types.RedelegationCooldownKey(
		redelegationCooldown.Address,
		redelegationCooldown.CreationDate,
	), []byte{1})
}

// GetRedelegationCooldownEntries ...
func (k Keeper) GetRedelegationCooldownEntries(ctx sdk.Context, delegatorAddress string) (creationDates []uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBuilder{Key: types.RedelegationCooldownPrefix}.AString(delegatorAddress).Key)
	iterator := sdk.KVStorePrefixIterator(store, nil)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		creationDates = append(creationDates, binary.BigEndian.Uint64(iterator.Key()[0:8]))
	}
	return
}

// RemoveRedelegationCooldown ...
func (k Keeper) RemoveRedelegationCooldown(ctx sdk.Context, delegatorAddress string, block uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationCooldownPrefix)
	store.Delete(types.RedelegationCooldownKey(delegatorAddress, block))
}

// GetAllRedelegationCooldownEntries ...
func (k Keeper) GetAllRedelegationCooldownEntries(ctx sdk.Context) (list []types.RedelegationCooldown) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationCooldownPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		val := types.RedelegationCooldown{
			Address:      string(iterator.Key()[0:43]),
			CreationDate: binary.BigEndian.Uint64(iterator.Key()[44:52]),
		}
		list = append(list, val)
	}

	return
}
