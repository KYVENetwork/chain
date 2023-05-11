package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/util"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ClaimCommissionRewards ...
func (k msgServer) ClaimCommissionRewards(goCtx context.Context, msg *types.MsgClaimCommissionRewards) (*types.MsgClaimCommissionRewardsResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if staker exists
	staker, found := k.GetStaker(ctx, msg.Creator)
	if !found {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrNoStaker.Error(), msg.Creator)
	}

	// Check if amount can be claimed
	if staker.CommissionRewards < msg.Amount {
		return nil, types.ErrNotEnoughRewards.Wrapf("%d > %d", msg.Amount, staker.CommissionRewards)
	}

	// send commission rewards from stakers module to claimer
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, stakertypes.ModuleName, msg.Creator, msg.Amount); err != nil {
		return nil, err
	}

	// calculate new commission rewards and save
	staker.CommissionRewards -= msg.Amount
	k.setStaker(ctx, staker)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventClaimCommissionRewards{
		Staker: msg.Creator,
		Amount: msg.Amount,
	})

	return &types.MsgClaimCommissionRewardsResponse{}, nil
}
