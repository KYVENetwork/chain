package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// FeeGrant
	feeGrantKeeper "cosmossdk.io/x/feegrant/keeper"
	// Global
	"github.com/KYVENetwork/chain/x/global/keeper"
	// Staking
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

// DeductFeeDecorator

// The DeductFeeDecorator is responsible for the
// consensus minimum gas price.
// Validators can still choose their own (higher) gas prices.
type DeductFeeDecorator struct {
	accountKeeper  authKeeper.AccountKeeper
	bankKeeper     bankKeeper.Keeper
	feeGrantKeeper feeGrantKeeper.Keeper
	globalKeeper   keeper.Keeper
	stakingKeeper  stakingKeeper.Keeper
}

func NewDeductFeeDecorator(ak authKeeper.AccountKeeper, bk bankKeeper.Keeper, fk feeGrantKeeper.Keeper, gk keeper.Keeper, sk stakingKeeper.Keeper) DeductFeeDecorator {
	return DeductFeeDecorator{
		accountKeeper:  ak,
		bankKeeper:     bk,
		feeGrantKeeper: fk,
		globalKeeper:   gk,
		stakingKeeper:  sk,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// NOTE: This is Tendermint specific.
	var tfc ante.TxFeeChecker
	if ctx.BlockHeight() > 1 {
		tfc = BuildTxFeeChecker(ctx, dfd.globalKeeper, dfd.stakingKeeper)
	}

	internalDfd := ante.NewDeductFeeDecorator(dfd.accountKeeper, dfd.bankKeeper, dfd.feeGrantKeeper, tfc)

	return internalDfd.AnteHandle(ctx, tx, simulate, next)
}

// GasAdjustmentDecorator

// The GasAdjustmentDecorator allows to add additional gas-consumption
// to message types, making transactions which should only be used rerely
// more expensive.
type GasAdjustmentDecorator struct {
	globalKeeper keeper.Keeper
}

func NewGasAdjustmentDecorator(gk keeper.Keeper) GasAdjustmentDecorator {
	return GasAdjustmentDecorator{
		globalKeeper: gk,
	}
}

func (gad GasAdjustmentDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	gasAdjustments := gad.globalKeeper.GetGasAdjustments(ctx)

	for _, msg := range tx.GetMsgs() {
		for _, adjustment := range gasAdjustments {
			if sdk.MsgTypeURL(msg) == adjustment.Type {
				ctx.GasMeter().ConsumeGas(adjustment.Amount, adjustment.Type)
				break
			}
		}
	}

	return next(ctx, tx, simulate)
}
