package v1_5

import (
	"context"
	"fmt"

	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/bundles"
	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/delegation"
	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/funders"
	"github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/stakers"
	delegationKeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	fundersKeeper "github.com/KYVENetwork/chain/x/funders/keeper"
	"github.com/KYVENetwork/chain/x/funders/types"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	stakersKeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"

	"cosmossdk.io/math"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/KYVENetwork/chain/x/bundles/keeper"
	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v1.5.0"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, cdc codec.Codec, storeKeys []storetypes.StoreKey, bundlesKeeper keeper.Keeper, delegationKeeper delegationKeeper.Keeper, fundersKeeper fundersKeeper.Keeper, stakersKeeper *stakersKeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		logger := sdkCtx.Logger().With("upgrade", UpgradeName)
		logger.Info(fmt.Sprintf("performing upgrade %v", UpgradeName))

		// TODO: migrate gov params

		// migrate fundings
		if storeKey, err := getStoreKey(storeKeys, fundersTypes.StoreKey); err == nil {
			migrateFundersModule(sdkCtx, cdc, storeKey, fundersKeeper)
		} else {
			return nil, err
		}

		// migrate delegations
		if storeKey, err := getStoreKey(storeKeys, delegationTypes.StoreKey); err == nil {
			migrateDelegationModule(sdkCtx, cdc, storeKey, delegationKeeper)
		} else {
			return nil, err
		}

		// migrate stakers
		if storeKey, err := getStoreKey(storeKeys, stakersTypes.StoreKey); err == nil {
			migrateStakersModule(sdkCtx, cdc, storeKey, stakersKeeper)
		} else {
			return nil, err
		}

		// migrate bundles
		if storeKey, err := getStoreKey(storeKeys, bundlestypes.StoreKey); err == nil {
			migrateBundlesModule(sdkCtx, cdc, storeKey, bundlesKeeper)
		} else {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

func getStoreKey(storeKeys []storetypes.StoreKey, storeName string) (storetypes.StoreKey, error) {
	for _, k := range storeKeys {
		if k.Name() == storeName {
			return k, nil
		}
	}

	return nil, fmt.Errorf("store key not found: %s", storeName)
}

func migrateFundersModule(sdkCtx sdk.Context, cdc codec.Codec, storeKey storetypes.StoreKey, fundersKeeper fundersKeeper.Keeper) {
	// migrate params
	// TODO: define final prices and initial whitelisted coins
	oldParams := funders.GetParams(sdkCtx, cdc, storeKey)
	fundersKeeper.SetParams(sdkCtx, fundersTypes.Params{
		CoinWhitelist: []*fundersTypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				MinFundingAmount:          oldParams.MinFundingAmount,
				MinFundingAmountPerBundle: oldParams.MinFundingAmountPerBundle,
				CoinWeight:                math.LegacyMustNewDecFromStr("0.06"),
			},
		},
		MinFundingMultiple: oldParams.MinFundingMultiple,
	})

	// migrate fundings
	oldFundings := funders.GetAllFundings(sdkCtx, cdc, storeKey)
	for _, f := range oldFundings {
		fundersKeeper.SetFunding(sdkCtx, &types.Funding{
			FunderAddress:    f.FunderAddress,
			PoolId:           f.PoolId,
			Amounts:          sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.Amount))),
			AmountsPerBundle: sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.AmountPerBundle))),
			TotalFunded:      sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(f.TotalFunded))),
		})
	}
}

func migrateDelegationModule(sdkCtx sdk.Context, cdc codec.Codec, storeKey storetypes.StoreKey, delegationKeeper delegationKeeper.Keeper) {
	// migrate delegation entries
	oldDelegationEntries := delegation.GetAllDelegationEntries(sdkCtx, cdc, storeKey)
	for _, d := range oldDelegationEntries {
		delegationKeeper.SetDelegationEntry(sdkCtx, delegationTypes.DelegationEntry{
			Staker: d.Staker,
			KIndex: d.KIndex,
			Value:  sdk.NewDecCoins(sdk.NewDecCoinFromDec(globalTypes.Denom, d.Value)),
		})
	}

	// migrate delegation data
	oldDelegationData := delegation.GetAllDelegationData(sdkCtx, cdc, storeKey)
	for _, d := range oldDelegationData {
		delegationKeeper.SetDelegationData(sdkCtx, delegationTypes.DelegationData{
			Staker:                     d.Staker,
			CurrentRewards:             sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(d.CurrentRewards))),
			TotalDelegation:            d.TotalDelegation,
			LatestIndexK:               d.LatestIndexK,
			DelegatorCount:             d.DelegatorCount,
			LatestIndexWasUndelegation: d.LatestIndexWasUndelegation,
		})
	}
}

func migrateStakersModule(sdkCtx sdk.Context, cdc codec.Codec, storeKey storetypes.StoreKey, stakersKeeper *stakersKeeper.Keeper) {
	// migrate stakers
	oldStakers := stakers.GetAllStakers(sdkCtx, cdc, storeKey)
	for _, s := range oldStakers {
		stakersKeeper.SetStaker(sdkCtx, stakersTypes.Staker{
			Address:           s.Address,
			Commission:        s.Commission,
			Moniker:           s.Moniker,
			Website:           s.Website,
			Identity:          s.Identity,
			SecurityContact:   s.SecurityContact,
			Details:           s.Details,
			CommissionRewards: sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(s.CommissionRewards))),
		})
	}
}

func migrateBundlesModule(sdkCtx sdk.Context, cdc codec.Codec, storeKey storetypes.StoreKey, bundlesKeeper keeper.Keeper) {
	oldParams := bundles.GetParams(sdkCtx, storeKey, cdc)

	// TODO: define final storage cost prices
	bundlesKeeper.SetParams(sdkCtx, bundlestypes.Params{
		UploadTimeout: oldParams.UploadTimeout,
		StorageCosts: []bundlestypes.StorageCost{
			// Arweave: https://arweave.net/price/1048576 -> 699 winston/byte * 40 USD/AR * 1.5 / 10**12
			{StorageProviderId: 1, Cost: math.LegacyMustNewDecFromStr("0.000000004194")},
			// Irys: https://node1.bundlr.network/price/1048576 -> 1048 winston/byte * 40 USD/AR * 1.5 / 10**12
			{StorageProviderId: 2, Cost: math.LegacyMustNewDecFromStr("0.000000006288")},
			// KYVE Storage Provider: zero since it is free to use for testnet participants
			{StorageProviderId: 3, Cost: math.LegacyMustNewDecFromStr("0.00")},
		},
		NetworkFee: oldParams.NetworkFee,
		MaxPoints:  oldParams.MaxPoints,
	})
}
