package host

import (
	"encoding/json"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// Capability
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	// IBC Core
	clientTypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	portTypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibcExported "github.com/cosmos/ibc-go/v6/modules/core/exported"
	// Oracle Host
	"github.com/KYVENetwork/chain/x/oracle/host/keeper"
	"github.com/KYVENetwork/chain/x/oracle/host/types"
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
	data, req, valid, err := types.ParseOraclePacket(packet)
	if !valid {
		if err == nil {
			return im.app.OnRecvPacket(ctx, packet, relayer)
		} else {
			return types.NewErrorAcknowledgement(err)
		}
	}

	// Execute underlying transfer.
	if ack := im.app.OnRecvPacket(ctx, packet, relayer); !ack.Success() {
		return ack
	}

	// Execute query.
	var response []byte
	var timestamp time.Time

	switch query := req.Query.(type) {
	case *types.OracleQuery_LatestSummary:
		// TODO(@john): Handle query error.
		latestSummary, finalisedAt, _ := im.keeper.GetLatestSummary(ctx, query.LatestSummary.PoolId)

		response = []byte(latestSummary)
		timestamp = finalisedAt
	}

	// Ensure fee is sufficient.
	cost := im.keeper.GetPricePerByte(ctx).MulInt64(int64(len(response))).TruncateInt()
	if amount, _ := math.NewIntFromString(data.Amount); amount.LT(cost) {
		return types.NewErrorAcknowledgement(types.ErrInsufficientOracleFee)
	}

	// Return.
	bz, _ := json.Marshal(types.OracleAcknowledgement{
		OracleResponse: response,
		Timestamp:      timestamp,
	})
	return channelTypes.NewResultAcknowledgement(bz)
}

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
