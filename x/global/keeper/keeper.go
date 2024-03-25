package keeper

import (
	"cosmossdk.io/log"
	"fmt"

	storeTypes "cosmossdk.io/store/types"
	"github.com/KYVENetwork/chain/x/global/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storeTypes.StoreKey
		logger   log.Logger

		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storeTypes.StoreKey,
	logger log.Logger,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		logger:   logger,

		authority: authority,
	}
}

func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
