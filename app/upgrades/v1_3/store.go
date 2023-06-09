package v1_3

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"

	// Consensus
	consensusTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	// Crisis
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	// PFM
	pfmTypes "github.com/strangelove-ventures/packet-forward-middleware/v7/router/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storeTypes.StoreUpgrades{
		Added: []string{
			consensusTypes.StoreKey, crisisTypes.StoreKey,
			pfmTypes.StoreKey,
		},
	}

	return upgradeTypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
