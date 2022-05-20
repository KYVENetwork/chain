package registry

import (
	"github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewRegistryProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CreatePoolProposal:
			return handleCreatePoolProposal(ctx, k, c)
		case *types.UpdatePoolProposal:
			return handleUpdatePoolProposal(ctx, k, c)
		case *types.PausePoolProposal:
			return handlePausePoolProposal(ctx, k, c)
		case *types.UnpausePoolProposal:
			return handleUnpausePoolProposal(ctx, k, c)
		case *types.SchedulePoolUpgradeProposal:
			return handleSchedulePoolUpgradeProposal(ctx, k, c)
		case *types.CancelPoolUpgradeProposal:
			return handleCancelPoolUpgradeProposal(ctx, k, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized registry proposal content type: %T", c)
		}
	}
}

func handleCreatePoolProposal(ctx sdk.Context, k keeper.Keeper, p *types.CreatePoolProposal) error {
	pool := types.Pool{
		Creator:        govtypes.ModuleName,
		Name:           p.Name,
		Runtime:        p.Runtime,
		Logo:           p.Logo,
		Config:         p.Config,
		HeightArchived: p.StartHeight,
		StartHeight:    p.StartHeight,
		UploadInterval: p.UploadInterval,
		OperatingCost:  p.OperatingCost,
		BundleProposal: &types.BundleProposal{
			FromHeight: p.StartHeight,
			ToHeight:   p.StartHeight,
		},
		MaxBundleSize: p.MaxBundleSize,
		Protocol: &types.Protocol{
			Version: p.Version,
			LastUpgrade: uint64(ctx.BlockTime().Unix()),
			Binaries: p.Binaries,
		},
		UpgradePlan: &types.UpgradePlan{},
	}

	k.AppendPool(ctx, pool)

	return nil
}

func handleUpdatePoolProposal(ctx sdk.Context, k keeper.Keeper, p *types.UpdatePoolProposal) error {
	pool, found := k.GetPool(ctx, p.Id)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, types.ErrPoolNotFound.Error(), p.Id)
	}

	pool.Name = p.Name
	pool.Runtime = p.Runtime
	pool.Logo = p.Logo
	pool.Config = p.Config
	pool.UploadInterval = p.UploadInterval
	pool.OperatingCost = p.OperatingCost
	pool.MaxBundleSize = p.MaxBundleSize

	k.SetPool(ctx, pool)

	return nil
}

func handlePausePoolProposal(ctx sdk.Context, k keeper.Keeper, p *types.PausePoolProposal) error {
	// Attempt to fetch the pool, throw an error if not found.
	pool, found := k.GetPool(ctx, p.Id)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, types.ErrPoolNotFound.Error(), p.Id)
	}

	// Throw an error if the pool is already paused.
	if pool.Paused {
		return sdkerrors.Wrapf(sdkerrors.ErrLogic, "Pool is already paused.")
	}

	// Pause the pool and return.
	pool.Paused = true
	k.SetPool(ctx, pool)

	return nil
}

func handleUnpausePoolProposal(ctx sdk.Context, k keeper.Keeper, p *types.UnpausePoolProposal) error {
	// Attempt to fetch the pool, throw an error if not found.
	pool, found := k.GetPool(ctx, p.Id)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrNotFound, types.ErrPoolNotFound.Error(), p.Id)
	}

	// Throw an error if the pool is already unpaused.
	if !pool.Paused {
		return sdkerrors.Wrapf(sdkerrors.ErrLogic, "Pool is already unpaused.")
	}

	// Unpause the pool and return.
	pool.Paused = false
	k.SetPool(ctx, pool)

	return nil
}

func handleSchedulePoolUpgradeProposal(ctx sdk.Context, k keeper.Keeper, p *types.SchedulePoolUpgradeProposal) error {
	// Check if upgrade version and binaries are not empty
	if p.Version == "" || p.Binaries == "" {
		return types.ErrInvalidArgs
	}

	var scheduledAt uint64

	// If upgrade time was already surpassed we upgrade immediately
	if (p.ScheduledAt < uint64(ctx.BlockTime().Unix())) {
		scheduledAt = uint64(ctx.BlockTime().Unix())
	} else {
		scheduledAt = p.ScheduledAt
	}

	// go through every pool and schedule the upgrade
	for _, pool := range k.GetAllPool(ctx) {
		// Skip if runtime does not match
		if pool.Runtime != p.Runtime {
			continue
		}

		// Skip if pool is currently upgrading
		if pool.UpgradePlan.ScheduledAt > 0 {
			continue
		}
	
		// register upgrade plan
		pool.UpgradePlan = &types.UpgradePlan{
			Version: p.Version,
			Binaries: p.Binaries,
			ScheduledAt: scheduledAt,
			Duration: p.Duration,
		}
	
		// Update the pool
		k.SetPool(ctx, pool)
	}

	return nil
}

func handleCancelPoolUpgradeProposal(ctx sdk.Context, k keeper.Keeper, p *types.CancelPoolUpgradeProposal) error {
	// go through every pool and cancel the upgrade
	for _, pool := range k.GetAllPool(ctx) {
		// Skip if runtime does not match
		if pool.Runtime != p.Runtime {
			continue
		}
		
		// Continue if there is no upgrade scheduled
		if pool.UpgradePlan.ScheduledAt == 0 {
			continue
		}

		// clear upgrade plan
		pool.UpgradePlan = &types.UpgradePlan{}

		// Update the pool
		k.SetPool(ctx, pool)
	}

	return nil
}
