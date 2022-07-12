package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeCreatePool  = "CreatePool"
	ProposalTypeUpdatePool  = "UpdatePool"
	ProposalTypePausePool   = "PausePool"
	ProposalTypeUnpausePool = "UnpausePool"
	ProposalTypeSchedulePoolUpgrade = "SchedulePoolUpgrade"
	ProposalTypeCancelPoolUpgrade = "CancelPoolUpgrade"
	ProposalTypeResetPool = "ResetPool"
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeCreatePool)
	govtypes.RegisterProposalTypeCodec(&CreatePoolProposal{}, "kyve/CreatePoolProposal")
	govtypes.RegisterProposalType(ProposalTypeUpdatePool)
	govtypes.RegisterProposalTypeCodec(&UpdatePoolProposal{}, "kyve/UpdatePoolProposal")
	govtypes.RegisterProposalType(ProposalTypePausePool)
	govtypes.RegisterProposalTypeCodec(&PausePoolProposal{}, "kyve/PausePoolProposal")
	govtypes.RegisterProposalType(ProposalTypeUnpausePool)
	govtypes.RegisterProposalTypeCodec(&UnpausePoolProposal{}, "kyve/UnpausePoolProposal")
	govtypes.RegisterProposalType(ProposalTypeSchedulePoolUpgrade)
	govtypes.RegisterProposalTypeCodec(&SchedulePoolUpgradeProposal{}, "kyve/SchedulePoolUpgradeProposal")
	govtypes.RegisterProposalType(ProposalTypeCancelPoolUpgrade)
	govtypes.RegisterProposalTypeCodec(&CancelPoolUpgradeProposal{}, "kyve/CancelPoolUpgradeProposal")
	govtypes.RegisterProposalType(ProposalTypeResetPool)
	govtypes.RegisterProposalTypeCodec(&ResetPoolProposal{}, "kyve/ResetPoolProposal")
}

var (
	_ govtypes.Content = &CreatePoolProposal{}
	_ govtypes.Content = &UpdatePoolProposal{}
	_ govtypes.Content = &PausePoolProposal{}
	_ govtypes.Content = &UnpausePoolProposal{}
	_ govtypes.Content = &SchedulePoolUpgradeProposal{}
	_ govtypes.Content = &CancelPoolUpgradeProposal{}
	_ govtypes.Content = &ResetPoolProposal{}
)

func NewCreatePoolProposal(title string, description string, name string, runtime string, logo string, config string, uploadInterval uint64, operatingCost uint64, maxBundleSize uint64, version string, binaries string, startKey string, minStake uint64) govtypes.Content {
	return &CreatePoolProposal{
		Title:         title,
		Description:   description,
		Name:          name,
		Runtime:       runtime,
		Logo:          logo,
		Config:        config,
		UploadInterval: uploadInterval,
		OperatingCost: operatingCost,
		MaxBundleSize: maxBundleSize,
		Version: version,
		Binaries: binaries,
		StartKey: startKey,
		MinStake: minStake,
	}
}

func (p *CreatePoolProposal) ProposalRoute() string { return RouterKey }

func (p *CreatePoolProposal) ProposalType() string {
	return ProposalTypeCreatePool
}

func (p *CreatePoolProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewUpdatePoolProposal(title string, description string, id uint64, name string, runtime string, logo string, config string, uploadInterval uint64, operatingCost uint64, maxBundleSize uint64, minStake uint64) govtypes.Content {
	return &UpdatePoolProposal{
		Title:         title,
		Description:   description,
		Id:            id,
		Name:          name,
		Runtime:       runtime,
		Logo:          logo,
		Config:        config,
		UploadInterval: uploadInterval,
		OperatingCost: operatingCost,
		MaxBundleSize: maxBundleSize,
		MinStake: minStake,
	}
}

func (p *UpdatePoolProposal) ProposalRoute() string { return RouterKey }

func (p *UpdatePoolProposal) ProposalType() string {
	return ProposalTypeUpdatePool
}

func (p *UpdatePoolProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewPausePoolProposal(title string, description string, id uint64) govtypes.Content {
	return &PausePoolProposal{
		Title:       title,
		Description: description,
		Id:          id,
	}
}

func (p *PausePoolProposal) ProposalRoute() string { return RouterKey }

func (p *PausePoolProposal) ProposalType() string {
	return ProposalTypePausePool
}

func (p *PausePoolProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewUnpausePoolProposal(title string, description string, id uint64) govtypes.Content {
	return &UnpausePoolProposal{
		Title:       title,
		Description: description,
		Id:          id,
	}
}

func (p *UnpausePoolProposal) ProposalRoute() string { return RouterKey }

func (p *UnpausePoolProposal) ProposalType() string {
	return ProposalTypeUnpausePool
}

func (p *UnpausePoolProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewSchedulePoolUpgradeProposal(title string, description string, runtime string, version string, scheduled_at uint64, duration uint64, binaries string) govtypes.Content {
	return &SchedulePoolUpgradeProposal{
		Title:       title,
		Description: description,
		Runtime:          runtime,
		Version: version,
		ScheduledAt: scheduled_at,
		Duration: duration,
		Binaries: binaries,
	}
}

func (p *SchedulePoolUpgradeProposal) ProposalRoute() string { return RouterKey }

func (p *SchedulePoolUpgradeProposal) ProposalType() string {
	return ProposalTypeUnpausePool
}

func (p *SchedulePoolUpgradeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewCancelPoolUpgradeProposal(title string, description string, runtime string) govtypes.Content {
	return &CancelPoolUpgradeProposal{
		Title:       title,
		Description: description,
		Runtime:          runtime,
	}
}

func (p *CancelPoolUpgradeProposal) ProposalRoute() string { return RouterKey }

func (p *CancelPoolUpgradeProposal) ProposalType() string {
	return ProposalTypeUnpausePool
}

func (p *CancelPoolUpgradeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}

func NewResetPoolProposal(title string, description string, id uint64, bundleId uint64) govtypes.Content {
	return &ResetPoolProposal{
		Title:       title,
		Description: description,
		Id:          id,
		BundleId: bundleId,
	}
}

func (p *ResetPoolProposal) ProposalRoute() string { return RouterKey }

func (p *ResetPoolProposal) ProposalType() string {
	return ProposalTypeUnpausePool
}

func (p *ResetPoolProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(p)
	if err != nil {
		return err
	}

	return nil
}
