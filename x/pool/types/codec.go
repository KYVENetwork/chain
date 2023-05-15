package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(_ *codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgFundPool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDefundPool{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDisablePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgEnablePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgScheduleRuntimeUpgrade{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCancelRuntimeUpgrade{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
