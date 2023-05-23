package app

import (
	_ "embed"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"cosmossdk.io/depinject"

	kyveDocs "github.com/KYVENetwork/chain/docs"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata_pulsar"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side effects
	// Authz
	authzKeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"
	// Bank
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Bundles
	"github.com/KYVENetwork/chain/x/bundles"
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	// Capability
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	// Consensus
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	// Crisis
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	// Delegation
	"github.com/KYVENetwork/chain/x/delegation"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	// Distribution
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	// Evidence
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidenceKeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	// FeeGrant
	feeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feeGrant "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	// GenUtil
	genUtil "github.com/cosmos/cosmos-sdk/x/genutil"
	genUtilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	// Global
	"github.com/KYVENetwork/chain/x/global"
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	// Governance
	"github.com/cosmos/cosmos-sdk/x/gov"
	govClient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	// Group
	groupKeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	group "github.com/cosmos/cosmos-sdk/x/group/module"
	// IBC Client
	ibcClientSolomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibcClientTendermint "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	// IBC Core
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcClient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcKeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	// IBC Fee
	ibcFee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcFeeKeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	// IBC Transfer
	ibcTransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	// ICA
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	// ICA Controller
	icaControllerKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	// ICA Host
	icaHostKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	// Mint
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	// NFT -- TODO(@john): Do we want to include this?
	nftKeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	nft "github.com/cosmos/cosmos-sdk/x/nft/module"
	// Params
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	// Pool
	"github.com/KYVENetwork/chain/x/pool"
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	// Query
	"github.com/KYVENetwork/chain/x/query"
	queryKeeper "github.com/KYVENetwork/chain/x/query/keeper"
	// Slashing
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	// Stakers
	"github.com/KYVENetwork/chain/x/stakers"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	// Staking
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	// Team
	"github.com/KYVENetwork/chain/x/team"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	// Upgrade
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	// Vesting
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		// Cosmos SDK Modules
		auth.AppModuleBasic{},
		genUtil.NewAppModuleBasic(genUtilTypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govClient.ProposalHandler{
				paramsClient.ProposalHandler,
				upgradeClient.LegacyProposalHandler,
				upgradeClient.LegacyCancelProposalHandler,
				ibcClient.UpdateClientProposalHandler,
				ibcClient.UpgradeProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feeGrant.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		authz.AppModuleBasic{},
		group.AppModuleBasic{},
		vesting.AppModuleBasic{},
		nft.AppModuleBasic{},
		consensus.AppModuleBasic{},

		// IBC Modules
		ibc.AppModuleBasic{},
		ibcClientSolomachine.AppModuleBasic{},
		ibcClientTendermint.AppModuleBasic{},
		ibcFee.AppModuleBasic{},
		ibcTransfer.AppModuleBasic{},
		ica.AppModuleBasic{},

		// KYVE Modules
		bundles.AppModuleBasic{},
		delegation.AppModuleBasic{},
		global.AppModuleBasic{},
		pool.AppModuleBasic{},
		query.AppModuleBasic{},
		stakers.AppModuleBasic{},
		team.AppModuleBasic{},
	)
)

var (
	_ runtime.AppI            = (*KYVEApp)(nil)
	_ serverTypes.Application = (*KYVEApp)(nil)
)

// KYVEApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type KYVEApp struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codecTypes.InterfaceRegistry

	// Cosmos SDK Keepers
	AccountKeeper         authKeeper.AccountKeeper
	BankKeeper            bankKeeper.Keeper
	CapabilityKeeper      *capabilityKeeper.Keeper
	StakingKeeper         *stakingKeeper.Keeper
	SlashingKeeper        slashingKeeper.Keeper
	MintKeeper            mintKeeper.Keeper
	DistributionKeeper    distributionKeeper.Keeper
	GovKeeper             *govKeeper.Keeper
	CrisisKeeper          *crisisKeeper.Keeper
	UpgradeKeeper         *upgradeKeeper.Keeper
	ParamsKeeper          paramsKeeper.Keeper
	AuthzKeeper           authzKeeper.Keeper
	EvidenceKeeper        evidenceKeeper.Keeper
	FeeGrantKeeper        feeGrantKeeper.Keeper
	GroupKeeper           groupKeeper.Keeper
	NFTKeeper             nftKeeper.Keeper
	ConsensusParamsKeeper consensusKeeper.Keeper

	// IBC Keepers
	IBCKeeper           *ibcKeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper        ibcFeeKeeper.Keeper
	IBCTransferKeeper   ibcTransferKeeper.Keeper
	ICAControllerKeeper icaControllerKeeper.Keeper
	ICAHostKeeper       icaHostKeeper.Keeper

	// KYVE Keepers
	BundlesKeeper    bundlesKeeper.Keeper
	DelegationKeeper delegationKeeper.Keeper
	GlobalKeeper     globalKeeper.Keeper
	PoolKeeper       poolKeeper.Keeper
	QueryKeeper      queryKeeper.Keeper
	StakersKeeper    stakersKeeper.Keeper
	TeamKeeper       teamKeeper.Keeper

	// Scoped Keepers
	ScopedIBCKeeper           capabilityKeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilityKeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilityKeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilityKeeper.ScopedKeeper
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".simapp")
}

// NewKYVEApp returns a reference to an initialized KYVEApp.
func NewKYVEApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts serverTypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *KYVEApp {
	var (
		app        = &KYVEApp{}
		appBuilder *runtime.AppBuilder

		// merge the AppConfig and other configuration in one config
		appConfig = depinject.Configs(
			AppConfig,
			depinject.Supply(
				// supply the application options
				appOpts,

				// ADVANCED CONFIGURATION

				//
				// AUTH
				//
				// For providing a custom function required in auth to generate custom account types
				// add it below. By default the auth module uses simulation.RandomGenesisAccounts.
				//
				// authtypes.RandomGenesisAccountsFn(simulation.RandomGenesisAccounts),

				// For providing a custom a base account type add it below.
				// By default the auth module uses authtypes.ProtoBaseAccount().
				//
				// func() authtypes.AccountI { return authtypes.ProtoBaseAccount() },

				//
				// MINT
				//

				// For providing a custom inflation function for x/mint add here your
				// custom function that implements the minttypes.InflationCalculationFn
				// interface.
			),
		)
	)

	if err := depinject.Inject(appConfig,
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.CapabilityKeeper,
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.MintKeeper,
		&app.DistributionKeeper,
		&app.GovKeeper,
		&app.CrisisKeeper,
		&app.UpgradeKeeper,
		&app.ParamsKeeper,
		&app.AuthzKeeper,
		&app.EvidenceKeeper,
		&app.FeeGrantKeeper,
		&app.GroupKeeper,
		&app.NFTKeeper,
		&app.ConsensusParamsKeeper,
		// KYVE Modules
		//&app.BundlesKeeper,
		//&app.DelegationKeeper,
		//&app.GlobalKeeper,
		//&app.PoolKeeper,
		//&app.QueryKeeper,
		//&app.StakersKeeper,
		//&app.TeamKeeper,
	); err != nil {
		panic(err)
	}

	// Below we could construct and set an application specific mempool and
	// ABCI 1.0 PrepareProposal and ProcessProposal handlers. These defaults are
	// already set in the SDK's BaseApp, this shows an example of how to override
	// them.
	//
	// Example:
	//
	// app.App = appBuilder.Build(...)
	// nonceMempool := mempool.NewSenderNonceMempool()
	// abciPropHandler := NewDefaultProposalHandler(nonceMempool, app.App.BaseApp)
	//
	// app.App.BaseApp.SetMempool(nonceMempool)
	// app.App.BaseApp.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// app.App.BaseApp.SetProcessProposal(abciPropHandler.ProcessProposalHandler())
	//
	// Alternatively, you can construct BaseApp options, append those to
	// baseAppOptions and pass them to the appBuilder.
	//
	// Example:
	//
	// prepareOpt = func(app *baseapp.BaseApp) {
	// 	abciPropHandler := baseapp.NewDefaultProposalHandler(nonceMempool, app)
	// 	app.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// }
	// baseAppOptions = append(baseAppOptions, prepareOpt)

	app.App = appBuilder.Build(logger, db, traceStore, baseAppOptions...)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(app.App.BaseApp, appOpts, app.appCodec, logger, app.kvStoreKeys()); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	/****  Module Options ****/

	app.RegisterLegacyModules()

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)

	// RegisterUpgradeHandlers is used for registering any on-chain upgrades.
	// app.RegisterUpgradeHandlers()

	// add test gRPC service for testing gRPC queries in isolation
	testdata_pulsar.RegisterQueryServer(app.GRPCQueryRouter(), testdata_pulsar.QueryImpl{})

	// A custom InitChainer can be set if extra pre-init-genesis logic is required.
	// By default, when using app wiring enabled module, this is not required.
	// For instance, the upgrade module will set automatically the module version map in its init genesis thanks to app wiring.
	// However, when registering a module manually (i.e. that does not support app wiring), the module version map
	// must be set manually as follow. The upgrade module will de-duplicate the module version map.
	//
	// app.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	// 	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap())
	// 	return app.App.InitChainer(ctx, req)
	// })

	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}

// Name returns the name of the App
func (app *KYVEApp) Name() string { return app.BaseApp.Name() }

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *KYVEApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *KYVEApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns SimApp's InterfaceRegistry
func (app *KYVEApp) InterfaceRegistry() codecTypes.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns SimApp's TxConfig
func (app *KYVEApp) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *KYVEApp) GetKey(storeKey string) *storeTypes.KVStoreKey {
	sk := app.UnsafeFindStoreKey(storeKey)
	kvStoreKey, ok := sk.(*storeTypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

func (app *KYVEApp) kvStoreKeys() map[string]*storeTypes.KVStoreKey {
	keys := make(map[string]*storeTypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storeTypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *KYVEApp) GetSubspace(moduleName string) paramsTypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *KYVEApp) SimulationManager() *module.SimulationManager {
	return nil
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *KYVEApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	apiSvr.Router.Handle("/swagger.yml", http.FileServer(http.FS(kyveDocs.Swagger)))
	apiSvr.Router.HandleFunc("/", kyveDocs.Handler("", "/swagger.yml"))
	//if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
	//	panic(err)
	//}
}

// GetMaccPerms returns a copy of the module account permissions
//
// NOTE: This is solely to be used for testing purposes.
func GetMaccPerms() map[string][]string {
	dup := make(map[string][]string)
	for _, perms := range moduleAccPerms {
		dup[perms.Account] = perms.Permissions
	}

	return dup
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	result := make(map[string]bool)

	if len(blockAccAddrs) > 0 {
		for _, addr := range blockAccAddrs {
			result[addr] = true
		}
	} else {
		for addr := range GetMaccPerms() {
			result[addr] = true
		}
	}

	return result
}
