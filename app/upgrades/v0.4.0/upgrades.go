package v0_4_0

import (
	"fmt"

	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
)

func migratePools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
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
				Version:  "1.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/evm/releases/download/v1.1.0/kyve-linux.zip?checksum=174eb1a3b6c1161959c64daca44076e1571f76eb1e2339cf36e4b71cced3b471\",\"macos\":\"https://github.com/kyve-org/evm/releases/download/v1.1.0/kyve-macos.zip?checksum=494357b0e292dc97e8a3436778115ed8b2e96a4001622fe575462468c8b39f35\"}",
			}
		case "@kyve/stacks":
			pool.Versions = ">=0.1.0"
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/stacks/releases/download/v0.1.0/kyve-linux.zip?checksum=1c8b5ec983bccb70cbe435ded2edd7c2e65e9d12c7b032973f4adea6fa68281b\",\"macos\":\"https://github.com/kyve-org/stacks/releases/download/v0.1.0/kyve-macos.zip?checksum=acea13fd91281d545b500c8be5c02f92c0b296fee921516b072ed57b2607ca54\"}",
			}
		case "@kyve/bitcoin":
			pool.Versions = ">=0.1.0"
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.1.0/kyve-linux.zip?checksum=0fbeaa64b22ab2ba34b6fa4571d2ff1fe7e573885abba48c15d878c19853bd79\",\"macos\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.1.0/kyve-macos.zip?checksum=de901be709378f6bec81ed3c31d54fdfd67f26215e2614ffc14e80f6f98a29d1\"}",
			}
		case "@kyve/solana":
			pool.Versions = ">=0.1.0"
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/solana/releases/download/v0.1.0/kyve-linux.zip?checksum=0e94e710624fb42947f538892160179050f8e99b286ae71668b7f5a48c285af0\",\"macos\":\"https://github.com/kyve-org/solana/releases/download/v0.1.0/kyve-macos.zip?checksum=bf25242be4a99c7b6181b835af7d90969bcdcd9a0473e850aa0ba065d723b038\"}",
			}
		case "@kyve/zilliqa":
			pool.Versions = ">=0.1.0"
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.1.0/kyve-linux.zip?checksum=895b090bb2746a29fd4de168280b224353eb33bf7dd3f72fca8a60c250cfef2a\",\"macos\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.1.0/kyve-macos.zip?checksum=5ca1330c922bfefacf1b94c8933930789cd314624aa29a22aaa55d2ab64e7d83\"}",
			}
		case "@kyve/near":
			pool.Versions = ">=0.1.0"
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/near/releases/download/v0.1.0/kyve-linux.zip?checksum=7d310651e19aedfff9e3360b5e6108a6271314eedf52dd86e8945fd1bb1d3793\",\"macos\":\"https://github.com/kyve-org/near/releases/download/v0.1.0/kyve-macos.zip?checksum=34082976222fd7e8feac57eaf9af828508398e0688b7a3af9e0caf02777ab51d\"}",
			}
		case "@kyve/celo":
			pool.Versions = ">=0.1.0"
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/celo/releases/download/v0.1.0/kyve-linux.zip?checksum=686c3bd436a3322f6bf09d2b3df465186360981aabd4da25e8576f0a3a867d66\",\"macos\":\"https://github.com/kyve-org/celo/releases/download/v0.1.0/kyve-macos.zip?checksum=1ddf112ea5e161108bfd3678fa7bb9a94be5cda8795ffb1cf8be904116518f05\"}",
			}
		case "@kyve/cosmos":
			pool.Versions = ">=0.1.0"
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/cosmos/releases/download/v0.1.0/kyve-linux.zip?checksum=be96d9befd3a1084af6d9de2b49d81834edbf5888fb7dbdcbcaafbe22c1975a5\",\"macos\":\"https://github.com/kyve-org/cosmos/releases/download/v0.1.0/kyve-macos.zip?checksum=fcdd4c561096f7f7ccbf7eeb625c2317c5166b615fa208f92f73a2c704e1966a\"}",
			}
		default:
			pool.UpgradePlan = &types.UpgradePlan{}
		}

		pool.UpgradePlan.ScheduledAt = uint64(ctx.BlockTime().Unix())
		pool.UpgradePlan.Duration = 86400 // 24h

		// save changes
		registryKeeper.SetPool(ctx, pool)
	}
}

func createDelegatorIndex(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {

	fmt.Printf("%sCreating second delegator index\n", MigrationLoggerPrefix)

	// Set all delegators again to create the index
	delegators := registryKeeper.GetAllDelegator(ctx)
	for index, delegator := range delegators {

		registryKeeper.SetDelegator(ctx, delegator)

		if index%1000 == 0 {
			fmt.Printf("%sDelegators processed: %d\n", MigrationLoggerPrefix, index)
		}
	}

	fmt.Printf("%sFinished index creation\n", MigrationLoggerPrefix)
}

func createProposalIndex(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	fmt.Printf("%sCreating second proposal index\n", MigrationLoggerPrefix)

	// Set all delegators again to create the index
	proposals := registryKeeper.GetAllProposal(ctx)
	for index, proposal := range proposals {

		registryKeeper.SetProposal(ctx, proposal)

		if index%1000 == 0 {
			fmt.Printf("%sProposals processed: %d\n", MigrationLoggerPrefix, index)
		}
	}

	fmt.Printf("%sFinished index creation\n", MigrationLoggerPrefix)
}

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		migratePools(registryKeeper, ctx)

		createDelegatorIndex(registryKeeper, ctx)

		createProposalIndex(registryKeeper, ctx)

		// init param
		registryKeeper.ParamStore().Set(ctx, types.KeyMaxPoints, types.DefaultMaxPoints)

		// Return.
		return vm, nil
	}
}
