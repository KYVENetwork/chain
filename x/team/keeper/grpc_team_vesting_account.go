package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/team/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TeamVestingAccounts(c context.Context, req *types.QueryTeamVestingAccountsRequest) (*types.QueryTeamVestingAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	accounts := k.GetTeamVestingAccounts(ctx)

	return &types.QueryTeamVestingAccountsResponse{Accounts: accounts}, nil
}

func (k Keeper) TeamVestingAccount(c context.Context, req *types.QueryTeamVestingAccountRequest) (*types.QueryTeamVestingAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	account, found := k.GetTeamVestingAccount(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryTeamVestingAccountResponse{Account: account}, nil
}
