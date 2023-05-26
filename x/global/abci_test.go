package global_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Global
	"github.com/KYVENetwork/chain/x/global"
	"github.com/KYVENetwork/chain/x/global/types"
)

/*

TEST CASES - DeductFeeDecorator

* BurnRatio = 0.0
* BurnRatio = 2/3 - test truncate
* BurnRatio = 0.5
* BurnRatio = 1.0

* TODO(@max): combine with refund

*/

var _ = Describe("AbciEndBlocker", Ordered, func() {
	s := i.NewCleanChain()
	encodingConfig := BuildEncodingConfig()
	dfd := global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, *s.App().GlobalKeeper, *s.App().StakingKeeper)

	accountBalanceBefore := s.GetBalanceFromAddress(i.DUMMY[0])
	totalSupplyBefore := s.App().BankKeeper.GetSupply(s.Ctx(), types.Denom).Amount.Uint64()

	BeforeEach(func() {
		s = i.NewCleanChain()

		mintParams := s.App().MintKeeper.GetParams(s.Ctx())
		mintParams.InflationMax = sdk.ZeroDec()
		mintParams.InflationMin = sdk.ZeroDec()
		_ = s.App().MintKeeper.SetParams(s.Ctx(), mintParams)

		accountBalanceBefore = s.GetBalanceFromAddress(i.DUMMY[0])
		totalSupplyBefore = s.App().BankKeeper.GetSupply(s.Ctx(), types.Denom).Amount.Uint64()
		dfd = global.NewDeductFeeDecorator(s.App().AccountKeeper, s.App().BankKeeper, s.App().FeeGrantKeeper, *s.App().GlobalKeeper, *s.App().StakingKeeper)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("BurnRatio = 0.0", func() {
		// ARRANGE
		// default burn ratio is zero
		denom := s.App().StakingKeeper.BondDenom(s.Ctx())
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(err).Should(Not(HaveOccurred()))

		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		accountBalanceDifference := accountBalanceBefore - accountBalanceAfter
		Expect(accountBalanceDifference).To(Equal(uint64(200_000)))

		totalSupplyAfter := s.App().BankKeeper.GetSupply(s.Ctx(), types.Denom).Amount.Uint64()
		totalSupplyDifference := totalSupplyBefore - totalSupplyAfter
		Expect(totalSupplyDifference).To(Equal(uint64(0)))
	})

	It("BurnRatio = 2/3 - test truncate", func() {
		// ARRANGE
		// set burn ratio to 0.3
		params := types.DefaultParams()
		params.BurnRatio = sdk.OneDec().MulInt64(2).QuoInt64(3)
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)

		// default burn ratio is zero
		denom := s.App().StakingKeeper.BondDenom(s.Ctx())
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(err).Should(Not(HaveOccurred()))

		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		accountBalanceDifference := accountBalanceBefore - accountBalanceAfter
		Expect(accountBalanceDifference).To(Equal(uint64(200_000)))

		totalSupplyAfter := s.App().BankKeeper.GetSupply(s.Ctx(), types.Denom).Amount.Uint64()
		totalSupplyDifference := totalSupplyBefore - totalSupplyAfter
		// Expect ..666 not ..667
		Expect(totalSupplyDifference).To(Equal(uint64(133_333)))
	})

	It("BurnRatio = 0.5", func() {
		// ARRANGE
		// set burn ratio to 0.5
		params := types.DefaultParams()
		params.BurnRatio = sdk.OneDec().QuoInt64(2)
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)

		denom := s.App().StakingKeeper.BondDenom(s.Ctx())
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(err).Should(Not(HaveOccurred()))

		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		accountBalanceDifference := accountBalanceBefore - accountBalanceAfter
		Expect(accountBalanceDifference).To(Equal(uint64(200_000)))

		totalSupplyAfter := s.App().BankKeeper.GetSupply(s.Ctx(), types.Denom).Amount.Uint64()
		totalSupplyDifference := totalSupplyBefore - totalSupplyAfter
		Expect(totalSupplyDifference).To(Equal(uint64(100_000)))
	})

	It("BurnRatio = 1.0", func() {
		// ARRANGE
		// set burn ratio to 0.5
		params := types.DefaultParams()
		params.BurnRatio = sdk.OneDec()
		s.App().GlobalKeeper.SetParams(s.Ctx(), params)

		denom := s.App().StakingKeeper.BondDenom(s.Ctx())
		tx := BuildTestTx(math.NewInt(1), denom, i.DUMMY[0], encodingConfig)

		// ACT
		_, err := dfd.AnteHandle(s.Ctx(), tx, false, AnteNextFn)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(err).Should(Not(HaveOccurred()))

		accountBalanceAfter := s.GetBalanceFromAddress(i.DUMMY[0])
		accountBalanceDifference := accountBalanceBefore - accountBalanceAfter
		Expect(accountBalanceDifference).To(Equal(uint64(200_000)))

		totalSupplyAfter := s.App().BankKeeper.GetSupply(s.Ctx(), types.Denom).Amount.Uint64()
		totalSupplyDifference := totalSupplyBefore - totalSupplyAfter
		Expect(totalSupplyDifference).To(Equal(uint64(200_000)))
	})
})
