package keeper

import (
	"github.com/KYVENetwork/chain/x/funders/types"
)

var _ types.QueryServer = Keeper{}
