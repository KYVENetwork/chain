package funders

import (
	"github.com/KYVENetwork/chain/x/funders/keeper"
	"github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	for _, entry := range genState.FunderList {
		k.SetFunder(ctx, &entry)
	}
	for _, entry := range genState.FundingList {
		k.SetFunding(ctx, &entry)
	}
	for _, entry := range genState.FundingStateList {
		k.SetFundingState(ctx, &entry)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.FunderList = k.GetAllFunders(ctx)
	genesis.FundingList = k.GetAllFundings(ctx)
	genesis.FundingStateList = k.GetAllFundingStates(ctx)
	// this line is used by starport scaffolding # genesis/module/export
	return genesis
}
