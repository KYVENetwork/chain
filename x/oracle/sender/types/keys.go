package types

import "fmt"

const (
	// ModuleName defines the module name.
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	RequestPrefix  = "request/"
	ResponsePrefix = "response/"
)

func RequestKey(sequence uint64) []byte {
	return []byte(fmt.Sprintf("%s%d", RequestPrefix, sequence))
}

func ResponseKey(sequence uint64) []byte {
	return []byte(fmt.Sprintf("%s%d", ResponsePrefix, sequence))
}
