package keeper

import (
	"context"
	"fmt"

	"github.com/KYVENetwork/chain/x/funders/types"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefundPool handles the logic to defund a pool.
// If the user is a funder, it will subtract the provided amount
// and send the tokens back. If there are no more funds left, the funding will get inactive.
func (k msgServer) DefundPool(goCtx context.Context, msg *types.MsgDefundPool) (*types.MsgDefundPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Funding has to exist
	funding, found := k.GetFunding(ctx, msg.Creator, msg.PoolId)
	if !found {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrFundingDoesNotExist.Error(), msg.PoolId, msg.Creator)
	}

	if funding.Amount == 0 {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrFundingIsUsedUp.Error(), msg.PoolId, msg.Creator)
	}

	// FundingState has to exist
	fundingState, found := k.GetFundingState(ctx, msg.PoolId)
	if !found {
		util.PanicHalt(k.upgradeKeeper, ctx, fmt.Sprintf("FundingState for pool %d does not exist", msg.PoolId))
	}

	// Subtract amount from funding
	amount := funding.SubtractAmount(msg.Amount)
	if funding.Amount == 0 {
		fundingState.SetInactive(&funding)
	} else {
		// If funding is not fully revoked, check if updated funding is still compatible with params.
		if err := k.ensureParamsCompatibility(ctx, funding); err != nil {
			return nil, err
		}
	}

	// Transfer tokens from this module to sender.
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Creator, amount); err != nil {
		return nil, err
	}

	// Save funding and funding state
	k.SetFunding(ctx, &funding)
	k.SetFundingState(ctx, &fundingState)

	// Emit a defund event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
		PoolId:  msg.PoolId,
		Address: msg.Creator,
		Amount:  amount,
	})

	return &types.MsgDefundPoolResponse{}, nil
}
