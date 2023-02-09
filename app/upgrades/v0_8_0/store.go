package v080

import (
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storeTypes.StoreUpgrades{
		Added: []string{
			// kyve
			globalTypes.StoreKey,
			teamTypes.StoreKey,
		},
		Deleted: []string{
			"registry",
			"fees",
		},
	}

	return upgradeTypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
