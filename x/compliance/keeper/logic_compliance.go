package keeper

import (
	"sort"

	"github.com/KYVENetwork/chain/x/compliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) addPendingComplianceRewards(ctx sdk.Context, address string, rewards sdk.Coins) {
	queueIndex := k.getNextQueueSlot(ctx, types.QUEUE_IDENTIFIER_MULTI_COIN_REWARDS)

	compliancePendingEntry := types.MultiCoinPendingRewardsEntry{
		Index:        queueIndex,
		Address:      address,
		Rewards:      rewards,
		CreationDate: ctx.BlockTime().Unix(),
	}

	k.SetMultiCoinPendingRewardsEntry(ctx, compliancePendingEntry)
}

// ProcessComplianceQueue ...
func (k Keeper) ProcessComplianceQueue(ctx sdk.Context) {
	collectedRewards := sdk.NewCoins()

	k.processQueue(ctx, types.QUEUE_IDENTIFIER_MULTI_COIN_REWARDS, func(index uint64) bool {
		// Get queue entry in question
		queueEntry, found := k.GetMultiCoinPendingRewardsEntry(ctx, index)
		if !found {
			// continue with the next entry
			return true
		}

		if queueEntry.CreationDate+int64(k.GetMultiCoinRefundPendingTime(ctx)) <= ctx.BlockTime().Unix() {
			k.RemoveMultiCoinPendingRewardsEntry(ctx, &queueEntry)
			collectedRewards = collectedRewards.Add(queueEntry.Rewards...)
			return true
		}

		return false
	})

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.MultiCoinRewardsRedistributionAccountName, collectedRewards); err != nil {
		panic(err)
	}
}

func (k Keeper) DistributeNonClaimedRewards(ctx sdk.Context) error {
	policy, err := k.MultiCoinRefundPolicy.Get(ctx)
	if err != nil {
		return err
	}

	complianceMap, err := types.ParseMultiCoinComplianceMap(policy)
	if err != nil {
		return err
	}

	type PoolRewards struct {
		account sdk.AccAddress
		rewards sdk.Coins
		poolId  uint64
	}

	rewards := k.bankKeeper.GetAllBalances(ctx, k.accountKeeper.GetModuleAddress(types.MultiCoinRewardsRedistributionAccountName))

	poolRewards := make(map[uint64]PoolRewards)

	for _, coin := range rewards {
		weightMap, ok := complianceMap[coin.Denom]
		if !ok {
			// Coin not registered in coin map, it will stay in the module account
			continue
		}

		for _, weight := range weightMap {
			// Check if pool is already in temporary map
			accounts, ok := poolRewards[weight.PoolId]
			if !ok {
				// if not, get pool from id
				pool, err := k.poolKeeper.GetPoolWithError(ctx, weight.PoolId)
				if err != nil {
					// if pool does not exist, the compliance map is incorrect, cancel the process
					return err
				}
				accounts.poolId = pool.Id
				accounts.account = pool.GetPoolAccount()
			}

			poolReward := sdk.NewCoin(coin.Denom, weight.NormalizedWeight.MulInt(rewards.AmountOf(coin.Denom)).TruncateInt())

			// Subtract reward from available rewards
			rewards = rewards.Sub(poolReward)
			// Add reward to pool
			accounts.rewards = accounts.rewards.Add(poolReward)
		}
	}

	accountList := make([]PoolRewards, 0)
	for _, account := range poolRewards {
		accountList = append(accountList, account)
	}
	sort.Slice(accountList, func(i, j int) bool { return accountList[i].poolId < accountList[j].poolId })

	for _, account := range accountList {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.MultiCoinRewardsRedistributionAccountName, account.account.String(), account.rewards)
		if err != nil {
			return err
		}
	}

	return nil
}
