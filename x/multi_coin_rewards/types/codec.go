package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateParams{}, "kyve/multi_coin_rewards/MsgUpdateParams", nil)
	cdc.RegisterConcrete(&MsgToggleMultiCoinRewards{}, "kyve/multi_coin_rewards/MsgToggleMultiCoinRewards", nil)
	cdc.RegisterConcrete(&MsgSetMultiCoinRewardsDistributionPolicy{}, "kyve/multi_coin_rewards/MsgSetMultiCoinRewardsDistributionPolicy", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgToggleMultiCoinRewards{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetMultiCoinRewardsDistributionPolicy{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
