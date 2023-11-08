package app

import (
	"github.com/KYVENetwork/chain/x/funders"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	// Authz
	authz "github.com/cosmos/cosmos-sdk/x/authz/module"
	// Bank
	"github.com/cosmos/cosmos-sdk/x/bank"
	// Bundles
	"github.com/KYVENetwork/chain/x/bundles"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Capability
	"github.com/cosmos/cosmos-sdk/x/capability"
	// Consensus
	"github.com/cosmos/cosmos-sdk/x/consensus"
	// Crisis
	"github.com/cosmos/cosmos-sdk/x/crisis"
	// Delegation
	"github.com/KYVENetwork/chain/x/delegation"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Distribution
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	// Evidence
	"github.com/cosmos/cosmos-sdk/x/evidence"
	// FeeGrant
	feeGrant "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	// GenUtil
	"github.com/cosmos/cosmos-sdk/x/genutil"
	// Global
	"github.com/KYVENetwork/chain/x/global"
	// Governance
	"github.com/cosmos/cosmos-sdk/x/gov"
	govClient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// Group
	group "github.com/cosmos/cosmos-sdk/x/group/module"
	// IBC Core
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcClient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	// IBC Light Clients
	ibcSm "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	ibcTm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	// IBC Fee
	ibcFee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcFeeTypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	// IBC Transfer
	ibcTransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibcTransferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	// ICA
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	// Mint
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	// Parameters
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	// PFM
	pfm "github.com/strangelove-ventures/packet-forward-middleware/v7/router"
	// Pool
	"github.com/KYVENetwork/chain/x/pool"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	// Query
	"github.com/KYVENetwork/chain/x/query"
	// Slashing
	"github.com/cosmos/cosmos-sdk/x/slashing"
	// Stakers
	"github.com/KYVENetwork/chain/x/stakers"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	// Staking
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// Team
	"github.com/KYVENetwork/chain/x/team"
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
	// Upgrade
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
)

// appModuleBasics returns ModuleBasics for the module BasicManager.
var appModuleBasics = []module.AppModuleBasic{
	// Cosmos SDK
	auth.AppModuleBasic{},
	authz.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	consensus.AppModuleBasic{},
	crisis.AppModuleBasic{},
	distribution.AppModuleBasic{},
	evidence.AppModuleBasic{},
	feeGrant.AppModuleBasic{},
	genutil.AppModuleBasic{},
	gov.NewAppModuleBasic([]govClient.ProposalHandler{
		paramsClient.ProposalHandler,
		upgradeClient.LegacyProposalHandler,
		upgradeClient.LegacyCancelProposalHandler,
		ibcClient.UpdateClientProposalHandler,
		ibcClient.UpgradeProposalHandler,
	}),
	group.AppModuleBasic{},
	mint.AppModuleBasic{},
	params.AppModuleBasic{},
	slashing.AppModuleBasic{},
	staking.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	vesting.AppModuleBasic{},

	// IBC
	ibc.AppModuleBasic{},
	ibcSm.AppModuleBasic{},
	ibcTm.AppModuleBasic{},
	ibcFee.AppModuleBasic{},
	ibcTransfer.AppModuleBasic{},
	ica.AppModuleBasic{},
	pfm.AppModuleBasic{},

	// KYVE
	bundles.AppModuleBasic{},
	delegation.AppModuleBasic{},
	global.AppModuleBasic{},
	pool.AppModuleBasic{},
	query.AppModuleBasic{},
	stakers.AppModuleBasic{},
	team.AppModuleBasic{},
	funders.AppModuleBasic{},
}

// moduleAccountPermissions ...
var moduleAccountPermissions = map[string][]string{
	// Cosmos SDK
	authTypes.FeeCollectorName:     {authTypes.Burner},
	distributionTypes.ModuleName:   nil,
	govTypes.ModuleName:            {authTypes.Burner},
	mintTypes.ModuleName:           {authTypes.Minter},
	stakingTypes.BondedPoolName:    {authTypes.Burner, authTypes.Staking},
	stakingTypes.NotBondedPoolName: {authTypes.Burner, authTypes.Staking},

	// IBC
	ibcTransferTypes.ModuleName: {authTypes.Minter, authTypes.Burner},
	ibcFeeTypes.ModuleName:      nil,
	icaTypes.ModuleName:         nil,

	// KYVE
	bundlesTypes.ModuleName:    nil,
	delegationTypes.ModuleName: nil,
	poolTypes.ModuleName:       nil,
	stakersTypes.ModuleName:    nil,
	teamTypes.ModuleName:       nil,
	fundersTypes.ModuleName:    nil,
}

// BlockedModuleAccountAddrs returns all the app's blocked module account addresses.
func (app *App) BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()
	delete(modAccAddrs, authTypes.NewModuleAddress(govTypes.ModuleName).String())

	return modAccAddrs
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		modAccAddrs[authTypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
