package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SkipUploaderRole handles the logic of an SDK message that allows protocol nodes to skip an upload.
func (k msgServer) SkipUploaderRole(
	goCtx context.Context, msg *types.MsgSkipUploaderRole,
) (*types.MsgSkipUploaderRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.AssertCanPropose(ctx, msg.PoolId, msg.Staker, msg.Creator, msg.FromIndex); err != nil {
		return nil, err
	}

	pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, msg.PoolId)

	// reset points of uploader as node has proven to be active
	k.resetPoints(ctx, msg.PoolId, msg.Staker)

	// Get next uploader from stakers voted
	stakers := make([]string, 0)
	nextUploader := ""

	// exclude the staker who skips the uploader role
	for _, staker := range k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, msg.PoolId) {
		if staker != msg.Staker {
			stakers = append(stakers, staker)
		}
	}

	if len(stakers) > 0 {
		nextUploader = k.chooseNextUploaderFromSelectedStakers(ctx, msg.PoolId, stakers)
	} else {
		nextUploader = k.chooseNextUploaderFromAllStakers(ctx, msg.PoolId)
	}

	bundleProposal.NextUploader = nextUploader
	bundleProposal.UpdatedAt = uint64(ctx.BlockTime().Unix())

	k.SetBundleProposal(ctx, bundleProposal)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventSkippedUploaderRole{
		PoolId:           msg.PoolId,
		Id:               pool.TotalBundles,
		PreviousUploader: msg.Staker,
		NewUploader:      nextUploader,
	})

	return &types.MsgSkipUploaderRoleResponse{}, nil
}
