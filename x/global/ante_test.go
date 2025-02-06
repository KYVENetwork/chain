package global_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// Global
	"github.com/KYVENetwork/chain/x/global"
	"github.com/KYVENetwork/chain/x/global/types"
	// Stakers
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	// Staking
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

/*

TEST CASES - DeductFeeDecorator

* Invalid transaction.
* consensusGasPrice = 0.0; validatorGasPrice = 0.0 - deliverTX
* consensusGasPrice = 0.0; validatorGasPrice = 0.0 - checkTX
* consensusGasPrice = 1.0; validatorGasPrice = 0.0 - deliverTX - not enough fees
* consensusGasPrice = 1.0; validatorGasPrice = 0.0 - deliverTX - enough fees
* consensusGasPrice = 1.0; validatorGasPrice = 0.0 - checkTx - not enough fees
* consensusGasPrice = 1.0; validatorGasPrice = 0.0 - checkTx - enough fees
* consensusGasPrice = 1.0; validatorGasPrice = 2.0 - deliverTX - not enough fees
* consensusGasPrice = 1.0; validatorGasPrice = 2.0 - deliverTX - not enough fees for validator but enough for consensus.
* consensusGasPrice = 1.0; validatorGasPrice = 2.0 - checkTx - not enough fees
* consensusGasPrice = 1.0; validatorGasPrice = 2.0 - checkTx - not enough fees for validator but enough for consensus.

*/

var _ = Describe("DeductFeeDecorator", Ordered, func() {
	s := i.NewCleanChain()
	encodingConfig := BuildEncodingConfig()
	dfd := global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper, s.App().StakingKeeper)
	denom, _ := s.App().StakingKeeper.BondDenom(s.Ctx())

	accountBalanceBefore := s.GetBalanceFromAddress(i.DUMMY[0])
	collectorBalanceBefore := s.GetBalanceFromModule(authTypes.FeeCollectorName)

	BeforeEach(func() {
		s = i.NewCleanChain()
		encodingConfig = BuildEncodingConfig()
		denom, _ = s.App().StakingKeeper.BondDenom(s.Ctx())
		dfd = global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper, s.App().StakingKeeper)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Invalid transaction.", func() {
		// ARRANGE
		dfd := global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper, s.App().StakingKeeper)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx(), &InvalidTx{}, false, AnteNextFn)

		// ASSERT
		Expect(err).Should(HaveOccurred())
	})

	It("consensusGasPrice = 0.0; validatorGasPrice = 0.0 - deliverTX", func() {
		// ARRANGE
		dfd := global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper, s.App().StakingKeeper)

		denom, _ := s.App().StakingKeeper.BondDenom(s.Ctx())
		tx := BuildTestTx(math.ZeroInt(), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx().WithIsCheckTx(false), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter))
	})

	It("consensusGasPrice = 0.0; validatorGasPrice = 0.0 - checkTX", func() {
		// ARRANGE
		dfd := global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper, s.App().StakingKeeper)

		denom, _ := s.App().StakingKeeper.BondDenom(s.Ctx())
		tx := BuildTestTx(math.ZeroInt(), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx().WithIsCheckTx(true), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 0.0 - deliverTX - not enough fees", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)
		tx := BuildTestTx(math.ZeroInt(), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx().WithIsCheckTx(false), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(HaveOccurred())
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 0.0 - deliverTX - enough fees", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx().WithIsCheckTx(false), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + 200_000))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter - 200_000))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 0.0 - checkTx - not enough fees", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)
		tx := BuildTestTx(math.ZeroInt(), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx().WithIsCheckTx(true), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(HaveOccurred())
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 0.0 - checkTx - enough fees", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx().WithIsCheckTx(true), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + 200_000))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter - 200_000))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 2.0 - deliverTX - not enough fees", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)

		ctx := s.Ctx().WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(denom, math.NewInt(2))))
		s.SetCtx(ctx)
		tx := BuildTestTx(math.ZeroInt(), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(ctx.WithIsCheckTx(false), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(HaveOccurred())
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 2.0 - deliverTX - not enough fees for validator but enough for consensus.", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)

		ctx := s.Ctx().WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(denom, math.NewInt(2))))
		s.SetCtx(ctx)
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(ctx.WithIsCheckTx(false), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + 200_000))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter - 200_000))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 2.0 - checkTx - not enough fees", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)

		ctx := s.Ctx().WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(denom, math.NewInt(2))))
		s.SetCtx(ctx)
		tx := BuildTestTx(math.ZeroInt(), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(ctx.WithIsCheckTx(true), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(HaveOccurred())
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter))
	})

	It("consensusGasPrice = 1.0; validatorGasPrice = 2.0 - checkTx - not enough fees for validator but enough for consensus.", func() {
		// ARRANGE
		params := types.DefaultParams()
		params.MinGasPrice = math.LegacyOneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)

		ctx := s.Ctx().WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(denom, math.NewInt(2))))
		s.SetCtx(ctx)
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(ctx.WithIsCheckTx(true), tx, false, AnteNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(err).Should(HaveOccurred())
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceBefore).To(Equal(collectorBalanceAfter))
	})
})

/*

TEST CASES - GasAdjustmentDecorator

* Empty transaction.
* Transaction with a normal message.
* Transaction with an adjusted message.
* Transaction with multiple adjusted messages.
* Transaction with multiple normal and multiple adjusted messages.

*/

var _ = Describe("GasAdjustmentDecorator", Ordered, func() {
	s := i.NewCleanChain()
	encodingConfig := BuildEncodingConfig()

	// NOTE: This will change as implementation changes.
	// TODO: Why does this change as the implementation changes?
	BaseCost := 62974

	BeforeEach(func() {
		s = i.NewCleanChain()

		params := types.DefaultParams()
		params.GasAdjustments = []types.GasAdjustment{
			{
				Type:   "/cosmos.staking.v1beta1.MsgCreateValidator",
				Amount: 2000,
			},
			{
				Type:   "/kyve.funders.v1beta1.MsgCreateFunder",
				Amount: 1000,
			},
		}
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Empty transaction.", func() {
		// ARRANGE
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		tx := txBuilder.GetTx()

		gad := global.NewGasAdjustmentDecorator(s.App().GlobalKeeper)

		// ACT
		_, err := gad.AnteHandle(s.Ctx(), tx, false, AnteNextFn)

		// ASSERT
		Expect(err).ToNot(HaveOccurred())
		Expect(s.Ctx().GasMeter().GasConsumed()).To(BeEquivalentTo(BaseCost))
	})

	It("Transaction with a normal message.", func() {
		// ARRANGE
		msg := bankTypes.MsgSend{}

		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		_ = txBuilder.SetMsgs(&msg)
		tx := txBuilder.GetTx()

		gad := global.NewGasAdjustmentDecorator(s.App().GlobalKeeper)

		// ACT
		_, err := gad.AnteHandle(s.Ctx(), tx, false, AnteNextFn)

		// ASSERT
		Expect(err).ToNot(HaveOccurred())
		Expect(s.Ctx().GasMeter().GasConsumed()).To(BeEquivalentTo(BaseCost))
	})

	It("Transaction with an adjusted message.", func() {
		// ARRANGE
		msg := stakingTypes.MsgCreateValidator{}

		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		_ = txBuilder.SetMsgs(&msg)
		tx := txBuilder.GetTx()

		gad := global.NewGasAdjustmentDecorator(s.App().GlobalKeeper)

		// ACT
		_, err := gad.AnteHandle(s.Ctx(), tx, false, AnteNextFn)

		// ASSERT
		Expect(err).ToNot(HaveOccurred())
		Expect(s.Ctx().GasMeter().GasConsumed()).To(BeEquivalentTo(BaseCost + 2000))
	})

	It("Transaction with multiple adjusted messages.", func() {
		// ARRANGE
		firstMsg := stakingTypes.MsgCreateValidator{}
		secondMsg := fundersTypes.MsgCreateFunder{}

		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		_ = txBuilder.SetMsgs(&firstMsg, &secondMsg)
		tx := txBuilder.GetTx()

		gad := global.NewGasAdjustmentDecorator(s.App().GlobalKeeper)

		// ACT
		_, err := gad.AnteHandle(s.Ctx(), tx, false, AnteNextFn)

		// ASSERT
		Expect(err).ToNot(HaveOccurred())
		Expect(s.Ctx().GasMeter().GasConsumed()).To(BeEquivalentTo(BaseCost + 3000))
	})

	It("Transaction with multiple normal and multiple adjusted messages.", func() {
		// ARRANGE
		firstMsg := stakersTypes.MsgJoinPool{}
		secondMsg := fundersTypes.MsgCreateFunder{}
		thirdMsg := stakingTypes.MsgCreateValidator{}

		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		_ = txBuilder.SetMsgs(&firstMsg, &secondMsg, &thirdMsg)
		tx := txBuilder.GetTx()

		gad := global.NewGasAdjustmentDecorator(s.App().GlobalKeeper)

		// ACT
		_, err := gad.AnteHandle(s.Ctx(), tx, false, AnteNextFn)

		// ASSERT
		Expect(err).ToNot(HaveOccurred())
		Expect(s.Ctx().GasMeter().GasConsumed()).To(BeEquivalentTo(BaseCost + 3000))
	})
})
