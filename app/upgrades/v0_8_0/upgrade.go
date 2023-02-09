package v080

import (
	"fmt"

	types2 "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	// Global
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	pk poolKeeper.Keeper,
	sk stakersKeeper.Keeper,
	govStoreKey *types2.KVStoreKey,
	globalKeeper globalKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		cosmosVersionmap, err := mm.RunMigrations(ctx, configurator, vm)
		// kyve params
		ctx.Logger().Info("Init global module.")
		InitGlobalParams(ctx, globalKeeper)

		ctx.Logger().Info("Remove all stakers from disabled pool.")
		RemoveStakersFromDisabledPools(ctx, pk, sk)

		ParseProposals(ctx, govStoreKey)

		return cosmosVersionmap, err
	}
}

func RemoveStakersFromDisabledPools(ctx sdk.Context, pk poolKeeper.Keeper, sk stakersKeeper.Keeper) {
	for _, pool := range pk.GetAllPools(ctx) {
		if pool.Disabled {
			for _, staker := range sk.GetAllStakerAddressesOfPool(ctx, pool.Id) {
				sk.LeavePool(ctx, staker, pool.Id)
				ctx.Logger().Info(fmt.Sprintf("Remove %s from pool %s.\n", staker, pool.Name))
			}
		}
	}
}

// InitGlobalParams ...
func InitGlobalParams(ctx sdk.Context, globalKeeper globalKeeper.Keeper) {
	params := globalTypes.DefaultParams()

	minInitialRatio, _ := sdk.NewDecFromStr("0.25")
	params.MinInitialDepositRatio = minInitialRatio

	params.MinGasPrice = sdk.NewDec(1)

	burnRatio, _ := sdk.NewDecFromStr("0.2")
	params.BurnRatio = burnRatio

	globalKeeper.SetParams(ctx, params)
}

func ParseProposals(ctx sdk.Context, storeKey *types2.KVStoreKey) {
	store := ctx.KVStore(storeKey)
	for i := 1; i <= 168; i++ {
		store.Delete(types.ProposalKey(uint64(i)))
	}
	store.Delete(types.ProposalKey(uint64(172)))
}
