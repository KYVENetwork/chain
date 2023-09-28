package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KYVENetwork/chain/x/pool/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) UpdatePool(goCtx context.Context, req *types.MsgUpdatePool) (*types.MsgUpdatePoolResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, req.Id)
	if !found {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrPoolNotFound.Error(), req.Id)
	}

	var update types.PoolUpdate
	if err := json.Unmarshal([]byte(req.Payload), &update); err != nil {
		return nil, err
	}

	if update.Name != nil {
		pool.Name = *update.Name
	}
	if update.Runtime != nil {
		pool.Runtime = *update.Runtime
	}
	if update.Logo != nil {
		pool.Logo = *update.Logo
	}
	if update.Config != nil {
		pool.Config = *update.Config
	}
	if update.UploadInterval != nil {
		pool.UploadInterval = *update.UploadInterval
	}
	if update.InflationShareWeight != nil {
		pool.InflationShareWeight = *update.InflationShareWeight
	}
	if update.MinDelegation != nil {
		pool.MinDelegation = *update.MinDelegation
	}
	if update.MaxBundleSize != nil {
		pool.MaxBundleSize = *update.MaxBundleSize
	}
	if update.StorageProviderId != nil {
		pool.CurrentStorageProviderId = *update.StorageProviderId
	}
	if update.CompressionId != nil {
		pool.CurrentCompressionId = *update.CompressionId
	}

	k.SetPool(ctx, pool)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventPoolUpdated{
		Id:                   pool.Id,
		RawUpdateString:      req.Payload,
		Name:                 pool.Name,
		Runtime:              pool.Runtime,
		Logo:                 pool.Logo,
		Config:               pool.Config,
		UploadInterval:       pool.UploadInterval,
		InflationShareWeight: pool.InflationShareWeight,
		MinDelegation:        pool.MinDelegation,
		MaxBundleSize:        pool.MaxBundleSize,
		StorageProviderId:    pool.CurrentStorageProviderId,
		CompressionId:        pool.CurrentCompressionId,
	})

	return &types.MsgUpdatePoolResponse{}, nil
}
