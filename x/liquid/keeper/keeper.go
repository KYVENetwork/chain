// This file was modified from its original version.
// https://github.com/cosmos/gaia/blob/b3583a755ad2947582f8e9c969806b187668ea23/x/liquid/keeper/keeper.go
// Modifications made by KYVE, 2025.
// - Changed import paths

package keeper

import (
	"context"

	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KYVENetwork/chain/x/liquid/types"
)

// Keeper of the x/liquid store
type Keeper struct {
	storeService  storetypes.KVStoreService
	cdc           codec.BinaryCodec
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	distKeeper    types.DistributionKeeper
	authority     string
}

// NewKeeper creates a new liquid Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	dk types.DistributionKeeper,
	authority string,
) *Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	return &Keeper{
		storeService:  storeService,
		cdc:           cdc,
		authKeeper:    ak,
		bankKeeper:    bk,
		stakingKeeper: sk,
		distKeeper:    dk,
		authority:     authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// GetAuthority returns the x/liquid module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
