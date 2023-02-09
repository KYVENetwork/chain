package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ClaimUnlocked(goCtx context.Context, msg *types.MsgClaimUnlocked) (*types.MsgClaimUnlockedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if types.AUTHORITY_ADDRESS != msg.Authority {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidAuthority.Error(), types.AUTHORITY_ADDRESS, msg.Authority)
	}

	account, found := k.GetTeamVestingAccount(ctx, msg.Id)
	if !found {
		return nil, sdkErrors.ErrNotFound
	}

	// get current claimable amount
	currentProgress := GetVestingStatus(account, uint64(ctx.BlockTime().Unix()))

	// throw error if the requested claim amount is bigger than the available unlocked amount
	if msg.Amount > currentProgress.CurrentClaimableAmount {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrClaimAmountTooHigh.Error(), msg.Amount, currentProgress.CurrentClaimableAmount)
	}

	// Transfer claim amount from this module to recipient.
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Recipient, msg.Amount); err != nil {
		return nil, err
	}

	// update claimed amount of unlocked $KYVE
	account.UnlockedClaimed += msg.Amount
	account.LastClaimedTime = uint64(ctx.BlockTime().Unix())

	k.SetTeamVestingAccount(ctx, account)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventClaimedUnlocked{
		Id:        account.Id,
		Amount:    msg.Amount,
		Recipient: msg.Recipient,
	})

	return &types.MsgClaimUnlockedResponse{}, nil
}
