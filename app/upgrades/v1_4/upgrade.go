package v1_4

import (
	"errors"

	"github.com/KYVENetwork/chain/app/upgrades/v1_4/v1_3_types"
	"github.com/KYVENetwork/chain/util"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	poolKeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibcTmMigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"

	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	consensusKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	ibcKeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

//nolint:all
//goland:noinspection GoDeprecation
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.BinaryCodec,
	consensusKeeper consensusKeeper.Keeper,
	globalKeeper globalKeeper.Keeper,
	govKeeper govKeeper.Keeper,
	ibcKeeper ibcKeeper.Keeper,
	paramsKeeper paramsKeeper.Keeper,
	poolKeeper poolKeeper.Keeper,
	fundersKeeper fundersKeeper.Keeper,
	bankKeeper bankKeeper.Keeper,
	accountKeeper authKeeper.AccountKeeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Info("Run v1.4 upgrade")

		distributionSpace, _ := paramsKeeper.GetSubspace("distribution")
		distributionSpace.WithKeyTable(distributionTypes.ParamKeyTable())

		stakingSpace, _ := paramsKeeper.GetSubspace("staking")
		stakingSpace.WithKeyTable(stakingTypes.ParamKeyTable())

		authSpace, _ := paramsKeeper.GetSubspace("auth")
		authSpace.WithKeyTable(authTypes.ParamKeyTable())

		bankSpace, _ := paramsKeeper.GetSubspace("bank")
		bankSpace.WithKeyTable(bankTypes.ParamKeyTable())

		crisisSpace, _ := paramsKeeper.GetSubspace("crisis")
		crisisSpace.WithKeyTable(crisisTypes.ParamKeyTable())

		govSpace, _ := paramsKeeper.GetSubspace("gov")
		govSpace.WithKeyTable(govTypes.ParamKeyTable())

		mintSpace, _ := paramsKeeper.GetSubspace("mint")
		mintSpace.WithKeyTable(mintTypes.ParamKeyTable())

		slashingSpace, _ := paramsKeeper.GetSubspace("slashing")
		slashingSpace.WithKeyTable(slashingTypes.ParamKeyTable())

		// Migrate consensus parameters from x/params to dedicated x/consensus module.
		baseAppSubspace := paramsKeeper.Subspace(baseapp.Paramspace).
			WithKeyTable(paramsTypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppSubspace, &consensusKeeper)

		var err error

		// ibc-go v7.0 to v7.1 upgrade
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := ibcKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		ibcKeeper.ClientKeeper.SetParams(ctx, params)

		// Run module migrations.
		vm, err = mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// Prune expired Tendermint consensus states.
		_, err = ibcTmMigrations.PruneExpiredConsensusStates(ctx, cdc, ibcKeeper.ClientKeeper)
		if err != nil {
			return vm, err
		}

		// Migrate initial deposit ratio.
		err = migrateInitialDepositRatio(ctx, globalKeeper, govKeeper)
		if err != nil {
			return vm, err
		}

		// Migrate funders.
		err = migrateFundersAndPools(ctx, cdc, poolKeeper, fundersKeeper, bankKeeper, accountKeeper)
		if err != nil {
			return vm, err
		}

		// Set min gas for funder creation in global module
		globalParams := globalKeeper.GetParams(ctx)
		globalParams.GasAdjustments = append(globalParams.GasAdjustments, globalTypes.GasAdjustment{
			Type:   "/kyve.funders.v1beta1.MsgCreateFunder",
			Amount: 50_000_000,
		})
		globalKeeper.SetParams(ctx, globalParams)

		return vm, nil
	}
}

// migrateInitialDepositRatio migrates the MinInitialDepositRatio parameter from
// our custom x/global module to the x/gov module.
func migrateInitialDepositRatio(
	ctx sdk.Context,
	globalKeeper globalKeeper.Keeper,
	govKeeper govKeeper.Keeper,
) error {
	minInitialDepositRatio := globalKeeper.GetMinInitialDepositRatio(ctx)

	params := govKeeper.GetParams(ctx)
	params.MinInitialDepositRatio = minInitialDepositRatio.String()

	return govKeeper.SetParams(ctx, params)
}

type FundingMigration struct {
	PoolId uint64
	Amount uint64
}

type FunderMigration struct {
	Address  string
	Fundings []FundingMigration
}

// migrateFunders migrates funders from x/pool to x/funders and creates funding states for pools.
func migrateFundersAndPools(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	poolKeeper poolKeeper.Keeper,
	fundersKeeper fundersKeeper.Keeper,
	bankKeeper bankKeeper.Keeper,
	accountKeeper authKeeper.AccountKeeper,
) error {
	pools, err := v1_3_types.GetAllPools(ctx, poolKeeper, cdc)
	if err != nil {
		return err
	}

	toBeCreatedFunders := make(map[string]*FunderMigration)
	amountToBeTransferred := uint64(0)

	// Get all funders and their funding from pools.
	for _, oldPool := range pools {
		checkTotalFunds := uint64(0)
		for _, funder := range oldPool.Funders {
			if funder.Amount > 0 {
				_, ok := toBeCreatedFunders[funder.Address]
				if ok {
					toBeCreatedFunders[funder.Address].Fundings = append(toBeCreatedFunders[funder.Address].Fundings, FundingMigration{PoolId: oldPool.Id, Amount: funder.Amount})
				} else {
					toBeCreatedFunders[funder.Address] = &FunderMigration{
						Address:  funder.Address,
						Fundings: []FundingMigration{{PoolId: oldPool.Id, Amount: funder.Amount}},
					}
				}
				checkTotalFunds += funder.Amount
			}
		}
		if checkTotalFunds != oldPool.TotalFunds {
			return errors.New("total funds is not equal to the sum of all funders amount")
		}
		amountToBeTransferred += oldPool.TotalFunds

		// Create funding state for pool.
		fundersKeeper.SetFundingState(ctx, &fundersTypes.FundingState{
			PoolId:                oldPool.Id,
			ActiveFunderAddresses: []string{},
		})

		poolKeeper.SetPool(ctx, poolTypes.Pool{
			Id:                   oldPool.Id,
			Name:                 oldPool.Name,
			Runtime:              oldPool.Runtime,
			Logo:                 oldPool.Logo,
			Config:               oldPool.Config,
			StartKey:             oldPool.StartKey,
			CurrentKey:           oldPool.CurrentKey,
			CurrentSummary:       oldPool.CurrentSummary,
			CurrentIndex:         oldPool.CurrentIndex,
			TotalBundles:         oldPool.TotalBundles,
			UploadInterval:       oldPool.UploadInterval,
			InflationShareWeight: oldPool.OperatingCost,
			MinDelegation:        oldPool.MinDelegation,
			MaxBundleSize:        oldPool.MaxBundleSize,
			Disabled:             oldPool.Disabled,
			Protocol: &poolTypes.Protocol{
				Version:     oldPool.Protocol.Version,
				Binaries:    oldPool.Protocol.Binaries,
				LastUpgrade: oldPool.Protocol.LastUpgrade,
			},
			UpgradePlan: &poolTypes.UpgradePlan{
				Version:     oldPool.UpgradePlan.Version,
				Binaries:    oldPool.UpgradePlan.Binaries,
				ScheduledAt: oldPool.UpgradePlan.ScheduledAt,
				Duration:    oldPool.UpgradePlan.Duration,
			},
			CurrentStorageProviderId: oldPool.CurrentStorageProviderId,
			CurrentCompressionId:     oldPool.CurrentCompressionId,
		})
	}

	// Create new funders and fundings.
	for _, funder := range toBeCreatedFunders {
		fundersKeeper.SetFunder(ctx, &fundersTypes.Funder{
			Address:     funder.Address,
			Moniker:     funder.Address,
			Identity:    "",
			Website:     "",
			Contact:     "",
			Description: "",
		})
		for _, funding := range funder.Fundings {
			fundersKeeper.SetFunding(ctx, &fundersTypes.Funding{
				FunderAddress:   funder.Address,
				PoolId:          funding.PoolId,
				Amount:          funding.Amount,
				AmountPerBundle: fundersTypes.DefaultMinFundingAmountPerBundle,
				// Previous funders will not be considered, as there is no way to calculate this on chain.
				// Although almost all funding was only provided by the Foundation itself.
				TotalFunded: 0,
			})
		}
	}

	// Check if pool module balance is equal to the sum of all pools total funds.
	poolModule := accountKeeper.GetModuleAddress(poolTypes.ModuleName)
	balance := bankKeeper.GetBalance(ctx, poolModule, globalTypes.Denom)
	if balance.Amount.Uint64() != amountToBeTransferred {
		return errors.New("pool module balance is not equal to the sum of all pools total funds")
	}

	// Transfer funds from pools to funders.
	if err := util.TransferFromModuleToModule(bankKeeper, ctx, poolTypes.ModuleName, fundersTypes.ModuleName, amountToBeTransferred); err != nil {
		return err
	}

	return nil
}
