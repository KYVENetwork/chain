package global

import (
	"bytes"
	sdkmath "cosmossdk.io/math"
	"math"

	sdkErrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"

	// Auth
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	// FeeGrant
	feeGrantKeeper "cosmossdk.io/x/feegrant/keeper"
	// Global
	"github.com/KYVENetwork/chain/x/global/keeper"
	// Staking
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

func GetFeeAccount(ctx sdk.Context, tx sdk.FeeTx, feeGrantKeeper feeGrantKeeper.Keeper) (sdk.AccAddress, error) {
	fee := tx.GetFee()
	feePayer := tx.FeePayer()
	feeGranter := tx.FeeGranter()

	account := feePayer
	if feeGranter != nil && !bytes.Equal(feeGranter, feePayer) {
		err := feeGrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, tx.GetMsgs())
		if err != nil {
			return nil, sdkErrors.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
		}

		account = feeGranter
	}

	return account, nil
}

// BuildTxFeeChecker ensures that the configured minimum gas price is met.
// In contrast to
// https://github.com/cosmos/cosmos-sdk/blob/release/v0.46.x/x/auth/ante/validator_tx_fee.go#L12
// this code runs within the consensus layer.
func BuildTxFeeChecker(ctx sdk.Context, fk keeper.Keeper, sk stakingKeeper.Keeper) ante.TxFeeChecker {
	return func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
		bondDenom, err := sk.BondDenom(ctx)
		if err != nil {
			return nil, 0, sdkErrors.Wrap(errorsTypes.ErrNotFound, "failed to get bond denom")
		}
		consensusMinGasPrices := sdk.NewDecCoins(sdk.NewDecCoinFromDec(bondDenom, fk.GetMinGasPrice(ctx)))

		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return nil, 0, sdkErrors.Wrap(errorsTypes.ErrTxDecode, "Tx must be a FeeTx")
		}

		feeCoins := feeTx.GetFee()
		gas := feeTx.GetGas()

		validatorMinGasPrices := ctx.MinGasPrices()

		requiredFees := make(sdk.Coins, len(consensusMinGasPrices))

		// Determine the required fees by multiplying each required minimum gas
		// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
		glDec := sdkmath.LegacyNewDec(int64(gas))
		for i, gp := range consensusMinGasPrices {
			fee := gp.Amount.Mul(glDec)
			requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		if ctx.IsCheckTx() {
			validatorFees := make(sdk.Coins, len(validatorMinGasPrices))
			for i, gp := range validatorMinGasPrices {
				fee := gp.Amount.Mul(glDec)
				validatorFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			requiredFees = requiredFees.Max(validatorFees)
		}

		if !requiredFees.IsZero() && !feeCoins.IsAnyGTE(requiredFees) {
			return nil, 0, sdkErrors.Wrapf(errorsTypes.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
		}

		priority := getTxPriority(feeCoins, int64(gas))
		return feeCoins, priority, nil
	}
}

// https://github.com/cosmos/cosmos-sdk/blob/release/v0.46.x/x/auth/ante/validator_tx_fee.go#L51
// As the default DeductFeeDecorator is overwritten, this is the place to add a custom priority.
// Although this is calculated within "consensus-code" the priority itself gets only used
// for mem-pool ordering.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
