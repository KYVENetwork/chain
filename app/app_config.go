package app

import (
	"time"

	runtime "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	app "cosmossdk.io/api/cosmos/app/v1alpha1"
	txConfig "cosmossdk.io/api/cosmos/tx/config/v1"
	"cosmossdk.io/core/appconfig"
	"google.golang.org/protobuf/types/known/durationpb"

	// Auth
	authModule "cosmossdk.io/api/cosmos/auth/module/v1"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Authz
	authzModule "cosmossdk.io/api/cosmos/authz/module/v1"
	"github.com/cosmos/cosmos-sdk/x/authz"
	// Bank
	bankModule "cosmossdk.io/api/cosmos/bank/module/v1"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// Bundles
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Capability
	capabilityModule "cosmossdk.io/api/cosmos/capability/module/v1"
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	// Consensus
	consensusModule "cosmossdk.io/api/cosmos/consensus/module/v1"
	consensusTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	// Crisis
	crisisModule "cosmossdk.io/api/cosmos/crisis/module/v1"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	// Delegation
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Distribution
	distributionModule "cosmossdk.io/api/cosmos/distribution/module/v1"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	// Evidence
	evidenceModule "cosmossdk.io/api/cosmos/evidence/module/v1"
	evidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	// FeeGrant
	feeGrantModule "cosmossdk.io/api/cosmos/feegrant/module/v1"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	// GenUtil
	genUtilModule "cosmossdk.io/api/cosmos/genutil/module/v1"
	genUtilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	// Global
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// Governance
	govModule "cosmossdk.io/api/cosmos/gov/module/v1"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// Group
	groupModule "cosmossdk.io/api/cosmos/group/module/v1"
	"github.com/cosmos/cosmos-sdk/x/group"
	// IBC Core
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/core/exported"
	// IBC Fee
	ibcFeeTypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	// IBC Transfer
	ibcTransferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	// ICA
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	// Mint
	mintModule "cosmossdk.io/api/cosmos/mint/module/v1"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	// NFT -- TODO(@john): Do we want to include this?
	nftModule "cosmossdk.io/api/cosmos/nft/module/v1"
	"github.com/cosmos/cosmos-sdk/x/nft"
	// Params
	paramsModule "cosmossdk.io/api/cosmos/params/module/v1"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	// Pool
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	// Query
	queryTypes "github.com/KYVENetwork/chain/x/query/types"
	// Slashing
	slashingModule "cosmossdk.io/api/cosmos/slashing/module/v1"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	// Stakers
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	// Staking
	stakingModule "cosmossdk.io/api/cosmos/staking/module/v1"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// Team
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
	// Upgrade
	upgradeModule "cosmossdk.io/api/cosmos/upgrade/module/v1"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	// Vesting
	vestingModule "cosmossdk.io/api/cosmos/vesting/module/v1"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

var (
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder = []string{
		capabilityTypes.ModuleName, authTypes.ModuleName, bankTypes.ModuleName,
		distributionTypes.ModuleName, stakingTypes.ModuleName, slashingTypes.ModuleName, govTypes.ModuleName,
		mintTypes.ModuleName, crisisTypes.ModuleName, genUtilTypes.ModuleName, evidenceTypes.ModuleName, authz.ModuleName,
		feegrant.ModuleName, nft.ModuleName, group.ModuleName, paramsTypes.ModuleName, upgradeTypes.ModuleName,
		vestingTypes.ModuleName, consensusTypes.ModuleName,

		ibcTransferTypes.ModuleName,
		ibcTypes.ModuleName,
		icaTypes.ModuleName,
		ibcFeeTypes.ModuleName,

		poolTypes.ModuleName,
		stakersTypes.ModuleName,
		delegationTypes.ModuleName,
		bundlesTypes.ModuleName,
		queryTypes.ModuleName,
		globalTypes.ModuleName,
		teamTypes.ModuleName,
	}

	// module account permissions
	moduleAccPerms = []*authModule.ModuleAccountPermission{
		{Account: authTypes.FeeCollectorName, Permissions: []string{authTypes.Burner}},
		{Account: distributionTypes.ModuleName},
		{Account: mintTypes.ModuleName, Permissions: []string{authTypes.Minter}},
		{Account: stakingTypes.BondedPoolName, Permissions: []string{authTypes.Burner, stakingTypes.ModuleName}},
		{Account: stakingTypes.NotBondedPoolName, Permissions: []string{authTypes.Burner, stakingTypes.ModuleName}},
		{Account: govTypes.ModuleName, Permissions: []string{authTypes.Burner}},
		{Account: nft.ModuleName},

		// IBC
		{Account: ibcFeeTypes.ModuleName},
		{Account: ibcTransferTypes.ModuleName, Permissions: []string{authTypes.Burner, authTypes.Minter}},
		{Account: icaTypes.ModuleName},

		// KYVE
		{Account: bundlesTypes.ModuleName},
		{Account: delegationTypes.ModuleName},
		{Account: poolTypes.ModuleName},
		{Account: globalTypes.ModuleName, Permissions: []string{authTypes.Burner}},
		{Account: stakersTypes.ModuleName},
		{Account: teamTypes.ModuleName},
	}

	// blocked account addresses
	blockAccAddrs = []string{
		authTypes.FeeCollectorName,
		distributionTypes.ModuleName,
		mintTypes.ModuleName,
		stakingTypes.BondedPoolName,
		stakingTypes.NotBondedPoolName,
		nft.ModuleName,
		// We allow the following module accounts to receive funds:
		// govTypes.ModuleName
	}

	// application configuration (used by depinject)
	AppConfig = appconfig.Compose(&app.Config{
		Modules: []*app.ModuleConfig{
			{
				Name: "runtime",
				Config: appconfig.WrapAny(&runtime.Module{
					AppName: "SimApp",
					// During begin block slashing happens after distr.BeginBlocker so that
					// there is nothing left over in the validator fee pool, so as to keep the
					// CanWithdrawInvariant invariant.
					// NOTE: staking module is required if HistoricalEntries param > 0
					// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
					BeginBlockers: []string{
						upgradeTypes.ModuleName,
						capabilityTypes.ModuleName,
						mintTypes.ModuleName,
						// NOTE: x/team must be run before x/distribution and after x/mint.
						teamTypes.ModuleName,
						distributionTypes.ModuleName,
						slashingTypes.ModuleName,
						evidenceTypes.ModuleName,
						stakingTypes.ModuleName,
						authTypes.ModuleName,
						bankTypes.ModuleName,
						govTypes.ModuleName,
						crisisTypes.ModuleName,
						genUtilTypes.ModuleName,
						authz.ModuleName,
						feegrant.ModuleName,
						nft.ModuleName,
						group.ModuleName,
						paramsTypes.ModuleName,
						vestingTypes.ModuleName,
						consensusTypes.ModuleName,

						ibcTransferTypes.ModuleName,
						ibcTypes.ModuleName,
						icaTypes.ModuleName,
						ibcFeeTypes.ModuleName,

						poolTypes.ModuleName,
						stakersTypes.ModuleName,
						delegationTypes.ModuleName,
						bundlesTypes.ModuleName,
						queryTypes.ModuleName,
						globalTypes.ModuleName,
					},
					EndBlockers: []string{
						crisisTypes.ModuleName,
						govTypes.ModuleName,
						stakingTypes.ModuleName,
						capabilityTypes.ModuleName,
						authTypes.ModuleName,
						bankTypes.ModuleName,
						distributionTypes.ModuleName,
						slashingTypes.ModuleName,
						mintTypes.ModuleName,
						genUtilTypes.ModuleName,
						evidenceTypes.ModuleName,
						authz.ModuleName,
						feegrant.ModuleName,
						nft.ModuleName,
						group.ModuleName,
						paramsTypes.ModuleName,
						consensusTypes.ModuleName,
						upgradeTypes.ModuleName,
						vestingTypes.ModuleName,

						ibcTransferTypes.ModuleName,
						ibcTypes.ModuleName,
						icaTypes.ModuleName,
						ibcFeeTypes.ModuleName,

						poolTypes.ModuleName,
						stakersTypes.ModuleName,
						delegationTypes.ModuleName,
						bundlesTypes.ModuleName,
						queryTypes.ModuleName,
						globalTypes.ModuleName,
						teamTypes.ModuleName,
					},
					OverrideStoreKeys: []*runtime.StoreKeyConfig{
						{
							ModuleName: authTypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					InitGenesis: genesisModuleOrder,
					// When ExportGenesis is not specified, the export genesis module order
					// is equal to the init genesis order
					// ExportGenesis: genesisModuleOrder,
					// Uncomment if you want to set a custom migration order here.
					// OrderMigrations: nil,
				}),
			},

			// ----- Cosmos SDK Modules -----

			{
				Name: authTypes.ModuleName,
				Config: appconfig.WrapAny(&authModule.Module{
					Bech32Prefix:             "kyve",
					ModuleAccountPermissions: moduleAccPerms,
					// By default modules authority is the governance module. This is configurable with the following:
					// Authority: "group", // A custom module authority can be set using a module name
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a specific address
				}),
			},
			{
				Name:   vestingTypes.ModuleName,
				Config: appconfig.WrapAny(&vestingModule.Module{}),
			},
			{
				Name: bankTypes.ModuleName,
				Config: appconfig.WrapAny(&bankModule.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   stakingTypes.ModuleName,
				Config: appconfig.WrapAny(&stakingModule.Module{}),
			},
			{
				Name:   slashingTypes.ModuleName,
				Config: appconfig.WrapAny(&slashingModule.Module{}),
			},
			{
				Name:   paramsTypes.ModuleName,
				Config: appconfig.WrapAny(&paramsModule.Module{}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txConfig.Config{}),
			},
			{
				Name:   genUtilTypes.ModuleName,
				Config: appconfig.WrapAny(&genUtilModule.Module{}),
			},
			{
				Name:   authz.ModuleName,
				Config: appconfig.WrapAny(&authzModule.Module{}),
			},
			{
				Name:   upgradeTypes.ModuleName,
				Config: appconfig.WrapAny(&upgradeModule.Module{}),
			},
			{
				Name:   distributionTypes.ModuleName,
				Config: appconfig.WrapAny(&distributionModule.Module{}),
			},
			{
				Name: capabilityTypes.ModuleName,
				Config: appconfig.WrapAny(&capabilityModule.Module{
					SealKeeper: true,
				}),
			},
			{
				Name:   evidenceTypes.ModuleName,
				Config: appconfig.WrapAny(&evidenceModule.Module{}),
			},
			{
				Name:   mintTypes.ModuleName,
				Config: appconfig.WrapAny(&mintModule.Module{}),
			},
			{
				Name: group.ModuleName,
				Config: appconfig.WrapAny(&groupModule.Module{
					MaxExecutionPeriod: durationpb.New(time.Second * 1209600),
					MaxMetadataLen:     255,
				}),
			},
			{
				Name:   nft.ModuleName,
				Config: appconfig.WrapAny(&nftModule.Module{}),
			},
			{
				Name:   feegrant.ModuleName,
				Config: appconfig.WrapAny(&feeGrantModule.Module{}),
			},
			{
				Name:   govTypes.ModuleName,
				Config: appconfig.WrapAny(&govModule.Module{}),
			},
			{
				Name:   crisisTypes.ModuleName,
				Config: appconfig.WrapAny(&crisisModule.Module{}),
			},
			{
				Name:   consensusTypes.ModuleName,
				Config: appconfig.WrapAny(&consensusModule.Module{}),
			},

			// ----- KYVE Modules -----

			// TODO(@john)

			//{
			//	Name:   bundlesTypes.ModuleName,
			//	Config: appconfig.WrapAny(&bundlesModule.Module{}),
			//},
			//{
			//	Name:   delegationTypes.ModuleName,
			//	Config: appconfig.WrapAny(&delegationModule.Module{}),
			//},
			//{
			//	Name:   globalTypes.ModuleName,
			//	Config: appconfig.WrapAny(&globalModule.Module{}),
			//},
			//{
			//	Name:   poolTypes.ModuleName,
			//	Config: appconfig.WrapAny(&poolModule.Module{}),
			//},
			//{
			//	Name:   queryTypes.ModuleName,
			//	Config: appconfig.WrapAny(&queryModule.Module{}),
			//},
			//{
			//	Name:   stakersTypes.ModuleName,
			//	Config: appconfig.WrapAny(&stakersModule.Module{}),
			//},
			//{
			//	Name:   teamTypes.ModuleName,
			//	Config: appconfig.WrapAny(&teamModule.Module{}),
			//},
		},
	})
)
