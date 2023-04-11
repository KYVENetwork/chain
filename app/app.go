package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	v11 "github.com/KYVENetwork/chain/app/upgrades/v1_1"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	authzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feeGrant "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	groupTypes "github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	group "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	icaController "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller"
	icaControllerKeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	icaControllerTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icaHost "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibcFee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee"
	ibcFeeKeeper "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/keeper"
	ibcFeeTypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	ibcTransfer "github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v6/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v6/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibcPortTypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	kyveDocs "github.com/KYVENetwork/chain/docs"
	// this line is used by starport scaffolding # stargate/app/moduleImport

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

const (
	AccountAddressPrefix = "kyve"
	Name                 = "kyve"
)

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.LegacyProposalHandler,
		upgradeclient.LegacyCancelProposalHandler,
		ibcclientclient.UpdateClientProposalHandler,
		ibcclientclient.UpgradeProposalHandler,
	)

	return govProposalHandlers
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(appModuleBasics...)

	// module account permissions
	maccPerms = moduleAccountPermissions
)

var _ servertypes.Application = (*App)(nil)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	Keepers

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// mm is the module manager
	mm *module.Manager

	configurator module.Configurator
}

// NewKYVEApp returns a reference to an initialized blockchain app
func NewKYVEApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, authzTypes.ModuleName, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey, govtypes.StoreKey,
		paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey,
		evidencetypes.StoreKey, ibcFeeTypes.StoreKey, ibcTransferTypes.StoreKey,
		icaControllerTypes.StoreKey, icaHostTypes.StoreKey, capabilitytypes.StoreKey,
		groupTypes.StoreKey,

		bundlesTypes.StoreKey,
		delegationTypes.StoreKey,
		globalTypes.StoreKey,
		poolTypes.StoreKey,
		queryTypes.StoreKey,
		stakersTypes.StoreKey,
		teamTypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, bundlesTypes.MemStoreKey, delegationTypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedIBCTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icaControllerTypes.SubModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	// this line is used by starport scaffolding # stargate/app/scopedKeeper

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32PrefixAccAddr,
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzTypes.ModuleName],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		app.BlockedModuleAccountAddrs(),
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		&app.StakersKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	app.DistributionKeeper = distrkeeper.NewKeeper(
		appCodec,
		keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)

	groupConfig := groupTypes.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[groupTypes.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupConfig,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// ... other modules keepers
	app.GlobalKeeper = *globalKeeper.NewKeeper(appCodec, keys[globalTypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String())

	app.TeamKeeper = *teamKeeper.NewKeeper(appCodec, keys[teamTypes.StoreKey], app.AccountKeeper, app.BankKeeper)

	app.PoolKeeper = *poolKeeper.NewKeeper(
		appCodec,
		keys[poolTypes.StoreKey],
		memKeys[poolTypes.MemStoreKey],

		authtypes.NewModuleAddress(govtypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.UpgradeKeeper,
	)

	app.StakersKeeper = *stakersKeeper.NewKeeper(
		appCodec,
		keys[stakersTypes.StoreKey],
		memKeys[stakersTypes.MemStoreKey],

		authtypes.NewModuleAddress(govtypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.UpgradeKeeper,
	)

	app.DelegationKeeper = *delegationKeeper.NewKeeper(
		appCodec,
		keys[delegationTypes.StoreKey],
		memKeys[delegationTypes.MemStoreKey],

		authtypes.NewModuleAddress(govtypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.UpgradeKeeper,
		app.StakersKeeper,
	)

	stakersKeeper.SetDelegationKeeper(&app.StakersKeeper, app.DelegationKeeper)
	poolKeeper.SetStakersKeeper(&app.PoolKeeper, app.StakersKeeper)

	app.BundlesKeeper = *bundlesKeeper.NewKeeper(
		appCodec,
		keys[bundlesTypes.StoreKey],
		memKeys[bundlesTypes.MemStoreKey],

		authtypes.NewModuleAddress(govtypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.StakersKeeper,
		app.DelegationKeeper,
	)

	// Create IBC Keepers
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	app.IBCFeeKeeper = ibcFeeKeeper.NewKeeper(
		appCodec, keys[ibcFeeTypes.StoreKey],
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	app.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibcTransferTypes.StoreKey],
		app.GetSubspace(ibcTransferTypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedIBCTransferKeeper,
	)

	app.ICAControllerKeeper = icaControllerKeeper.NewKeeper(
		appCodec, keys[icaControllerTypes.StoreKey],
		app.GetSubspace(icaControllerTypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icaHostTypes.StoreKey],
		app.GetSubspace(icaHostTypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		keys[evidencetypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(app.DistributionKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	govConfig := govtypes.DefaultConfig()
	app.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		&app.StakersKeeper,
		govRouter,
		app.MsgServiceRouter(),
		govConfig,
	)

	app.QueryKeeper = *queryKeeper.NewKeeper(
		appCodec,
		keys[queryTypes.StoreKey],
		keys[queryTypes.MemStoreKey],
		app.GetSubspace(queryTypes.ModuleName),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.StakersKeeper,
		app.DelegationKeeper,
		app.BundlesKeeper,
		app.GlobalKeeper,
		app.GovKeeper,
		app.TeamKeeper,
	)
	// this line is used by starport scaffolding # stargate/app/keeperDefinition

	// Create static IBC router, add transfer route, then set and seal it
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

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		// Cosmos SDK
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		authz.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		distribution.NewAppModule(appCodec, app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feeGrant.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		group.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, minttypes.DefaultInflationCalculationFn),
		params.NewAppModule(app.ParamsKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),

		// IBC
		ibc.NewAppModule(app.IBCKeeper),
		ibcFee.NewAppModule(app.IBCFeeKeeper),
		ibcTransfer.NewAppModule(app.IBCTransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),

		// KYVE
		bundles.NewAppModule(appCodec, app.BundlesKeeper, app.AccountKeeper, app.BankKeeper),
		delegation.NewAppModule(appCodec, app.DelegationKeeper, app.AccountKeeper, app.BankKeeper),
		global.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.GlobalKeeper, app.UpgradeKeeper),
		pool.NewAppModule(appCodec, app.PoolKeeper, app.AccountKeeper, app.BankKeeper),
		query.NewAppModule(appCodec, app.QueryKeeper, app.AccountKeeper, app.BankKeeper),
		stakers.NewAppModule(appCodec, app.StakersKeeper, app.AccountKeeper, app.BankKeeper),
		team.NewAppModule(appCodec, app.BankKeeper, app.MintKeeper, app.TeamKeeper, app.UpgradeKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		// NOTE: x/team must be run before x/distribution and after x/mint.
		teamTypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibcFeeTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		authzTypes.ModuleName,
		feegrant.ModuleName,
		groupTypes.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,

		// this line is used by starport scaffolding # stargate/app/beginBlockers
		poolTypes.ModuleName,
		stakersTypes.ModuleName,
		delegationTypes.ModuleName,
		bundlesTypes.ModuleName,
		queryTypes.ModuleName,
		globalTypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibcFeeTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authzTypes.ModuleName,
		feegrant.ModuleName,
		groupTypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		// this line is used by starport scaffolding # stargate/app/endBlockers
		poolTypes.ModuleName,
		stakersTypes.ModuleName,
		delegationTypes.ModuleName,
		bundlesTypes.ModuleName,
		queryTypes.ModuleName,
		globalTypes.ModuleName,
		teamTypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		ibcFeeTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		evidencetypes.ModuleName,
		authzTypes.ModuleName,
		feegrant.ModuleName,
		groupTypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		// this line is used by starport scaffolding # stargate/app/initGenesis
		poolTypes.ModuleName,
		stakersTypes.ModuleName,
		delegationTypes.ModuleName,
		bundlesTypes.ModuleName,
		queryTypes.ModuleName,
		globalTypes.ModuleName,
		teamTypes.ModuleName,
	)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	var err error

	anteHandler, err := NewAnteHandler(
		app.AccountKeeper,
		app.BankKeeper,
		app.FeeGrantKeeper,
		app.GlobalKeeper,
		app.GovKeeper,
		app.IBCKeeper,
		app.StakingKeeper,
		ante.DefaultSigVerificationGasConsumer,
		encodingConfig.TxConfig.SignModeHandler(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	postHandler, err := NewPostHandler(
		app.BankKeeper,
		app.FeeGrantKeeper,
		app.GlobalKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create PostHandler: %s", err))
	}

	app.SetAnteHandler(anteHandler)
	app.SetPostHandler(postHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.UpgradeKeeper.SetUpgradeHandler(
		v11.UpgradeName,
		v11.CreateUpgradeHandler(
			app.mm,
			app.configurator,
			appCodec,
			keys[stakersTypes.StoreKey],
			keys[bundlesTypes.StoreKey],
			keys[delegationTypes.StoreKey],
			app.AccountKeeper,
		),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == v11.UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		app.SetStoreLoader(v11.CreateStoreLoader(upgradeInfo.Height))
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedIBCTransferKeeper = scopedIBCTransferKeeper
	app.ScopedICAControllerKeeper = scopedICAControllerKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper

	return app
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// GetBaseApp returns the base app of the application
func (app App) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, _ config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API
	apiSvr.Router.Handle("/swagger.yml", http.FileServer(http.FS(kyveDocs.Swagger)))
	apiSvr.Router.HandleFunc("/", kyveDocs.Handler(Name, "/swagger.yml"))
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// SimulationManager implements the SimulationApp interface.
// NOTE: We simply return nil as we don't use the simulation manager anywhere.
func (app *App) SimulationManager() *module.SimulationManager { return nil }
