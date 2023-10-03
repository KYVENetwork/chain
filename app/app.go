package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	v1p4 "github.com/KYVENetwork/chain/app/upgrades/v1_4"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	cmtOs "github.com/cometbft/cometbft/libs/os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cast"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	// Authz
	authzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	authzKeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"
	// Bank
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// Bundles
	"github.com/KYVENetwork/chain/x/bundles"
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Capability
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	// Consensus
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	// Crisis
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	// Delegation
	"github.com/KYVENetwork/chain/x/delegation"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Distribution
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	// Evidence
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidenceKeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	// FeeGrant
	feeGrantTypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	feeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feeGrant "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	// GenUtil
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genUtilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	// Global
	"github.com/KYVENetwork/chain/x/global"
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// Gov
	"github.com/cosmos/cosmos-sdk/x/gov"
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	// Group
	groupTypes "github.com/cosmos/cosmos-sdk/x/group"
	groupKeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	group "github.com/cosmos/cosmos-sdk/x/group/module"
	// IBC Core
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcClientHandler "github.com/cosmos/ibc-go/v7/modules/core/02-client" // TODO
	ibcClientTypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcPortTypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcExported "github.com/cosmos/ibc-go/v7/modules/core/exported"
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
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	// ICA Controller
	icaController "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icaControllerKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icaControllerTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	// ICA Host
	icaHost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icaHostKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	// Mint
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	// Params
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsProposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	// Pool
	"github.com/KYVENetwork/chain/x/pool"
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	// PFM
	pfm "github.com/strangelove-ventures/packet-forward-middleware/v7/router"
	pfmKeeper "github.com/strangelove-ventures/packet-forward-middleware/v7/router/keeper"
	pfmTypes "github.com/strangelove-ventures/packet-forward-middleware/v7/router/types"
	// Query
	"github.com/KYVENetwork/chain/x/query"
	queryKeeper "github.com/KYVENetwork/chain/x/query/keeper"
	queryTypes "github.com/KYVENetwork/chain/x/query/types"
	// Slashing
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	// Stakers
	"github.com/KYVENetwork/chain/x/stakers"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	// Staking
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// Team
	"github.com/KYVENetwork/chain/x/team"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
	// Upgrade
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	// Funders
	"github.com/KYVENetwork/chain/x/funders"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
)

const (
	AccountAddressPrefix = "kyve"
	Name                 = "kyve"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(appModuleBasics...)
)

var (
	// TODO(@john): Ask if this is needed for a "v1" app.
	_ runtime.AppI            = (*App)(nil)
	_ serverTypes.Application = (*App)(nil)
)

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
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	keys    map[string]*storeTypes.KVStoreKey
	tkeys   map[string]*storeTypes.TransientStoreKey
	memKeys map[string]*storeTypes.MemoryStoreKey

	Keepers
	mm           *module.Manager
	configurator module.Configurator
}

// NewKYVEApp returns a reference to an initialized blockchain app
func NewKYVEApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts serverTypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	encodingConfig := MakeEncodingConfig()

	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	// Below we could construct and set an application specific mempool and
	// ABCI 1.0 PrepareProposal and ProcessProposal handlers. These defaults are
	// already set in the SDK's BaseApp, this shows an example of how to override
	// them.
	//
	// Example:
	//
	// bApp := baseapp.NewBaseApp(...)
	// nonceMempool := mempool.NewSenderNonceMempool()
	// abciPropHandler := NewDefaultProposalHandler(nonceMempool, bApp)
	//
	// bApp.SetMempool(nonceMempool)
	// bApp.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// bApp.SetProcessProposal(abciPropHandler.ProcessProposalHandler())
	//
	// Alternatively, you can construct BaseApp options, append those to
	// baseAppOptions and pass them to NewBaseApp.
	//
	// Example:
	//
	// prepareOpt = func(app *baseapp.BaseApp) {
	// 	abciPropHandler := baseapp.NewDefaultProposalHandler(nonceMempool, app)
	// 	app.SetPrepareProposal(abciPropHandler.PrepareProposalHandler())
	// }
	// baseAppOptions = append(baseAppOptions, prepareOpt)

	bApp := baseapp.NewBaseApp(Name, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := sdk.NewKVStoreKeys(
		authTypes.StoreKey,
		authzTypes.ModuleName,
		bankTypes.StoreKey,
		capabilityTypes.StoreKey,
		consensusTypes.StoreKey,
		crisisTypes.StoreKey,
		distributionTypes.StoreKey,
		evidenceTypes.StoreKey,
		feeGrantTypes.StoreKey,
		govTypes.StoreKey,
		groupTypes.StoreKey,
		mintTypes.StoreKey,
		paramsTypes.StoreKey,
		slashingTypes.StoreKey,
		stakingTypes.StoreKey,
		upgradeTypes.StoreKey,

		ibcExported.StoreKey,
		ibcFeeTypes.StoreKey,
		ibcTransferTypes.StoreKey,
		icaControllerTypes.StoreKey,
		icaHostTypes.StoreKey,
		pfmTypes.StoreKey,

		bundlesTypes.StoreKey,
		delegationTypes.StoreKey,
		globalTypes.StoreKey,
		poolTypes.StoreKey,
		queryTypes.StoreKey,
		stakersTypes.StoreKey,
		teamTypes.StoreKey,
		fundersTypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(
		capabilityTypes.MemStoreKey,

		bundlesTypes.MemStoreKey, delegationTypes.MemStoreKey,
	)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, keys); err != nil {
		logger.Error("failed to load state streaming", "err", err)
		os.Exit(1)
	}

	app := &App{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		legacyAmino,
		keys[paramsTypes.StoreKey],
		tkeys[paramsTypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	app.ConsensusKeeper = consensusKeeper.NewKeeper(
		appCodec,
		keys[consensusTypes.StoreKey],
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)
	bApp.SetParamStore(&app.ConsensusKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilityKeeper.NewKeeper(
		appCodec,
		keys[capabilityTypes.StoreKey],
		memKeys[capabilityTypes.MemStoreKey],
	)

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcExported.ModuleName)
	scopedIBCTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icaControllerTypes.SubModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)

	// TODO(@john): Seal x/capability keeper.

	// add keepers
	app.AccountKeeper = authKeeper.NewAccountKeeper(
		appCodec,
		keys[authTypes.StoreKey],
		authTypes.ProtoBaseAccount,
		moduleAccountPermissions,
		sdk.Bech32MainPrefix,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	app.AuthzKeeper = authzKeeper.NewKeeper(
		keys[authzTypes.ModuleName],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	app.BankKeeper = bankKeeper.NewBaseKeeper(
		appCodec,
		keys[bankTypes.StoreKey],
		app.AccountKeeper,
		app.BlockedModuleAccountAddrs(),
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	app.StakingKeeper = stakingKeeper.NewKeeper(
		appCodec,
		keys[stakingTypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	app.MintKeeper = mintKeeper.NewKeeper(
		appCodec,
		keys[mintTypes.StoreKey],
		app.StakingKeeper,
		&app.StakersKeeper, // TODO(@john)
		app.AccountKeeper,
		app.BankKeeper,
		authTypes.FeeCollectorName,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	app.DistributionKeeper = distributionKeeper.NewKeeper(
		appCodec,
		keys[distributionTypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		authTypes.FeeCollectorName,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingKeeper.NewKeeper(
		appCodec,
		legacyAmino,
		keys[slashingTypes.StoreKey],
		app.StakingKeeper,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	app.CrisisKeeper = crisisKeeper.NewKeeper(
		appCodec,
		keys[crisisTypes.StoreKey],
		invCheckPeriod,
		app.BankKeeper,
		authTypes.FeeCollectorName,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	app.FeeGrantKeeper = feeGrantKeeper.NewKeeper(
		appCodec,
		keys[feeGrantTypes.StoreKey],
		app.AccountKeeper,
	)

	app.FeeGrantKeeper = feeGrantKeeper.NewKeeper(
		appCodec,
		keys[feeGrantTypes.StoreKey],
		app.AccountKeeper,
	)

	app.GroupKeeper = groupKeeper.NewKeeper(
		keys[groupTypes.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
		groupTypes.DefaultConfig(),
	)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	// set the governance module account as the authority for conducting upgrades
	app.UpgradeKeeper = upgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradeTypes.StoreKey],
		appCodec,
		homePath,
		app.BaseApp,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingTypes.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// ... other modules keepers
	app.GlobalKeeper = *globalKeeper.NewKeeper(appCodec, keys[globalTypes.StoreKey], authTypes.NewModuleAddress(govTypes.ModuleName).String())

	app.TeamKeeper = *teamKeeper.NewKeeper(appCodec, keys[teamTypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.MintKeeper, *app.UpgradeKeeper)

	app.PoolKeeper = *poolKeeper.NewKeeper(
		appCodec,
		keys[poolTypes.StoreKey],
		memKeys[poolTypes.MemStoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.MintKeeper,
		app.UpgradeKeeper,
		app.TeamKeeper,
	)

	app.StakersKeeper = *stakersKeeper.NewKeeper(
		appCodec,
		keys[stakersTypes.StoreKey],
		memKeys[stakersTypes.MemStoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

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

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.UpgradeKeeper,
		app.StakersKeeper,
	)

	app.FundersKeeper = *fundersKeeper.NewKeeper(
		appCodec,
		keys[fundersTypes.StoreKey],
		memKeys[fundersTypes.MemStoreKey],

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
		appCodec,
		keys[bundlesTypes.StoreKey],
		memKeys[bundlesTypes.MemStoreKey],

		authTypes.NewModuleAddress(govTypes.ModuleName).String(),

		app.AccountKeeper,
		app.BankKeeper,
		app.DistributionKeeper,
		app.PoolKeeper,
		app.StakersKeeper,
		app.DelegationKeeper,
		app.FundersKeeper,
	)

	app.IBCKeeper = ibcKeeper.NewKeeper(
		appCodec,
		keys[ibcExported.StoreKey],
		app.GetSubspace(ibcExported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	app.IBCFeeKeeper = ibcFeeKeeper.NewKeeper(
		appCodec,
		keys[ibcFeeTypes.StoreKey],
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)

	app.IBCTransferKeeper = ibcTransferKeeper.NewKeeper(
		appCodec,
		keys[ibcTransferTypes.StoreKey],
		app.GetSubspace(ibcTransferTypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedIBCTransferKeeper,
	)

	app.ICAControllerKeeper = icaControllerKeeper.NewKeeper(
		appCodec,
		keys[icaControllerTypes.StoreKey],
		app.GetSubspace(icaControllerTypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper,
		app.MsgServiceRouter(),
	)

	app.ICAHostKeeper = icaHostKeeper.NewKeeper(
		appCodec,
		keys[icaHostTypes.StoreKey],
		app.GetSubspace(icaHostTypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)

	app.PFMKeeper = pfmKeeper.NewKeeper(
		appCodec, keys[pfmTypes.StoreKey],
		app.GetSubspace(pfmTypes.ModuleName),
		app.IBCTransferKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.DistributionKeeper,
		app.BankKeeper,
		app.IBCKeeper.ChannelKeeper,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	app.EvidenceKeeper = evidenceKeeper.NewKeeper(
		appCodec,
		keys[evidenceTypes.StoreKey],
		app.StakingKeeper,
		app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	// app.EvidenceKeeper = *evidenceKeeper

	govRouter := v1beta1.NewRouter()
	govRouter.
		AddRoute(govTypes.RouterKey, v1beta1.ProposalHandler).
		AddRoute(paramsProposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		// AddRoute(distrtypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(app.DistributionKeeper)).
		AddRoute(upgradeTypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcClientTypes.RouterKey, ibcClientHandler.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	govConfig := govTypes.DefaultConfig()
	app.GovKeeper = govKeeper.NewKeeper(
		appCodec,
		keys[govTypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.StakersKeeper,
		app.MsgServiceRouter(),
		govConfig,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)
	app.GovKeeper.SetLegacyRouter(govRouter)

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
		*app.GovKeeper,
		app.TeamKeeper,
		app.FundersKeeper,
	)
	// this line is used by starport scaffolding # stargate/app/keeperDefinition

	// Create static IBC router, add transfer route, then set and seal it
	var ibcTransferStack ibcPortTypes.IBCModule
	ibcTransferStack = ibcTransfer.NewIBCModule(app.IBCTransferKeeper)
	ibcTransferStack = ibcFee.NewIBCMiddleware(ibcTransferStack, app.IBCFeeKeeper)
	ibcTransferStack = pfm.NewIBCMiddleware(
		ibcTransferStack,
		app.PFMKeeper,
		0,
		pfmKeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		pfmKeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

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
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authTypes.ModuleName)),
		authz.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(bankTypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		consensus.NewAppModule(appCodec, app.ConsensusKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisisTypes.ModuleName)),
		distribution.NewAppModule(appCodec, app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distributionTypes.ModuleName)),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feeGrant.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govTypes.ModuleName)),
		group.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, mintTypes.DefaultInflationCalculationFn, app.GetSubspace(mintTypes.ModuleName)),
		params.NewAppModule(app.ParamsKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingTypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingTypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),

		// IBC
		ibc.NewAppModule(app.IBCKeeper),
		ibcFee.NewAppModule(app.IBCFeeKeeper),
		ibcTransfer.NewAppModule(app.IBCTransferKeeper),
		ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),

		// KYVE
		bundles.NewAppModule(appCodec, app.BundlesKeeper, app.AccountKeeper, app.BankKeeper, app.DistributionKeeper, app.MintKeeper, *app.UpgradeKeeper, app.PoolKeeper, app.TeamKeeper),
		delegation.NewAppModule(appCodec, app.DelegationKeeper, app.AccountKeeper, app.BankKeeper),
		global.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.GlobalKeeper, *app.UpgradeKeeper),
		pool.NewAppModule(appCodec, app.PoolKeeper, app.AccountKeeper, app.BankKeeper, *app.UpgradeKeeper),
		query.NewAppModule(appCodec, app.QueryKeeper, app.AccountKeeper, app.BankKeeper),
		stakers.NewAppModule(appCodec, app.StakersKeeper, app.AccountKeeper, app.BankKeeper),
		team.NewAppModule(appCodec, app.BankKeeper, app.MintKeeper, app.TeamKeeper, *app.UpgradeKeeper),
		funders.NewAppModule(appCodec, app.FundersKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		// upgrades should be run first
		upgradeTypes.ModuleName,
		capabilityTypes.ModuleName,
		mintTypes.ModuleName,
		// NOTE: x/team must be run before x/distribution and after x/mint.
		teamTypes.ModuleName,
		// NOTE: x/bundles must be run before x/distribution and after x/team.
		bundlesTypes.ModuleName,
		distributionTypes.ModuleName,
		slashingTypes.ModuleName,
		evidenceTypes.ModuleName,
		stakingTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		govTypes.ModuleName,
		crisisTypes.ModuleName,
		ibcFeeTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcExported.ModuleName,
		icaTypes.ModuleName,
		genUtilTypes.ModuleName,
		authzTypes.ModuleName,
		feeGrantTypes.ModuleName,
		groupTypes.ModuleName,
		paramsTypes.ModuleName,
		vestingTypes.ModuleName,
		consensusTypes.ModuleName,

		// this line is used by starport scaffolding # stargate/app/beginBlockers
		poolTypes.ModuleName,
		stakersTypes.ModuleName,
		delegationTypes.ModuleName,
		queryTypes.ModuleName,
		globalTypes.ModuleName,
		fundersTypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisisTypes.ModuleName,
		govTypes.ModuleName,
		stakingTypes.ModuleName,
		ibcFeeTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcExported.ModuleName,
		icaTypes.ModuleName,
		capabilityTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		slashingTypes.ModuleName,
		mintTypes.ModuleName,
		genUtilTypes.ModuleName,
		evidenceTypes.ModuleName,
		authzTypes.ModuleName,
		feeGrantTypes.ModuleName,
		groupTypes.ModuleName,
		paramsTypes.ModuleName,
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,
		consensusTypes.ModuleName,

		// this line is used by starport scaffolding # stargate/app/endBlockers
		poolTypes.ModuleName,
		stakersTypes.ModuleName,
		delegationTypes.ModuleName,
		bundlesTypes.ModuleName,
		queryTypes.ModuleName,
		globalTypes.ModuleName,
		teamTypes.ModuleName,
		fundersTypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilityTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		stakingTypes.ModuleName,
		slashingTypes.ModuleName,
		govTypes.ModuleName,
		mintTypes.ModuleName,
		crisisTypes.ModuleName,
		genUtilTypes.ModuleName,
		ibcFeeTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcExported.ModuleName,
		icaTypes.ModuleName,
		evidenceTypes.ModuleName,
		authzTypes.ModuleName,
		feeGrantTypes.ModuleName,
		groupTypes.ModuleName,
		paramsTypes.ModuleName,
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,
		consensusTypes.ModuleName,

		// this line is used by starport scaffolding # stargate/app/initGenesis
		poolTypes.ModuleName,
		stakersTypes.ModuleName,
		delegationTypes.ModuleName,
		bundlesTypes.ModuleName,
		queryTypes.ModuleName,
		globalTypes.ModuleName,
		teamTypes.ModuleName,
		fundersTypes.ModuleName,
	)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

	app.mm.RegisterInvariants(app.CrisisKeeper)
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
		app.IBCKeeper,
		*app.StakingKeeper,
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
		v1p4.UpgradeName,
		v1p4.CreateUpgradeHandler(
			app.mm,
			app.configurator,
			appCodec,
			app.ConsensusKeeper,
			app.GlobalKeeper,
			*app.GovKeeper,
			*app.IBCKeeper,
			app.ParamsKeeper,
		),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == v1p4.UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		app.SetStoreLoader(v1p4.CreateStoreLoader(upgradeInfo.Height))
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			cmtOs.Exit(err.Error())
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

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *App) Configurator() module.Configurator {
	return app.configurator
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
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
	return app.legacyAmino
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

// TxConfig returns a TxConfig
func (app *App) TxConfig() client.TxConfig {
	return app.txConfig
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *App) DefaultGenesis() map[string]json.RawMessage {
	return ModuleBasics.DefaultGenesis(app.appCodec)
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *storeTypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *storeTypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *storeTypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramsTypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authTx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway.
	node.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authTx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
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

func (app *App) RegisterNodeService(clientCtx client.Context) {
	node.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// SimulationManager implements the SimulationApp interface.
// NOTE: We simply return nil as we don't use the simulation manager anywhere.
func (app *App) SimulationManager() *module.SimulationManager { return nil }
