package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgFundPool{}, "kyve/pool/MsgFundPool", nil)
	cdc.RegisterConcrete(&MsgDefundPool{}, "kyve/pool/MsgDefundPool", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgFundPool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDefundPool{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDisablePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgEnablePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgScheduleRuntimeUpgrade{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCancelRuntimeUpgrade{})
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
