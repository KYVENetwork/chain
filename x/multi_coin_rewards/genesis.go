package multi_coin_rewards

import (
	"github.com/KYVENetwork/chain/x/multi_coin_rewards/keeper"
	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, entry := range genState.MultiCoinPendingRewardsEntries {
		k.SetMultiCoinPendingRewardsEntry(ctx, entry)
	}

	for _, entry := range genState.MultiCoinEnabled {
		err := k.MultiCoinRewardsEnabled.Set(ctx, sdk.MustAccAddressFromBech32(entry))
		if err != nil {
			panic(err)
		}
	}

	k.SetQueueState(ctx, types.QUEUE_IDENTIFIER_MULTI_COIN_REWARDS, genState.QueueStatePendingRewards)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.MultiCoinPendingRewardsEntries = k.GetAllMultiCoinPendingRewardsEntries(ctx)

	genesis.QueueStatePendingRewards = k.GetQueueState(ctx, types.QUEUE_IDENTIFIER_MULTI_COIN_REWARDS)

	policy, err := k.MultiCoinDistributionPolicy.Get(ctx)
	if err != nil {
		panic(err)
	}
	genesis.MultiCoinDistributionPolicy = &policy

	genesis.MultiCoinEnabled = k.GetAllEnabledMultiCoinAddresses(ctx)

	return genesis
}
