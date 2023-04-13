package types

import (
	"cosmossdk.io/errors"
	channelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
)

var (
	ErrInvalidOracleToken    = errors.Register(ModuleName, 0, "invalid token sent to oracle")
	ErrInvalidOracleMemo     = errors.Register(ModuleName, 1, "invalid memo sent to oracle")
	ErrInsufficientOracleFee = errors.Register(ModuleName, 2, "insufficient fee sent to oracle")
	QueryError               = errors.Register(ModuleName, 3, "an error occurred while executing the query")
)

// NewErrorAcknowledgement is inspired by channelTypes.NewNewErrorAcknowledgement,
// however returns the error in a more readable format.
func NewErrorAcknowledgement(err error) channelTypes.Acknowledgement {
	return channelTypes.Acknowledgement{
		Response: &channelTypes.Acknowledgement_Error{
			Error: err.Error(),
		},
	}
}
