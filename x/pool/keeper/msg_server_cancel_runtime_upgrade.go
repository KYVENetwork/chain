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

func (k msgServer) CancelRuntimeUpgrade(
	goCtx context.Context, req *types.MsgCancelRuntimeUpgrade,
) (*types.MsgCancelRuntimeUpgradeResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	affectedPools := make([]uint64, 0)
	for _, pool := range k.GetAllPools(ctx) {
		if pool.Runtime != req.Runtime {
			continue
		}
		if pool.UpgradePlan.ScheduledAt == 0 {
			continue
		}

		affectedPools = append(affectedPools, pool.Id)

		pool.UpgradePlan = &types.UpgradePlan{}
		k.SetPool(ctx, pool)
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventRuntimeUpgradeCancelled{
		Runtime:       req.Runtime,
		AffectedPools: affectedPools,
	})

	return &types.MsgCancelRuntimeUpgradeResponse{}, nil
}
