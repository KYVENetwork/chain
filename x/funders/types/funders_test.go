package types_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/funders/types"
	globaltypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - funders.go

* Funding.ChargeOneBundle
* Funding.ChargeOneBundle - charge more than available
* Funding.ChargeOneBundle - charge with multiple coins
* Funding.ChargeOneBundle - charge with no coins
* FundintState.SetActive
* FundintState.SetActive - add same funder twice
* FundintState.SetInactive
* FundintState.SetInactive - with multiple funders

*/

var _ = Describe("logic_funders.go", Ordered, func() {
	s := i.NewCleanChain()

	funding := types.Funding{}
	fundingState := types.FundingState{}
	var whitelist []*types.WhitelistCoinEntry

	BeforeEach(func() {
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

		funding = types.Funding{
			FunderAddress:    i.ALICE,
			PoolId:           0,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(1 * i.T_KYVE),
			TotalFunded:      sdk.NewCoins(),
		}
		fundingState = types.FundingState{
			PoolId:                0,
			ActiveFunderAddresses: []string{i.ALICE, i.BOB},
		}
	})

	It("Funding.ChargeOneBundle", func() {
		// ACT
		payouts := funding.ChargeOneBundle()

		// ASSERT
		Expect(payouts.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
		Expect(funding.Amounts.String()).To(Equal(i.ACoins(99 * i.T_KYVE).String()))
		Expect(funding.TotalFunded.String()).To(Equal(i.ACoins(1 * i.T_KYVE).String()))
	})

	It("Funding.ChargeOneBundle - charge more than available", func() {
		// ARRANGE
		funding.Amounts = i.ACoins(1 * i.T_KYVE / 2)

		// ACT
		payouts := funding.ChargeOneBundle()

		// ASSERT
		Expect(payouts.String()).To(Equal(i.ACoins(1 * i.T_KYVE / 2).String()))
		Expect(funding.Amounts.IsZero()).To(BeTrue())
		Expect(funding.TotalFunded.String()).To(Equal(i.ACoins(1 * i.T_KYVE / 2).String()))
	})

	It("Funding.ChargeOneBundle - charge with multiple coins", func() {
		// ARRANGE
		funding.Amounts = sdk.NewCoins(i.ACoin(100*i.T_KYVE), i.BCoin(100*i.T_KYVE))
		funding.AmountsPerBundle = sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(2*i.T_KYVE))

		// ACT
		payouts := funding.ChargeOneBundle()

		// ASSERT
		Expect(payouts.String()).To(Equal(sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(2*i.T_KYVE)).String()))
		Expect(funding.Amounts.String()).To(Equal(sdk.NewCoins(i.ACoin(99*i.T_KYVE), i.BCoin(98*i.T_KYVE)).String()))
		Expect(funding.TotalFunded.String()).To(Equal(sdk.NewCoins(i.ACoin(1*i.T_KYVE), i.BCoin(2*i.T_KYVE)).String()))
	})

	It("Funding.ChargeOneBundle - charge with no coins", func() {
		// ARRANGE
		funding.Amounts = sdk.NewCoins()
		funding.AmountsPerBundle = sdk.NewCoins()

		// ACT
		payouts := funding.ChargeOneBundle()

		// ASSERT
		Expect(payouts.IsZero()).To(BeTrue())
		Expect(funding.Amounts.IsZero()).To(BeTrue())
		Expect(funding.TotalFunded.IsZero()).To(BeTrue())
	})

	It("FundintState.SetActive", func() {
		// ARRANGE
		fundingState.ActiveFunderAddresses = []string{}

		// ACT
		fundingState.SetActive(&funding)

		// ASSERT
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
	})

	It("FundintState.SetActive - add same funder twice", func() {
		// ACT
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))

		fundingState.SetActive(&funding)

		// ASSERT
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
	})

	It("FundintState.SetInactive", func() {
		// ACT
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
		fundingState.SetInactive(&funding)

		// ASSERT
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("FundintState.SetInactive - with multiple funders", func() {
		// ARRANGE
		fundingState.ActiveFunderAddresses = []string{i.ALICE, i.BOB, i.CHARLIE}

		// ACT
		fundingState.SetInactive(&funding)

		// ASSERT
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(2))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.CHARLIE))
		Expect(fundingState.ActiveFunderAddresses[1]).To(Equal(i.BOB))
	})
})
