package keeper

import (
	"time"

	"cosmossdk.io/errors"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLatestSummary(ctx sdk.Context, poolId uint64) (string, time.Time, error) {
	pool, found := k.poolKeeper.GetPool(ctx, poolId)
	if !found {
		return "", ctx.BlockTime(), errors.Wrapf(
			errorsTypes.ErrNotFound, poolTypes.ErrPoolNotFound.Error(), poolId,
		)
	}

	bundle, _ := k.bundlesKeeper.GetFinalizedBundle(ctx, poolId, pool.TotalBundles-1)
	return bundle.BundleSummary, time.Unix(int64(bundle.FinalizedAt.Timestamp), 0), nil
}
