package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"
	"github.com/KYVENetwork/chain/util"

	"cosmossdk.io/log"
	storeTypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		memService   store.MemoryStoreService
		logger       log.Logger

		authority string

		stakersKeeper types.StakersKeeper
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		distrkeeper   util.DistributionKeeper
		upgradeKeeper util.UpgradeKeeper
		fundersKeeper types.FundersKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	memService store.MemoryStoreService,
	logger log.Logger,

	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrKeeper util.DistributionKeeper,
	upgradeKeeper util.UpgradeKeeper,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		memService:   memService,
		logger:       logger,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrkeeper:   distrKeeper,
		upgradeKeeper: upgradeKeeper,
	}
}

func (k Keeper) EnsurePoolAccount(ctx sdk.Context, id uint64) {
	name := fmt.Sprintf("%s/%d", types.ModuleName, id)

	address := authTypes.NewModuleAddress(name)
	account := k.accountKeeper.GetAccount(ctx, address)

	if account == nil {
		// account doesn't exist, initialise a new module account.
		newAcc := authTypes.NewEmptyModuleAccount(name)
		account = k.accountKeeper.NewAccountWithAddress(ctx, newAcc.GetAddress())
	} else {
		// account exists, adjust it to a module account.
		baseAccount := authTypes.NewBaseAccount(address, nil, account.GetAccountNumber(), 0)

		account = authTypes.NewModuleAccount(baseAccount, name)
	}

	k.accountKeeper.SetAccount(ctx, account)
}

func SetStakersKeeper(k *Keeper, stakersKeeper types.StakersKeeper) {
	k.stakersKeeper = stakersKeeper
}

func SetFundersKeeper(k *Keeper, fundersKeeper types.FundersKeeper) {
	k.fundersKeeper = fundersKeeper
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storeTypes.StoreKey {
	// TODO: Check this
	return storeTypes.NewKVStoreKey(types.StoreKey)
}
