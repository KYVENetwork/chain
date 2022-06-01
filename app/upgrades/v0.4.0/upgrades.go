package v0_4_0

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
)

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		for _, pool := range registryKeeper.GetAllPool(ctx) {
			// set max_bundle_size
			pool.MaxBundleSize = 100

			// init protocol
			pool.Protocol = &types.Protocol{}

			// schedule upgrades for each runtime
			switch pool.Runtime {
			case "@kyve/evm":
				pool.Versions = ">=1.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "1.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			case "@kyve/stacks":
				pool.Versions = ">=0.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "0.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			case "@kyve/bitcoin":
				pool.Versions = ">=0.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "0.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			case "@kyve/solana":
				pool.Versions = ">=0.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "0.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			case "@kyve/zilliqa":
				pool.Versions = ">=0.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "0.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			case "@kyve/near":
				pool.Versions = ">=0.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "0.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			case "@kyve/celo":
				pool.Versions = ">=0.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "0.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			case "@kyve/cosmos":
				pool.Versions = ">=0.1.0"
				pool.UpgradePlan = &types.UpgradePlan{
					Version: "0.1.0",
					Binaries: "{\"linux\":\"todo\",\"macos\":\"todo\"}",
				}
			default:
				pool.UpgradePlan = &types.UpgradePlan{}
			}

			pool.UpgradePlan.ScheduledAt = uint64(ctx.BlockTime().Unix())
			pool.UpgradePlan.Duration = 1200 // 20 min

			// save changes
			registryKeeper.SetPool(ctx, pool)
		}

		// init param
		registryKeeper.ParamStore().Set(ctx, types.KeyMaxPoints, types.DefaultMaxPoints)

		// Return.
		return vm, nil
	}
}
