package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - abci.go

* pool should always receive inflation funds
* TODO: pool should split inflation funds depending on operating cost
* TODO: pools with zero operating cost should receive nothing
* TODO: every pool has zero operating cost

*/

var _ = Describe("abci.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		s = i.NewCleanChain()

		s.App().PoolKeeper.SetParams(s.Ctx(), poolTypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0.1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		s.App().PoolKeeper.AppendPool(s.Ctx(), poolTypes.Pool{
			Name:           "PoolTest",
			MaxBundleSize:  100,
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  1_000_000,
			Protocol: &poolTypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &poolTypes.UpgradePlan{},
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("pool should always receive inflation funds", func() {
		// ARRANGE
		b1, b2 := uint64(0), uint64(0)

		for t := 0; t < 100; t++ {
			// ACT
			pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
			b1 = uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool.GetPoolAccount(), globalTypes.Denom).Amount.Int64())
			s.Commit()
			b2 = uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

			// ASSERT
			Expect(b2).To(BeNumerically(">", b1))
		}
	})
})
