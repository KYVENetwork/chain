package query

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"

	"github.com/KYVENetwork/chain/util"
	bundlekeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	// this line is used by starport scaffolding # 1

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	modulev1 "github.com/KYVENetwork/chain/api/kyve/query/module"
	"github.com/KYVENetwork/chain/x/query/client/cli"
	"github.com/KYVENetwork/chain/x/query/keeper"
	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModuleBasic      = (*AppModule)(nil)
	_ module.HasGenesis          = (*AppModule)(nil)
	_ module.HasInvariants       = (*AppModule)(nil)
	_ module.HasConsensusVersion = (*AppModule)(nil)

	_ appmodule.AppModule = (*AppModule)(nil)
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface that defines the independent methods a Cosmos SDK module needs to implement.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the name of the module as a string
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used to marshal and unmarshal structs to/from []byte in order to persist them in the module's KVStore
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message
func (a AppModuleBasic) RegisterInterfaces(_ cdctypes.InterfaceRegistry) {}

// DefaultGenesis returns a default GenesisState for the module, marshalled to json.RawMessage. The default GenesisState need to be defined by the module developer and is primarily used for testing
func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	return nil
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form
func (AppModuleBasic) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryAccountHandlerClient(context.Background(), mux, types.NewQueryAccountClient(clientCtx))
	_ = types.RegisterQueryPoolHandlerClient(context.Background(), mux, types.NewQueryPoolClient(clientCtx))
	_ = types.RegisterQueryStakersHandlerClient(context.Background(), mux, types.NewQueryStakersClient(clientCtx))
	_ = types.RegisterQueryBundlesHandlerClient(context.Background(), mux, types.NewQueryBundlesClient(clientCtx))
	_ = types.RegisterQueryParamsHandlerClient(context.Background(), mux, types.NewQueryParamsClient(clientCtx))
	_ = types.RegisterQueryFundersHandlerClient(context.Background(), mux, types.NewQueryFundersClient(clientCtx))
}

// GetTxCmd returns the root Tx command for the module. The subcommands of this root command are used by end-users to generate new transactions containing messages defined in the module
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the root query command for the module. The subcommands of this root command are used by end-users to generate new queries to the subset of the state defined by the module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface that defines the inter-dependent methods that modules need to implement
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// Deprecated: use RegisterServices
func (AppModule) QuerierRoute() string { return types.RouterKey }

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryAccountServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryPoolServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryStakersServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryBundlesServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryParamsServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryFundersServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(_ sdk.Context, _ codec.JSONCodec, _ json.RawMessage) {
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONCodec) json.RawMessage {
	return nil
}

// ConsensusVersion is a sequence number for state-breaking change of the module. It should be incremented on each consensus-breaking change introduced by the module. To avoid wrong/empty versions, the initial version should be set to 1
func (AppModule) ConsensusVersion() uint64 { return 1 }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Cdc    codec.Codec
	Config *modulev1.Module
	Logger log.Logger

	AccountKeeper      authkeeper.AccountKeeper
	BankKeeper         bankKeeper.Keeper
	DistributionKeeper util.DistributionKeeper
	UpgradeKeeper      util.UpgradeKeeper
	PoolKeeper         *poolKeeper.Keeper
	TeamKeeper         teamKeeper.Keeper
	StakersKeeper      *stakersKeeper.Keeper
	BundlesKeeper      bundlekeeper.Keeper
	GovKeeper          *govkeeper.Keeper
	GlobalKeeper       globalKeeper.Keeper
	FundersKeeper      fundersKeeper.Keeper
	StakingKeeper      util.StakingKeeper
}

type ModuleOutputs struct {
	depinject.Out

	QueryKeeper keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Cdc,
		in.Logger,
		in.AccountKeeper,
		in.BankKeeper,
		in.DistributionKeeper,
		in.PoolKeeper,
		in.StakersKeeper,
		in.BundlesKeeper,
		in.GlobalKeeper,
		in.GovKeeper,
		in.TeamKeeper,
		in.FundersKeeper,
		in.StakingKeeper,
	)
	m := NewAppModule(
		in.Cdc,
		k,
		in.AccountKeeper,
		in.BankKeeper,
	)

	return ModuleOutputs{QueryKeeper: k, Module: m}
}
