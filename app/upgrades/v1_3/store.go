package v1_3

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"

	// Consensus
	consensusTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	// Crisis
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storeTypes.StoreUpgrades{
		Added: []string{
			consensusTypes.StoreKey, crisisTypes.StoreKey,
		},
	}

	return upgradeTypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
