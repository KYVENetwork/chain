package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDisablePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgEnablePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgScheduleRuntimeUpgrade{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCancelRuntimeUpgrade{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
