package keeper

import (
	"context"
	"time"

	"github.com/KYVENetwork/chain/x/team/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TeamVestingStatusByTime(c context.Context, req *types.QueryTeamVestingStatusByTimeRequest) (*types.QueryTeamVestingStatusByTimeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	account, found := k.GetTeamVestingAccount(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	vestingStatus := GetVestingStatus(account, req.Time)

	queryVestingStatus := types.QueryVestingStatus{
		TotalVestedAmount:       vestingStatus.TotalVestedAmount,
		TotalUnlockedAmount:     vestingStatus.TotalUnlockedAmount,
		CurrentClaimableAmount:  vestingStatus.CurrentClaimableAmount,
		LockedVestedAmount:      vestingStatus.LockedVestedAmount,
		RemainingUnvestedAmount: vestingStatus.RemainingUnvestedAmount,
		ClaimedAmount:           account.UnlockedClaimed,
		TotalRewards:            account.TotalRewards,
		ClaimedRewards:          account.RewardsClaimed,
		AvailableRewards:        account.TotalRewards - account.RewardsClaimed,
	}

	vestingPlan := GetVestingPlan(account)

	queryVestingPlan := types.QueryVestingPlan{
		Commencement:         time.Unix(int64(account.Commencement), 0).String(),
		TokenVestingStart:    time.Unix(int64(vestingPlan.TokenVestingStart), 0).String(),
		TokenVestingFinished: time.Unix(int64(vestingPlan.TokenVestingFinished), 0).String(),
		TokenUnlockStart:     time.Unix(int64(vestingPlan.TokenUnlockStart), 0).String(),
		TokenUnlockFinished:  time.Unix(int64(vestingPlan.TokenUnlockFinished), 0).String(),
		Clawback:             account.Clawback,
		ClawbackAmount:       vestingPlan.ClawbackAmount,
		MaximumVestingAmount: vestingPlan.MaximumVestingAmount,
	}

	return &types.QueryTeamVestingStatusByTimeResponse{
		RequestDate: time.Unix(ctx.BlockTime().Unix(), 0).String(),
		Plan:        &queryVestingPlan,
		Status:      &queryVestingStatus,
	}, nil
}
