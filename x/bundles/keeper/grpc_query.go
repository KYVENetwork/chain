package keeper

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
)

var _ types.QueryServer = Keeper{}
