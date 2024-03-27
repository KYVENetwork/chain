package util

import (
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TransferFromAddressToAddress sends tokens from the given address to a specified address.
func TransferFromAddressToAddress(
	bankKeeper BankKeeper,
	ctx sdk.Context,
	fromAddress string,
	toAddress string,
	amount uint64,
) error {
	sender, errSenderAddress := sdk.AccAddressFromBech32(fromAddress)
	if errSenderAddress != nil {
		return errSenderAddress
	}

	recipient, errRecipientAddress := sdk.AccAddressFromBech32(toAddress)
	if errRecipientAddress != nil {
		return errRecipientAddress
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(amount)))
	err := bankKeeper.SendCoins(ctx, sender, recipient, coins)
	return err
}

// TransferFromModuleToAddress sends tokens from the given module to a specified address.
func TransferFromModuleToAddress(
	bankKeeper BankKeeper,
	ctx sdk.Context,
	module string,
	address string,
	amount uint64,
) error {
	recipient, errAddress := sdk.AccAddressFromBech32(address)
	if errAddress != nil {
		return errAddress
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(amount)))
	err := bankKeeper.SendCoinsFromModuleToAccount(ctx, module, recipient, coins)
	return err
}

// TransferFromAddressToModule sends tokens from a specified address to the given module.
func TransferFromAddressToModule(
	bankKeeper BankKeeper,
	ctx sdk.Context,
	address string,
	module string,
	amount uint64,
) error {
	sender, errAddress := sdk.AccAddressFromBech32(address)
	if errAddress != nil {
		return errAddress
	}
	coins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(amount)))

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, sender, module, coins)
	return err
}

// TransferFromModuleToModule sends tokens from a specified module to the given module.
func TransferFromModuleToModule(
	bankKeeper BankKeeper,
	ctx sdk.Context,
	fromModule string,
	toModule string,
	amount uint64,
) error {
	coins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(amount)))
	err := bankKeeper.SendCoinsFromModuleToModule(ctx, fromModule, toModule, coins)
	return err
}

// TransferFromAddressToTreasury sends tokens from a given address to the treasury (community spend pool).
func TransferFromAddressToTreasury(distrKeeper DistributionKeeper, ctx sdk.Context, address string, amount uint64) error {
	sender, errAddress := sdk.AccAddressFromBech32(address)
	if errAddress != nil {
		return errAddress
	}
	coins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(amount)))

	if err := distrKeeper.FundCommunityPool(ctx, coins, sender); err != nil {
		return err
	}

	return nil
}

// TransferFromModuleToTreasury sends tokens from a module to the treasury (community spend pool).
func TransferFromModuleToTreasury(
	accountKeeper AccountKeeper,
	distrKeeper DistributionKeeper,
	ctx sdk.Context,
	module string,
	amount uint64,
) error {
	sender := accountKeeper.GetModuleAddress(module)
	coins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(amount)))

	if err := distrKeeper.FundCommunityPool(ctx, coins, sender); err != nil {
		return err
	}

	return nil
}
