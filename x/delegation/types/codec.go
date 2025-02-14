package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
