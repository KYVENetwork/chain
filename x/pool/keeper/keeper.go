package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Distribution
	distributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	// Mint
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
	// Team
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storeTypes.StoreKey
		memKey   storeTypes.StoreKey

		authority string

		stakersKeeper types.StakersKeeper
		accountKeeper types.AccountKeeper
		bankKeeper    bankKeeper.Keeper
		distrkeeper   distributionKeeper.Keeper
		mintKeeper    mintKeeper.Keeper
		upgradeKeeper types.UpgradeKeeper
		teamKeeper    teamKeeper.Keeper
		fundersKeeper types.FundersKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storeTypes.StoreKey,
	memKey storeTypes.StoreKey,

	authority string,

	accountKeeper types.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	distrKeeper distributionKeeper.Keeper,
	mintKeeper mintKeeper.Keeper,
	upgradeKeeper types.UpgradeKeeper,
	teamKeeper teamKeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrkeeper:   distrKeeper,
		mintKeeper:    mintKeeper,
		upgradeKeeper: upgradeKeeper,
		teamKeeper:    teamKeeper,
	}
}

func (k Keeper) EnsurePoolAccount(ctx sdk.Context, id uint64) {
	name := fmt.Sprintf("%s/%d", types.ModuleName, id)

	address := authTypes.NewModuleAddress(name)
	account := k.accountKeeper.GetAccount(ctx, address)

	if account == nil {
		// account doesn't exist, initialise a new module account.
		account = authTypes.NewEmptyModuleAccount(name)
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

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() storeTypes.StoreKey {
	return k.storeKey
}
