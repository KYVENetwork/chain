package types

import (
	"github.com/KYVENetwork/chain/util"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type StakersKeeper interface {
	util.StakersKeeper

	GetValaccountsFromStaker(sdk.Context, string) []*stakersTypes.Valaccount
}
