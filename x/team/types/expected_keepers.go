package types

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	util.BankKeeper
	GetSupply(ctx context.Context, denom string) sdk.Coin
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
