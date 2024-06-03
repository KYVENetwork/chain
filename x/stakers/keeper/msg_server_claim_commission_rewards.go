package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
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
	if !msg.Amount.IsAllLTE(staker.CommissionRewards) {
		return nil, types.ErrNotEnoughRewards
	}

	// send commission rewards from stakers module to claimer
	recipient := sdk.MustAccAddressFromBech32(msg.Creator)
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, msg.Amount); err != nil {
		return nil, err
	}

	// calculate new commission rewards and save
	staker.CommissionRewards = staker.CommissionRewards.Sub(msg.Amount...)
	k.setStaker(ctx, staker)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventClaimCommissionRewards{
		Staker: msg.Creator,
		Amount: msg.Amount,
	})

	return &types.MsgClaimCommissionRewardsResponse{}, nil
}
