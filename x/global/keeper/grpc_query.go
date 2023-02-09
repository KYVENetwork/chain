package keeper

import (
	"github.com/KYVENetwork/chain/x/global/types"
)

var _ types.QueryServer = Keeper{}
