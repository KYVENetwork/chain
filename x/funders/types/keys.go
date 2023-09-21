package types

const (
	// ModuleName defines the module name
	ModuleName = "funders"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_funders"
)

var (
	// ParamsKey is the prefix for all module params defined in params.proto
	ParamsKey = []byte{0x00}
)
