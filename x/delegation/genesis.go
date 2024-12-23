package delegation

import (
	"github.com/KYVENetwork/chain/x/delegation/keeper"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Params = k.GetParams(ctx)

	genesis.DelegatorList = k.GetAllDelegators(ctx)

	genesis.DelegationEntryList = k.GetAllDelegationEntries(ctx)

	genesis.DelegationDataList = k.GetAllDelegationData(ctx)

	genesis.DelegationSlashList = k.GetAllDelegationSlashEntries(ctx)

	genesis.UndelegationQueueEntryList = k.GetAllUnbondingDelegationQueueEntries(ctx)

	genesis.QueueStateUndelegation = k.GetQueueState(ctx)

	genesis.RedelegationCooldownList = k.GetAllRedelegationCooldownEntries(ctx)

	return genesis
}
