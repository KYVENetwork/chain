package migration

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	_ "embed"
	"encoding/hex"
	"fmt"
	"github.com/KYVENetwork/chain/util"
	bundleskeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var logger log.Logger

const (
	BundlesMigrationStepSizePerPool uint64 = 100
	WaitingBlockPeriod                     = 1
)

//go:embed files/merkle_roots_pool_1
var merkleRootsPool1 []byte

//go:embed files/merkle_roots_pool_2
var merkleRootsPool2 []byte

type BundlesMigrationEntry struct {
	merkleRoots []byte
	poolId      uint64
}

// bundlesMigration includes the poolId and max bundleId to determine which bundles are migrated
var bundlesMigration = []BundlesMigrationEntry{
	{
		merkleRoots: merkleRootsPool1,
		poolId:      1,
	},
	{
		merkleRoots: merkleRootsPool2,
		poolId:      2,
	},
}

func MigrateBundlesModule(sdkCtx sdk.Context, bundlesKeeper bundleskeeper.Keeper, upgradeHeight int64) {
	if sdkCtx.BlockHeight()-upgradeHeight < WaitingBlockPeriod {
		logger.Info("sdkCtx.BlockHeight()-upgradeHeight < WaitingBlockPeriod > return")
		return
	}

	for _, bundlesMigrationEntry := range bundlesMigration {
		step := uint64(sdkCtx.BlockHeight()-upgradeHeight) - WaitingBlockPeriod
		offset := step * BundlesMigrationStepSizePerPool

		var maxBundleId uint64
		switch bundlesMigrationEntry.poolId {
		case 1:
			maxBundleId = uint64(len(merkleRootsPool1)) / 32
		case 2:
			maxBundleId = uint64(len(merkleRootsPool2)) / 32
		}

		// Exit if all bundles have already been migrated
		if offset > maxBundleId {
			logger.Info("offset > maxBundleId > return")
			return
		}

		if err := migrateFinalizedBundles(sdkCtx, bundlesKeeper, offset, bundlesMigrationEntry.poolId, maxBundleId); err != nil {
			// TODO: Error handling
			panic(err)
		}
	}
}

// MigrateFinalizedBundles ...
// maxBundleId -> inclusive
func migrateFinalizedBundles(ctx sdk.Context, bundlesKeeper bundleskeeper.Keeper, offset uint64, poolId uint64, maxBundleId uint64) error {
	// Init Bundles Store
	storeAdapter := runtime.KVStoreAdapter(bundlesKeeper.Migration_GetStoreService().OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, util.GetByteKey(types.FinalizedBundlePrefix, poolId))

	var iterator storeTypes.Iterator
	iterator = store.Iterator(util.GetByteKey(offset), util.GetByteKey(offset+BundlesMigrationStepSizePerPool))

	var migratedBundles []types.FinalizedBundle

	for ; iterator.Valid(); iterator.Next() {
		var rawFinalizedBundle types.FinalizedBundle
		if err := bundlesKeeper.Migration_GetCodec().Unmarshal(iterator.Value(), &rawFinalizedBundle); err != nil {
			return err
		}

		if rawFinalizedBundle.Id >= maxBundleId {
			return nil
		}

		var merkleRoot []byte
		switch rawFinalizedBundle.PoolId {
		case 1:
			merkleRoot = merkleRootsPool1[rawFinalizedBundle.Id*32 : rawFinalizedBundle.Id*32+32]
		case 2:
			merkleRoot = merkleRootsPool2[rawFinalizedBundle.Id*32 : rawFinalizedBundle.Id*32+32]
		}
		rawFinalizedBundle.BundleSummary = fmt.Sprintf("{\"merkle_root\":\"%v\"}", hex.EncodeToString(merkleRoot))

		migratedBundles = append(migratedBundles, rawFinalizedBundle)
	}
	iterator.Close()

	for _, migratedBundle := range migratedBundles {
		bundlesKeeper.SetFinalizedBundle(ctx, migratedBundle)
	}
	return nil
}
