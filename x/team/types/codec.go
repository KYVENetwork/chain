package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(_ *codec.LegacyAmino) {}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreateTeamVestingAccount{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimUnlocked{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClawback{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimAccountRewards{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimAuthorityRewards{})
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
