package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateParams{}, "kyve/compliance/MsgUpdateParams", nil)
	cdc.RegisterConcrete(&MsgToggleMultiCoinRewards{}, "kyve/compliance/MsgToggleMultiCoinRewards", nil)
	cdc.RegisterConcrete(&MsgSetMultiCoinRewardsRefundPolicy{}, "kyve/compliance/MsgSetMultiCoinRewardsRefundPolicy", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgToggleMultiCoinRewards{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetMultiCoinRewardsRefundPolicy{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
