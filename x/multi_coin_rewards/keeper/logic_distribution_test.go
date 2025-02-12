package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	multicoinrewardstypes "github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_leave_pool.go

* Do not redistribute if no policy is set
* Redistribute single
* Redistribute multi
* Redistribute partial

*/

var _ = Describe("logic_distribution_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var gov string
	var validator1 i.TestValidatorAddress

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
		gov = s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

		// create pool
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		// Create second pool
		s.RunTxPoolSuccess(msg)

		// create staker
		validator1 = s.CreateNewValidator("MyValidator-1", 1000*i.KYVE)

		params := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())
		params.MultiCoinDistributionPolicyAdminAddress = validator1.Address
		s.App().MultiCoinRewardsKeeper.SetParams(s.Ctx(), params)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Do not redistribute if no policy is set", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(50))
		payoutRewards(s, validator1.Address, i.BCoins(70))

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		s.CommitAfterSeconds(s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx()).MultiCoinDistributionPendingTime)
		for i := 0; i < 100; i++ {
			s.CommitAfterSeconds(1)
		}

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.MultiCoinRewardsRedistributionAccountName).String()).To(Equal("50acoin,70bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("Redistribute single coin, same weight", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(50))
		payoutRewards(s, validator1.Address, i.BCoins(70))

		s.RunTxSuccess(&multicoinrewardstypes.MsgSetMultiCoinRewardsDistributionPolicy{
			Creator: validator1.Address,
			Policy: &multicoinrewardstypes.MultiCoinDistributionPolicy{
				Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{
					{Denom: "acoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
				},
			},
		})

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		s.CommitAfterSeconds(s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx()).MultiCoinDistributionPendingTime)
		for i := 0; i < 100; i++ {
			s.CommitAfterSeconds(1)
		}

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.MultiCoinRewardsRedistributionAccountName).String()).To(Equal("70bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())

		p0, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		p1, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)

		Expect(s.GetCoinsFromAddress(p0.GetPoolAccount().String()).String()).To(Equal("25acoin"))
		Expect(s.GetCoinsFromAddress(p1.GetPoolAccount().String()).String()).To(Equal("25acoin"))
	})

	It("Redistribute multiple coins, same weight", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(50))
		payoutRewards(s, validator1.Address, i.BCoins(70))
		payoutRewards(s, validator1.Address, i.CCoins(110))

		s.RunTxSuccess(&multicoinrewardstypes.MsgSetMultiCoinRewardsDistributionPolicy{
			Creator: validator1.Address,
			Policy: &multicoinrewardstypes.MultiCoinDistributionPolicy{
				Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{
					{Denom: "acoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
					{Denom: "bcoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
					{Denom: "ccoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
				},
			},
		})

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		s.CommitAfterSeconds(s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx()).MultiCoinDistributionPendingTime)
		for i := 0; i < 100; i++ {
			s.CommitAfterSeconds(1)
		}

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.MultiCoinRewardsRedistributionAccountName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())

		p0, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		p1, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)

		Expect(s.GetCoinsFromAddress(p0.GetPoolAccount().String()).String()).To(Equal("25acoin,35bcoin,55ccoin"))
		Expect(s.GetCoinsFromAddress(p1.GetPoolAccount().String()).String()).To(Equal("25acoin,35bcoin,55ccoin"))
	})

	It("Redistribute multiple coins, different weights", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(50))
		payoutRewards(s, validator1.Address, i.BCoins(70))
		payoutRewards(s, validator1.Address, i.CCoins(110))

		s.RunTxSuccess(&multicoinrewardstypes.MsgSetMultiCoinRewardsDistributionPolicy{
			Creator: validator1.Address,
			Policy: &multicoinrewardstypes.MultiCoinDistributionPolicy{
				Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{
					{Denom: "acoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"), // -> 16acoin
						},
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("2"), // -> 33acoin
						},
					}},
					{Denom: "bcoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"), // -> 70bcoin
						},
					}},
					{Denom: "ccoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"), // -> 110bcoin
						},
					}},
				},
			},
		})

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		s.CommitAfterSeconds(s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx()).MultiCoinDistributionPendingTime)
		for i := 0; i < 100; i++ {
			s.CommitAfterSeconds(1)
		}

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.MultiCoinRewardsRedistributionAccountName).String()).To(Equal("1acoin")) // left over acoin
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())

		p0, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		p1, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)

		Expect(s.GetCoinsFromAddress(p0.GetPoolAccount().String()).String()).To(Equal("16acoin,70bcoin"))
		Expect(s.GetCoinsFromAddress(p1.GetPoolAccount().String()).String()).To(Equal("33acoin,110ccoin"))
	})
})
