package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Authz
	authzKeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// Bundles
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	// Capability
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	// Crisis
	crisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	// Delegation
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	// Distribution
	distributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	// Evidence
	evidenceKeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	// FeeGrant
	feeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	// Global
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	// Governance
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govV1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	// Group
	groupKeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	// IBC
	ibcHost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibcKeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	// IBC Fee
	ibcFeeKeeper "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/keeper"
	// IBC Transfer
	ibcTransferKeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	// ICA Controller
	icaControllerKeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	icaControllerTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	// ICA Host
	icaHostKeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	// Mint
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	// Parameters
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	// Pool
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	// Query
	queryKeeper "github.com/KYVENetwork/chain/x/query/keeper"
	// Router
	routerKeeper "github.com/strangelove-ventures/packet-forward-middleware/v6/router/keeper"
	routerTypes "github.com/strangelove-ventures/packet-forward-middleware/v6/router/types"
	// Slashing
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	// Stakers
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	// Staking
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// Team
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	// Upgrade
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
)

type Keepers struct {
	// Cosmos SDK
	AccountKeeper      authKeeper.AccountKeeper
	AuthzKeeper        authzKeeper.Keeper
	BankKeeper         bankKeeper.Keeper
	CapabilityKeeper   *capabilityKeeper.Keeper
	CrisisKeeper       crisisKeeper.Keeper
	DistributionKeeper distributionKeeper.Keeper
	EvidenceKeeper     evidenceKeeper.Keeper
	FeeGrantKeeper     feeGrantKeeper.Keeper
	GovKeeper          govKeeper.Keeper
	GroupKeeper        groupKeeper.Keeper
	MintKeeper         mintKeeper.Keeper
	ParamsKeeper       paramsKeeper.Keeper
	SlashingKeeper     slashingKeeper.Keeper
	StakingKeeper      stakingKeeper.Keeper
	UpgradeKeeper      upgradeKeeper.Keeper

	// IBC
	IBCKeeper           *ibcKeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper        ibcFeeKeeper.Keeper
	IBCTransferKeeper   ibcTransferKeeper.Keeper
	ICAControllerKeeper icaControllerKeeper.Keeper
	ICAHostKeeper       icaHostKeeper.Keeper
	RouterKeeper        *routerKeeper.Keeper

	// KYVE
	BundlesKeeper    bundlesKeeper.Keeper
	DelegationKeeper delegationKeeper.Keeper
	GlobalKeeper     globalKeeper.Keeper
	PoolKeeper       poolKeeper.Keeper
	QueryKeeper      queryKeeper.Keeper
	StakersKeeper    stakersKeeper.Keeper
	TeamKeeper       teamKeeper.Keeper

	// ----- Scoped Keepers -----
	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilityKeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilityKeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilityKeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilityKeeper.ScopedKeeper
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storeTypes.StoreKey) paramsKeeper.Keeper {
	keeper := paramsKeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	keeper.Subspace(authTypes.ModuleName)
	keeper.Subspace(bankTypes.ModuleName)
	keeper.Subspace(crisisTypes.ModuleName)
	keeper.Subspace(distributionTypes.ModuleName)
	keeper.Subspace(govTypes.ModuleName).WithKeyTable(govV1.ParamKeyTable())
	keeper.Subspace(mintTypes.ModuleName)
	keeper.Subspace(slashingTypes.ModuleName)
	keeper.Subspace(stakingTypes.ModuleName)

	keeper.Subspace(ibcHost.ModuleName)
	keeper.Subspace(ibcTransferTypes.ModuleName)
	keeper.Subspace(icaControllerTypes.SubModuleName)
	keeper.Subspace(icaHostTypes.SubModuleName)
	keeper.Subspace(routerTypes.ModuleName)

	return keeper
}
