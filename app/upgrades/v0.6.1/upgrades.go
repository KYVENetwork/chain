package v0_6_1

import (
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func migrateStakers(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, staker := range registryKeeper.GetAllStaker(ctx) {
		staker.Status = types.STAKER_STATUS_ACTIVE
		registryKeeper.GetStaker(ctx, staker.Account, staker.PoolId)
	}
}

func createRedelegationParameters(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	registryKeeper.ParamStore().Set(ctx, types.KeyRedelegationCooldown, types.DefaultRedelegationCooldown)

	registryKeeper.ParamStore().Set(ctx, types.KeyRedelegationMaxAmount, types.DefaultRedelegationMaxAmount)

	registryKeeper.ParamStore().Set(ctx, types.KeyCommissionChangeTime, types.DefaultCommissionChangeTime)
}

func migratePools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, pool := range registryKeeper.GetAllPool(ctx) {
		// schedule upgrades for each runtime
		switch pool.Runtime {
		case "@kyve/evm":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "1.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/evm/releases/download/v1.3.0/kyve-linux.zip?checksum=3959be874794a2f2ea8b56a5b63952880930c995460b0500230db19de516078a\",\"macos\":\"https://github.com/kyve-org/evm/releases/download/v1.3.0/kyve-macos.zip?checksum=713ebe89def5cfb4a79895d88923c82f18e3fa41d54083805ff8a7a4651ac0c0\"}\n",
			}
		case "@kyve/stacks":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/stacks/releases/download/v0.3.0/kyve-linux.zip?checksum=bb2c8708f7cc27321bc4f9397d8c0616a88650c846fd6a0005b91ac2dbacb2e1\",\"macos\":\"https://github.com/kyve-org/stacks/releases/download/v0.3.0/kyve-macos.zip?checksum=29ea76cbb618416883f671fd27790ec41c648dfa528739748701529d72d8f2ed\"}\n",
			}
		case "@kyve/bitcoin":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.3.0/kyve-linux.zip?checksum=b5591251ed539b8b37291e843fc861ac61694a1e4e2780cb934a78ce5ee854e2\",\"macos\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.3.0/kyve-macos.zip?checksum=5001c14e0b797823a93ddb538e1bd22e10bf0380f6b565c4ec25342c558b6bd3\"}\n",
			}
		case "@kyve/solana":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/solana/releases/download/v0.3.0/kyve-linux.zip?checksum=d28f6ed953e5b1c1770315bedec0573f5b60465b11f96fe493a6768cf12fb7cf\",\"macos\":\"https://github.com/kyve-org/solana/releases/download/v0.3.0/kyve-macos.zip?checksum=ca87985b382598c321474eaae7ffa2ad4abad0c1cb1e3ba0f0f0c92ac10fabd2\"}\n",
			}
		case "@kyve/zilliqa":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.3.0/kyve-linux.zip?checksum=00c765def29f9c223790a383793d660c06d6e26ef3c8f09cf002f30429dc24d3\",\"macos\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.3.0/kyve-macos.zip?checksum=051886f4940dcfdb5de4f8d40153bf49c2ead6e97c1b6ec5d2c3e013df0314e8\"}\n",
			}
		case "@kyve/near":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/near/releases/download/v0.3.0/kyve-linux.zip?checksum=28a252381b53a4f89100fc16edda78b7b3a2a8581d3d6bf556f84f2e03a95edb\",\"macos\":\"https://github.com/kyve-org/near/releases/download/v0.3.0/kyve-macos.zip?checksum=9a47f0dff125aa526d04a1d204f7c92ff38411880e5ad32e24d5f89717c7b659\"}\n",
			}
		case "@kyve/celo":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/celo/releases/download/v0.3.0/kyve-linux.zip?checksum=1fdcccb9c296157174d9588cf18b5b65d3dc9b4ddd907e21a21dbb843e5bc6fe\",\"macos\":\"https://github.com/kyve-org/celo/releases/download/v0.3.0/kyve-macos.zip?checksum=def7225c32e7e09b47f71be1038087add02613da86ce697a6e048b4cc29427c0\"}\n",
			}
		case "@kyve/cosmos":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/cosmos/releases/download/v0.3.0/kyve-linux.zip?checksum=f7a30d56a9058ba7c183c2a926cb6d875bf08a02a55d123d15c65ef25d8250b2\",\"macos\":\"https://github.com/kyve-org/cosmos/releases/download/v0.3.0/kyve-macos.zip?checksum=899ad7152d2a8b64f84130a14a1cf83629437a9219e3060ad66d97f7c8460966\"}\n",
			}
		case "@kyve/substrate":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/substrate/releases/download/v0.3.0/kyve-linux.zip?checksum=ecc4d9ee3e80e4a4347bfcd03b9a21ae12f98000974cb1a73231237ffb189d97\",\"macos\":\"https://github.com/kyve-org/substrate/releases/download/v0.3.0/kyve-macos.zip?checksum=39a4c121b6e6040f5628d09f98ac1e92292673355392687eed468cc6ab81ba8c\"}\n",
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

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		registryKeeper.UpgradeHelperV060MigrateSecondIndex(ctx)

		migrateStakers(registryKeeper, ctx)

		createRedelegationParameters(registryKeeper, ctx)

		migratePools(registryKeeper, ctx)

		return vm, nil
	}
}
