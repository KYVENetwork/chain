package types

import (
	"context"

	"cosmossdk.io/x/upgrade/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
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
	GetValaccountsFromStaker(ctx sdk.Context, stakerAddress string) (val []*stakerstypes.Valaccount)
	GetPoolCount(ctx sdk.Context, stakerAddress string) (poolCount uint64)
}
