package types

import (
	"context"

	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

type BundlesKeeper interface {
	AssertCanVote(sdk.Context, uint64, string, string, string) error
	AssertCanPropose(sdk.Context, uint64, string, string, uint64) error
	GetBundleVersionMap(sdk.Context) bundlesTypes.BundleVersionMap
	GetFinalizedBundleByIndex(sdk.Context, uint64, uint64) (FinalizedBundle, bool)
	GetBundleProposal(sdk.Context, uint64) (bundlesTypes.BundleProposal, bool)
	GetFinalizedBundle(sdk.Context, uint64, uint64) (bundlesTypes.FinalizedBundle, bool)
	GetPaginatedFinalizedBundleQuery(sdk.Context, *query.PageRequest, uint64) ([]FinalizedBundle, *query.PageResponse, error)
	GetParams(sdk.Context) bundlesTypes.Params
	GetVoteDistribution(sdk.Context, uint64) bundlesTypes.VoteDistribution
}
