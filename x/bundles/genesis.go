package bundles

import (
	"github.com/KYVENetwork/chain/x/bundles/keeper"
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	for _, entry := range genState.BundleProposalList {
		k.SetBundleProposal(ctx, entry)
	}

	for _, entry := range genState.FinalizedBundleList {
		k.SetFinalizedBundle(ctx, entry)
	}

	for _, entry := range genState.RoundRobinProgressList {
		k.SetRoundRobinProgress(ctx, entry)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.BundleProposalList = k.GetAllBundleProposals(ctx)

	genesis.FinalizedBundleList = k.GetAllFinalizedBundles(ctx)

	genesis.RoundRobinProgressList = k.GetAllRoundRobinProgress(ctx)

	return genesis
}
