package keeper_test

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/funders/types"
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_defund_pool.go

* Defund 50 coins from a funder who has previously funded 100 coins
* Defund more than actually funded
* Defund full funding amount from a funder who has previously funded 100 coins
* Defund as highest funder 75 coins in order to be the lowest funder afterward
* Try to defund zero amounts
* Try to defund nonexistent fundings
* Try to defund a funding twice
* Try to defund below minimum funding params (but not full defund)
* Try to partially defund after a coin has been removed from the whitelist
* Try to fully defund after a coin has been removed from the whitelist

*/

var _ = Describe("msg_server_defund_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	var initialBalance sdk.Coins
	var whitelist []*types.WhitelistCoinEntry

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		initialBalance = s.GetBalancesFromAddress(i.ALICE)

		// create clean pool for every test case
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10_000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		// set whitelist
		whitelist = []*types.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.A_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(3),
			},
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), types.NewParams(whitelist, 20))

		// create funder
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "moniker",
		})

		// fund pool
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Defund 50 coins from a funder who has previously funded 100 coins", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(50 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)

		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		// assert if the funder can still fund afterwards
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
	})

	It("Defund more than actually funded", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(101 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceAfter...).IsZero()).To(BeTrue())

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.Amounts).To(BeEmpty())
		Expect(funding.AmountsPerBundle).To(BeEmpty())
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))

		// assert if the funder can still fund afterwards
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
	})

	It("Defund full funding amount from a funder who has previously funded 100 coins", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceAfter...).IsZero()).To(BeTrue())

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.Amounts).To(BeEmpty())
		Expect(funding.AmountsPerBundle).To(BeEmpty())
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))

		// assert if the funder can still fund afterwards
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
	})

	It("Defund as highest funder 75 coins in order to be the lowest funder afterwards", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.BOB,
			Moniker: "moniker",
		})
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(50 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.BOB))

		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(75 * i.T_KYVE),
		})

		// ASSERT
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		activeFundings = s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err = s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))

		// assert if the funder can still fund afterwards
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
	})

	It("Try to defund zero amounts", func() {
		// ACT
		_, err := s.RunTx(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: sdk.NewCoins(),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(errors.Wrapf(errorsTypes.ErrInvalidRequest, "empty amount").Error()))
	})

	It("Try to defund nonexistent fundings", func() {
		// ASSERT
		s.RunTxFundersError(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  1,
			Amounts: i.ACoins(1 * i.T_KYVE),
		})

		s.RunTxFundersError(&types.MsgDefundPool{
			Creator: i.BOB,
			PoolId:  0,
			Amounts: i.ACoins(1 * i.T_KYVE),
		})
	})

	It("Try to defund a funding twice", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})

		// ACT
		_, err := s.RunTx(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(errors.Wrapf(errorsTypes.ErrInvalidRequest, types.ErrFundsTooLow.Error()).Error()))
	})

	It("Try to defund below minimum funding params (but not full defund)", func() {
		// ACT
		_, err := s.RunTx(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(95 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(types.ErrMinFundingAmount.Error()))
	})

	It("Try to partially defund after a coin has been removed from the whitelist", func() {
		// ARRANGE
		whitelist = []*types.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(3),
			},
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), types.NewParams(whitelist, 20))

		// ACT
		_, err := s.RunTx(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(50 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(types.ErrCoinNotWhitelisted.Error()))
	})

	It("Try to fully defund after a coin has been removed from the whitelist", func() {
		// ARRANGE
		whitelist = []*types.WhitelistCoinEntry{
			{
				CoinDenom:                 globaltypes.Denom,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          10 * i.KYVE,
				MinFundingAmountPerBundle: 1 * i.KYVE,
				CoinWeight:                math.LegacyNewDec(3),
			},
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), types.NewParams(whitelist, 20))

		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceAfter...).IsZero()).To(BeTrue())

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.Amounts).To(BeEmpty())
		Expect(funding.AmountsPerBundle).To(BeEmpty())
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))

		// assert if the funder can still fund afterwards
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.BCoins(100 * i.T_KYVE),
			AmountsPerBundle: i.BCoins(1 * i.T_KYVE),
		})
	})
})
