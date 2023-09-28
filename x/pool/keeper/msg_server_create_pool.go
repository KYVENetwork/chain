package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

func (k msgServer) CreatePool(goCtx context.Context, req *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	if !json.Valid([]byte(req.Binaries)) {
		return nil, errors.Wrapf(errorsTypes.ErrLogic, types.ErrInvalidJson.Error(), req.Binaries)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	id := k.AppendPool(ctx, types.Pool{
		Name:                 req.Name,
		Runtime:              req.Runtime,
		Logo:                 req.Logo,
		Config:               req.Config,
		StartKey:             req.StartKey,
		UploadInterval:       req.UploadInterval,
		InflationShareWeight: req.InflationShareWeight,
		MinDelegation:        req.MinDelegation,
		MaxBundleSize:        req.MaxBundleSize,
		Protocol: &types.Protocol{
			Version:     req.Version,
			Binaries:    req.Binaries,
			LastUpgrade: uint64(ctx.BlockTime().Unix()),
		},
		UpgradePlan:              &types.UpgradePlan{},
		CurrentStorageProviderId: req.StorageProviderId,
		CurrentCompressionId:     req.CompressionId,
	})

	k.EnsurePoolAccount(ctx, id)
	k.fundersKeeper.CreateFundingState(ctx, id)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventCreatePool{
		Id:                   k.GetPoolCount(ctx) - 1,
		Name:                 req.Name,
		Runtime:              req.Runtime,
		Logo:                 req.Logo,
		Config:               req.Config,
		StartKey:             req.StartKey,
		UploadInterval:       req.UploadInterval,
		InflationShareWeight: req.InflationShareWeight,
		MinDelegation:        req.MinDelegation,
		MaxBundleSize:        req.MaxBundleSize,
		Version:              req.Version,
		Binaries:             req.Binaries,
		StorageProviderId:    req.StorageProviderId,
		CompressionId:        req.CompressionId,
	})

	return &types.MsgCreatePoolResponse{}, nil
}
