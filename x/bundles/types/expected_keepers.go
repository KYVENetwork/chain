package types

import (
	"cosmossdk.io/math"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
	GetPoolWithError(ctx sdk.Context, poolId uint64) (pooltypes.Pool, error)
	GetPool(ctx sdk.Context, id uint64) (val pooltypes.Pool, found bool)

	IncrementBundleInformation(ctx sdk.Context, poolId uint64, currentHeight uint64, currentKey string, currentValue string)

	GetAllPools(ctx sdk.Context) (list []pooltypes.Pool)
	ChargeInflationPool(ctx sdk.Context, poolId uint64) (payout uint64, err error)
	GetProtocolInflationShare(ctx sdk.Context) (res math.LegacyDec)
	GetMaxVotingPowerPerPool(ctx sdk.Context) (res math.LegacyDec)
}

type StakerKeeper interface {
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
	AssertPoolAccountAuthorized(ctx sdk.Context, stakerAddress string, poolId uint64, poolAddress string) error

	GetPoolAccount(ctx sdk.Context, stakerAddress string, poolId uint64) (valaccount stakersTypes.PoolAccount, active bool)

	LeavePool(ctx sdk.Context, stakerAddress string, poolId uint64)

	IncrementPoints(ctx sdk.Context, stakerAddress string, poolId uint64) (newPoints uint64)
	ResetPoints(ctx sdk.Context, stakerAddress string, poolId uint64) (previousPoints uint64)

	GetValidator(ctx sdk.Context, stakerAddress string) (stakingtypes.Validator, bool)
	GetValidatorPoolCommission(ctx sdk.Context, staker string, poolId uint64) math.LegacyDec
	GetValidatorPoolStake(ctx sdk.Context, staker string, poolId uint64) uint64
	GetTotalStakeOfPool(ctx sdk.Context, poolId uint64) (totalStake uint64)
	GetValidatorPoolStakes(ctx sdk.Context, poolId uint64, mustIncludeStakers ...string) map[string]uint64
	IsVotingPowerTooHigh(ctx sdk.Context, poolId uint64) bool
	Slash(ctx sdk.Context, poolId uint64, staker string, slashType stakersTypes.SlashType)
	PayoutRewards(ctx sdk.Context, staker string, amount sdk.Coins, payerModuleName string) error
	PayoutAdditionalCommissionRewards(ctx sdk.Context, validator string, payerModuleName string, amount sdk.Coins) error
}

type FundersKeeper interface {
	GetCoinWhitelistMap(ctx sdk.Context) (whitelist map[string]fundersTypes.WhitelistCoinEntry)
	ChargeFundersOfPool(ctx sdk.Context, poolId uint64, recipient string) (payout sdk.Coins, err error)
}

type TeamKeeper interface {
	GetTeamBlockProvision(ctx sdk.Context) int64
}
