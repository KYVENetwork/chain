package util

import (
	"context"

	"cosmossdk.io/math"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	upgradeTypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AccountKeeper interface {
	GetModuleAddress(string) sdk.AccAddress
}

type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type DistributionKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
	AllocateTokensToValidator(ctx context.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) error
	IncrementValidatorPeriod(ctx context.Context, val stakingtypes.ValidatorI) (uint64, error)
	CalculateDelegationRewards(ctx context.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, err error)

	GetValidatorAccumulatedCommission(ctx context.Context, val sdk.ValAddress) (commission distributionTypes.ValidatorAccumulatedCommission, err error)
	SetValidatorAccumulatedCommission(ctx context.Context, val sdk.ValAddress, commission distributionTypes.ValidatorAccumulatedCommission) error
	GetValidatorOutstandingRewards(ctx context.Context, val sdk.ValAddress) (rewards distributionTypes.ValidatorOutstandingRewards, err error)
	SetValidatorOutstandingRewards(ctx context.Context, val sdk.ValAddress, rewards distributionTypes.ValidatorOutstandingRewards) error
}

type StakingKeeper interface {
	Slash(ctx context.Context, consAddr sdk.ConsAddress, infractionHeight, power int64, slashFactor math.LegacyDec) (math.Int, error)
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error)
	PowerReduction(ctx context.Context) math.Int
	SetHooks(sh stakingtypes.StakingHooks)
	Delegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.DelegationI, error)
}

type UpgradeKeeper interface {
	ScheduleUpgrade(context.Context, upgradeTypes.Plan) error
}
