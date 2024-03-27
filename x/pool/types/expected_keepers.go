package types

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

type BankKeeper interface {
	util.BankKeeper
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx context.Context, plan upgradetypes.Plan) error
}

type StakersKeeper interface {
	LeavePool(ctx sdk.Context, staker string, poolId uint64)
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
}

type FundersKeeper interface {
	CreateFundingState(ctx sdk.Context, poolId uint64)
}
