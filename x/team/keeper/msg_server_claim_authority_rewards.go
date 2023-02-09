package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ClaimAuthorityRewards(goCtx context.Context, msg *types.MsgClaimAuthorityRewards) (*types.MsgClaimAuthorityRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if types.AUTHORITY_ADDRESS != msg.Authority {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidAuthority.Error(), types.AUTHORITY_ADDRESS, msg.Authority)
	}

	authority := k.GetAuthority(ctx)

	// check if authority has enough available rewards to claim
	if authority.TotalRewards-authority.RewardsClaimed < msg.Amount {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrClaimAmountTooHigh.Error(), authority.TotalRewards-authority.RewardsClaimed, msg.Amount)
	}

	// send authority inflation rewards to recipient
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Recipient, msg.Amount); err != nil {
		return nil, err
	}

	// increase claimed inflation rewards
	authority.RewardsClaimed += msg.Amount
	k.SetAuthority(ctx, authority)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventClaimAuthorityRewards{
		Amount:    msg.Amount,
		Recipient: msg.Recipient,
	})

	return &types.MsgClaimAuthorityRewardsResponse{}, nil
}
