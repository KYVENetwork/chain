package v0_6_0

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
}

func migratePools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	for _, pool := range registryKeeper.GetAllPool(ctx) {
		// schedule upgrades for each runtime
		switch pool.Runtime {
		case "@kyve/evm":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "1.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/evm/releases/download/v1.2.1/kyve-linux.zip?checksum=add69311973079efd092cc960c857f60546b9e7193888f6421428338a2b24bb0\",\"macos\":\"https://cdn.discordapp.com/attachments/889827445132374036/996778129030926436/kyve-macos.zip?checksum=90f7b40d2db1185ddc83d289726c2be95f883dc9cc03159f5be5f8957e242017\"}",
			}
		case "@kyve/stacks":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/stacks/releases/download/v0.2.1/kyve-linux.zip?checksum=690f479472b535ae842d098ac4e1f6e2a183f6c8ee2f539eb9faa3dd0b0a4b7b\",\"macos\":\"https://github.com/kyve-org/stacks/releases/download/v0.2.1/kyve-macos.zip?checksum=d213a61798969c338a6e86ddcff9b77f2f894f17c84896c0e1df0918280b7774\"}",
			}
		case "@kyve/bitcoin":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.2.1/kyve-linux.zip?checksum=ced62e0808c4e6fae78daef28052405bc60491328902d416ae3edd717e89ab84\",\"macos\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.2.1/kyve-macos.zip?checksum=ff91ff97729e61ba621a82945e89aaf8b80bfd36346b80c3971d48a11fbc9db9\"}",
			}
		case "@kyve/solana":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/solana/releases/download/v0.2.1/kyve-linux.zip?checksum=19903b6a5bd621219025a90d059dfcd3a0b522344365eaaa20feadb64ca80ab9\",\"macos\":\"https://github.com/kyve-org/solana/releases/download/v0.2.1/kyve-macos.zip?checksum=4a0b403e7f1196dadd2cafec11529f0c4036ba6048b6fe1a86dadd6835f84680\"}",
			}
		case "@kyve/zilliqa":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.2.1/kyve-linux.zip?checksum=5d32003201c7b22f580f9e6552139a32bc2775386f8f9ccbda363fe0602b44a8\",\"macos\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.2.1/kyve-macos.zip?checksum=dd35df7a88c1517c3fc355612ad7bb1eda596b6b25d1ad7d108f26ff3d2e8084\"}",
			}
		case "@kyve/near":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/near/releases/download/v0.2.1/kyve-linux.zip?checksum=1b5697e65277247e6c7a819618d517c0c38ff5ae8869c0a54252a31179dabab5\",\"macos\":\"https://github.com/kyve-org/near/releases/download/v0.2.1/kyve-macos.zip?checksum=fd4ef500849c625ecfcc80afbd8d3b8a4a2c74386e6759aa8bd17f7bde54a28d\"}",
			}
		case "@kyve/celo":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/celo/releases/download/v0.2.1/kyve-linux.zip?checksum=3c4d819183d7adc4fc614b5d78f5665ee38f81d78f4991d34af333b04b225594\",\"macos\":\"https://github.com/kyve-org/celo/releases/download/v0.2.1/kyve-macos.zip?checksum=713707e4b19efe9c1570f086e1f699cd4d1217c921a5dbc50aed55aa9c14a71a\"}",
			}
		case "@kyve/cosmos":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/cosmos/releases/download/v0.2.1/kyve-linux.zip?checksum=1bbfd864fb787eb4e360c2d40f5d49966e0004448eb0baa00354ffb269b61615\",\"macos\":\"https://github.com/kyve-org/cosmos/releases/download/v0.2.1/kyve-macos.zip?checksum=1e01d7074b462110cf8ff5263c32c5f1369c04a85cbcb6d2b99eb3d09efbf64f\"}",
			}
		case "@kyve/substrate":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.3.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/substrate/releases/download/v0.1.1/kyve-linux.zip?checksum=d4120fa124f8eeb8f4f346b952cfdad419dc65083e26c9906b7138c123b4614f\",\"macos\":\"https://github.com/kyve-org/substrate/releases/download/v0.1.1/kyve-macos.zip?checksum=1e29724c45f7c8fd81c7840a7d3d66136861e213545920c84c4d47e2c3e619c7\"}",
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
