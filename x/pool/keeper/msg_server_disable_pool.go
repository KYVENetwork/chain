package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

func (k msgServer) DisablePool(
	goCtx context.Context,
	req *types.MsgDisablePool,
) (*types.MsgDisablePoolResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, req.Id)

	if !found {
		return nil, errors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.Id)
	}

	if pool.Disabled {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, "Pool is already disabled.")
	}

	pool.Disabled = true
	k.SetPool(ctx, pool)

	// remove all stakers from pool in order to "reset" it
	poolMembers := k.stakersKeeper.GetAllStakerAddressesOfPool(ctx, pool.Id)
	for _, staker := range poolMembers {
		k.stakersKeeper.LeavePool(ctx, staker, pool.Id)
	}

	if err := k.fundersKeeper.DefundFundingState(ctx, pool.Id); err != nil {
		return nil, err
	}

	// send remaining pool assets to treasury
	if balance := k.bankKeeper.GetBalance(ctx, pool.GetPoolAccount(), globalTypes.Denom).Amount.Uint64(); balance > 0 {
		if err := util.TransferFromAddressToTreasury(k.distrkeeper, ctx, pool.GetPoolAccount().String(), balance); err != nil {
			return nil, err
		}
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventPoolDisabled{Id: req.Id})

	return &types.MsgDisablePoolResponse{}, nil
}
