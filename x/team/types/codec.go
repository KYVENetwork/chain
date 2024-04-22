package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateTeamVestingAccount{}, "kyve/team/MsgCreateTeamVestingAccount", nil)
	cdc.RegisterConcrete(&MsgClaimUnlocked{}, "kyve/team/MsgClaimUnlocked", nil)
	cdc.RegisterConcrete(&MsgClawback{}, "kyve/team/MsgClawback", nil)
	cdc.RegisterConcrete(&MsgClaimAccountRewards{}, "kyve/team/MsgClaimAccountRewards", nil)
	cdc.RegisterConcrete(&MsgClaimAuthorityRewards{}, "kyve/team/MsgClaimAuthorityRewards", nil)
}

func RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreateTeamVestingAccount{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimUnlocked{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClawback{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimAccountRewards{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimAuthorityRewards{})
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
