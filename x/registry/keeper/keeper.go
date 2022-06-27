package keeper

import (
	"fmt"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		distrKeeper   types.DistrKeeper
		upgradeKeeper types.UpgradeKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper types.DistrKeeper,
	upgradeKeeper types.UpgradeKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func (k Keeper) StoreKey() sdk.StoreKey {
	return k.storeKey
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) IterateProtocolBonding(ctx sdk.Context, address sdk.AccAddress, fn func(poolId uint64, amount sdk.Int) (stop bool)) {
	for _, pool := range k.GetAllPool(ctx) {
		total := uint64(0)

		//
		staker, isStaker := k.GetStaker(ctx, address.String(), pool.Id)
		if isStaker {
			total += staker.Amount
		}

		//
		delegatorPrefix := types.KeyPrefixBuilder{Key: types.DelegatorKeyPrefixIndex2}.AString(address.String()).AInt(pool.Id).Key
		delegatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), delegatorPrefix)
		delegatorIterator := sdk.KVStorePrefixIterator(delegatorStore, nil)

		for ; delegatorIterator.Valid(); delegatorIterator.Next() {
			key := delegatorIterator.Key()
			delegator, _ := k.GetDelegator(ctx, pool.Id, string(key[0:43]), address.String())

			total += delegator.DelegationAmount
		}

		delegatorIterator.Close()

		//
		stop := fn(pool.Id, sdk.NewIntFromUint64(total))
		if stop {
			break
		}
	}
}

func (k Keeper) TotalProtocolBonding(ctx sdk.Context) sdk.Int {
	total := uint64(0)

	for _, pool := range k.GetAllPool(ctx) {
		total += pool.TotalStake
		total += pool.TotalDelegation
	}

	return sdk.NewIntFromUint64(total)
}
