package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateTeamVestingAccount(goCtx context.Context, msg *types.MsgCreateTeamVestingAccount) (*types.MsgCreateTeamVestingAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if types.AUTHORITY_ADDRESS != msg.Authority {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidAuthority.Error(), types.AUTHORITY_ADDRESS, msg.Authority)
	}

	if msg.TotalAllocation == 0 || msg.Commencement == 0 {
		return nil, errors.Wrapf(sdkErrors.ErrLogic, "total allocation %v or commencement %v invalid", msg.TotalAllocation, msg.Commencement)
	}

	// check if new team vesting account still has allocation left
	if k.GetIssuedTeamAllocation(ctx)+msg.TotalAllocation > types.TEAM_ALLOCATION {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrAvailableFundsTooLow.Error(), types.TEAM_ALLOCATION-k.GetIssuedTeamAllocation(ctx), msg.TotalAllocation)
	}

	id := k.AppendTeamVestingAccount(ctx, types.TeamVestingAccount{
		TotalAllocation: msg.TotalAllocation,
		Commencement:    msg.Commencement,
	})

	_ = ctx.EventManager().EmitTypedEvent(&types.EventCreateTeamVestingAccount{
		Id:              id,
		TotalAllocation: msg.TotalAllocation,
		Commencement:    msg.Commencement,
	})

	return &types.MsgCreateTeamVestingAccountResponse{}, nil
}
