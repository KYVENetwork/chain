package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SkipUploaderRole handles the logic of an SDK message that allows protocol nodes to skip an upload.
func (k msgServer) SkipUploaderRole(goCtx context.Context, msg *types.MsgSkipUploaderRole) (*types.MsgSkipUploaderRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.AssertCanPropose(ctx, msg.PoolId, msg.Staker, msg.Creator, msg.FromIndex); err != nil {
		return nil, err
	}

	bundleProposal, _ := k.GetBundleProposal(ctx, msg.PoolId)

	// reset points of uploader as node has proven to be active
	k.resetPoints(ctx, msg.Staker, msg.PoolId)

	// If previous bundle was dropped just skip uploader role
	// No previous round needs to be evaluated
	if bundleProposal.StorageId == "" {
		nextUploader := k.chooseNextUploader(ctx, msg.PoolId)

		// Register empty bundle with next uploader
		bundleProposal = types.BundleProposal{
			PoolId:       msg.PoolId,
			NextUploader: nextUploader,
			UpdatedAt:    uint64(ctx.BlockTime().Unix()),
		}
		k.SetBundleProposal(ctx, bundleProposal)

		return &types.MsgSkipUploaderRoleResponse{}, nil
	}

	// Previous round contains a bundle which needs to be validated now
	result, err := k.tallyBundleProposal(ctx, bundleProposal, msg.PoolId)
	if err != nil {
		return nil, err
	}

	// Get next uploader, except the one who skipped
	nextUploader := k.chooseNextUploader(ctx, msg.PoolId, msg.Staker)

	switch result.Status {
	case types.TallyResultValid:
		// Finalize bundle by adding it to the store
		k.finalizeCurrentBundleProposal(ctx, msg.PoolId, result.VoteDistribution, result.FundersPayout, result.InflationPayout, result.BundleReward, nextUploader)

		// Register empty bundle with next uploader
		bundleProposal = types.BundleProposal{
			PoolId:       msg.PoolId,
			NextUploader: nextUploader,
			UpdatedAt:    uint64(ctx.BlockTime().Unix()),
		}
		k.SetBundleProposal(ctx, bundleProposal)
	case types.TallyResultInvalid:
		// Drop current bundle.
		k.dropCurrentBundleProposal(ctx, msg.PoolId, result.VoteDistribution, nextUploader)
	case types.TallyResultNoQuorum:
		// Set next uploader and update the bundle proposal
		bundleProposal.NextUploader = nextUploader
		bundleProposal.UpdatedAt = uint64(ctx.BlockTime().Unix())
		k.SetBundleProposal(ctx, bundleProposal)
	}

	pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventSkippedUploaderRole{
		PoolId:           msg.PoolId,
		Id:               pool.TotalBundles,
		PreviousUploader: msg.Staker,
		NewUploader:      nextUploader,
	})

	return &types.MsgSkipUploaderRoleResponse{}, nil
}
