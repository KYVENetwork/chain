package v1_5

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storetypes.StoreUpgrades{
		//	TODO: packet forward middleware -> Delete
	}

	return upgradetypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
