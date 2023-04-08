package types

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState creates the default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams())
}

// ValidateGenesis validates the provided genesis state to ensure the expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
