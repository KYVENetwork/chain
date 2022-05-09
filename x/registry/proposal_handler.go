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
		Versions:       p.Versions,
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
	pool.Versions = p.Versions
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
