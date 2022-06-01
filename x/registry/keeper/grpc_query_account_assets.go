package keeper

import (
	"context"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountAssets returns an overview of the balances of the given user regarding the protocol nodes
// This includes the current balance, funding, staking, and delegation.
// Supports Pagination
func (k Keeper) AccountAssets(goCtx context.Context, req *types.QueryAccountAssetsRequest) (*types.QueryAccountAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	response := types.QueryAccountAssetsResponse{
		Balance:            0,
		ProtocolStaking:    0,
		ProtocolDelegation: 0,
		ProtocolRewards:    0,
		ProtocolFunding:    0,
	}

	// Fetch account balance
	account, _ := sdk.AccAddressFromBech32(req.Address)
	balance := k.bankKeeper.GetBalance(ctx, account, "tkyve")
	response.Balance = balance.Amount.Uint64()

	// Iterate all Delegator entries
	// Fetches the total delegation and calculates the outstanding rewards
	delegatorPrefix := types.KeyPrefixBuilder{Key: types.DelegatorKeyPrefixIndex2}.AString(req.Address).Key
	delegatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), delegatorPrefix)
	delegatorIterator := sdk.KVStorePrefixIterator(delegatorStore, nil)

	defer delegatorIterator.Close()

	for ; delegatorIterator.Valid(); delegatorIterator.Next() {

		key := delegatorIterator.Key()
		staker := string(key[9:52])
		poolId := binary.BigEndian.Uint64(key[0:8])
		var delegator, found = k.GetDelegator(ctx, poolId, staker, req.Address)
		if !found {
			k.Logger(ctx).Error("Delegator entry does not exist: {delegator: %s, staker: %s, poolId: %d}",
				req.Address, staker, poolId)
			continue
		}

		f1 := F1Distribution{
			k:                k,
			ctx:              ctx,
			poolId:           delegator.Id,
			stakerAddress:    delegator.Staker,
			delegatorAddress: delegator.Delegator,
		}

		response.ProtocolRewards += f1.getCurrentReward()
		response.ProtocolDelegation += delegator.DelegationAmount
	}

	// Iterate all Staker entries
	// Fetches the total delegation and calculates the outstanding rewards
	stakerStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StakerKeyPrefix))
	var stakerPrefix []byte
	stakerPrefix = append(stakerPrefix, []byte(req.Address)...)
	stakerPrefix = append(stakerPrefix, []byte("/")...)
	stakerIterator := sdk.KVStorePrefixIterator(stakerStore, stakerPrefix)

	defer stakerIterator.Close()

	for ; stakerIterator.Valid(); stakerIterator.Next() {
		var val types.Staker
		k.cdc.MustUnmarshal(stakerIterator.Value(), &val)

		response.ProtocolStaking += val.Amount
	}

	// Iterate all funding entries
	funderStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FunderKeyPrefix))
	funderIterator := sdk.KVStorePrefixIterator(funderStore, []byte(req.Address))

	defer funderIterator.Close()

	for ; funderIterator.Valid(); funderIterator.Next() {
		var val types.Funder
		k.cdc.MustUnmarshal(funderIterator.Value(), &val)

		response.ProtocolFunding += val.Amount
	}

	return &response, nil
}
