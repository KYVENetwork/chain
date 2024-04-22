package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitBundleProposal{}, "kyve/bundles/MsgSubmitBundleProposal", nil)
	cdc.RegisterConcrete(&MsgVoteBundleProposal{}, "kyve/bundles/MsgVoteBundleProposal", nil)
	cdc.RegisterConcrete(&MsgClaimUploaderRole{}, "kyve/bundles/MsgClaimUploaderRole", nil)
	cdc.RegisterConcrete(&MsgSkipUploaderRole{}, "kyve/bundles/MsgSkipUploaderRole", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitBundleProposal{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgVoteBundleProposal{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimUploaderRole{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSkipUploaderRole{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
