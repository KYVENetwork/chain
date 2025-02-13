package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	multicoinrewardstypes "github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_leave_pool.go

* Valid Policy: Update simple valid policy
* Invalid Policy: Duplicate denom entries
* Invalid Policy: Duplicate pool entries for same denom
* Invalid Policy: Negative weights
* Invalid Policy: Zero weights

*/

var _ = Describe("msg_server_distribution_policy_test.go", Ordered, func() {
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

		// create staker
		validator1 = s.CreateNewValidator("MyValidator-1", 1000*i.KYVE)

		params := s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx())
		params.MultiCoinDistributionPolicyAdminAddress = validator1.Address
		s.App().MultiCoinRewardsKeeper.SetParams(s.Ctx(), params)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Valid Policy: Update simple valid policy", func() {
		// Arrange
		// ACT
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
							Weight: math.LegacyMustNewDecFromStr("2"),
						},
					}},
					{Denom: "bcoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
					{Denom: "ccoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
				},
			},
		})

		// ASSERT
		policy, _ := s.App().MultiCoinRewardsKeeper.MultiCoinDistributionPolicy.Get(s.Ctx())
		distributionMap, err := multicoinrewardstypes.ParseAndNormalizeMultiCoinDistributionMap(policy)
		Expect(err).To(BeNil())
		Expect(distributionMap).To(HaveLen(3))
	})

	It("Invalid Policy: Duplicate denom entries", func() {
		// Arrange
		// ACT
		_, err := s.RunTx(&multicoinrewardstypes.MsgSetMultiCoinRewardsDistributionPolicy{
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
							Weight: math.LegacyMustNewDecFromStr("2"),
						},
					}},
					{Denom: "acoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
					{Denom: "ccoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
				},
			},
		})

		// ASSERT
		Expect(err.Error()).To(Equal("duplicate entry for denom acoin"))
		_, err = s.App().MultiCoinRewardsKeeper.MultiCoinDistributionPolicy.Get(s.Ctx())
		Expect(err).NotTo(BeNil())
	})

	It("Invalid Policy: Duplicate pool entries for same denom", func() {
		// Arrange
		// ACT
		_, err := s.RunTx(&multicoinrewardstypes.MsgSetMultiCoinRewardsDistributionPolicy{
			Creator: validator1.Address,
			Policy: &multicoinrewardstypes.MultiCoinDistributionPolicy{
				Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{
					{Denom: "acoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("2"),
						},
					}},
					{Denom: "bcoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
					{Denom: "ccoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("1"),
						},
					}},
				},
			},
		})

		// ASSERT
		Expect(err.Error()).To(Equal("duplicate distribution weight for pool id 0"))
		_, err = s.App().MultiCoinRewardsKeeper.MultiCoinDistributionPolicy.Get(s.Ctx())
		Expect(err).NotTo(BeNil())
	})

	It("Invalid Policy: Negative weights", func() {
		// Arrange
		// ACT
		_, err := s.RunTx(&multicoinrewardstypes.MsgSetMultiCoinRewardsDistributionPolicy{
			Creator: validator1.Address,
			Policy: &multicoinrewardstypes.MultiCoinDistributionPolicy{
				Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{
					{Denom: "acoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("-1"),
						},
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("2"),
						},
					}},
				},
			},
		})

		// ASSERT
		Expect(err.Error()).To(Equal("invalid weight for pool id 0"))
		_, err = s.App().MultiCoinRewardsKeeper.MultiCoinDistributionPolicy.Get(s.Ctx())
		Expect(err).NotTo(BeNil())
	})

	It("Invalid Policy: Zero weights", func() {
		// Arrange
		// ACT
		_, err := s.RunTx(&multicoinrewardstypes.MsgSetMultiCoinRewardsDistributionPolicy{
			Creator: validator1.Address,
			Policy: &multicoinrewardstypes.MultiCoinDistributionPolicy{
				Entries: []*multicoinrewardstypes.MultiCoinDistributionDenomEntry{
					{Denom: "acoin", PoolWeights: []*multicoinrewardstypes.MultiCoinDistributionPoolWeightEntry{
						{
							PoolId: 0,
							Weight: math.LegacyMustNewDecFromStr("-1"),
						},
						{
							PoolId: 1,
							Weight: math.LegacyMustNewDecFromStr("2"),
						},
					}},
				},
			},
		})

		// ASSERT
		Expect(err.Error()).To(Equal("invalid weight for pool id 0"))
		_, err = s.App().MultiCoinRewardsKeeper.MultiCoinDistributionPolicy.Get(s.Ctx())
		Expect(err).NotTo(BeNil())
	})
})
