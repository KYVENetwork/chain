package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

func (k msgServer) ScheduleRuntimeUpgrade(goCtx context.Context, req *types.MsgScheduleRuntimeUpgrade) (*types.MsgScheduleRuntimeUpgradeResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	if req.Version == "" || req.Binaries == "" {
		return nil, types.ErrInvalidArgs
	}

	var scheduledAt uint64
	ctx := sdk.UnwrapSDKContext(goCtx)

	// if upgrade was scheduled in the past we reschedule it to now
	if req.ScheduledAt < uint64(ctx.BlockTime().Unix()) {
		scheduledAt = uint64(ctx.BlockTime().Unix())
	} else {
		scheduledAt = req.ScheduledAt
	}

	for _, pool := range k.GetAllPools(ctx) {
		// only schedule upgrade if the runtime matches
		if pool.Runtime != req.Runtime {
			continue
		}

		// only schedule upgrade if there is no upgrade already
		if pool.UpgradePlan.ScheduledAt != 0 {
			continue
		}

		pool.UpgradePlan = &types.UpgradePlan{
			Version:     req.Version,
			Binaries:    req.Binaries,
			ScheduledAt: scheduledAt,
			Duration:    req.Duration,
		}

		k.SetPool(ctx, pool)
	}

	return &types.MsgScheduleRuntimeUpgradeResponse{}, nil
}
