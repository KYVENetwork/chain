package migration

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storeTypes "cosmossdk.io/store/types"
	"embed"
	_ "embed"
	"encoding/hex"
	"fmt"
	"github.com/KYVENetwork/chain/util"
	bundleskeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"
)

var logger log.Logger

const (
	BundlesMigrationStepSizePerPool uint64 = 100
	WaitingBlockPeriod                     = 1
)

//go:embed files/*
var merkelRoots embed.FS

type BundlesMigrationEntry struct {
	merkleRoots []byte
	poolId      uint64
	maxBundleId uint64
}

// bundlesMigration includes the poolId and maxBundleId (exclusive) to determine which bundles are migrated
var bundlesMigration []BundlesMigrationEntry

func init() {
	dir, err := merkelRoots.ReadDir("files")
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		readFile, err := merkelRoots.ReadFile(fmt.Sprintf("files/%s", file.Name()))
		if err != nil {
			panic(err)
		}

		poolId, err := strconv.ParseUint(strings.ReplaceAll(file.Name(), "merkle_roots_pool_", ""), 10, 64)
		if err != nil {
			panic(err)
		}

		bundlesMigration = append(bundlesMigration, BundlesMigrationEntry{
			merkleRoots: readFile,
			poolId:      poolId,
			maxBundleId: uint64(len(readFile)) / 32,
		})
	}
}

func MigrateBundlesModule(sdkCtx sdk.Context, bundlesKeeper bundleskeeper.Keeper, upgradeHeight int64) {
	logger = sdkCtx.Logger().With("upgrade", "bundles-migration")

	if sdkCtx.BlockHeight()-upgradeHeight < WaitingBlockPeriod {
		logger.Info("sdkCtx.BlockHeight()-upgradeHeight < WaitingBlockPeriod > return")
		return
	}

	for _, bundlesMigrationEntry := range bundlesMigration {
		step := uint64(sdkCtx.BlockHeight()-upgradeHeight) - WaitingBlockPeriod
		offset := step * BundlesMigrationStepSizePerPool

		// Skip if all bundles have already been migrated
		if offset > bundlesMigrationEntry.maxBundleId+BundlesMigrationStepSizePerPool {
			continue
		}

		if err := migrateFinalizedBundles(sdkCtx, bundlesKeeper, offset, bundlesMigrationEntry); err != nil {
			// TODO: Error handling
			panic(err)
		}
	}
}

// MigrateFinalizedBundles ...
// maxBundleId -> inclusive
func migrateFinalizedBundles(ctx sdk.Context, bundlesKeeper bundleskeeper.Keeper, offset uint64, bundlesMigrationEntry BundlesMigrationEntry) error {
	// Init Bundles Store
	storeAdapter := runtime.KVStoreAdapter(bundlesKeeper.Migration_GetStoreService().OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, util.GetByteKey(types.FinalizedBundlePrefix, bundlesMigrationEntry.poolId))

	var iterator storeTypes.Iterator
	iterator = store.Iterator(util.GetByteKey(offset), util.GetByteKey(offset+BundlesMigrationStepSizePerPool))

	var migratedBundles []types.FinalizedBundle

	for ; iterator.Valid(); iterator.Next() {
		var rawFinalizedBundle types.FinalizedBundle
		if err := bundlesKeeper.Migration_GetCodec().Unmarshal(iterator.Value(), &rawFinalizedBundle); err != nil {
			return err
		}

		if rawFinalizedBundle.Id >= bundlesMigrationEntry.maxBundleId {
			return nil
		}

		merkleRoot := bundlesMigrationEntry.merkleRoots[rawFinalizedBundle.Id*32 : rawFinalizedBundle.Id*32+32]

		rawFinalizedBundle.BundleSummary = fmt.Sprintf("{\"merkle_root\":\"%v\"}", hex.EncodeToString(merkleRoot))

		migratedBundles = append(migratedBundles, rawFinalizedBundle)
	}
	iterator.Close()

	for _, migratedBundle := range migratedBundles {
		bundlesKeeper.SetFinalizedBundle(ctx, migratedBundle)
	}
	return nil
}
