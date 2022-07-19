package v0_6_3

import (
	"fmt"

	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func migratePools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, pool := range registryKeeper.GetAllPool(ctx) {

		// set minimum total stake to 10k $KYVE
		pool.MinStake = 10_000_000_000_000

		// schedule upgrades for each runtime
		switch pool.Runtime {
		case "@kyve/evm":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "1.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/evm/releases/download/v1.3.5/kyve-linux.zip?checksum=f37eb5178890f74cdd6ba272cc783b25e59b3abc2fb13bd0c939736425e09123\",\"macos\":\"https://github.com/kyve-org/evm/releases/download/v1.3.5/kyve-macos.zip?checksum=ddccbe416c79b1e58a76813f14f83f41571c0c72cb3b913d5ca4d32d0fb4c8c9\"}\n",
			}
		case "@kyve/stacks":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/stacks/releases/download/v0.3.5/kyve-linux.zip?checksum=df4c6f66a05505a4b72d6f71d62aadbf6bcc740ae8ea2250bd9d3a2cc17d1c7f\",\"macos\":\"https://github.com/kyve-org/stacks/releases/download/v0.3.5/kyve-macos.zip?checksum=f1cc1d3f9f5f873685eea6b631a52a4efc0644dd2eb3408c942f31dcb4e54131\"}\n",
			}
		case "@kyve/bitcoin":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.3.5/kyve-linux.zip?checksum=815a035f737fb7c5d92017f697317e44159d30233b01548a08777b637db600e9\",\"macos\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.3.5/kyve-macos.zip?checksum=4e6879d25c5f05c7c9410f8ac65a4226d4793c0881b18541bd122cb41e72685a\"}\n",
			}
		case "@kyve/solana":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/solana/releases/download/v0.3.5/kyve-linux.zip?checksum=1faafe36451b0cba24dada4e655408d897ab4ca9e417ec7568abdcee7c095e01\",\"macos\":\"https://github.com/kyve-org/solana/releases/download/v0.3.5/kyve-macos.zip?checksum=bc518f685ef53f01ecfc8b036460871895718e8c0e4b0c4d3b87a35ea4207b9c\"}\n",
			}
		case "@kyve/zilliqa":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.3.5/kyve-linux.zip?checksum=06cedc13e6dd450ee3dc739503c55ed079341ee911da36aeaf5fea902186f332\",\"macos\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.3.5/kyve-macos.zip?checksum=e6914c721a29b9d7e38475eff495189884328e1492a72eab1b2f15ad0ebd9238\"}\n",
			}
		case "@kyve/near":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/near/releases/download/v0.3.5/kyve-linux.zip?checksum=b4614324223799113e30e9d382776ff6bd59ded12c73ceb60ab79ec0d6515ae1\",\"macos\":\"https://github.com/kyve-org/near/releases/download/v0.3.5/kyve-macos.zip?checksum=880de70dcfbc027cd0a454b51bcea742b5bea2bf3dc2f5da643d2a8f5830f748\"}\n",
			}
		case "@kyve/celo":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/celo/releases/download/v0.3.5/kyve-linux.zip?checksum=5e8059d8aeffeda66ee833e75ab2b7b3e627f508620aeb196f70cb2ba46a776b\",\"macos\":\"https://github.com/kyve-org/celo/releases/download/v0.3.5/kyve-macos.zip?checksum=eafd247d88062219070697b69e1f77f589a0e63b594807a0cc7221d01ff9b577\"}\n",
			}
		case "@kyve/cosmos":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/cosmos/releases/download/v0.3.5/kyve-linux.zip?checksum=86cdd1bad6461f00b831cb8976c104a2f230c4a7d8139a0e27dcb227c617f499\",\"macos\":\"https://github.com/kyve-org/cosmos/releases/download/v0.3.5/kyve-macos.zip?checksum=5c47e939a4c406a759d8c107f141e09ad73b78979978b711292e0833566b06b0\"}\n",
			}
		case "@kyve/substrate":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.5",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/substrate/releases/download/v0.3.5/kyve-linux.zip?checksum=0c42aba35d1b6d3bf00a2a42d2aa8b9bd8887fa8960e91d998648146b18d4721\",\"macos\":\"https://github.com/kyve-org/substrate/releases/download/v0.3.5/kyve-macos.zip?checksum=039f0400dc2cd86a5e1568ad7e0d18e08cbaa8d4a5be84b3a548889595237155\"}\n",
			}
		default:
			pool.UpgradePlan = &types.UpgradePlan{}
		}

		// add pool upgrade info
		pool.UpgradePlan.ScheduledAt = uint64(ctx.BlockTime().Unix())
		pool.UpgradePlan.Duration = 300 // 5min

		// save changes
		registryKeeper.SetPool(ctx, pool)
	}
}

func removeStringFromList(list []string, el string) []string {
	for i, other := range list {
		if other == el {
			return append(list[0:i], list[i+1:]...)
		}
	}
	return list
}

func deactivateStaker(pool *types.Pool, staker *types.Staker) {
	// make user an inactive staker
	pool.Stakers = removeStringFromList(pool.Stakers, staker.Account)
	pool.InactiveStakers = append(pool.InactiveStakers, staker.Account)
	pool.TotalStake -= staker.Amount
	pool.TotalInactiveStake += staker.Amount
	staker.Status = types.STAKER_STATUS_INACTIVE
}

func deactivateStakers(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, pool := range registryKeeper.GetAllPool(ctx) {

		poolStakers := make([]string, len(pool.Stakers))
		copy(poolStakers, pool.Stakers)

		for _, stakerAddress := range poolStakers {
			staker, _ := registryKeeper.GetStaker(ctx, stakerAddress, pool.Id)

			deactivateStaker(&pool, &staker)
			registryKeeper.SetStaker(ctx, staker)
		}

		pool.TotalStake = 0

		pool.BundleProposal = &types.BundleProposal{
			CreatedAt: uint64(ctx.BlockTime().Unix()),
		}

		registryKeeper.UpdateLowestStaker(ctx, &pool)
		registryKeeper.SetPool(ctx, pool)
	}
}

func correctPools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {

	// Check for duplicated
	for _, pool := range registryKeeper.GetAllPool(ctx) {

		stakersList := make([]string, 0)
		visitedStakers := make(map[string]bool, 0)
		var totalInactiveStake uint64 = 0
		for _, staker := range pool.InactiveStakers {
			if visitedStakers[staker] == false {

				inactiveStaker, found := registryKeeper.GetStaker(ctx, staker, pool.Id)
				if !found {
					fmt.Printf("Error: Staker does not exist. Skipping staker %s\n", staker)
					continue
				}

				stakersList = append(stakersList, staker)
				visitedStakers[staker] = true
				totalInactiveStake += inactiveStaker.Amount

			} else {
				fmt.Printf("Duplicate staker: %s\n", staker)
			}
		}
		pool.InactiveStakers = stakersList
		pool.TotalInactiveStake = totalInactiveStake
		pool.TotalStake = 0

		registryKeeper.SetPool(ctx, pool)
	}

}

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		deactivateStakers(registryKeeper, ctx)
		migratePools(registryKeeper, ctx)
		correctPools(registryKeeper, ctx)

		return vm, nil
	}
}
