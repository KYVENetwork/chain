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

type Update struct {
	Name              *string
	Runtime           *string
	Logo              *string
	Config            *string
	UploadInterval    *uint64
	OperatingCost     *uint64
	MinDelegation     *uint64
	MaxBundleSize     *uint64
	StorageProviderId *uint32
	CompressionId     *uint32
}

func (k msgServer) UpdatePool(goCtx context.Context, req *types.MsgUpdatePool) (*types.MsgUpdatePoolResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, req.Id)
	if !found {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrPoolNotFound.Error(), req.Id)
	}

	var update Update
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
	if update.OperatingCost != nil {
		pool.OperatingCost = *update.OperatingCost
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

	return &types.MsgUpdatePoolResponse{}, nil
}
