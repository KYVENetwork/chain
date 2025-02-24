package v2_0

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	multicoinrewardstypes "github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

func CreateStoreLoader(upgradeHeight int64) baseapp.StoreLoader {
	storeUpgrades := storetypes.StoreUpgrades{
		Added:   []string{multicoinrewardstypes.ModuleName},
		Deleted: []string{"group"},
	}

	return upgradetypes.UpgradeStoreLoader(upgradeHeight, &storeUpgrades)
}
