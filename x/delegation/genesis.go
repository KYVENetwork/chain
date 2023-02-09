package delegation

import (
	"github.com/KYVENetwork/chain/x/delegation/keeper"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, delegator := range genState.DelegatorList {
		k.SetDelegator(ctx, delegator)
	}

	for _, entry := range genState.DelegationEntryList {
		k.SetDelegationEntry(ctx, entry)
	}

	for _, entry := range genState.DelegationDataList {
		k.SetDelegationData(ctx, entry)
	}

	for _, entry := range genState.DelegationSlashList {
		k.SetDelegationSlashEntry(ctx, entry)
	}

	for _, entry := range genState.UndelegationQueueEntryList {
		k.SetUndelegationQueueEntry(ctx, entry)
	}

	k.SetQueueState(ctx, genState.QueueStateUndelegation)

	for _, entry := range genState.RedelegationCooldownList {
		k.SetRedelegationCooldown(ctx, entry)
	}

	k.InitMemStore(ctx)
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
