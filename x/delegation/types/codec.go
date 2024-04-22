package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDelegate{}, "kyve/delegation/MsgDelegate", nil)
	cdc.RegisterConcrete(&MsgWithdrawRewards{}, "kyve/delegation/MsgWithdrawRewards", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "kyve/delegation/MsgUndelegate", nil)
	cdc.RegisterConcrete(&MsgRedelegate{}, "kyve/delegation/MsgRedelegate", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgDelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgWithdrawRewards{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUndelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRedelegate{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
