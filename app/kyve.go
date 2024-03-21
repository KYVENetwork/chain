package app

import (
	"cosmossdk.io/core/appmodule"
	bundlesModule "github.com/KYVENetwork/chain/x/bundles"
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	delegationModule "github.com/KYVENetwork/chain/x/delegation"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	fundersModule "github.com/KYVENetwork/chain/x/funders"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	globalModule "github.com/KYVENetwork/chain/x/global"
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	poolModule "github.com/KYVENetwork/chain/x/pool"
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	queryModule "github.com/KYVENetwork/chain/x/query"
	queryTypes "github.com/KYVENetwork/chain/x/query/types"
	stakersModule "github.com/KYVENetwork/chain/x/stakers"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	teamModule "github.com/KYVENetwork/chain/x/team"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (app *App) registerKyveModules() {
	app.GlobalKeeper = *globalKeeper.NewKeeper(
		app.appCodec,
		app.GetKey(globalTypes.StoreKey),
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	app.TeamKeeper = *teamKeeper.NewKeeper(
		app.appCodec,
		app.GetKey(teamTypes.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.MintKeeper,
		*app.UpgradeKeeper,
	)

	app.PoolKeeper = *poolKeeper.NewKeeper(
		app.appCodec,
		app.GetKey(poolTypes.StoreKey),
		app.GetMemKey(poolTypes.MemStoreKey),

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.MintKeeper,
		app.UpgradeKeeper,
		app.TeamKeeper,
	)

	app.StakersKeeper = *stakersKeeper.NewKeeper(
		app.appCodec,
		app.GetKey(stakersTypes.StoreKey),
		app.GetMemKey(stakersTypes.MemStoreKey),

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.UpgradeKeeper,
	)

	app.DelegationKeeper = *delegationKeeper.NewKeeper(
		app.appCodec,
		app.GetKey(delegationTypes.StoreKey),
		app.GetMemKey(delegationTypes.MemStoreKey),

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.UpgradeKeeper,
		app.StakersKeeper,
	)

	app.FundersKeeper = *fundersKeeper.NewKeeper(
		app.appCodec,
		app.GetKey(fundersTypes.StoreKey),
		app.GetMemKey(fundersTypes.MemStoreKey),

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.PoolKeeper,
		app.UpgradeKeeper,
	)

	stakersKeeper.SetDelegationKeeper(&app.StakersKeeper, app.DelegationKeeper)
	poolKeeper.SetStakersKeeper(&app.PoolKeeper, app.StakersKeeper)
	poolKeeper.SetFundersKeeper(&app.PoolKeeper, app.FundersKeeper)

	app.BundlesKeeper = *bundlesKeeper.NewKeeper(
		app.appCodec,
		app.GetKey(bundlesTypes.StoreKey),
		app.GetMemKey(bundlesTypes.MemStoreKey),

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.StakersKeeper,
		app.DelegationKeeper,
		app.FundersKeeper,
	)
}

func RegisterKyveModules(registry cdctypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		bundlesTypes.ModuleName:    bundlesModule.AppModule{},
		delegationTypes.ModuleName: delegationModule.AppModule{},
		globalTypes.ModuleName:     globalModule.AppModule{},
		poolTypes.ModuleName:       poolModule.AppModule{},
		stakersTypes.ModuleName:    stakersModule.AppModule{},
		teamTypes.ModuleName:       teamModule.AppModule{},
		fundersTypes.ModuleName:    fundersModule.AppModule{},
		queryTypes.ModuleName:      queryModule.AppModule{},
	}
	for _, module := range modules {
		if mod, ok := module.(interface {
			RegisterInterfaces(registry cdctypes.InterfaceRegistry)
		}); ok {
			mod.RegisterInterfaces(registry)
		}
	}

	return modules
}
