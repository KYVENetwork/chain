package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	"github.com/KYVENetwork/chain/x/team/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - abci.go

* total_supply
* team_balance
* community_pool
* distribution

*/

var _ = Describe("abci.go", Ordered, func() {
	s := i.NewCleanChainAtTime(int64(types.TGE))

	BeforeEach(func() {
		s = i.NewCleanChainAtTime(int64(types.TGE))
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("total_supply", func() {
		// ARRANGE
		b1, b2 := int64(0), int64(0)

		for t := 0; t < 100; t++ {
			// ACT
			b1 = s.App().BankKeeper.GetSupply(s.Ctx(), globalTypes.Denom).Amount.Int64()
			s.Commit()
			b2 = s.App().BankKeeper.GetSupply(s.Ctx(), globalTypes.Denom).Amount.Int64()

			// ASSERT
			Expect(b2).To(BeNumerically(">", b1))
		}
	})

	It("team_balance", func() {
		// ARRANGE
		b1, b2 := uint64(0), uint64(0)

		for t := 0; t < 100; t++ {
			// ACT
			b1 = s.App().TeamKeeper.GetTeamInfo(s.Ctx()).RequiredModuleBalance
			s.Commit()
			b2 = s.App().TeamKeeper.GetTeamInfo(s.Ctx()).RequiredModuleBalance

			// ASSERT
			Expect(b2).To(BeNumerically(">", b1))
		}
	})

	It("community_pool", func() {
		// ARRANGE
		b1, b2 := int64(0), int64(0)

		for t := 0; t < 100; t++ {
			// ACT
			feePool, _ := s.App().DistributionKeeper.FeePool.Get(s.Ctx())
			b1 = feePool.CommunityPool.AmountOf(globalTypes.Denom).TruncateInt64()
			s.Commit()
			feePool, _ = s.App().DistributionKeeper.FeePool.Get(s.Ctx())
			b2 = feePool.CommunityPool.AmountOf(globalTypes.Denom).TruncateInt64()

			// ASSERT
			Expect(b2).To(BeNumerically(">", b1))
		}
	})

	It("distribution", func() {
		for t := 0; t < 100; t++ {
			// ARRANGE

			// get the team balance and total supply at current block which will
			// be used to calculate distribution in BeginBlock of next block
			teamBalance := math.LegacyNewDec(int64(s.GetBalanceFromModule(types.ModuleName)))
			totalSupply := math.LegacyNewDec(s.App().BankKeeper.GetSupply(s.Ctx(), globalTypes.Denom).Amount.Int64())

			// get current team and validators reward for this block
			r1 := s.App().TeamKeeper.GetTeamInfo(s.Ctx()).TotalAuthorityRewards
			feePool, _ := s.App().DistributionKeeper.FeePool.Get(s.Ctx())
			c1 := uint64(feePool.CommunityPool.AmountOf(globalTypes.Denom).TruncateInt64())

			// ACT

			// inflation is minted and distributed here
			s.Commit()

			// calculate delta for team and community rewards in order to verify distribution
			r2 := s.App().TeamKeeper.GetTeamInfo(s.Ctx()).TotalAuthorityRewards
			feePool, _ = s.App().DistributionKeeper.FeePool.Get(s.Ctx())
			c2 := uint64(feePool.CommunityPool.AmountOf(globalTypes.Denom).TruncateInt64())

			teamReward := r2 - r1
			communityReward := c2 - c1

			// get block reward for this block
			minter, _ := s.App().MintKeeper.Minter.Get(s.Ctx())
			params, _ := s.App().MintKeeper.Params.Get(s.Ctx())
			blockProvision := minter.BlockProvision(params)

			// ASSERT

			// calculate if team and community distribution add up to total inflation reward
			Expect(teamReward + communityReward).To(Equal(blockProvision.Amount.Uint64()))

			// calculate if distribution share matches with team balance and total supply
			Expect(teamBalance.Mul(math.LegacyNewDec(blockProvision.Amount.Int64())).Quo(totalSupply).TruncateInt64()).To(Equal(int64(teamReward)))
		}
	})
})
