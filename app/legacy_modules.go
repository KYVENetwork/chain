package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bundles
	"github.com/KYVENetwork/chain/x/bundles"
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Delegation
	"github.com/KYVENetwork/chain/x/delegation"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Global
	"github.com/KYVENetwork/chain/x/global"
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// Gov
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// IBC Client
	ibcClientSolomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibcClientTendermint "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	// IBC Core
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcPortTypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibcKeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	// IBC Fee
	ibcFee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcFeeKeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcFeeTypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	// IBC Transfer
	ibcTransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	// ICA
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	// ICA Controller
	icaController "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icaControllerKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icaControllerTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	// ICA Host
	icaHost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icaHostKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	// Pool
	"github.com/KYVENetwork/chain/x/pool"
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	// Query
	"github.com/KYVENetwork/chain/x/query"
	queryKeeper "github.com/KYVENetwork/chain/x/query/keeper"
	queryTypes "github.com/KYVENetwork/chain/x/query/types"
	// Stakers
	"github.com/KYVENetwork/chain/x/stakers"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	// Team
	"github.com/KYVENetwork/chain/x/team"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
)

var (
	keys = sdk.NewKVStoreKeys(
		ibcTypes.StoreKey, ibcFeeTypes.StoreKey, ibcTransferTypes.StoreKey,
		icaControllerTypes.StoreKey, icaHostTypes.StoreKey,

		bundlesTypes.StoreKey, delegationTypes.StoreKey, globalTypes.StoreKey,
		poolTypes.StoreKey, queryTypes.StoreKey, stakersTypes.StoreKey,
		teamTypes.StoreKey,
	)

	memKeys = sdk.NewMemoryStoreKeys(
		bundlesTypes.MemStoreKey, delegationTypes.MemStoreKey,
		poolTypes.MemStoreKey, queryTypes.MemStoreKey, stakersTypes.MemStoreKey,
	)

	subspaces = []string{
		ibcTypes.ModuleName, ibcFeeTypes.ModuleName, ibcTransferTypes.ModuleName,
		icaControllerTypes.SubModuleName, icaHostTypes.SubModuleName,
	}
)

func (app *KYVEApp) RegisterLegacyModules() {
	// Register the additional keys.
	app.MountKVStores(keys)
	app.MountMemoryStores(memKeys)

	// Initialise the additional param subspaces.
	for _, subspace := range subspaces {
		app.ParamsKeeper.Subspace(subspace)
	}

	// Keeper: IBC
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcTypes.ModuleName)
	app.IBCKeeper = ibcKeeper.NewKeeper(
		app.appCodec,
		keys[ibcTypes.StoreKey],
		app.GetSubspace(ibcTypes.ModuleName),

		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)
	app.ScopedIBCKeeper = scopedIBCKeeper

	// Keeper: IBC Fee
	app.IBCFeeKeeper = ibcFeeKeeper.NewKeeper(
		app.appCodec,
		keys[ibcFeeTypes.ModuleName],

		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,

		app.AccountKeeper,
		app.BankKeeper,
	)

	// Keeper: IBC Transfer
	scopedIBCTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	app.IBCTransferKeeper = ibcTransferKeeper.NewKeeper(
		app.appCodec,
		keys[ibcTransferTypes.StoreKey],
		app.GetSubspace(ibcTransferTypes.ModuleName),

		app.IBCFeeKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,

		app.AccountKeeper,
		app.BankKeeper,
		scopedIBCTransferKeeper,
	)
	app.ScopedIBCTransferKeeper = scopedIBCTransferKeeper

	// Keeper: ICA Controller
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icaControllerTypes.SubModuleName)
	app.ICAControllerKeeper = icaControllerKeeper.NewKeeper(
		app.appCodec,
		keys[icaControllerTypes.StoreKey],
		app.GetSubspace(icaControllerTypes.SubModuleName),

		app.IBCFeeKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,

		scopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)
	app.ScopedICAControllerKeeper = scopedICAControllerKeeper

	// Keeper: ICA Host
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	app.ICAHostKeeper = icaHostKeeper.NewKeeper(
		app.appCodec,
		keys[icaHostTypes.StoreKey],
		app.GetSubspace(icaHostTypes.SubModuleName),

		app.IBCFeeKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,

		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)
	app.ScopedICAHostKeeper = scopedICAHostKeeper

	// IBC: Create a router.
	var ibcTransferStack ibcPortTypes.IBCModule
	ibcTransferStack = ibcTransfer.NewIBCModule(app.IBCTransferKeeper)
	ibcTransferStack = ibcFee.NewIBCMiddleware(ibcTransferStack, app.IBCFeeKeeper)

	var icaControllerStack ibcPortTypes.IBCModule
	icaControllerStack = icaController.NewIBCMiddleware(icaControllerStack, app.ICAControllerKeeper)
	icaControllerStack = ibcFee.NewIBCMiddleware(icaControllerStack, app.IBCFeeKeeper)

	var icaHostStack ibcPortTypes.IBCModule
	icaHostStack = icaHost.NewIBCModule(app.ICAHostKeeper)
	icaHostStack = ibcFee.NewIBCMiddleware(icaHostStack, app.IBCFeeKeeper)

	ibcRouter := ibcPortTypes.NewRouter()
	ibcRouter.AddRoute(ibcTransferTypes.ModuleName, ibcTransferStack).
		AddRoute(icaControllerTypes.SubModuleName, icaControllerStack).
		AddRoute(icaHostTypes.SubModuleName, icaHostStack)
	app.IBCKeeper.SetRouter(ibcRouter)

	// Keeper: Global
	app.GlobalKeeper = *globalKeeper.NewKeeper(
		app.appCodec,
		keys[globalTypes.StoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)
	// Keeper: Pool
	app.PoolKeeper = *poolKeeper.NewKeeper(
		app.appCodec,
		keys[poolTypes.StoreKey],
		memKeys[poolTypes.MemStoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.UpgradeKeeper,
	)
	// Keeper: Stakers
	app.StakersKeeper = *stakersKeeper.NewKeeper(
		app.appCodec,
		keys[stakersTypes.StoreKey],
		memKeys[stakersTypes.MemStoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.UpgradeKeeper,
	)
	// Keeper: Delegation
	app.DelegationKeeper = *delegationKeeper.NewKeeper(
		app.appCodec,
		keys[delegationTypes.StoreKey],
		memKeys[delegationTypes.MemStoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.UpgradeKeeper,
		app.StakersKeeper,
	)
	// Keeper: Bundles
	app.BundlesKeeper = *bundlesKeeper.NewKeeper(
		app.appCodec,
		keys[bundlesTypes.StoreKey],
		memKeys[bundlesTypes.MemStoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.StakersKeeper,
		app.DelegationKeeper,
	)
	// Keeper: Team
	app.TeamKeeper = *teamKeeper.NewKeeper(
		app.appCodec,
		keys[teamTypes.StoreKey],

		app.AccountKeeper,
		app.BankKeeper,
	)
	// Keeper: Query
	app.QueryKeeper = *queryKeeper.NewKeeper(
		app.appCodec,
		keys[queryTypes.StoreKey],
		memKeys[queryTypes.MemStoreKey],

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.StakersKeeper,
		app.DelegationKeeper,
		app.BundlesKeeper,
		app.GlobalKeeper,
		app.GovKeeper,
	)

	app.StakersKeeper.SetDelegationKeeper(app.DelegationKeeper)
	app.PoolKeeper.SetStakersKeeper(app.StakersKeeper)
	app.GovKeeper.SetProtocolStakingKeeper(app.StakersKeeper)

	// Register modules and interfaces/services.
	legacyModules := []module.AppModule{
		ibc.NewAppModule(app.IBCKeeper),
		ibcFee.NewAppModule(app.IBCFeeKeeper),
		ibcTransfer.NewAppModule(app.IBCTransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),

		bundles.NewAppModule(app.appCodec, app.BundlesKeeper),
		delegation.NewAppModule(app.appCodec, app.DelegationKeeper, app.AccountKeeper, app.BankKeeper),
		global.NewAppModule(app.appCodec, app.AccountKeeper, app.BankKeeper, app.GlobalKeeper, app.UpgradeKeeper),
		pool.NewAppModule(app.appCodec, app.PoolKeeper, app.AccountKeeper, app.BankKeeper),
		query.NewAppModule(app.appCodec, app.QueryKeeper),
		stakers.NewAppModule(app.appCodec, app.StakersKeeper, app.AccountKeeper, app.BankKeeper),
		team.NewAppModule(app.appCodec, app.BankKeeper, app.MintKeeper, app.TeamKeeper, app.UpgradeKeeper),
	}
	if err := app.RegisterModules(legacyModules...); err != nil {
		panic(err)
	}

	for _, m := range legacyModules {
		if s, ok := m.(module.HasServices); ok {
			s.RegisterServices(app.Configurator())
		}
	}

	ibcClientSolomachine.AppModuleBasic{}.RegisterInterfaces(app.interfaceRegistry)
	ibcClientTendermint.AppModuleBasic{}.RegisterInterfaces(app.interfaceRegistry)
}
