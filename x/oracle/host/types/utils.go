package types

import (
	"encoding/json"

	"github.com/cosmos/gogoproto/jsonpb"

	// Global
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// IBC Core
	channelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	// IBC Transfer
	transferTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

// TODO(@john): Is hard-coding this the only solution?
const OracleAddress = "kyve1jgp27m8fykex4e4jtt0l7ze8q528ux2lxl4pxd"

func ParseOracleAcknowledgement(raw []byte) (ack *OracleAcknowledgement, valid bool) {
	var data channelTypes.Acknowledgement_Result
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, false
	}

	if err := json.Unmarshal(data.Result, &ack); err != nil {
		return nil, false
	}

	return ack, true
}

func ParseOraclePacket(packet channelTypes.Packet) (
	data *transferTypes.FungibleTokenPacketData, req *OracleQuery, valid bool, err error,
) {
	if err := json.Unmarshal(packet.Data, &data); err != nil {
		return nil, nil, false, nil
	}
	if data.Receiver != OracleAddress {
		return data, nil, false, nil
	}

	// Check -- Token Denom
	trace := transferTypes.ParseDenomTrace(data.Denom)
	isNativeToken := transferTypes.ReceiverChainIsSource(packet.SourcePort, packet.SourceChannel, data.Denom)
	isNativeKYVE := isNativeToken && trace.BaseDenom == globalTypes.Denom

	if !isNativeKYVE {
		return data, nil, false, ErrInvalidOracleToken
	}

	// Check -- Oracle Instructions
	var memo OracleMemo
	if err := jsonpb.UnmarshalString(data.GetMemo(), &memo); err != nil {
		return nil, nil, false, ErrInvalidOracleMemo
	}

	return data, memo.Query, true, nil
}
