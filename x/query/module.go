package query

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"

	"github.com/KYVENetwork/chain/util"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	// this line is used by starport scaffolding # 1

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	moduleV1 "github.com/KYVENetwork/chain/pulsar/kyve/query/module/v1"
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
	_ appmodule.AppModule   = AppModule{}
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
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
	_ = types.RegisterQueryDelegationHandlerClient(context.Background(), mux, types.NewQueryDelegationClient(clientCtx))
	_ = types.RegisterQueryBundlesHandlerClient(context.Background(), mux, types.NewQueryBundlesClient(clientCtx))
	_ = types.RegisterQueryParamsHandlerClient(context.Background(), mux, types.NewQueryParamsClient(clientCtx))
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

	keeper keeper.Keeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryAccountServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryPoolServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryStakersServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryDelegationServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryBundlesServer(cfg.QueryServer(), am.keeper)
	types.RegisterQueryParamsServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(_ sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONCodec) json.RawMessage {
	return nil
}

// ConsensusVersion is a sequence number for state-breaking change of the module. It should be incremented on each consensus-breaking change introduced by the module. To avoid wrong/empty versions, the initial version should be set to 1
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// App Wiring Setup

func init() {
	appmodule.Register(&moduleV1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type QueryInputs struct {
	depinject.In

	Config *moduleV1.Module
	Cdc    codec.Codec
	Key    *storeTypes.KVStoreKey
	MemKey *storeTypes.MemoryStoreKey

	AccountKeeper      util.AccountKeeper
	BankKeeper         util.BankKeeper
	BundlesKeeper      types.BundlesKeeper
	DelegationKeeper   types.DelegationKeeper
	DistributionKeeper util.DistributionKeeper
	GlobalKeeper       types.GlobalKeeper
	GovKeeper          util.GovKeeper
	PoolKeeper         types.PoolKeeper
	StakersKeeper      types.StakersKeeper
	UpgradeKeeper      util.UpgradeKeeper
}

type QueryOutputs struct {
	depinject.Out

	QueryKeeper keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in QueryInputs) QueryOutputs {
	queryKeeper := *keeper.NewKeeper(
		in.Cdc,
		in.Key,
		in.MemKey,
		in.AccountKeeper,
		in.BankKeeper,
		in.DistributionKeeper,
		in.PoolKeeper,
		in.StakersKeeper,
		in.DelegationKeeper,
		in.BundlesKeeper,
		in.GlobalKeeper,
		in.GovKeeper,
	)
	m := NewAppModule(in.Cdc, queryKeeper)

	return QueryOutputs{QueryKeeper: queryKeeper, Module: m}
}
