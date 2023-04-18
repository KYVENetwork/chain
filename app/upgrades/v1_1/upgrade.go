package v1_1

import (
	"github.com/KYVENetwork/chain/app/upgrades/v1_1/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Auth
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingExported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	// Bundles
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Delegation
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// IBC Transfer
	transferKeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	transferTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	// ICA Controller
	icaControllerKeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	icaControllerTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	// ICA Host
	icaHostKeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	// Stakers
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	// Upgrade
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.BinaryCodec,
	bundlesStoreKey storeTypes.StoreKey,
	delegationStoreKey storeTypes.StoreKey,
	stakersStoreKey storeTypes.StoreKey,
	accountKeeper authKeeper.AccountKeeper,
	icaControllerKeeper icaControllerKeeper.Keeper,
	icaHostKeeper icaHostKeeper.Keeper,
	transferKeeper transferKeeper.Keeper,
) upgradeTypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradeTypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		if ctx.ChainID() == MainnetChainID {
			for _, address := range InvestorAccounts {
				AdjustInvestorVesting(ctx, accountKeeper, sdk.MustAccAddressFromBech32(address))
			}
		}

		EnableIBCTransfers(ctx, transferKeeper)
		if ctx.ChainID() == TestnetChainID {
			InitialiseICAControllerParams(ctx, icaControllerKeeper)
			UpdateICAHostParams(ctx, icaHostKeeper)
		}

		MigrateBundlesParameters(ctx, cdc, bundlesStoreKey)
		MigrateDelegationParameters(ctx, cdc, delegationStoreKey)

		MigrateStakerMetadata(ctx, cdc, stakersStoreKey)
		MigrateStakerCommissionEntries(ctx, cdc, stakersStoreKey)

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

// EnableIBCTransfers updates the parameters of the IBC Transfer module to
// allow both sending and receiving of IBC tokens. Since the default parameters
// of the module have everything enabled, we simply switch to the defaults.
func EnableIBCTransfers(ctx sdk.Context, keeper transferKeeper.Keeper) {
	params := transferTypes.DefaultParams()
	keeper.SetParams(ctx, params)
}

// InitialiseICAControllerParams initialises the parameters of the ICA
// Controller module. This is because when upgrading from IBC v5 to v6 on Kaon
// (our testnet), we didn't initialise the parameters.
func InitialiseICAControllerParams(ctx sdk.Context, keeper icaControllerKeeper.Keeper) {
	params := icaControllerTypes.DefaultParams()
	keeper.SetParams(ctx, params)
}

// UpdateICAHostParams updates the parameters of the ICA Host module to allow
// all messages to be called. This configuration is the new default, which
// wasn't updated on Kaon (our testnet).
func UpdateICAHostParams(ctx sdk.Context, keeper icaHostKeeper.Keeper) {
	params := icaHostTypes.DefaultParams()
	keeper.SetParams(ctx, params)
}

// MigrateBundlesParameters re-encodes all parameters of the Bundles module
// that were converted to sdk.Dec.
func MigrateBundlesParameters(ctx sdk.Context, cdc codec.BinaryCodec, bundlesStoreKey storeTypes.StoreKey) {
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

// MigrateDelegationParameters re-encodes all parameters of the Delegation
// module that were converted to sdk.Dec.
func MigrateDelegationParameters(ctx sdk.Context, cdc codec.BinaryCodec, delegationStoreKey storeTypes.StoreKey) {
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

// MigrateStakerMetadata migrates all existing staker metadata. The `Logo`
// field has been deprecated and replaced by the `Identity` field. This new
// field must be a valid hex string; therefore, must be set to empty for now.
func MigrateStakerMetadata(ctx sdk.Context, cdc codec.BinaryCodec, stakerStoreKey storeTypes.StoreKey) {
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

// MigrateStakerCommissionEntries re-encodes all CommissionChangeEntry
// entries from the Stakers module that were converted to sdk.Dec.
func MigrateStakerCommissionEntries(ctx sdk.Context, cdc codec.BinaryCodec, stakerStoreKey storeTypes.StoreKey) {
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
