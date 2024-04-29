package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_defund_pool.go

* Defund 50 coins from a funder who has previously funded 100 coins
* Defund more than actually funded
* Defund full funding amount from a funder who has previously funded 100 coins
* Defund as highest funder 75 coins in order to be the lowest funder afterward
* Try to defund nonexistent fundings
* Try to defund a funding twice
* Try to defund below minimum funding params (but not full defund)

*/

var _ = Describe("msg_server_defund_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalance := s.GetBalancesFromAddress(i.ALICE)
	whitelist := []*types.WhitelistCoinEntry{
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

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

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
		Expect(funding.Amounts.IsZero()).To(BeTrue())
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))
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
		Expect(funding.Amounts.IsZero()).To(BeTrue())
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))
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
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})

		// ASSERT
		s.RunTxFundersError(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})
	})

	It("Try to defund below minimum funding params (but not full defund)", func() {
		// ACT
		_, err := s.RunTx(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.ACoins(95 * i.T_KYVE),
		})

		// ASSERT
		Expect(err.Error()).To(Equal("minimum funding amount of 1000000000kyve not reached: invalid request"))
	})
})
