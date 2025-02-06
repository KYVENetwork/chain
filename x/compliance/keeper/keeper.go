package keeper

import (
	"fmt"

	"cosmossdk.io/collections"

	"github.com/KYVENetwork/chain/x/compliance/types"

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

		MultiCoinRewardsEnabled collections.KeySet[sdk.AccAddress]
		MultiCoinRefundPolicy   collections.Item[types.MultiCoinRefundPolicy]

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

		MultiCoinRewardsEnabled: collections.NewKeySet(sb, types.MultiCoinRewardsEnabledKeyPrefix,
			"compliance_multi_coin_enabled", sdk.AccAddressKey),
		MultiCoinRefundPolicy: collections.NewItem(sb, types.MultiCoinRefundPolicyKeyPrefix,
			"compliance_multi_coin_policy", codec.CollValue[types.MultiCoinRefundPolicy](cdc)),
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
