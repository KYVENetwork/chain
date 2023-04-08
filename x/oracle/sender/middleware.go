package sender

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Capability
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	// IBC Core
	clientTypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	portTypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibcExported "github.com/cosmos/ibc-go/v6/modules/core/exported"
	// Oracle Host
	hostTypes "github.com/KYVENetwork/chain/x/oracle/host/types"
	// Oracle Sender
	"github.com/KYVENetwork/chain/x/oracle/sender/keeper"
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

func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channelTypes.Packet,
	relayer sdk.AccAddress,
) ibcExported.Acknowledgement {
	return im.app.OnRecvPacket(ctx, packet, relayer)
}

// OnAcknowledgementPacket implements the IBCMiddleware interface.
func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channelTypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	_, req, valid, _ := hostTypes.ParseOraclePacket(packet)
	if !valid {
		return im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	}

	im.keeper.SetRequest(ctx, packet.Sequence, *req)

	// TODO(@john): Save response.

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
