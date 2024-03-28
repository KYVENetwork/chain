package global_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Global
	"github.com/KYVENetwork/chain/x/global"
	"github.com/KYVENetwork/chain/x/global/types"
)

/*

TEST CASES - RefundFeeDecorator

* Non-refundable message
* Refund 0%
* Refund 10%
* Refund 2/3 %
* Refund 100%
* Don't refund multiple

*/

var _ = Describe("RefundFeeDecorator", Ordered, func() {
	s := i.NewCleanChain()
	encodingConfig := BuildEncodingConfig()
	rfd := global.NewRefundFeeDecorator(s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper)
	dfd := global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper, *s.App().StakingKeeper)
	denom, _ := s.App().StakingKeeper.BondDenom(s.Ctx())

	accountBalanceBefore := s.GetBalanceFromAddress(i.DUMMY[0])

	BeforeEach(func() {
		s = i.NewCleanChain()

		accountBalanceBefore = s.GetBalanceFromAddress(i.DUMMY[0])
		rfd = global.NewRefundFeeDecorator(s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper)
		dfd = global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, s.App().GlobalKeeper, *s.App().StakingKeeper)

		denom, _ = s.App().StakingKeeper.BondDenom(s.Ctx())

		params := types.DefaultParams()
		params.GasRefunds = []types.GasRefund{
			{
				Type:     "/kyve.bundles.v1beta1.MsgSubmitBundleProposal",
				Fraction: math.LegacyNewDec(1).QuoInt64(10),
			},
			{
				Type:     "/kyve.bundles.v1beta1.MsgVoteBundleProposal",
				Fraction: math.LegacyOneDec(),
			},
			{
				Type:     "/kyve.bundles.v1beta1.MsgSkipUploaderRole",
				Fraction: math.LegacyZeroDec(),
			},
			{
				Type:     "/kyve.stakers.v1beta1.MsgCreateStaker",
				Fraction: math.LegacyNewDec(2).QuoInt64(3),
			},
		}
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Non-refundable message", func() {
		// ARRANGE
		msg := bundlesTypes.MsgClaimUploaderRole{Creator: i.ALICE}
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		gasLimit := uint64(200_000)
		txBuilder.SetGasLimit(gasLimit)
		fees := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1).MulRaw(int64(gasLimit))))
		txBuilder.SetFeeAmount(fees)
		_ = txBuilder.SetMsgs(&msg)
		tx := txBuilder.GetTx()

		// ACT
		_, errAnte := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		_, errPost := rfd.PostHandle(s.Ctx(), tx, false, true, PostNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.ALICE)
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(errAnte).Should(Not(HaveOccurred()))
		Expect(errPost).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + 200_000))
		Expect(collectorBalanceAfter).To(Equal(uint64(200_000)))
	})

	It("Refund 0%", func() {
		// ARRANGE
		msg := bundlesTypes.MsgSkipUploaderRole{Creator: i.ALICE}
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		gasLimit := uint64(200_000)
		txBuilder.SetGasLimit(gasLimit)
		fees := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1).MulRaw(int64(gasLimit))))
		txBuilder.SetFeeAmount(fees)
		_ = txBuilder.SetMsgs(&msg)
		tx := txBuilder.GetTx()

		// ACT
		_, errAnte := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		_, errPost := rfd.PostHandle(s.Ctx(), tx, false, true, PostNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.ALICE)
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(errAnte).Should(Not(HaveOccurred()))
		Expect(errPost).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + 200_000))
		Expect(collectorBalanceAfter).To(Equal(uint64(200_000)))
	})

	It("Refund 10%", func() {
		// ARRANGE
		msg := bundlesTypes.MsgSubmitBundleProposal{Creator: i.ALICE}
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		gasLimit := uint64(200_000)
		txBuilder.SetGasLimit(gasLimit)
		fees := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1).MulRaw(int64(gasLimit))))
		txBuilder.SetFeeAmount(fees)
		_ = txBuilder.SetMsgs(&msg)
		tx := txBuilder.GetTx()

		// ACT
		_, errAnte := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		_, errPost := rfd.PostHandle(s.Ctx(), tx, false, true, PostNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.ALICE)
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(errAnte).Should(Not(HaveOccurred()))
		Expect(errPost).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + 180_000))
		Expect(collectorBalanceAfter).To(Equal(uint64(180_000)))
	})

	It("Refund 2/3 %", func() {
		// ARRANGE
		msg := stakersTypes.MsgCreateStaker{Creator: i.ALICE}
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		gasLimit := uint64(200_000)
		txBuilder.SetGasLimit(gasLimit)
		fees := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1).MulRaw(int64(gasLimit))))
		txBuilder.SetFeeAmount(fees)
		_ = txBuilder.SetMsgs(&msg)
		tx := txBuilder.GetTx()

		// ACT
		_, errAnte := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		_, errPost := rfd.PostHandle(s.Ctx(), tx, false, true, PostNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.ALICE)
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(errAnte).Should(Not(HaveOccurred()))
		Expect(errPost).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + (66_667)))
		Expect(collectorBalanceAfter).To(Equal(uint64(66_667)))
	})

	It("Refund 100%", func() {
		// ARRANGE
		msg := bundlesTypes.MsgVoteBundleProposal{Creator: i.ALICE}
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		gasLimit := uint64(200_000)
		txBuilder.SetGasLimit(gasLimit)
		fees := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1).MulRaw(int64(gasLimit))))
		txBuilder.SetFeeAmount(fees)
		_ = txBuilder.SetMsgs(&msg)
		tx := txBuilder.GetTx()

		// ACT
		_, errAnte := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		_, errPost := rfd.PostHandle(s.Ctx(), tx, false, true, PostNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.ALICE)
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(errAnte).Should(Not(HaveOccurred()))
		Expect(errPost).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter))
		Expect(collectorBalanceAfter).To(Equal(uint64(0)))
	})

	It("Don't refund multiple", func() {
		// ARRANGE
		msg1 := bundlesTypes.MsgVoteBundleProposal{Creator: i.ALICE}
		msg2 := bundlesTypes.MsgSkipUploaderRole{Creator: i.ALICE}
		msg3 := stakersTypes.MsgJoinPool{Creator: i.ALICE}
		txBuilder := encodingConfig.TxConfig.NewTxBuilder()
		gasLimit := uint64(200_000)
		txBuilder.SetGasLimit(gasLimit)
		fees := sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(1).MulRaw(int64(gasLimit))))
		txBuilder.SetFeeAmount(fees)
		_ = txBuilder.SetMsgs(&msg1, &msg2, &msg3)
		tx := txBuilder.GetTx()

		// ACT
		_, errAnte := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		_, errPost := rfd.PostHandle(s.Ctx(), tx, false, true, PostNextFn)

		// ASSERT
		accountBalanceAfter := s.GetBalanceFromAddress(i.ALICE)
		collectorBalanceAfter := s.GetBalanceFromModule(authTypes.FeeCollectorName)

		Expect(errAnte).Should(Not(HaveOccurred()))
		Expect(errPost).Should(Not(HaveOccurred()))
		Expect(accountBalanceBefore).To(Equal(accountBalanceAfter + 200_000))
		Expect(collectorBalanceAfter).To(Equal(uint64(200_000)))
	})
})
