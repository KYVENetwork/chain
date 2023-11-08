package keeper

import (
	"github.com/KYVENetwork/chain/x/query/types"
)

var (
	_ types.QueryAccountServer    = Keeper{}
	_ types.QueryPoolServer       = Keeper{}
	_ types.QueryStakersServer    = Keeper{}
	_ types.QueryDelegationServer = Keeper{}
	_ types.QueryBundlesServer    = Keeper{}
	_ types.QueryParamsServer     = Keeper{}
	_ types.QueryFundersServer    = Keeper{}
)
