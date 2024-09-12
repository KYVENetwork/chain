package keeper

import (
	"encoding/binary"
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetBundlesMigrationUpgradeHeight(ctx sdk.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	return binary.BigEndian.Uint64(storeAdapter.Get(types.BundlesMigrationHeightKey))
}

// SetBundlesMigrationUpgradeHeight stores the upgrade height of the v1.6 bundles migration
// upgrade in the KV-Store.
func (k Keeper) SetBundlesMigrationUpgradeHeight(ctx sdk.Context, upgradeHeight uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, upgradeHeight)

	storeAdapter.Set(types.BundlesMigrationHeightKey, bz)
}
