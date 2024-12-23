package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - stakers leave

* Staker leaves, although he is the next uploader and runs into the upload timeout
* Staker leaves, although he was the uploader of the previous round and should receive the uploader reward
* Staker leaves, although he was the uploader of the previous round and should get slashed
* Staker leaves, although he was a voter in the previous round and should get slashed
* Staker leaves, although he was a voter in the previous round and should get a point
* Staker leaves, although he was a voter who did not vote max points in a row should not get slashed

*/

var _ = Describe("stakers leave", Ordered, func() {
	var s *i.KeeperTestSuite
	var initialBalanceStaker0, initialBalanceStaker1 uint64

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)

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
			InflationShareWeight: math.LegacyNewDec(10_000),
			MinDelegation:        0 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		// create funders
		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		params := funderstypes.DefaultParams()
		params.CoinWhitelist[0].MinFundingAmountPerBundle = math.NewInt(10_000)
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(10_000),
		})

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Staker leaves, although he is the next uploader and runs into the upload timeout", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.CommitAfterSeconds(60)

		// leave pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_0, 0)

		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		// check if next uploader is still removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers[0]).To(Equal(i.STAKER_1))

		_, valaccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(valaccountActive).To(BeFalse())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader got not slashed
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Staker leaves, although he was the uploader of the previous round and should receive the uploader reward", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		// leave pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_0, 0)

		// overwrite next uploader for test purposes
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_1
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		// ACT
		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_2))
		Expect(bundleProposal.StorageId).To(Equal("18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4"))

		// check if next uploader is still removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers[0]).To(Equal(i.STAKER_1))

		_, valaccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(valaccountActive).To(BeFalse())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader got not slashed
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))

		// check if next uploader received the uploader reward
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		balanceUploader := s.GetBalanceFromAddress(i.STAKER_0)

		// calculate uploader rewards
		networkFee := s.App().BundlesKeeper.GetNetworkFee(s.Ctx())
		treasuryReward := pool.InflationShareWeight.Mul(networkFee)
		storageReward := s.App().BundlesKeeper.GetStorageCost(s.Ctx(), pool.CurrentStorageProviderId).MulInt64(100)
		totalUploaderReward := pool.InflationShareWeight.Sub(treasuryReward).Sub(storageReward)

		uploaderPayoutReward := totalUploaderReward.Mul(math.LegacyMustNewDecFromStr("0.1"))
		uploaderDelegationReward := totalUploaderReward.Sub(uploaderPayoutReward)

		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0))

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		// assert commission rewards
		// TODO: fix
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), i.STAKER_0).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(uploaderPayoutReward.Add(storageReward).TruncateInt64())))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(uploaderDelegationReward.TruncateInt64())))
	})

	It("Staker leaves, although he was the uploader of the previous round and should get slashed", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		// leave pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_0, 0)

		// overwrite next uploader for test purposes
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_1
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		// ACT
		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		// check if next uploader is still removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers[0]).To(Equal(i.STAKER_1))

		_, valaccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(valaccountActive).To(BeFalse())

		// check if next uploader got slashed
		fraction := s.App().StakersKeeper.GetUploadSlash(s.Ctx())
		slashAmount := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100*i.KYVE - slashAmount))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader did not receive the uploader reward
		balanceUploader := s.GetBalanceFromAddress(i.STAKER_0)

		// assert payout transfer
		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeEmpty())
	})

	It("Staker leaves, although he was a voter in the previous round and should get slashed", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		// leave pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_1, 0)

		// overwrite next uploader for test purposes
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_2
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		// ACT
		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_2_A,
			Staker:        i.STAKER_2,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal("18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4"))

		// check if next uploader is still removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers[0]).To(Equal(i.STAKER_0))

		_, valaccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(valaccountActive).To(BeFalse())

		// check if voter got slashed
		fraction := s.App().StakersKeeper.GetVoteSlash(s.Ctx())
		slashAmount := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(Equal(100*i.KYVE - slashAmount))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader did not receive any rewards
		balanceVoter := s.GetBalanceFromAddress(i.STAKER_1)

		// assert payout transfer
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())
	})

	It("Staker leaves, although he was a voter in the previous round and should get a point", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)

		// leave pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_1, 0)

		// do not vote

		// overwrite next uploader for test purposes
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_0
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		// ACT
		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal("18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4"))

		// check if next uploader is still removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers[0]).To(Equal(i.STAKER_0))

		_, valaccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(valaccountActive).To(BeFalse())

		// check if voter status

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader did not receive any rewards
		balanceVoter := s.GetBalanceFromAddress(i.STAKER_1)

		// assert payout transfer
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())
	})

	It("Staker leaves, although he was a voter who did not vote max points in a row should not get slashed", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CommitAfterSeconds(60)

		maxPoints := int(s.App().BundlesKeeper.GetMaxPoints(s.Ctx()))

		for r := 0; r < maxPoints; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.VALADDRESS_0_A,
				Staker:        i.STAKER_0,
				PoolId:        0,
				StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				DataSize:      100,
				DataHash:      "test_hash",
				FromIndex:     uint64(r * 100),
				BundleSize:    100,
				FromKey:       "test_key",
				ToKey:         "test_key",
				BundleSummary: "test_value",
			})

			s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
				Creator:   i.VALADDRESS_2_A,
				Staker:    i.STAKER_2,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)

		// leave pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_1, 0)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     2400,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_2))
		Expect(bundleProposal.StorageId).To(Equal("18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4"))

		// check if next uploader is still removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers[0]).To(Equal(i.STAKER_0))

		_, valaccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(valaccountActive).To(BeFalse())

		// check if voter not got slashed
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader did not receive any rewards
		balanceVoter := s.GetBalanceFromAddress(i.STAKER_1)

		// assert payout transfer
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		// assert uploader self delegation rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeEmpty())
	})
})
