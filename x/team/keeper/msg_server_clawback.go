package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Clawback(goCtx context.Context, msg *types.MsgClawback) (*types.MsgClawbackResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if types.AUTHORITY_ADDRESS != msg.Authority {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidAuthority.Error(), types.AUTHORITY_ADDRESS, msg.Authority)
	}

	account, found := k.GetTeamVestingAccount(ctx, msg.Id)
	if !found {
		return nil, sdkErrors.ErrNotFound
	}

	// can not clawback before commencement
	if msg.Clawback > 0 && msg.Clawback < account.Commencement {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidClawbackDate.Error())
	}

	// can not clawback before last claim time because claimed $KYVE can not be returned to team module
	if msg.Clawback > 0 && msg.Clawback < account.LastClaimedTime {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidClawbackDate.Error())
	}

	account.Clawback = msg.Clawback
	k.SetTeamVestingAccount(ctx, account)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventClawback{
		Id:       account.Id,
		Clawback: msg.Clawback,
		Amount:   account.TotalAllocation - getVestingMaxAmount(account),
	})

	return &types.MsgClawbackResponse{}, nil
}
