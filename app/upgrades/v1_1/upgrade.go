package v1_1

import (
	"github.com/KYVENetwork/chain/app/upgrades/v1_1/types"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingExported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.BinaryCodec,
	stakersStoreKey storetypes.StoreKey,
	bundlesStoreKey storetypes.StoreKey,
	delegationStoreKey storetypes.StoreKey,
	accountKeeper authKeeper.AccountKeeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if ctx.ChainID() == MainnetChainID {
			for _, address := range InvestorAccounts {
				AdjustInvestorVesting(ctx, accountKeeper, sdk.MustAccAddressFromBech32(address))
			}
		}

		MigrateStakerMetadata(ctx, cdc, stakersStoreKey)
		MigrateStakerCommissionEntries(ctx, cdc, stakersStoreKey)

		MigrateBundleParameters(ctx, cdc, bundlesStoreKey)
		MigrateDelegationParameters(ctx, cdc, delegationStoreKey)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// AdjustInvestorVesting correctly adjusts the vesting schedules of investors
// from our second funding round. In genesis, the accounts were set up with an
// 18-month cliff instead of a 6-month cliff.
func AdjustInvestorVesting(ctx sdk.Context, keeper authKeeper.AccountKeeper, address sdk.AccAddress) {
	rawAccount := keeper.GetAccount(ctx, address)
	account := rawAccount.(vestingExported.VestingAccount)

	baseAccount := authTypes.NewBaseAccount(
		account.GetAddress(), account.GetPubKey(), account.GetAccountNumber(), account.GetSequence(),
	)
	updatedAccount := vestingTypes.NewContinuousVestingAccount(
		baseAccount, account.GetOriginalVesting(), StartTime, EndTime,
	)

	keeper.SetAccount(ctx, updatedAccount)
}

// MigrateStakerMetadata migrates all existing staker metadata. The `Logo`
// field has been deprecated and replaced by the `Identity` field. This new
// field must be a valid hex string; therefore, must be set to empty for now.
func MigrateStakerMetadata(ctx sdk.Context, cdc codec.BinaryCodec, stakerStoreKey storetypes.StoreKey) {
	store := prefix.NewStore(ctx.KVStore(stakerStoreKey), stakersTypes.StakerKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	oldStakers := make([]types.OldStaker, 0)

	for ; iterator.Valid(); iterator.Next() {
		var val types.OldStaker
		cdc.MustUnmarshal(iterator.Value(), &val)
		oldStakers = append(oldStakers, val)
	}

	for _, oldStaker := range oldStakers {

		commission, err := sdk.NewDecFromStr(oldStaker.Commission)
		if err != nil {
			commission = stakersTypes.DefaultCommission
		}

		newStaker := stakersTypes.Staker{
			Address:         oldStaker.Address,
			Commission:      commission,
			Moniker:         oldStaker.Moniker,
			Website:         oldStaker.Website,
			Identity:        "",
			SecurityContact: "",
			Details:         "",
		}

		b := cdc.MustMarshal(&newStaker)
		store.Set(stakersTypes.StakerKey(newStaker.Address), b)
	}
}

// MigrateStakerCommissionEntries re-encodes the CommissionChangeEntry fields which got converted to sdk.Dec
func MigrateStakerCommissionEntries(ctx sdk.Context, cdc codec.BinaryCodec, stakerStoreKey storetypes.StoreKey) {
	store := prefix.NewStore(ctx.KVStore(stakerStoreKey), stakersTypes.CommissionChangeEntryKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	oldCommissionChangeEntries := make([]types.OldCommissionChangeEntry, 0)

	for ; iterator.Valid(); iterator.Next() {
		var val types.OldCommissionChangeEntry
		cdc.MustUnmarshal(iterator.Value(), &val)
		oldCommissionChangeEntries = append(oldCommissionChangeEntries, val)
	}

	for _, oldCommissionEntry := range oldCommissionChangeEntries {

		commission, err := sdk.NewDecFromStr(oldCommissionEntry.Commission)
		if err != nil {
			commission = stakersTypes.DefaultCommission
		}

		newCommissionChangeEntry := stakersTypes.CommissionChangeEntry{
			Index:        oldCommissionEntry.Index,
			Staker:       oldCommissionEntry.Staker,
			Commission:   commission,
			CreationDate: oldCommissionEntry.CreationDate,
		}

		b := cdc.MustMarshal(&newCommissionChangeEntry)
		store.Set(stakersTypes.CommissionChangeEntryKey(newCommissionChangeEntry.Index), b)
	}
}

// MigrateBundleParameters re-encodes the params fields which got converted to sdk.Dec
func MigrateBundleParameters(ctx sdk.Context, cdc codec.BinaryCodec, bundlesStoreKey storetypes.StoreKey) {
	store := ctx.KVStore(bundlesStoreKey)
	bz := store.Get(bundlesTypes.ParamsKey)

	var oldParams types.OldBundlesParams
	cdc.MustUnmarshal(bz, &oldParams)

	newParams := bundlesTypes.Params{
		UploadTimeout: oldParams.UploadTimeout,
		StorageCost:   oldParams.StorageCost,
		NetworkFee:    sdk.MustNewDecFromStr(oldParams.NetworkFee),
		MaxPoints:     oldParams.MaxPoints,
	}

	bz = cdc.MustMarshal(&newParams)
	store.Set(bundlesTypes.ParamsKey, bz)
}

// MigrateDelegationParameters re-encodes the params fields which got converted to sdk.Dec
func MigrateDelegationParameters(ctx sdk.Context, cdc codec.BinaryCodec, delegationStoreKey storetypes.StoreKey) {
	store := ctx.KVStore(delegationStoreKey)
	bz := store.Get(delegationTypes.ParamsKey)

	var oldParams types.OldDelegationParams
	cdc.MustUnmarshal(bz, &oldParams)

	newParams := delegationTypes.Params{
		UnbondingDelegationTime: oldParams.UnbondingDelegationTime,
		RedelegationCooldown:    oldParams.RedelegationCooldown,
		RedelegationMaxAmount:   oldParams.RedelegationMaxAmount,
		VoteSlash:               sdk.MustNewDecFromStr(oldParams.VoteSlash),
		UploadSlash:             sdk.MustNewDecFromStr(oldParams.UploadSlash),
		TimeoutSlash:            sdk.MustNewDecFromStr(oldParams.TimeoutSlash),
	}

	bz = cdc.MustMarshal(&newParams)
	store.Set(delegationTypes.ParamsKey, bz)
}
