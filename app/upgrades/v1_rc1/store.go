package v1rc1

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"

	// ICA Controller
	icaControllerTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storeTypes.StoreUpgrades{
		Added: []string{
			icaControllerTypes.StoreKey,
		},
	}

	return upgradeTypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
