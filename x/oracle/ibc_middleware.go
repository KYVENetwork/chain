package oracle

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/jsonpb"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Capability
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	// Global
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// IBC Core
	clientTypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	portTypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibcExported "github.com/cosmos/ibc-go/v6/modules/core/exported"
	// IBC Transfer
	transferTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	// Oracle
	"github.com/KYVENetwork/chain/x/oracle/keeper"
	"github.com/KYVENetwork/chain/x/oracle/types"
)

var _ portTypes.Middleware = &IBCMiddleware{}

type IBCMiddleware struct {
	app    portTypes.IBCModule
	keeper keeper.Keeper
}

func NewIBCMiddleware(app portTypes.IBCModule, keeper keeper.Keeper) IBCMiddleware {
	return IBCMiddleware{app: app, keeper: keeper}
}

func (im IBCMiddleware) OnChanOpenInit(
	ctx sdk.Context,
	order channelTypes.Order,
	hops []string,
	portID string,
	channelID string,
	capability *capabilityTypes.Capability,
	counterparty channelTypes.Counterparty,
	version string,
) (string, error) {
	return im.app.OnChanOpenInit(ctx, order, hops, portID, channelID, capability, counterparty, version)
}

func (im IBCMiddleware) OnChanOpenTry(
	ctx sdk.Context,
	order channelTypes.Order,
	hops []string,
	portID string,
	channelID string,
	capability *capabilityTypes.Capability,
	counterparty channelTypes.Counterparty,
	counterpartyVersion string,
) (version string, err error) {
	return im.app.OnChanOpenTry(ctx, order, hops, portID, channelID, capability, counterparty, counterpartyVersion)
}

func (im IBCMiddleware) OnChanOpenAck(
	ctx sdk.Context,
	portID string,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (im IBCMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID string, channelID string) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (im IBCMiddleware) OnChanCloseInit(ctx sdk.Context, portID string, channelID string) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

func (im IBCMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID string, channelID string) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the IBCMiddleware interface.
//
// This middleware is intended to be used alongside ICS-20 ("IBC Transfer").
// When receiving an ICS-20 packet, we first check if the tokens are being sent
// to the x/oracle module. If they are, we ensure that the denom sent is the
// native KYVE token (ukyve). Note that if the denom is not correct, we simply
// throw an error, returning the funds. Next, we utilise the memo field to
// trigger an interchain query request.
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channelTypes.Packet,
	relayer sdk.AccAddress,
) ibcExported.Acknowledgement {
	var parsedPacket transferTypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.Data, &parsedPacket); err != nil {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}

	if parsedPacket.Receiver != authTypes.NewModuleAddress(types.ModuleName).String() {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}

	// Ensure that the tokens being transferred are native KYVE.
	// TODO: Is this better than doing a regex check? -- ^transfer/channel-\d+$
	trace := transferTypes.ParseDenomTrace(parsedPacket.Denom)
	isNativeToken := transferTypes.ReceiverChainIsSource(packet.SourcePort, packet.SourceChannel, parsedPacket.Denom)
	isNativeKYVE := isNativeToken && trace.BaseDenom == globalTypes.Denom

	if !isNativeKYVE {
		return types.NewErrorAcknowledgement(types.ErrInvalidOracleToken)
	}

	// Attempt to parse the memo included in the transfer, error if unrecognisable.
	var parsedMemo types.OracleMemo
	if err := jsonpb.UnmarshalString(parsedPacket.Memo, &parsedMemo); err != nil {
		return types.NewErrorAcknowledgement(types.ErrInvalidOracleMemo)
	}

	// First, let's execute the underlying token transfer.
	rawIBCAck := im.app.OnRecvPacket(ctx, packet, relayer)
	if !rawIBCAck.Success() {
		return rawIBCAck
	}

	var ibcAck channelTypes.Acknowledgement_Result
	_ = json.Unmarshal(rawIBCAck.Acknowledgement(), &ibcAck)

	// Next, let's execute the query and ensure the included fee is sufficient.
	var response []byte
	var timestamp time.Time

	switch query := parsedMemo.Query.Query.(type) {
	case *types.OracleQuery_LatestSummary:
		// TODO(@john): Handle query error.
		latestSummary, finalisedAt, _ := im.keeper.GetLatestSummary(ctx, query.LatestSummary.PoolId)

		response = []byte(latestSummary)
		timestamp = finalisedAt
	default:
		// TODO: Do we even need this? This is in theory unreachable.
		return types.NewErrorAcknowledgement(types.ErrInvalidOracleMemo)
	}

	cost := im.keeper.GetPricePerByte(ctx).MulInt64(int64(len(response))).TruncateInt()
	if amount, _ := math.NewIntFromString(parsedPacket.Amount); amount.LT(cost) {
		return types.NewErrorAcknowledgement(types.ErrInsufficientOracleFee)
	}

	// Return an acknowledgement with the query response.
	bz, _ := json.Marshal(types.OracleAcknowledgement{
		IBCAcknowledgement: ibcAck.Result,
		OracleResponse:     response,
		Timestamp:          timestamp,
	})
	return channelTypes.NewResultAcknowledgement(bz)
}

// OnAcknowledgementPacket implements the IBCMiddleware interface.
//
// TODO(@john): Save query responses.
func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channelTypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channelTypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) SendPacket(
	_ sdk.Context,
	_ *capabilityTypes.Capability,
	_ string,
	_ string,
	_ clientTypes.Height,
	_ uint64,
	_ []byte,
) (sequence uint64, err error) {
	panic("SendPacket is not supported by the KYVE Oracle module.")
}

func (im IBCMiddleware) WriteAcknowledgement(
	_ sdk.Context,
	_ *capabilityTypes.Capability,
	_ ibcExported.PacketI,
	_ ibcExported.Acknowledgement,
) error {
	panic("WriteAcknowledgement is not supported by the KYVE Oracle module.")
}

func (im IBCMiddleware) GetAppVersion(_ sdk.Context, _ string, _ string) (string, bool) {
	panic("GetAppVersion is currently not implemented.")
}
