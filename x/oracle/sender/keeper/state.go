package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// Oracle Host
	hostTypes "github.com/KYVENetwork/chain/x/oracle/host/types"
	// Oracle Sender
	"github.com/KYVENetwork/chain/x/oracle/sender/types"
)

func (k Keeper) GetRequests(ctx sdk.Context) map[uint64]hostTypes.OracleQuery {
	requests := make(map[uint64]hostTypes.OracleQuery)
	itr := ctx.KVStore(k.storeKey).Iterator([]byte(types.RequestPrefix), nil)

	for ; itr.Valid(); itr.Next() {
		var req hostTypes.OracleQuery
		k.cdc.MustUnmarshal(itr.Value(), &req)

		// TODO(@john): Is it better to use LittleEndian here?
		requests[binary.BigEndian.Uint64(itr.Key())] = req
	}

	return requests
}

func (k Keeper) SetRequest(ctx sdk.Context, seq uint64, req hostTypes.OracleQuery) {
	bz := k.cdc.MustMarshal(&req)
	ctx.KVStore(k.storeKey).Set(types.RequestKey(seq), bz)
}

func (k Keeper) GetResponses(ctx sdk.Context) map[uint64]hostTypes.OracleAcknowledgement {
	responses := make(map[uint64]hostTypes.OracleAcknowledgement)
	itr := ctx.KVStore(k.storeKey).Iterator([]byte(types.ResponsePrefix), nil)

	for ; itr.Valid(); itr.Next() {
		var res hostTypes.OracleAcknowledgement
		k.cdc.MustUnmarshal(itr.Value(), &res)

		// TODO(@john): Is it better to use LittleEndian here?
		responses[binary.BigEndian.Uint64(itr.Key())] = res
	}

	return responses
}

func (k Keeper) SetResponse(ctx sdk.Context, seq uint64, res hostTypes.OracleAcknowledgement) {
	bz := k.cdc.MustMarshal(&res)
	ctx.KVStore(k.storeKey).Set(types.ResponseKey(seq), bz)
}
