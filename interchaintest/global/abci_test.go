package global_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
)

/*

TEST CASES - x/global/abci.go

* BurnRatio = 0.0
* BurnRatio = 2/3 - test truncate
* BurnRatio = 0.5
* BurnRatio = 1.0

*/

func TestProposalHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "interchaintest/x/global Test Suite")
}

var _ = Describe("x/global/abci.go - Endblocker", func() {
	It("BurnRatio = 0.0", func() {
		// ARRANGE
		ctx := context.Background()
		burnRatio := math.LegacyZeroDec()

		chain, interchain, broadcaster, wallet := startNewChainWithCustomBurnRatio(ctx, burnRatio)
		DeferCleanup(func() {
			_ = chain.StopAllNodes(context.Background())
			_ = interchain.Close()
		})

		balanceBefore, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())
		result, err := chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupplyBefore := result.AmountOf(chain.Config().Denom)

		// ACT
		msgSend := &banktypes.MsgSend{
			FromAddress: wallet.FormattedAddress(),
			ToAddress:   wallet.FormattedAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(chain.Config().Denom, math.NewInt(i.T_KYVE))),
		}
		tx, err := cosmos.BroadcastTx(ctx, broadcaster, wallet, msgSend)

		// ASSERT
		Expect(err).To(BeNil())
		Expect(tx.Code).To(Equal(uint32(0)))

		balance, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())

		accountBalanceDifference := balanceBefore.Sub(balance)
		Expect(accountBalanceDifference.Int64()).To(Equal(int64(200_000)))

		result, err = chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupply := result.AmountOf(chain.Config().Denom)
		Expect(totalSupply).To(Equal(totalSupplyBefore))
	})

	It("BurnRatio = 2/3 - test truncate", func() {
		// ARRANGE
		ctx := context.Background()
		burnRatio := math.LegacyOneDec().MulInt64(2).QuoInt64(3)

		chain, interchain, broadcaster, wallet := startNewChainWithCustomBurnRatio(ctx, burnRatio)
		DeferCleanup(func() {
			_ = chain.StopAllNodes(context.Background())
			_ = interchain.Close()
		})

		balanceBefore, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())
		result, err := chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupplyBefore := result.AmountOf(chain.Config().Denom)

		// ACT
		msgSend := &banktypes.MsgSend{
			FromAddress: wallet.FormattedAddress(),
			ToAddress:   wallet.FormattedAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(chain.Config().Denom, math.NewInt(i.T_KYVE))),
		}
		tx, err := cosmos.BroadcastTx(ctx, broadcaster, wallet, msgSend)

		// ASSERT
		Expect(err).To(BeNil())
		Expect(tx.Code).To(Equal(uint32(0)))

		balance, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())

		accountBalanceDifference := balanceBefore.Sub(balance)
		Expect(accountBalanceDifference.Int64()).To(Equal(int64(200_000)))

		result, err = chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupply := result.AmountOf(chain.Config().Denom)
		totalSupplyDifference := totalSupplyBefore.Sub(totalSupply)

		// Expect ..666 not ..667
		Expect(totalSupplyDifference).To(Equal(math.NewInt(133_333)))
	})

	It("BurnRatio = 0.5", func() {
		// ARRANGE
		ctx := context.Background()
		burnRatio := math.LegacyMustNewDecFromStr("0.5")

		chain, interchain, broadcaster, wallet := startNewChainWithCustomBurnRatio(ctx, burnRatio)
		DeferCleanup(func() {
			_ = chain.StopAllNodes(context.Background())
			_ = interchain.Close()
		})

		balanceBefore, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())
		result, err := chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupplyBefore := result.AmountOf(chain.Config().Denom)

		// ACT
		msgSend := &banktypes.MsgSend{
			FromAddress: wallet.FormattedAddress(),
			ToAddress:   wallet.FormattedAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(chain.Config().Denom, math.NewInt(i.T_KYVE))),
		}
		tx, err := cosmos.BroadcastTx(ctx, broadcaster, wallet, msgSend)

		// ASSERT
		Expect(err).To(BeNil())
		Expect(tx.Code).To(Equal(uint32(0)))

		balance, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())

		accountBalanceDifference := balanceBefore.Sub(balance)
		Expect(accountBalanceDifference.Int64()).To(Equal(int64(200_000)))

		result, err = chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupply := result.AmountOf(chain.Config().Denom)
		totalSupplyDifference := totalSupplyBefore.Sub(totalSupply)

		Expect(totalSupplyDifference).To(Equal(math.NewInt(100_000)))
	})

	It("BurnRatio = 1.0", func() {
		// ARRANGE
		ctx := context.Background()
		burnRatio := math.LegacyMustNewDecFromStr("1.0")

		chain, interchain, broadcaster, wallet := startNewChainWithCustomBurnRatio(ctx, burnRatio)
		DeferCleanup(func() {
			_ = chain.StopAllNodes(context.Background())
			_ = interchain.Close()
		})

		balanceBefore, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())
		result, err := chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupplyBefore := result.AmountOf(chain.Config().Denom)

		// ACT
		msgSend := &banktypes.MsgSend{
			FromAddress: wallet.FormattedAddress(),
			ToAddress:   wallet.FormattedAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(chain.Config().Denom, math.NewInt(i.T_KYVE))),
		}
		tx, err := cosmos.BroadcastTx(ctx, broadcaster, wallet, msgSend)

		// ASSERT
		Expect(err).To(BeNil())
		Expect(tx.Code).To(Equal(uint32(0)))

		balance, err := chain.GetBalance(ctx, wallet.FormattedAddress(), chain.Config().Denom)
		Expect(err).To(BeNil())

		accountBalanceDifference := balanceBefore.Sub(balance)
		Expect(accountBalanceDifference.Int64()).To(Equal(int64(200_000)))

		result, err = chain.BankQueryTotalSupply(ctx)
		Expect(err).To(BeNil())
		totalSupply := result.AmountOf(chain.Config().Denom)
		totalSupplyDifference := totalSupplyBefore.Sub(totalSupply)

		Expect(totalSupplyDifference).To(Equal(math.NewInt(200_000)))
	})
})
