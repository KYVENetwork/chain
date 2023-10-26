package types_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/funders/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - funders.go

* Funding.AddAmount
* Funding.SubtractAmount
* Funding.SubtractAmount - subtract more than available
* Funding.ChargeOneBundle
* Funding.ChargeOneBundle - charge more than available
* FundintState.SetActive
* FundintState.SetActive - add same funder twice
* FundintState.SetInactive
* FundintState.SetInactive - with multiple funders

*/

var _ = Describe("logic_funders.go", Ordered, func() {
	funding := types.Funding{}
	fundingState := types.FundingState{}

	BeforeEach(func() {
		funding = types.Funding{
			FunderAddress:   i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
			TotalFunded:     0,
		}
		fundingState = types.FundingState{
			PoolId:                0,
			ActiveFunderAddresses: []string{i.ALICE, i.BOB},
		}
	})

	It("Funding.AddAmount", func() {
		// ACT
		funding.AddAmount(100 * i.KYVE)

		// ASSERT
		Expect(funding.Amount).To(Equal(200 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(uint64(0)))
	})

	It("Funding.SubtractAmount", func() {
		// ACT
		funding.SubtractAmount(50 * i.KYVE)

		// ASSERT
		Expect(funding.Amount).To(Equal(50 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(uint64(0)))
	})

	It("Funding.SubtractAmount - subtract more than available", func() {
		// ACT
		subtracted := funding.SubtractAmount(200 * i.KYVE)

		// ASSERT
		Expect(subtracted).To(Equal(100 * i.KYVE))
		Expect(funding.Amount).To(Equal(uint64(0)))
	})

	It("Funding.ChargeOneBundle", func() {
		// ACT
		amount := funding.ChargeOneBundle()

		// ASSERT
		Expect(amount).To(Equal(1 * i.KYVE))
		Expect(funding.Amount).To(Equal(99 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(1 * i.KYVE))
	})

	It("Funding.ChargeOneBundle - charge more than available", func() {
		// ARRANGE
		funding.Amount = 1 * i.KYVE / 2

		// ACT
		amount := funding.ChargeOneBundle()

		// ASSERT
		Expect(amount).To(Equal(1 * i.KYVE / 2))
		Expect(funding.Amount).To(Equal(uint64(0)))
		Expect(funding.TotalFunded).To(Equal(1 * i.KYVE / 2))
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
