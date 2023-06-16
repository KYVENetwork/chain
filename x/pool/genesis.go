package pool

import (
	"github.com/KYVENetwork/chain/x/pool/keeper"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the pool module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, elem := range genState.PoolList {
		k.SetPool(ctx, elem)
	}

	k.SetPoolCount(ctx, genState.PoolCount)
}

// ExportGenesis returns the pool module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Params = k.GetParams(ctx)

	genesis.PoolList = k.GetAllPools(ctx)
	genesis.PoolCount = k.GetPoolCount(ctx)

	return genesis
}
