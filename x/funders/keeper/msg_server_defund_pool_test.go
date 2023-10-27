package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_defund_pool.go

* Defund 50 KYVE from a funder who has previously funded 100 KYVE
* Defund more than actually funded
* Defund full funding amount from a funder who has previously funded 100 KYVE
* Defund as highest funder 75 KYVE in order to be the lowest funder afterward
* Try to defund nonexistent fundings
* Try to defund a funding twice
* Try to defund below minimum funding params (but not full defund)

*/

var _ = Describe("msg_server_defund_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalance := s.GetBalanceFromAddress(i.ALICE)

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

		// create funder
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "moniker",
		})

		// fund pool
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Defund 50 KYVE from a funder who has previously funded 100 KYVE", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		Expect(initialBalance - balanceAfter).To(Equal(50 * i.KYVE))

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.Amount).To(Equal(50 * i.KYVE))
		Expect(funding.AmountPerBundle).To(Equal(1 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(0 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(1))
		Expect(fundingState.ActiveFunderAddresses[0]).To(Equal(i.ALICE))
	})

	It("Defund more than actually funded", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  101 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceAfter).To(BeZero())

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.Amount).To(Equal(0 * i.KYVE))
		Expect(funding.AmountPerBundle).To(Equal(1 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(0 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))
	})

	It("Defund full funding amount from a funder who has previously funded 100 KYVE", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  100 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)
		Expect(initialBalance - balanceAfter).To(BeZero())

		funding, _ := s.App().FundersKeeper.GetFunding(s.Ctx(), i.ALICE, 0)
		Expect(funding.Amount).To(Equal(0 * i.KYVE))
		Expect(funding.AmountPerBundle).To(Equal(1 * i.KYVE))
		Expect(funding.TotalFunded).To(Equal(0 * i.KYVE))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(len(fundingState.ActiveFunderAddresses)).To(Equal(0))
	})

	It("Defund as highest funder 75 KYVE in order to be the lowest funder afterwards", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.BOB,
			Moniker: "moniker",
		})
		s.RunTxFundersSuccess(&types.MsgFundPool{
			Creator:         i.BOB,
			PoolId:          0,
			Amount:          50 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err := s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.BOB))

		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  75 * i.KYVE,
		})

		// ASSERT
		fundingState, _ = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		activeFundings = s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
		lowestFunding, err = s.App().FundersKeeper.GetLowestFunding(activeFundings)
		Expect(err).To(BeNil())
		Expect(lowestFunding.FunderAddress).To(Equal(i.ALICE))
	})

	It("Try to defund nonexistent fundings", func() {
		// ASSERT
		s.RunTxFundersError(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  1,
			Amount:  1 * i.KYVE,
		})

		s.RunTxFundersError(&types.MsgDefundPool{
			Creator: i.BOB,
			PoolId:  0,
			Amount:  1 * i.KYVE,
		})
	})

	It("Try to defund a funding twice", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  100 * i.KYVE,
		})

		// ASSERT
		s.RunTxFundersError(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  100 * i.KYVE,
		})
	})

	It("Try to defund below minimum funding params (but not full defund)", func() {
		// ACT
		_, err := s.RunTx(&types.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  100*i.KYVE - types.DefaultMinFundingAmount/2,
		})

		// ASSERT
		Expect(err.Error()).To(Equal("minimum funding amount of 1000000000kyve not reached: invalid request"))
	})
})
