package keeper

import (
	"sort"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) addPendingRewards(ctx sdk.Context, address string, rewards sdk.Coins) {
	queueIndex := k.getNextQueueSlot(ctx, types.QUEUE_IDENTIFIER_MULTI_COIN_REWARDS)

	pendingEntry := types.MultiCoinPendingRewardsEntry{
		Index:        queueIndex,
		Address:      address,
		Rewards:      rewards,
		CreationDate: ctx.BlockTime().Unix(),
	}

	k.SetMultiCoinPendingRewardsEntry(ctx, pendingEntry)
}

// ProcessPendingRewardsQueue ...
func (k Keeper) ProcessPendingRewardsQueue(ctx sdk.Context) error {
	collectedRewards := sdk.NewCoins()

	k.processQueue(ctx, types.QUEUE_IDENTIFIER_MULTI_COIN_REWARDS, func(index uint64) bool {
		// Get queue entry in question
		queueEntry, found := k.GetMultiCoinPendingRewardsEntry(ctx, index)
		if !found {
			// continue with the next entry
			return true
		}

		if queueEntry.CreationDate+int64(k.GetMultiCoinDistributionPendingTime(ctx)) <= ctx.BlockTime().Unix() {
			k.RemoveMultiCoinPendingRewardsEntry(ctx, &queueEntry)
			collectedRewards = collectedRewards.Add(queueEntry.Rewards...)
			return true
		}

		return false
	})

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.MultiCoinRewardsRedistributionAccountName, collectedRewards); err != nil {
		return err
	}

	return nil
}

// DistributeNonClaimedRewards takes all non-claimed rewards which have exceeding the claim time
// and re-distribute them to the pools according to the redistribution-policy.
func (k Keeper) DistributeNonClaimedRewards(ctx sdk.Context) error {
	policy, err := k.MultiCoinDistributionPolicy.Get(ctx)
	if err != nil {
		return err
	}

	distributionMap, err := types.ParseAndNormalizeMultiCoinDistributionMap(policy)
	if err != nil {
		return err
	}

	// Store rewards for all pools. There could be multiple rules which re-direct coins to the same pool
	type PoolRewardsBasket struct {
		account sdk.AccAddress
		rewards sdk.Coins
		poolId  uint64
	}
	poolRewardBaskets := make(map[uint64]PoolRewardsBasket)

	// Get all rewards
	rewards := k.bankKeeper.GetAllBalances(ctx, k.accountKeeper.GetModuleAddress(types.MultiCoinRewardsRedistributionAccountName))

	// Iterate every coin denom and re-distribute accordingly
	for _, coin := range rewards {
		weightMap, ok := distributionMap[coin.Denom]
		if !ok {
			// Coin not registered in coin map, it will stay in the module account
			continue
		}

		// weight-map contains for every denom the destination pools together with a weight.
		for _, weight := range weightMap {
			// Check if pool is already in temporary map
			poolBasket, ok := poolRewardBaskets[weight.PoolId]
			// If pool is not registered in map yet, initialize new pool basket.
			if !ok {
				// if not, get pool from id
				pool, err := k.poolKeeper.GetPoolWithError(ctx, weight.PoolId)
				if err != nil {
					// if pool does not exist, the distribution map is incorrect, cancel the process
					return err
				}

				poolBasket.poolId = pool.Id
				poolBasket.account = pool.GetPoolAccount()
				poolBasket.rewards = sdk.NewCoins()
				poolRewardBaskets[weight.PoolId] = poolBasket
			}

			// Truncate int ensures that there are never more tokens distributed than available
			poolReward := sdk.NewCoin(coin.Denom, weight.NormalizedWeight.MulInt(rewards.AmountOf(coin.Denom)).TruncateInt())

			// Add reward to pool
			poolBasket.rewards = poolBasket.rewards.Add(poolReward)

			// Update map
			poolRewardBaskets[weight.PoolId] = poolBasket
		}
	}

	// Sort PoolRewards for determinism
	accountList := make([]PoolRewardsBasket, 0)
	for _, account := range poolRewardBaskets {
		accountList = append(accountList, account)
	}
	sort.Slice(accountList, func(i, j int) bool { return accountList[i].poolId < accountList[j].poolId })

	// Redistribute all tokens
	for _, account := range accountList {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.MultiCoinRewardsRedistributionAccountName, account.account, account.rewards); err != nil {
			return err
		}
	}

	return nil
}
