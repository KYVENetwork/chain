package types

import (
	"context"

	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"cosmossdk.io/x/upgrade/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type DistrKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
	AllocateTokensToValidator(ctx context.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) error
	GetValidatorAccumulatedCommission(ctx context.Context, val sdk.ValAddress) (commission distributionTypes.ValidatorAccumulatedCommission, err error)
	IncrementValidatorPeriod(ctx context.Context, val stakingtypes.ValidatorI) (uint64, error)
	CalculateDelegationRewards(ctx context.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, err error)
}

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx context.Context, plan types.Plan) error
}

type StakersKeeper interface {
	DoesStakerExist(ctx sdk.Context, staker string) bool
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
	GetValaccountsFromStaker(ctx sdk.Context, stakerAddress string) (val []*stakerstypes.PoolAccount)
	GetPoolCount(ctx sdk.Context, stakerAddress string) (poolCount uint64)
	GetActiveStakers(ctx sdk.Context) []string
}
