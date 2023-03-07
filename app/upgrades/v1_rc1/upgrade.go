package v1rc1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Bundles
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Capability
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	// ICA Controller
	v6 "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/migrations/v6"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	bundlesKeeper bundlesKeeper.Keeper,
	cdc codec.BinaryCodec,
	capabilityStoreKey *storeTypes.KVStoreKey,
	capabilityKeeper *capabilityKeeper.Keeper,
	moduleName string,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if err := v6.MigrateICS27ChannelCapability(ctx, cdc, capabilityStoreKey, capabilityKeeper, moduleName); err != nil {
			return nil, err
		}

		ReinitialiseBundlesParams(ctx, bundlesKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// ReinitialiseBundlesParams ...
func ReinitialiseBundlesParams(ctx sdk.Context, bundlesKeeper bundlesKeeper.Keeper) {
	params := bundlesTypes.DefaultParams()
	bundlesKeeper.SetParams(ctx, params)
}
