package keeper

import (
	"time"

	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLatestSummary(ctx sdk.Context, poolId uint64) (string, time.Time, error) {
	pool, found := k.poolKeeper.GetPool(ctx, poolId)
	if !found {
		return "", ctx.BlockTime(), poolTypes.ErrPoolNotFound
	}

	bundle, _ := k.bundlesKeeper.GetFinalizedBundle(ctx, poolId, pool.TotalBundles-1)
	return bundle.BundleSummary, time.Unix(int64(bundle.FinalizedAt.Timestamp), 0), nil
}
