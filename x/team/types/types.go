package types

import (
	"errors"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// VestingPlan contains basic information for one member
type VestingPlan struct {
	MaximumVestingAmount uint64

	ClawbackAmount uint64

	TokenVestingStart uint64

	TokenVestingFinished uint64

	TokenUnlockStart uint64

	TokenUnlockFinished uint64
}

// VestingStatus contains computed vesting values for a member at a given time
type VestingStatus struct {
	// vested_amount ...
	TotalVestedAmount uint64
	// unlocked_amount ...
	TotalUnlockedAmount uint64

	// i.e. U(t) - C
	CurrentClaimableAmount uint64

	// unvested_amount ...
	LockedVestedAmount uint64

	// unvested_amount ...
	RemainingUnvestedAmount uint64
}

var (
	TEAM_FOUNDATION_STRING = "kyve1u7ukf2nv6v5j5y2yqprm8yqruue2rlmrkx4xgq"
	TEAM_BCP_STRING        = "kyve1ruxaec07ca3dh0amkzxjap7av3xjt5vjgnd424"
	TEAM_ALLOCATION_STRING = "165000000000000"
	TGE_STRING             = "2023-02-01T10:34:15"
)

// Convert passed build variables (string) to the corresponding int values
func init() {
	// Foundation needs to be a valid Bech32 address
	foundationPrefix, _, err := bech32.DecodeAndConvert(TEAM_FOUNDATION_STRING)
	if err != nil {
		panic(err)
	}
	if foundationPrefix != "kyve" {
		panic(errors.New("team foundation address is not a KYVE address"))
	}
	FOUNDATION_ADDRESS = TEAM_FOUNDATION_STRING

	// BCP needs to be a valid Bech32 address
	bcpPrefix, _, err := bech32.DecodeAndConvert(TEAM_BCP_STRING)
	if err != nil {
		panic(err)
	}
	if bcpPrefix != "kyve" {
		panic(errors.New("team foundation address is not a KYVE address"))
	}
	BCP_ADDRESS = TEAM_BCP_STRING

	// TEAM_ALLOCATION must be a valid integer
	parsedAllocation, err := strconv.ParseUint(TEAM_ALLOCATION_STRING, 10, 64)
	if err != nil {
		panic(err)
	}
	TEAM_ALLOCATION = parsedAllocation

	// TGE must be a valid unix date-string
	tge, err := time.Parse("2006-01-02T15:04:05", TGE_STRING)
	if err != nil {
		panic(err)
	}
	TGE = uint64(tge.Unix())
}
