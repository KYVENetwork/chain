package keeper

import (
	"fmt"

	"cosmossdk.io/collections"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	"github.com/KYVENetwork/chain/util"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		memService   store.MemoryStoreService
		logger       log.Logger

		authority string

		accountKeeper util.AccountKeeper
		bankKeeper    util.BankKeeper
		poolKeeper    types.PoolKeeper

		MultiCoinRewardsEnabled     collections.KeySet[sdk.AccAddress]
		MultiCoinDistributionPolicy collections.Item[types.MultiCoinDistributionPolicy]

		Schema collections.Schema
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	memService store.MemoryStoreService,
	logger log.Logger,

	authority string,

	accountKeeper util.AccountKeeper,
	bankKeeper util.BankKeeper,
	poolKeeper types.PoolKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		memService:   memService,
		logger:       logger,

		authority: authority,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		poolKeeper:    poolKeeper,

		MultiCoinRewardsEnabled: collections.NewKeySet(sb, types.MultiCoinRewardsEnabledKey,
			"multi_coin_rewards_enabled", sdk.AccAddressKey),
		MultiCoinDistributionPolicy: collections.NewItem(sb, types.MultiCoinDistributionPolicyKey,
			"multi_coin_rewards_policy", codec.CollValue[types.MultiCoinDistributionPolicy](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
