package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(_ *codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitBundleProposal{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgVoteBundleProposal{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimUploaderRole{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSkipUploaderRole{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateParams{})
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
