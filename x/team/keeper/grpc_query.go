package keeper

import (
	"github.com/KYVENetwork/chain/x/team/types"
)

var _ types.QueryServer = Keeper{}
