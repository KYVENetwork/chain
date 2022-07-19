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
				Version:  "1.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/evm/releases/download/v1.3.4/kyve-linux.zip?checksum=223c4099afcf98c9736167952ba45b13553fc53e4f21d26c7e593e0e94c926c6\",\"macos\":\"https://github.com/kyve-org/evm/releases/download/v1.3.4/kyve-macos.zip?checksum=63138cbdfa3dc1ec38981f110464240e2b981fe28bff20cc957fe8386c8b6650\"}\n",
			}
		case "@kyve/stacks":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/stacks/releases/download/v0.3.4/kyve-linux.zip?checksum=479467a8fbd0e0c48c4e659622947e90bd5745a53f0d5338f4b34c63de6b16fc\",\"macos\":\"https://github.com/kyve-org/stacks/releases/download/v0.3.4/kyve-macos.zip?checksum=64a6dd6fa70d6750eefc02759af0991f1f768e6cd072f01e3f716e0d8a897173\"}\n",
			}
		case "@kyve/bitcoin":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.3.4/kyve-linux.zip?checksum=b9255671b7d74c726db971a74049312692c246edbebf369d255d09d000fa8ca7\",\"macos\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.3.4/kyve-macos.zip?checksum=89dce78e5f7f8f531538adf333522b04810f3e36e65b2cfe55598a26d9b1a53f\"}\n",
			}
		case "@kyve/solana":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/solana/releases/download/v0.3.4/kyve-linux.zip?checksum=06a2fd9271aaacaed89308d858c68d7aea0bb0b60a5ad0e6997f3065acd91cf0\",\"macos\":\"https://github.com/kyve-org/solana/releases/download/v0.3.4/kyve-macos.zip?checksum=c82415f9513c4893df5a578db57671768761a6d4374f4791156c427925a41240\"}\n",
			}
		case "@kyve/zilliqa":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.3.4/kyve-linux.zip?checksum=18a8f70a6c64151c2e8428e0d2ed95e13c822638180478976f624a32602ab1f8\",\"macos\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.3.4/kyve-macos.zip?checksum=c617b9e48a2a6eb2b8e3bfe1132d08f55ccf254c5900dc8211ca301a544fc44e\"}\n",
			}
		case "@kyve/near":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/near/releases/download/v0.3.4/kyve-linux.zip?checksum=bbf036c6c368ca5d11439ed73cd272c161798c4fd52ff096007417d5fe6a3e5f\",\"macos\":\"https://github.com/kyve-org/near/releases/download/v0.3.4/kyve-macos.zip?checksum=52eb5e59863967dbfb21f7f16eab9da61472923f931df2ef9112e6f27e7b223f\"}\n",
			}
		case "@kyve/celo":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/celo/releases/download/v0.3.4/kyve-linux.zip?checksum=bb7760bde70b93d9cad7ec0b9d484fad6b4b7b814b1ed8b49278f0b46d281f5f\",\"macos\":\"https://github.com/kyve-org/celo/releases/download/v0.3.4/kyve-macos.zip?checksum=f67a284770c6b51900f29dc02e08ec2ddec83e7b7fe517e32b094529d693d48e\"}\n",
			}
		case "@kyve/cosmos":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/cosmos/releases/download/v0.3.4/kyve-linux.zip?checksum=449edd89d2a74b64aca86915c1fb41f13ad4e068ff3247f41f2405cd360452ad\",\"macos\":\"https://github.com/kyve-org/cosmos/releases/download/v0.3.4/kyve-macos.zip?checksum=9c064d9dc2cc96b5665ae46bdddcc11b6c676a6bec33143a59e4a81462977e85\"}\n",
			}
		case "@kyve/substrate":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.4",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/substrate/releases/download/v0.3.4/kyve-linux.zip?checksum=2891626d214f4175ed213c39b414e4f16c17a7d4122020371280bde55cde9077\",\"macos\":\"https://github.com/kyve-org/substrate/releases/download/v0.3.4/kyve-macos.zip?checksum=d117a737d07e1aa2de3eda9f93252f4fc6869f3b86b84d8cf45413bf0aad15c9\"}\n",
			}
		default:
			pool.UpgradePlan = &types.UpgradePlan{}
		}

		// add pool upgrade info
		pool.UpgradePlan.ScheduledAt = uint64(ctx.BlockTime().Unix())
		pool.UpgradePlan.Duration = 1800 // 30min

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
