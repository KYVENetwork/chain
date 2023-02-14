package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// FeeGrant
	feeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	// Global
	"github.com/KYVENetwork/chain/x/global"
	globalKeeper "github.com/KYVENetwork/chain/x/global/keeper"
	// Gov
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	// IBC
	ibcAnte "github.com/cosmos/ibc-go/v6/modules/core/ante"
	ibcKeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	// Staking
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// https://github.com/cosmos/cosmos-sdk/blob/release/v0.46.x/x/auth/ante/ante.go#L25

func NewAnteHandler(
	accountKeeper authKeeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	feeGrantKeeper feeGrantKeeper.Keeper,
	globalKeeper globalKeeper.Keeper,
	govKeeper govKeeper.Keeper,
	ibcKeeper *ibcKeeper.Keeper,
	stakingKeeper stakingKeeper.Keeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) (sdk.AnteHandler, error) {
	deductFeeDecorator := global.NewDeductFeeDecorator(accountKeeper, bankKeeper, feeGrantKeeper, globalKeeper, stakingKeeper)

	gasAdjustmentDecorator := global.NewGasAdjustmentDecorator(globalKeeper)

	initialDepositDecorator := global.NewInitialDepositDecorator(globalKeeper, govKeeper)

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		gasAdjustmentDecorator,
		ante.NewExtensionOptionsDecorator(nil),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(accountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(accountKeeper),
		deductFeeDecorator,
		ante.NewSetPubKeyDecorator(accountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(accountKeeper),
		ante.NewSigGasConsumeDecorator(accountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(accountKeeper, signModeHandler),
		ante.NewIncrementSequenceDecorator(accountKeeper),
		ibcAnte.NewRedundantRelayDecorator(ibcKeeper),
		initialDepositDecorator,
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}

//

func NewPostHandler(
	bankKeeper bankKeeper.Keeper,
	feeGrantKeeper feeGrantKeeper.Keeper,
	globalKeeper globalKeeper.Keeper,
) (sdk.AnteHandler, error) {
	refundFeeDecorator := global.NewRefundFeeDecorator(bankKeeper, feeGrantKeeper, globalKeeper)

	postDecorators := []sdk.AnteDecorator{
		refundFeeDecorator,
	}

	return sdk.ChainAnteDecorators(postDecorators...), nil
}
