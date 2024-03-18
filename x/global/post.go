package global

import (
	sdkErrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// FeeGrant
	feeGrantKeeper "cosmossdk.io/x/feegrant/keeper"
	// Global
	"github.com/KYVENetwork/chain/x/global/keeper"
)

// RefundFeeDecorator

type RefundFeeDecorator struct {
	bankKeeper     bankKeeper.Keeper
	feeGrantKeeper feeGrantKeeper.Keeper
	globalKeeper   keeper.Keeper
}

func NewRefundFeeDecorator(bk bankKeeper.Keeper, fk feeGrantKeeper.Keeper, gk keeper.Keeper) RefundFeeDecorator {
	return RefundFeeDecorator{
		bankKeeper:     bk,
		feeGrantKeeper: fk,
		globalKeeper:   gk,
	}
}

func (rfd RefundFeeDecorator) PostHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, success bool, next sdk.PostHandler) (newCtx sdk.Context, err error) {
	// Ensure that this is a fee transaction.
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkErrors.Wrap(errorsTypes.ErrTxDecode, "Tx must be a FeeTx")
	}

	// Return early if the transaction fee is zero (nothing to refund)
	// or there are more than one message (can't refund).
	fee := feeTx.GetFee()
	msgs := feeTx.GetMsgs()
	if fee.IsZero() || len(msgs) != 1 {
		return next(ctx, tx, simulate, success)
	}

	// Find the refund percentage based on the transaction message type.
	refundPercentage := sdk.ZeroDec()
	gasRefunds := rfd.globalKeeper.GetGasRefunds(ctx)
	for _, refund := range gasRefunds {
		if sdk.MsgTypeURL(msgs[0]) == refund.Type {
			refundPercentage = refund.Fraction
			break
		}
	}

	// Return early if the refund percentage is zero.
	if refundPercentage.IsZero() {
		return next(ctx, tx, simulate, success)
	}

	// Calculate the refund amount.
	refund := sdk.NewCoins()
	for _, coin := range fee {
		amount := sdk.NewDecFromInt(coin.Amount).Mul(refundPercentage)
		refund = refund.Add(sdk.NewCoin(coin.Denom, amount.TruncateInt()))
	}

	// Send the refund back to this transaction's fee payer.
	account, err := GetFeeAccount(ctx, feeTx, rfd.feeGrantKeeper)
	if err != nil {
		return ctx, err
	}
	err = rfd.bankKeeper.SendCoinsFromModuleToAccount(ctx, authTypes.FeeCollectorName, account, refund)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate, success)
}
