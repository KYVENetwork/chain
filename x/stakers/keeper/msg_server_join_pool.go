package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// JoinPool handles the SDK message of joining a pool.
// For joining a pool the staker needs to exist and must not
// be in that pool (even with a different valaccount)
// Second, there must be free slots available or the staker
// must have more stake than the lowest staker in that pool.
// After the staker joined the pool he is subject to slashing.
// The protocol node should be configured and running before
// submitting this transaction
func (k msgServer) JoinPool(goCtx context.Context, msg *types.MsgJoinPool) (*types.MsgJoinPoolResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, poolErr := k.poolKeeper.GetPoolWithError(ctx, msg.PoolId)
	if poolErr != nil {
		return nil, poolErr
	}
	if pool.Disabled {
		return nil, errors.Wrapf(errorsTypes.ErrLogic, types.ErrCanNotJoinDisabledPool.Error())
	}

	// throw error if staker was not found
	staker, stakerFound := k.GetStaker(ctx, msg.Creator)
	if !stakerFound {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrNoStaker.Error())
	}

	// Stakers are not allowed to use their own address, to prevent
	// users from putting their staker private key on the protocol node server.
	if msg.Creator == msg.Valaddress {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrValaddressSameAsStaker.Error())
	}

	// Stakers are not allowed to join a pool twice.
	if _, valaccountFound := k.GetValaccount(ctx, msg.PoolId, msg.Creator); valaccountFound {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrAlreadyJoinedPool.Error())
	}

	// Only join if it is possible
	if errFreeSlot := k.ensureFreeSlot(ctx, msg.PoolId, staker.Address); errFreeSlot != nil {
		return nil, errFreeSlot
	}

	// Every valaddress can only be used for one pool. It is not allowed
	// to use the same valaddress for multiple pools. (to avoid account sequence errors,
	// when two processes try so submit transactions simultaneously)
	for _, valaccount := range k.GetValaccountsFromStaker(ctx, msg.Creator) {
		if valaccount.Valaddress == msg.Valaddress {
			return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ValaddressAlreadyUsed.Error())
		}
	}

	// It is not allowed to use the valaddress of somebody else.
	for _, poolStaker := range k.GetAllStakerAddressesOfPool(ctx, msg.PoolId) {
		valaccount, _ := k.GetValaccount(ctx, msg.PoolId, poolStaker)

		if valaccount.Valaddress == msg.Valaddress {
			return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ValaddressAlreadyUsed.Error())
		}
	}

	k.AddValaccountToPool(ctx, msg.PoolId, msg.Creator, msg.Valaddress)

	if err := util.TransferFromAddressToAddress(k.bankKeeper, ctx, msg.Creator, msg.Valaddress, msg.Amount); err != nil {
		return nil, err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventJoinPool{
		PoolId:     msg.PoolId,
		Staker:     msg.Creator,
		Valaddress: msg.Valaddress,
		Amount:     msg.Amount,
	})

	return &types.MsgJoinPoolResponse{}, nil
}
