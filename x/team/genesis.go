package team

import (
	"github.com/KYVENetwork/chain/x/team/keeper"
	"github.com/KYVENetwork/chain/x/team/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the team module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetAuthority(ctx, genState.Authority)

	for _, elem := range genState.AccountList {
		k.SetTeamVestingAccount(ctx, elem)
	}

	k.SetTeamVestingAccountCount(ctx, genState.AccountCount)
}

// ExportGenesis returns the team module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Authority = k.GetAuthority(ctx)
	genesis.AccountList = k.GetTeamVestingAccounts(ctx)
	genesis.AccountCount = k.GetTeamVestingAccountCount(ctx)

	return genesis
}
