package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govTypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	oldParams := k.GetParams(ctx)

	newParams := oldParams
	_ = json.Unmarshal([]byte(msg.Payload), &newParams)
	k.SetParams(ctx, newParams)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventUpdateParams{
		OldParams: oldParams,
		NewParams: newParams,
		Payload:   msg.Payload,
	})

	return &types.MsgUpdateParamsResponse{}, nil
}
