package v1_3

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/log"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Pool
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	accountKeeper authKeeper.AccountKeeper,
	poolKeeper poolKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)

		CheckPoolAddresses(ctx, logger, accountKeeper, poolKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// CheckPoolAddresses ensures that each pool account exists post upgrade.
func CheckPoolAddresses(ctx sdk.Context, logger log.Logger, ak authKeeper.AccountKeeper, pk poolKeeper.Keeper) {
	pools := pk.GetAllPools(ctx)

	for _, pool := range pools {
		name := fmt.Sprintf("pool/%d", pool.Id)

		address := authTypes.NewModuleAddress(name)
		account := ak.GetAccount(ctx, address)

		if account == nil {
			// account doesn't exist, initialise a new module account.
			account = authTypes.NewEmptyModuleAccount(name)
			ak.SetAccount(ctx, account)
		} else { //nolint:staticcheck
			// TODO: account exists, adjust it to a module account.
		}

		logger.Info("successfully initialised pool account", "name", name)
	}
}
