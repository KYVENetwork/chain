package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// Oracle Host
	hostTypes "github.com/KYVENetwork/chain/x/oracle/host/types"
	// Oracle Sender
	"github.com/KYVENetwork/chain/x/oracle/sender/types"
)

func (k Keeper) GetRequest(ctx sdk.Context, seq uint64) (req hostTypes.OracleQuery, found bool) {
	key := types.RequestKey(seq)

	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return hostTypes.OracleQuery{}, false
	}

	k.cdc.MustUnmarshal(bz, &req)
	return req, true
}

func (k Keeper) GetRequests(ctx sdk.Context) map[uint64]hostTypes.OracleQuery {
	requests := make(map[uint64]hostTypes.OracleQuery)
	itr := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.RequestPrefix))

	for ; itr.Valid(); itr.Next() {
		var req hostTypes.OracleQuery
		k.cdc.MustUnmarshal(itr.Value(), &req)

		requests[binary.BigEndian.Uint64(itr.Key())] = req
	}

	_ = itr.Close()

	return requests
}

func (k Keeper) SetRequest(ctx sdk.Context, seq uint64, req hostTypes.OracleQuery) {
	bz := k.cdc.MustMarshal(&req)
	ctx.KVStore(k.storeKey).Set(types.RequestKey(seq), bz)
}

func (k Keeper) GetResponse(ctx sdk.Context, seq uint64) (res hostTypes.OracleAcknowledgement, found bool) {
	key := types.ResponseKey(seq)

	bz := ctx.KVStore(k.storeKey).Get(key)
	if bz == nil {
		return hostTypes.OracleAcknowledgement{}, false
	}

	k.cdc.MustUnmarshal(bz, &res)
	return res, true
}

func (k Keeper) GetResponses(ctx sdk.Context) map[uint64]hostTypes.OracleAcknowledgement {
	responses := make(map[uint64]hostTypes.OracleAcknowledgement)
	itr := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte(types.ResponsePrefix))

	for ; itr.Valid(); itr.Next() {
		var res hostTypes.OracleAcknowledgement
		k.cdc.MustUnmarshal(itr.Value(), &res)

		responses[binary.BigEndian.Uint64(itr.Key())] = res
	}

	_ = itr.Close()

	return responses
}

func (k Keeper) SetResponse(ctx sdk.Context, seq uint64, res hostTypes.OracleAcknowledgement) {
	bz := k.cdc.MustMarshal(&res)
	ctx.KVStore(k.storeKey).Set(types.ResponseKey(seq), bz)
}
