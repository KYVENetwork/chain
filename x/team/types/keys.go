package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "team"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

// Team module account address
// kyve1e29j95xmsw3zmvtrk4st8e89z5n72v7nf70ma4

// VESTING_DURATION 3 years
const VESTING_DURATION uint64 = 3 * 365 * 24 * 3600 // 3 * 365 * 24 * 3600

// UNLOCK_DURATION 2 years
const UNLOCK_DURATION uint64 = 2 * 365 * 24 * 3600 // 2 * 365 * 24 * 3600

// CLIFF_DURATION 1 year
const CLIFF_DURATION uint64 = 1 * 365 * 24 * 3600 // 1 * 365 * 24 * 3600

// FOUNDATION_ADDRESS is initialised in types.go by the init function which uses linker flags
var FOUNDATION_ADDRESS = ""

// BCP_ADDRESS is initialised in types.go by the init function which uses linker flags
var BCP_ADDRESS = ""

// TEAM_ALLOCATION is initialised in types.go by the init function which uses linker flags
var TEAM_ALLOCATION uint64 = 0

// TGE is initialised in types.go by the init function which uses linker flags
var TGE uint64 = 0

var (
	ParamsKey                  = []byte{0x00}
	AuthorityKey               = []byte{0x01}
	TeamVestingAccountKey      = []byte{0x02}
	TeamVestingAccountCountKey = []byte{0x03}
)

func TeamVestingAccountKeyPrefix(id uint64) []byte {
	return util.GetByteKey(id)
}
