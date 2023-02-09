package global

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// FeeGrant
	feeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	// Global
	"github.com/KYVENetwork/chain/x/global/keeper"
	// Gov
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	legacyGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
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

// InitialDepositDecorator

// The InitialDepositDecorator is responsible for checking
// if the submit-proposal message also provides the required
// minimum deposit. Otherwise, the message is rejected.
type InitialDepositDecorator struct {
	globalKeeper keeper.Keeper
	govKeeper    govKeeper.Keeper
}

func NewInitialDepositDecorator(globalKeeper keeper.Keeper, govKeeper govKeeper.Keeper) InitialDepositDecorator {
	return InitialDepositDecorator{
		globalKeeper: globalKeeper,
		govKeeper:    govKeeper,
	}
}

func (idd InitialDepositDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// NOTE: This is Tendermint specific.
	if ctx.BlockHeight() <= 1 {
		return next(ctx, tx, simulate)
	}

	minInitialDepositRatio := idd.globalKeeper.GetMinInitialDepositRatio(ctx)
	depositParams := idd.govKeeper.GetDepositParams(ctx)

	requiredDeposit := sdk.NewCoins()
	for _, coin := range depositParams.MinDeposit {
		amount := sdk.NewDecFromInt(coin.Amount).Mul(minInitialDepositRatio).TruncateInt()
		requiredDeposit = requiredDeposit.Add(sdk.NewCoin(coin.Denom, amount))
	}

	for _, rawMsg := range tx.GetMsgs() {
		initialDeposit := sdk.NewCoins()
		throwError := false

		if sdk.MsgTypeURL(rawMsg) == sdk.MsgTypeURL(&legacyGovTypes.MsgSubmitProposal{}) {
			// cosmos.gov.v1beta1.MsgSubmitProposal
			if legacyMsg, ok := rawMsg.(*legacyGovTypes.MsgSubmitProposal); ok {
				initialDeposit = legacyMsg.GetInitialDeposit()
				throwError = !initialDeposit.IsAllGTE(requiredDeposit)
			}
		} else if sdk.MsgTypeURL(rawMsg) == sdk.MsgTypeURL(&govTypes.MsgSubmitProposal{}) {
			// cosmos.gov.v1.MsgSubmitProposal
			if msg, ok := rawMsg.(*govTypes.MsgSubmitProposal); ok {
				initialDeposit = msg.GetInitialDeposit()
				throwError = !initialDeposit.IsAllGTE(requiredDeposit)
			}
		}

		if throwError {
			return ctx, errors.Wrapf(
				errorsTypes.ErrLogic, "minimum deposit is too small - was (%s), need (%s)",
				initialDeposit, requiredDeposit,
			)
		}
	}

	return next(ctx, tx, simulate)
}
