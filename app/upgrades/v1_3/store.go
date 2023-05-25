package v1_3

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"

	// PFM
	pfmTypes "github.com/strangelove-ventures/packet-forward-middleware/v6/router/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storeTypes.StoreUpgrades{
		Added: []string{
			pfmTypes.StoreKey,
		},
	}

	return upgradeTypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
