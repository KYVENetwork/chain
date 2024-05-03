package keeper

import (
	"context"

	sdkErrors "cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WithdrawRewards calculates the current rewards of a delegator and transfers the balance to
// the delegator's wallet. Only the delegator himself can call this transaction.
func (k msgServer) WithdrawRewards(
	goCtx context.Context,
	msg *types.MsgWithdrawRewards,
) (*types.MsgWithdrawRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the sender has delegated to the given staker
	if !k.DoesDelegatorExist(ctx, msg.Staker, msg.Creator) {
		return nil, sdkErrors.Wrapf(types.ErrNotADelegator, "%s does not delegate to %s", msg.Creator, msg.Staker)
	}

	// Withdraw all rewards of the sender and send them back
	_, err := k.performWithdrawal(ctx, msg.Staker, msg.Creator)
	if err != nil {
		return nil, err
	}

	return &types.MsgWithdrawRewardsResponse{}, nil
}
