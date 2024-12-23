package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/stakers/types"
)

// UpdateStakeFraction updates the stake fraction of a validator in the specified pool.
// If the validator wants to increase their stake fraction we can do this immediately
// since there are no security risks involved there. If the validator wants
// to decrease it however we do that only after the stake fraction change time
// so validators can not decrease their stake before e.g. doing something maliciously
func (k msgServer) UpdateStakeFraction(goCtx context.Context, msg *types.MsgUpdateStakeFraction) (*types.MsgUpdateStakeFractionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolAccount, active := k.GetPoolAccount(ctx, msg.Creator, msg.PoolId)
	if !active {
		return nil, errors.Wrap(errorsTypes.ErrUnauthorized, types.ErrNoPoolAccount.Error())
	}

	// if the validator wants to decrease their stake fraction in a pool we have
	// to do that in the bonding time
	if msg.StakeFraction.LT(poolAccount.StakeFraction) {
		// Insert commission change into queue
		k.orderNewStakeFractionChange(ctx, msg.Creator, msg.PoolId, msg.StakeFraction)
		return &types.MsgUpdateStakeFractionResponse{}, nil
	}

	// if the validator wants to increase their stake fraction we can do this immediately.
	// Before we clear any change entries if there are currently bonding
	queueEntry, found := k.GetStakeFractionChangeEntryByIndex2(ctx, msg.Creator, msg.PoolId)
	if found {
		k.RemoveStakeFractionEntry(ctx, &queueEntry)
	}

	poolAccount.StakeFraction = msg.StakeFraction
	k.SetPoolAccount(ctx, poolAccount)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventUpdateStakeFraction{
		Staker:        msg.Creator,
		PoolId:        msg.PoolId,
		StakeFraction: msg.StakeFraction,
	})

	return &types.MsgUpdateStakeFractionResponse{}, nil
}
