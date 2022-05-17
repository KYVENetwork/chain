package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgFundPool{}, "registry/FundPool", nil)
	cdc.RegisterConcrete(&MsgDefundPool{}, "registry/DefundPool", nil)
	cdc.RegisterConcrete(&MsgStakePool{}, "registry/StakePool", nil)
	cdc.RegisterConcrete(&MsgUnstakePool{}, "registry/UnstakePool", nil)
	cdc.RegisterConcrete(&MsgSubmitBundleProposal{}, "registry/SubmitBundleProposal", nil)
	cdc.RegisterConcrete(&MsgVoteProposal{}, "registry/VoteProposal", nil)
	cdc.RegisterConcrete(&MsgClaimUploaderRole{}, "registry/ClaimUploaderRole", nil)
	cdc.RegisterConcrete(&MsgDelegatePool{}, "registry/DelegatePool", nil)
	cdc.RegisterConcrete(&MsgWithdrawPool{}, "registry/WithdrawPool", nil)
	cdc.RegisterConcrete(&MsgUndelegatePool{}, "registry/UndelegatePool", nil)
	cdc.RegisterConcrete(&MsgUpdateMetadata{}, "registry/UpdateMetadata", nil)
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&CreatePoolProposal{}, "kyve/CreatePoolProposal", nil)
	cdc.RegisterConcrete(&UpdatePoolProposal{}, "kyve/UpdatePoolProposal", nil)
	cdc.RegisterConcrete(&PausePoolProposal{}, "kyve/PausePoolProposal", nil)
	cdc.RegisterConcrete(&UnpausePoolProposal{}, "kyve/UnpausePoolProposal", nil)
	cdc.RegisterConcrete(&SchedulePoolUpgradeProposal{}, "kyve/SchedulePoolUpgradeProposal", nil)
	cdc.RegisterConcrete(&CancelPoolUpgradeProposal{}, "kyve/CancelPoolUpgradeProposal", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFundPool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDefundPool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgStakePool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnstakePool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitBundleProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVoteProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaimUploaderRole{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDelegatePool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawPool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUndelegatePool{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMetadata{},
	)
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&CreatePoolProposal{},
		&UpdatePoolProposal{},
		&PausePoolProposal{},
		&UnpausePoolProposal{},
		&SchedulePoolUpgradeProposal{},
		&CancelPoolUpgradeProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
