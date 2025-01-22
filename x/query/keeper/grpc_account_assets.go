package keeper

import (
	"context"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountAssets returns an overview of the balances of the given user regarding the protocol nodes
// This includes the current balance, funding, staking, and delegation.
func (k Keeper) AccountAssets(goCtx context.Context, req *types.QueryAccountAssetsRequest) (*types.QueryAccountAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryAccountAssetsResponse{}

	// =======
	// Balance
	// =======
	account, _ := sdk.AccAddressFromBech32(req.Address)
	balance := k.bankKeeper.GetBalance(ctx, account, globalTypes.Denom)
	response.Balance = balance.Amount.Uint64()

	// ================================================
	// OutstandingRewards + Delegation + Unbonding
	// ================================================

	delegatorAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	validators, err := k.stakingKeeper.GetDelegatorValidators(ctx, delegatorAddr, 1000)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	for _, validator := range validators.Validators {
		response.DelegationRewards = response.DelegationRewards.Add(
			k.stakerKeeper.GetOutstandingRewards(ctx, util.MustAccountAddressFromValAddress(validator.OperatorAddress), req.Address)...,
		)
	}

	response.CommissionRewards = k.stakerKeeper.GetOutstandingCommissionRewards(ctx, util.MustAccountAddressFromValAddress(util.MustValaddressFromOperatorAddress(req.Address)))

	delegatorBonded, err := k.stakingKeeper.GetDelegatorBonded(ctx, delegatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	response.Delegation = delegatorBonded.Uint64()

	delegatorUnbonding, err := k.stakingKeeper.GetDelegatorUnbonding(ctx, delegatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	response.DelegationUnbonding = delegatorUnbonding.Uint64()

	// ===============
	// ProtocolFunding
	// ===============

	// Iterate all fundings of the user to get total funding amount
	for _, funding := range k.fundersKeeper.GetFundingsOfFunder(ctx, req.Address) {
		response.ProtocolFunding = response.ProtocolFunding.Add(funding.Amounts...)
	}

	return &response, nil
}
