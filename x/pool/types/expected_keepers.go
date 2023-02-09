package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type AccountKeeper interface{}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx sdk.Context, plan upgradetypes.Plan) error
}

type StakersKeeper interface {
	LeavePool(ctx sdk.Context, staker string, poolId uint64)
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
}
