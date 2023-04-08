package types

import hostTypes "github.com/KYVENetwork/chain/x/oracle/host/types"

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(
	requests map[uint64]hostTypes.OracleQuery,
	responses map[uint64]hostTypes.OracleAcknowledgement,
) *GenesisState {
	return &GenesisState{
		Requests:  requests,
		Responses: responses,
	}
}

// DefaultGenesisState creates the default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		map[uint64]hostTypes.OracleQuery{},
		map[uint64]hostTypes.OracleAcknowledgement{},
	)
}

// ValidateGenesis validates the provided genesis state to ensure the expected invariants holds.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
