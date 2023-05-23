package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateStaker handles the logic of an SDK message that allows protocol nodes to create
// a staker with an initial self delegation.
// Every user can create a staker object with some stake. However,
// only if self_delegation + delegation is large enough to join a pool the staker
// is able to participate in the protocol
func (k msgServer) CreateStaker(
	goCtx context.Context,
	msg *types.MsgCreateStaker,
) (*types.MsgCreateStakerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Only create new stakers
	if k.DoesStakerExist(ctx, msg.Creator) {
		return nil, types.ErrStakerAlreadyCreated
	}

	// Create and append new staker to store
	k.AppendStaker(ctx, types.Staker{
		Address:    msg.Creator,
		Commission: msg.Commission,
	})

	// Perform initial self delegation
	// TODO(@john)
	//delegationMsgServer := delegationKeeper.NewMsgServerImpl(k.delegationKeeper)
	//if _, err := delegationMsgServer.Delegate(ctx, &delegationTypes.MsgDelegate{
	//	Creator: msg.Creator,
	//	Staker:  msg.Creator,
	//	Amount:  msg.Amount,
	//}); err != nil {
	//	return nil, err
	//}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventCreateStaker{
		Staker:     msg.Creator,
		Amount:     msg.Amount,
		Commission: msg.Commission,
	})

	return &types.MsgCreateStakerResponse{}, nil
}
