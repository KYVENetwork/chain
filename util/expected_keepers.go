package util

import (
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	// Mint
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type AccountKeeper interface {
	GetModuleAddress(string) sdk.AccAddress
}

type BankKeeper interface {
	BurnCoins(sdk.Context, string, sdk.Coins) error
	GetAllBalances(sdk.Context, sdk.AccAddress) sdk.Coins
	GetBalance(sdk.Context, sdk.AccAddress, string) sdk.Coin
	GetSupply(sdk.Context, string) sdk.Coin
	SendCoins(sdk.Context, sdk.AccAddress, sdk.AccAddress, sdk.Coins) error
	SendCoinsFromAccountToModule(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
	SendCoinsFromModuleToModule(sdk.Context, string, string, sdk.Coins) error
}

type BundlesKeeper interface {
	AssertCanPropose(sdk.Context, uint64, string, string, uint64) error
	AssertCanVote(sdk.Context, uint64, string, string, string) error
}

type DelegationKeeper interface {
	GetDelegationAmount(sdk.Context, string) uint64
	GetDelegationAmountOfDelegator(sdk.Context, string, string) uint64
	GetDelegationOfPool(sdk.Context, uint64) uint64
	GetOutstandingRewards(sdk.Context, string, string) uint64
	GetPaginatedActiveStakersByDelegation(sdk.Context, *query.PageRequest, func(string, bool) bool) (*query.PageResponse, error)
	GetPaginatedActiveStakersByPoolCountAndDelegation(sdk.Context, *query.PageRequest) ([]string, *query.PageResponse, error)
	GetPaginatedInactiveStakersByDelegation(sdk.Context, *query.PageRequest, func(string, bool) bool) (*query.PageResponse, error)
	GetPaginatedStakersByDelegation(sdk.Context, *query.PageRequest, func(string, bool) bool) (*query.PageResponse, error)
	GetRedelegationCooldown(sdk.Context) uint64
	GetRedelegationCooldownEntries(sdk.Context, string) []uint64
	GetRedelegationMaxAmount(sdk.Context) uint64
	GetStakersByDelegator(sdk.Context, string) []string
	PayoutRewards(sdk.Context, string, uint64, string) bool
	StoreKey() storeTypes.StoreKey
}

type DistributionKeeper interface {
	FundCommunityPool(sdk.Context, sdk.Coins, sdk.AccAddress) error
}

type GlobalKeeper interface{}

type GovKeeper interface {
	GetParams(sdk.Context) govTypes.Params
}

type MintKeeper interface {
	GetMinter(sdk.Context) mintTypes.Minter
	GetParams(sdk.Context) mintTypes.Params
}

type PoolKeeper interface {
	ChargeFundersOfPool(sdk.Context, uint64, uint64) error
	IncrementBundleInformation(sdk.Context, uint64, uint64, string, string)
}

type StakersKeeper interface {
	AssertValaccountAuthorized(sdk.Context, uint64, string, string) error
	DoesStakerExist(sdk.Context, string) bool
	DoesValaccountExist(sdk.Context, uint64, string) bool
	GetActiveStakers(sdk.Context) []string
	GetAllStakerAddressesOfPool(sdk.Context, uint64) []string
	GetCommission(sdk.Context, string) sdk.Dec
	GetPoolCount(sdk.Context, string) uint64
	IncrementPoints(sdk.Context, uint64, string) uint64
	LeavePool(sdk.Context, string, uint64)
	ResetPoints(sdk.Context, uint64, string) uint64
}

type UpgradeKeeper interface {
	ScheduleUpgrade(sdk.Context, upgradeTypes.Plan) error
}
