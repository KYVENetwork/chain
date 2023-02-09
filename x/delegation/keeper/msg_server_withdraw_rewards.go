package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"

	sdkErrors "cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WithdrawRewards calculates the current rewards of a delegator and transfers the balance to
// the delegator's wallet. Only the delegator himself can call this transaction.
func (k msgServer) WithdrawRewards(goCtx context.Context, msg *types.MsgWithdrawRewards) (*types.MsgWithdrawRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the sender has delegated to the given staker
	if !k.DoesDelegatorExist(ctx, msg.Staker, msg.Creator) {
		return nil, sdkErrors.WithType(types.ErrNotADelegator, msg.Creator)
	}

	// Withdraw all rewards of the sender.
	reward := k.f1WithdrawRewards(ctx, msg.Staker, msg.Creator)

	// Transfer reward $KYVE from this module to sender.
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Creator, reward); err != nil {
		return nil, err
	}

	// Emit a delegation event.
	_ = ctx.EventManager().EmitTypedEvent(&types.EventWithdrawRewards{
		Address: msg.Creator,
		Staker:  msg.Staker,
		Amount:  reward,
	})

	return &types.MsgWithdrawRewardsResponse{}, nil
}
