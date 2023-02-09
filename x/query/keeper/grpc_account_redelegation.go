package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountRedelegation(goCtx context.Context, req *types.QueryAccountRedelegationRequest) (*types.QueryAccountRedelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var redelegationEntries []types.RedelegationEntry
	usedSlots := uint64(0)

	for _, creationDate := range k.delegationKeeper.GetRedelegationCooldownEntries(ctx, req.Address) {

		finishDate := creationDate + k.delegationKeeper.GetRedelegationCooldown(ctx)

		if finishDate >= uint64(ctx.BlockTime().Unix()) {
			redelegationEntries = append(redelegationEntries, types.RedelegationEntry{
				CreationDate: creationDate,
				FinishDate:   finishDate,
			})
			usedSlots += 1
		}
	}

	return &types.QueryAccountRedelegationResponse{
		RedelegationCooldownEntries: redelegationEntries,
		AvailableSlots:              k.delegationKeeper.GetRedelegationMaxAmount(ctx) - usedSlots,
	}, nil
}
