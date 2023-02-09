package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(_ *codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreateStaker{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateCommission{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateMetadata{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgJoinPool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgLeavePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
