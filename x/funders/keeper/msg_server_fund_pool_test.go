package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_fund_pool.go

* Fund a pool with 100 coins
* Fund additional 50 coins to an existing funding with 100 coins
* Fund additional 50 different coins to an existing funding with 100 coins
* Change amount per bundle after funding 100 coins
* Change multiple amounts per bundle and also fund 50 additional coins after funding 100 coins
* Try to fund more coins than available in balance
* Try to fund multiple coins where one coin exceeds balance
* Fund with a new funder less coins than the existing one
* Fund with a new funder more coins than the existing one
* Try funding with no amounts and no amounts per bundle
* Try funding coin which is not in the whitelist
* Try funding multiple coins where one coin is not in the whitelist
* Try funding 100 coins but amount per bundle is not set yet and empty
* Try funding 100 coins but amount per bundle is a different coin
* Try changing the amount per bundle of a coin which is not funded
* Try changing the amount per bundle of a coin which is not funded and whitelisted
* Try to fund with a non-existent funder
* Try to fund less coins than the lowest funder with full funding slots
* Fund more coins than the lowest funder with full funding slots
* Refund a funding as the lowest funder
* Try to fund a non-existent pool
* Try to fund below the minimum amount
* Try to fund multiple coins where one is below the minimum amount
* Try to fund below the minimum amount per bundle
* Try to fund multiple coins where one is below the minimum amount per bundle
* Try to fund without fulfilling min funding multiple
* Try to fund multiple coins where one is not fulfilling min funding multiple

*/

var _ = Describe("msg_server_fund_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalance := s.GetBalancesFromAddress(i.ALICE)
	whitelist := []*funderstypes.WhitelistCoinEntry{
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
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams(whitelist, 20))

		// create funder
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.BOB,
			Moniker: "Bob",
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Fund a pool with 100 coins", func() {
		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)

		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(i.ACoins(100 * i.T_KYVE).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.FunderAddress).To(Equal(i.ALICE))
		Expect(funding.PoolId).To(Equal(uint64(0)))
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(100 * i.T_KYVE).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
	})

	It("Fund additional 50 coins to an existing funding with 100 coins", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(50 * i.T_KYVE),
			AmountsPerBundle: sdk.NewCoins(),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)

		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(i.ACoins(150 * i.T_KYVE).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.FunderAddress).To(Equal(i.ALICE))
		Expect(funding.PoolId).To(Equal(uint64(0)))
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(150 * i.T_KYVE).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Fund additional 50 different coins to an existing funding with 100 coins", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.BCoins(50 * i.T_KYVE),
			AmountsPerBundle: i.BCoins(1 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)

		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.FunderAddress).To(Equal(i.ALICE))
		Expect(funding.PoolId).To(Equal(uint64(0)))
		Expect(funding.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(1*i.T_KYVE)).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Change amount per bundle after funding 100 coins", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(),
			AmountsPerBundle: i.ACoins(2 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)

		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(i.ACoins(100 * i.T_KYVE).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.FunderAddress).To(Equal(i.ALICE))
		Expect(funding.PoolId).To(Equal(uint64(0)))
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(100 * i.T_KYVE).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(2 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Change multiple amounts per bundle and also fund 50 additional coins after funding 100 coins", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.BCoins(50 * i.T_KYVE),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(2*i.T_KYVE), i.BCoin(1*i.T_KYVE)),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)

		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.FunderAddress).To(Equal(i.ALICE))
		Expect(funding.PoolId).To(Equal(uint64(0)))
		Expect(funding.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(50*i.T_KYVE)).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(sdk.NewCoins(i.ACoin(2*i.T_KYVE), i.BCoin(1*i.T_KYVE)).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Try to fund more coins than available in balance", func() {
		// ACT
		_, currentBalance := s.GetBalancesFromAddress(i.ALICE).Find(i.A_DENOM)

		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(currentBalance.Add(i.ACoin(1))),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceAfter...).IsZero()).To(BeTrue())

		_, found := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(found).To(BeFalse())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))
	})

	It("Try to fund multiple coins where one coin exceeds balance", func() {
		// ACT
		currentBalance := s.GetBalancesFromAddress(i.ALICE)

		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          currentBalance.Add(i.ACoin(1)),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceAfter...).IsZero()).To(BeTrue())

		_, found := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(found).To(BeFalse())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))
	})

	It("Fund with a new funder less coins than the existing one", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(50 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.BOB)
		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(50 * i.T_KYVE).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(2))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.BOB))
	})

	It("Fund with a new funder more coins than the existing one", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(200 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		balanceAfter := s.GetBalancesFromAddress(i.BOB)
		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(i.ACoins(200 * i.T_KYVE).String()))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.BOB, 0)
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(200 * i.T_KYVE).String()))
		Expect(funding.AmountsPerBundle.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState.PoolId).To(Equal(uint64(0)))
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(2))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Try funding with no amounts and no amounts per bundle", func() {
		// ASSERT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(),
			AmountsPerBundle: sdk.NewCoins(),
		})
	})

	It("Try funding coin which is not in the whitelist", func() {
		// ARRANGE
		whitelist = []*funderstypes.WhitelistCoinEntry{
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
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams(whitelist, 20))

		// ASSERT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.CCoins(100 * i.T_KYVE),
			AmountsPerBundle: i.CCoins(1 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrCoinNotWhitelisted.Error()))
	})

	It("Try funding multiple coins where one coin is not in the whitelist", func() {
		// ARRANGE
		whitelist = []*funderstypes.WhitelistCoinEntry{
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
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams(whitelist, 20))

		// ACT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.CCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.CCoin(1*i.T_KYVE)),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrCoinNotWhitelisted.Error()))
	})

	It("Try funding 100 coins but amount per bundle is not set yet and empty", func() {
		// ACT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: sdk.NewCoins(),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrInvalidAmountPerBundleCoin.Error()))
	})

	It("Try funding 100 coins but amount per bundle is a different coin", func() {
		// ACT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.BCoins(1 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrInvalidAmountPerBundleCoin.Error()))
	})

	It("Try changing the amount per bundle of a coin which is not funded", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ACT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(),
			AmountsPerBundle: i.BCoins(1 * i.T_KYVE),
		})

		// ASSERT
		// we actually allow this since the amount per bundle will not be automatically removed
		// if the coin in amounts is empty. The funder can only specify amount per bundles which
		// are whitelisted
		Expect(err).ToNot(HaveOccurred())
	})

	It("Try changing the amount per bundle of a coin which is not funded and whitelisted", func() {
		// ARRANGE
		whitelist = []*funderstypes.WhitelistCoinEntry{
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
		}
		s.App().FundersKeeper.SetParams(s.Ctx(), funderstypes.NewParams(whitelist, 20))

		// ASSERT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(),
			AmountsPerBundle: i.CCoins(1 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrCoinNotWhitelisted.Error()))
	})

	It("Try to fund with a non-existent funder", func() {
		// ASSERT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:          i.CHARLIE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
	})

	It("Try to fund less coins than the lowest funder with full funding slots", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		for a := 0; a < funderstypes.MaxFunders-1; a++ {
			s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
				Creator: i.DUMMY[a],
				Moniker: i.DUMMY[a],
			})
			// fill remaining funding slots
			s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
				Creator:          i.DUMMY[a],
				PoolId:           0,
				Amounts:          i.ACoins(1000 * i.T_KYVE),
				AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
			})
		}

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		balanceAfter := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceAfter...)).To(Equal(i.ACoins(100 * i.T_KYVE)))

		// ACT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(50 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		Expect(activeFundings).To(HaveLen(50))
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Fund more coins than the lowest funder with full funding slots", func() {
		// ARRANGE
		initialBalanceBob := s.GetBalancesFromAddress(i.BOB)
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		for a := 0; a < funderstypes.MaxFunders-1; a++ {
			s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
				Creator: i.DUMMY[a],
				Moniker: i.DUMMY[a],
			})
			// fill remaining funding slots
			s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
				Creator:          i.DUMMY[a],
				PoolId:           0,
				Amounts:          i.ACoins(1000 * i.T_KYVE),
				AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
			})
		}
		balanceAfter := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceAfter...).String()).To(Equal(i.ACoins(100 * i.T_KYVE).String()))

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.BOB,
			PoolId:           0,
			Amounts:          i.ACoins(200 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
		x := s.GetBalancesFromAddress(i.BOB)
		Expect(initialBalanceBob.Sub(x...).String()).To(Equal(i.ACoins(200 * i.T_KYVE).String()))
		// ASSERT
		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		Expect(activeFundings).To(HaveLen(50))
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.BOB))

		balanceEnd := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceEnd...).IsZero()).To(BeTrue())

		balanceAfterBob := s.GetBalancesFromAddress(i.BOB)
		Expect(initialBalanceBob.Sub(balanceAfterBob...).String()).To(Equal(i.ACoins(200 * i.T_KYVE).String()))
	})

	It("Refund a funding as the lowest funder", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		for a := 0; a < funderstypes.MaxFunders-1; a++ {
			s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
				Creator: i.DUMMY[a],
				Moniker: i.DUMMY[a],
			})
			// fill remaining funding slots
			s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
				Creator:          i.DUMMY[a],
				PoolId:           0,
				Amounts:          i.ACoins(1000 * i.T_KYVE),
				AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
			})
		}

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		Expect(activeFundings).To(HaveLen(50))
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings, whitelist)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))

		// ACT
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(50 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(funderstypes.MaxFunders))

		balanceEnd := s.GetBalancesFromAddress(i.ALICE)
		Expect(initialBalance.Sub(balanceEnd...).String()).To(Equal(i.ACoins(150 * i.T_KYVE).String()))
	})

	It("Try to fund a non-existent pool", func() {
		// ASSERT
		s.RunTxFundersError(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           1,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})
	})

	It("Try to fund below the minimum amount", func() {
		// ACT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(1 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrMinFundingAmount.Error()))
	})

	It("Try to fund multiple coins where one is below the minimum amount", func() {
		// ACT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(1*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(1*i.T_KYVE)),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrMinFundingAmount.Error()))
	})

	It("Try to fund below the minimum amount per bundle", func() {
		// ASSERT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrMinAmountPerBundle.Error()))
	})

	It("Try to fund multiple coins where one is below the minimum amount per bundle", func() {
		// ASSERT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(1)),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrMinAmountPerBundle.Error()))
	})

	It("Try to fund without fulfilling min funding multiple", func() {
		// ASSERT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(50 * i.T_KYVE),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrMinFundingMultiple.Error()))
	})

	It("Try to fund multiple coins where one is not fulfilling min funding multiple", func() {
		// ASSERT
		_, err := s.RunTx(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(50*i.T_KYVE)),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(funderstypes.ErrMinFundingMultiple.Error()))
	})
})
