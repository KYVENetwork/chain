package keeper

import (
	"context"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"github.com/KYVENetwork/chain/util"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

	// ======================
	// ProtocolSelfDelegation
	// ======================

	response.ProtocolSelfDelegation = k.delegationKeeper.GetDelegationAmountOfDelegator(ctx, req.Address, req.Address)

	// ================================================
	// ProtocolDelegation + ProtocolDelegationUnbonding
	// ================================================

	// Iterate all Delegator entries
	delegatorStore := prefix.NewStore(
		ctx.KVStore(k.delegationKeeper.StoreKey()),
		util.GetByteKey(delegationtypes.DelegatorKeyPrefixIndex2, req.Address))
	delegatorIterator := sdk.KVStorePrefixIterator(delegatorStore, nil)
	defer delegatorIterator.Close()

	for ; delegatorIterator.Valid(); delegatorIterator.Next() {

		staker := string(delegatorIterator.Key()[0:43])

		response.ProtocolDelegation += k.delegationKeeper.GetDelegationAmountOfDelegator(ctx, staker, req.Address)
		response.ProtocolRewards += k.delegationKeeper.GetOutstandingRewards(ctx, staker, req.Address)
	}

	// ======================================================
	// Delegation Unbonding + ProtocolSelfDelegationUnbonding
	// ======================================================

	// Iterate all UnbondingDelegation entries to get total delegation unbonding amount
	for _, entry := range k.delegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(ctx, req.Address) {
		response.ProtocolDelegationUnbonding += entry.Amount
		if entry.Staker == req.Address {
			response.ProtocolSelfDelegationUnbonding += entry.Amount
		}
	}

	// ===============
	// ProtocolFunding
	// ===============

	// Iterate all fundings of the user to get total funding amount
	for _, funding := range k.fundersKeeper.GetFundingsOfFunder(ctx, req.Address) {
		response.ProtocolFunding += funding.Amount
	}

	return &response, nil
}
