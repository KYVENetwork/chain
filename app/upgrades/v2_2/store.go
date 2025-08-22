package v2_2

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storetypes.StoreUpgrades{
		Added: []string{""},
	}

	return upgradetypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
