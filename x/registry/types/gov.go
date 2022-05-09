package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeCreatePool  = "CreatePool"
	ProposalTypeUpdatePool  = "UpdatePool"
	ProposalTypePausePool   = "PausePool"
	ProposalTypeUnpausePool = "UnpausePool"
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
}

var (
	_ govtypes.Content = &CreatePoolProposal{}
	_ govtypes.Content = &UpdatePoolProposal{}
	_ govtypes.Content = &PausePoolProposal{}
	_ govtypes.Content = &UnpausePoolProposal{}
)

func NewCreatePoolProposal(title string, description string, name string, runtime string, logo string, versions string, config string, startHeight uint64, uploadInterval uint64, operatingCost uint64, maxBundleSize uint64) govtypes.Content {
	return &CreatePoolProposal{
		Title:         title,
		Description:   description,
		Name:          name,
		Runtime:       runtime,
		Logo:          logo,
		Versions:      versions,
		Config:        config,
		StartHeight:   startHeight,
		UploadInterval: uploadInterval,
		OperatingCost: operatingCost,
		MaxBundleSize: maxBundleSize,
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

func NewUpdatePoolProposal(title string, description string, id uint64, name string, runtime string, logo string, versions string, config string, uploadInterval uint64, operatingCost uint64, maxBundleSize uint64) govtypes.Content {
	return &UpdatePoolProposal{
		Title:         title,
		Description:   description,
		Id:            id,
		Name:          name,
		Runtime:       runtime,
		Logo:          logo,
		Versions:      versions,
		Config:        config,
		UploadInterval: uploadInterval,
		OperatingCost: operatingCost,
		MaxBundleSize: maxBundleSize,
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
