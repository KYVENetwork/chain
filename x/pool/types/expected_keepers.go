package types

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx sdk.Context, plan upgradetypes.Plan) error
}

type StakersKeeper interface {
	LeavePool(ctx sdk.Context, staker string, poolId uint64)
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
}

type FundersKeeper interface {
	CreateFundingState(ctx sdk.Context, poolId uint64)
}
