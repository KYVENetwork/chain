package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAuthority get the authority
func (k Keeper) GetAuthority(ctx sdk.Context) (authority types.Authority) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.AuthorityKey
	bz := store.Get(byteKey)

	// Authority doesn't exist: no element
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshal(bz, &authority)
	return
}

// SetAuthority set the authority
func (k Keeper) SetAuthority(ctx sdk.Context, authority types.Authority) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.AuthorityKey
	b := k.cdc.MustMarshal(&authority)
	store.Set(byteKey, b)
}

// GetTeamVestingAccountCount get the total number of team vesting accounts
func (k Keeper) GetTeamVestingAccountCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.TeamVestingAccountCountKey
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetTeamVestingAccountCount set the total number of team vesting accounts
func (k Keeper) SetTeamVestingAccountCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.TeamVestingAccountCountKey
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendTeamVestingAccount appends a team vesting account in the store with a new id and update the count
func (k Keeper) AppendTeamVestingAccount(
	ctx sdk.Context,
	tva types.TeamVestingAccount,
) uint64 {
	// Create the pool
	count := k.GetTeamVestingAccountCount(ctx)

	// Set the ID of the appended value
	tva.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TeamVestingAccountKey)
	appendedValue := k.cdc.MustMarshal(&tva)
	store.Set(types.TeamVestingAccountKeyPrefix(tva.Id), appendedValue)

	// Update team vesting account count
	k.SetTeamVestingAccountCount(ctx, count+1)

	return count
}

// GetTeamVestingAccount returns a team vesting account given its address.
func (k Keeper) GetTeamVestingAccount(ctx sdk.Context, id uint64) (tva types.TeamVestingAccount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TeamVestingAccountKey)
	b := store.Get(types.TeamVestingAccountKeyPrefix(id))

	if b == nil {
		return tva, false
	}

	k.cdc.MustUnmarshal(b, &tva)
	return tva, true
}

// GetTeamVestingAccounts returns all team vesting accounts
func (k Keeper) GetTeamVestingAccounts(ctx sdk.Context) (teamVestingAccounts []types.TeamVestingAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TeamVestingAccountKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var tva types.TeamVestingAccount
		k.cdc.MustUnmarshal(iterator.Value(), &tva)
		teamVestingAccounts = append(teamVestingAccounts, tva)
	}

	return
}

// SetTeamVestingAccount sets a specific team vesting account in the store.
func (k Keeper) SetTeamVestingAccount(ctx sdk.Context, tva types.TeamVestingAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TeamVestingAccountKey)
	b := k.cdc.MustMarshal(&tva)
	store.Set(types.TeamVestingAccountKeyPrefix(tva.Id), b)
}
