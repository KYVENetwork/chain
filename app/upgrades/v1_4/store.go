package v1_4

import (
	storeTypes "cosmossdk.io/store/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/baseapp"

	// Consensus
	consensusTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	// Crisis
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	// Upgrade
	upgradeTypes "cosmossdk.io/x/upgrade/types"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storeTypes.StoreUpgrades{
		Added: []string{
			consensusTypes.StoreKey, crisisTypes.StoreKey, funderstypes.StoreKey,
		},
	}

	return upgradeTypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
