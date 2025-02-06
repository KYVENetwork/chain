package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		memService   store.MemoryStoreService
		logger       log.Logger

		authority string

		bankKeeper    util.BankKeeper
		poolKeeper    types.PoolKeeper
		stakingKeeper util.StakingKeeper
		distKeeper    util.DistributionKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	memService store.MemoryStoreService,
	logger log.Logger,

	authority string,

	bankKeeper util.BankKeeper,
	poolKeeper types.PoolKeeper,
	stakingKeeper util.StakingKeeper,
	distributionKeeper util.DistributionKeeper,
) *Keeper {
	k := &Keeper{
		cdc:          cdc,
		storeService: storeService,
		memService:   memService,
		logger:       logger,

		authority: authority,

		bankKeeper:    bankKeeper,
		poolKeeper:    poolKeeper,
		stakingKeeper: stakingKeeper,
		distKeeper:    distributionKeeper,
	}

	return k
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// HOOKS

func (k Keeper) AfterValidatorCreated(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) BeforeValidatorModified(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) AfterValidatorBonded(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) AfterValidatorBeginUnbonding(goCtx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	for _, v := range k.GetPoolAccountsFromStaker(ctx, util.MustAccountAddressFromValAddress(valAddr.String())) {
		k.LeavePool(ctx, v.Staker, v.PoolId)
	}
	return nil
}

func (k Keeper) BeforeDelegationCreated(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) BeforeDelegationRemoved(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) BeforeValidatorSlashed(goCtx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	return nil
}

func (k Keeper) AfterUnbondingInitiated(ctx context.Context, id uint64) error {
	return nil
}
