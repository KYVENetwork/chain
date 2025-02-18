package types

import (
	"github.com/KYVENetwork/chain/x/stakers/types_v1beta1"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateCommission{}, "kyve/stakers/MsgUpdateCommission", nil)
	cdc.RegisterConcrete(&MsgUpdateStakeFraction{}, "kyve/stakers/MsgUpdateStakeFraction", nil)
	cdc.RegisterConcrete(&MsgJoinPool{}, "kyve/stakers/MsgJoinPool", nil)
	cdc.RegisterConcrete(&MsgLeavePool{}, "kyve/stakers/MsgLeavePool", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "kyve/stakers/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateCommission{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateStakeFraction{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgJoinPool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgLeavePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &types_v1beta1.MsgUpdateParams{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
