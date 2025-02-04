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
// be in that pool (even with a different pool account)
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

	// Validator must exist.
	validator, validatorFound := k.GetValidator(ctx, msg.Creator)
	if !validatorFound {
		return nil, errors.Wrapf(errorsTypes.ErrNotFound, types.ErrNoStaker.Error())
	}

	// Validator must be in the active set.
	if !validator.IsBonded() {
		return nil, types.ErrValidatorNotInActiveSet
	}

	// Validators are not allowed to use their own address, to prevent
	// users from putting their validator private key on the protocol node server.
	if msg.Creator == msg.PoolAddress {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrPoolAddressSameAsStaker.Error())
	}

	// Validators are not allowed to join a pool twice.
	if _, poolAccountFound := k.GetPoolAccount(ctx, msg.Creator, msg.PoolId); poolAccountFound {
		return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrAlreadyJoinedPool.Error())
	}

	// Only join if it is possible
	if errFreeSlot := k.ensureFreeSlot(ctx, msg.PoolId, msg.Creator, msg.StakeFraction); errFreeSlot != nil {
		return nil, errFreeSlot
	}

	// Every pool address can only be used for one pool. It is not allowed
	// to use the same pool address for multiple pools. (to avoid account sequence errors,
	// when two processes try so submit transactions simultaneously)
	for _, poolAccount := range k.GetPoolAccountsFromStaker(ctx, msg.Creator) {
		if poolAccount.PoolAddress == msg.PoolAddress {
			return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.PoolAddressAlreadyUsed.Error())
		}
	}

	// It is not allowed to use the pool address of somebody else.
	for _, poolStaker := range k.GetAllStakerAddressesOfPool(ctx, msg.PoolId) {
		poolAccount, _ := k.GetPoolAccount(ctx, poolStaker, msg.PoolId)

		if poolAccount.PoolAddress == msg.PoolAddress {
			return nil, errors.Wrapf(errorsTypes.ErrInvalidRequest, types.PoolAddressAlreadyUsed.Error())
		}
	}

	k.AddPoolAccountToPool(ctx, msg.Creator, msg.PoolId, msg.PoolAddress, msg.Commission, msg.StakeFraction)

	if err := util.TransferFromAddressToAddress(k.bankKeeper, ctx, msg.Creator, msg.PoolAddress, msg.Amount); err != nil {
		return nil, err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventJoinPool{
		PoolId:        msg.PoolId,
		Staker:        msg.Creator,
		PoolAddress:   msg.PoolAddress,
		Amount:        msg.Amount,
		Commission:    msg.Commission,
		StakeFraction: msg.StakeFraction,
	})

	return &types.MsgJoinPoolResponse{}, nil
}
