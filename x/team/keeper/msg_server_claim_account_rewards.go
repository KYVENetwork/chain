package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ClaimAccountRewards(goCtx context.Context, msg *types.MsgClaimAccountRewards) (*types.MsgClaimAccountRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if types.FOUNDATION_ADDRESS != msg.Authority && types.BCP_ADDRESS != msg.Authority {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidAuthority.Error(), types.FOUNDATION_ADDRESS, types.BCP_ADDRESS, msg.Authority)
	}

	account, found := k.GetTeamVestingAccount(ctx, msg.Id)
	if !found {
		return nil, sdkErrors.ErrNotFound
	}

	// check if account has available inflation rewards which can be claimed
	if account.TotalRewards-account.RewardsClaimed < msg.Amount {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrClaimAmountTooHigh.Error(), account.TotalRewards-account.RewardsClaimed, msg.Amount)
	}

	// send inflation rewards to recipient
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Recipient, msg.Amount); err != nil {
		return nil, err
	}

	// increase claimed inflation rewards
	account.RewardsClaimed += msg.Amount
	k.SetTeamVestingAccount(ctx, account)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventClaimInflationRewards{
		Authority: msg.Authority,
		Id:        account.Id,
		Amount:    msg.Amount,
		Recipient: msg.Recipient,
	})

	return &types.MsgClaimAccountRewardsResponse{}, nil
}
