package types

import (
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type PoolKeeper interface {
	AssertPoolExists(ctx sdk.Context, poolId uint64) error
}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx sdk.Context, plan types.Plan) error
}

type StakersKeeper interface {
	DoesStakerExist(ctx sdk.Context, staker string) bool
	GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string)
	GetValaccountsFromStaker(ctx sdk.Context, stakerAddress string) (val []*stakerstypes.Valaccount)
	GetPoolCount(ctx sdk.Context, stakerAddress string) (poolCount uint64)
	GetActiveStakers(ctx sdk.Context) []string
}
