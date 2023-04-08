package sender

import (
	"github.com/KYVENetwork/chain/x/oracle/sender/keeper"
	"github.com/KYVENetwork/chain/x/oracle/sender/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the x/oracle sender module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for seq, req := range genState.Requests {
		k.SetRequest(ctx, seq, req)
	}

	for seq, res := range genState.Responses {
		k.SetResponse(ctx, seq, res)
	}
}

// ExportGenesis returns the x/oracle sender module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	reqs := k.GetRequests(ctx)
	ress := k.GetResponses(ctx)

	return types.NewGenesisState(reqs, ress)
}
