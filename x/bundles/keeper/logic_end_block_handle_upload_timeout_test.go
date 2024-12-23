package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/util"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_end_block_handle_upload_timeout.go

* Next uploader can stay although pool ran out of funds
* First staker who joins gets automatically chosen as next uploader
* Next uploader gets removed due to pool upgrading
* Next uploader gets removed due to pool being disabled
* Next uploader gets removed due to pool not reaching min delegation
* Next uploader gets not removed although pool having one node with more than 50% voting power
* Staker is next uploader of genesis bundle and upload interval and timeout does not pass
* Staker is next uploader of genesis bundle and upload timeout does not pass but upload interval passes
* Staker is next uploader of genesis bundle and upload timeout does pass together with upload interval
* Staker is next uploader of bundle proposal and upload interval does not pass
* Staker is next uploader of bundle proposal and upload timeout does not pass
* Staker is next uploader of bundle proposal and upload timeout passes with the previous bundle being valid
* Staker is next uploader of bundle proposal and upload timeout passes with the previous bundle not reaching quorum
* Staker is next uploader of bundle proposal and upload timeout passes with the previous bundle being invalid
* Staker with already max points is next uploader of bundle proposal and upload timeout passes
* A bundle proposal with no quorum does not reach the upload interval
* A bundle proposal with no quorum does reach the upload interval
* Staker who just left the pool is next uploader of dropped bundle proposal and upload timeout passes
* Staker who just left the pool is next uploader of valid bundle proposal and upload timeout passes
* Staker who just left the pool is next uploader of invalid bundle proposal and upload timeout passes
* Staker with already max points is next uploader of bundle proposal in a second pool and upload timeout passes

*/

var _ = Describe("logic_end_block_handle_upload_timeout.go", Ordered, func() {
	var s *i.KeeperTestSuite

	var originalRoundRobinProgress bundletypes.RoundRobinProgress
	var originalBundleProposal bundletypes.BundleProposal

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
			InflationShareWeight: math.LegacyNewDec(10_000),
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(1 * i.T_KYVE),
		})

		originalRoundRobinProgress, _ = s.App().BundlesKeeper.GetRoundRobinProgress(s.Ctx(), 0)
		originalBundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		s.CreateValidatorWithoutCommit(i.STAKER_0, "Staker-0", int64(100*i.KYVE))
		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgress)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposal)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Next uploader can stay although pool ran out of funds", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.RunTxPoolSuccess(&funderstypes.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amounts: i.KYVECoins(100 * i.T_KYVE),
		})

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Last staker who joins gets automatically chosen as next uploader", func() {
		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Next uploader gets removed due to pool upgrading", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		pool.UpgradePlan = &pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix()),
			Duration:    3600,
		}

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Next uploader gets removed due to pool being disabled", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Next uploader gets removed due to pool not reaching min delegation", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.RunTxSuccess(stakingTypes.NewMsgUndelegate(
			i.STAKER_0,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(80*i.KYVE)),
		))

		s.RunTxSuccess(stakingTypes.NewMsgUndelegate(
			i.STAKER_1,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(80*i.KYVE)),
		))

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(20 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(20 * i.KYVE))
	})

	It("Next uploader gets removed due to pool having validators with too much voting power", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.SetMaxVotingPower("0.2")

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(BeZero())

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(BeZero())
	})

	It("Staker is next uploader of genesis bundle and upload interval and timeout does not pass", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of genesis bundle and upload timeout does not pass but upload interval passes", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ACT
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of genesis bundle and upload timeout does pass together with upload interval", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		// check if next uploader got not removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		// check if next uploader received a point
		valaccount, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(valaccount.Points).To(Equal(uint64(1)))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader not got slashed
		expectedBalance := 100 * i.KYVE
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
	})

	It("Staker is next uploader of bundle proposal and upload interval does not pass", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of bundle proposal and upload timeout does not pass", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		// ACT
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of bundle proposal and upload timeout passes with the previous bundle being valid", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal(""))

		// check that previous bundle got finalized
		finalizedBundle, _ := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		// check that nobody got removed from the pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		// check that staker 0 has no points
		valaccount, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(valaccount.Points).To(Equal(uint64(0)))

		// check that staker 1 (next uploader) received a point for not uploading
		valaccount, _ = s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(valaccount.Points).To(Equal(uint64(1)))

		// check that nobody got slashed
		expectedBalance := 100 * i.KYVE
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)))

		// pool delegations equals delegations of staker 0 & 1
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))
	})

	It("Staker is next uploader of bundle proposal and upload timeout passes with the previous bundle not reaching quorum", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal(""))

		// check that bundle didn't get finalized
		_, found := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check that nobody got removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers).To(ContainElements(i.STAKER_0, i.STAKER_1))

		// check that staker 0 (uploader) has no points
		valaccount, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(valaccount.Points).To(Equal(uint64(0)))

		// check that staker 1 received a point for not voting
		valaccount, _ = s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(valaccount.Points).To(Equal(uint64(1)))

		// check that nobody got slashed
		expectedBalance := 100 * i.KYVE
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())

		// pool delegations equals delegations of staker 0 & 1
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))
	})

	It("Staker is next uploader of bundle proposal and upload timeout passes with the previous bundle being invalid", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgress)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposal)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
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

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(300 * i.KYVE))

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_2))
		Expect(bundleProposal.StorageId).To(Equal(""))

		// check that bundle didn't get finalized
		_, found := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check that staker 0 (uploader) got removed from pool because his proposal was voted invalid
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers).To(ContainElements(i.STAKER_1, i.STAKER_2))

		// check that staker 1 (next uploader) received a point for missing the upload
		valaccount, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(valaccount.Points).To(Equal(uint64(1)))

		// check that staker 2 has a no points
		valaccount, _ = s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(valaccount.Points).To(Equal(uint64(0)))

		// check that staker 0 (uploader) got slashed
		expectedBalance := 100*i.KYVE - uint64(s.App().StakersKeeper.GetUploadSlash(s.Ctx()).Mul(math.LegacyNewDec(int64(100*i.KYVE))).TruncateInt64())
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))

		// check that staker 1 (next uploader) didn't get slashed
		expectedBalance = 100 * i.KYVE
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)))

		// check that staker 2 didn't get slashed
		expectedBalance = 100 * i.KYVE
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.STAKER_2)))

		// pool delegations equals delegations of staker 1 & 2
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))
	})

	It("Staker with already max points is next uploader of bundle proposal and upload timeout passes", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgress)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposal)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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

		s.CommitAfterSeconds(60)

		maxPoints := int(s.App().BundlesKeeper.GetMaxPoints(s.Ctx())) - 1

		for r := 1; r <= maxPoints; r++ {
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
				Creator:   i.VALADDRESS_1_A,
				Staker:    i.STAKER_1,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// overwrite next uploader with staker_2
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_2
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal(""))

		// check that previous bundle got finalized
		finalizedBundle, _ := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		// check that staker 2 (next uploader) got removed from pool because he didn't upload
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))
		Expect(poolStakers).To(ContainElements(i.STAKER_0, i.STAKER_1))

		// check that staker 0 (uploader) has no points
		valaccount, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(valaccount.Points).To(Equal(uint64(0)))

		// check that staker 2 has a no points
		valaccount, _ = s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(valaccount.Points).To(Equal(uint64(0)))

		// check that staker 0 (uploader) didn't get slashed
		expectedBalance := 100 * i.KYVE
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))

		// check that staker 1 didn't get slashed
		expectedBalance = 100 * i.KYVE
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)))

		// check that staker 2 (next uploader) got slashed
		expectedBalance = 100*i.KYVE - uint64(math.LegacyNewDec(int64(100*i.KYVE)).Mul(s.App().StakersKeeper.GetTimeoutSlash(s.Ctx())).TruncateInt64())
		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.STAKER_2)))

		// pool delegations equals delegations of staker 0 & 1
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))
	})

	It("A bundle proposal with no quorum does not reach the upload interval", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("A bundle proposal with no quorum does reach the upload interval", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		// ACT
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.DataSize).To(BeZero())
		Expect(bundleProposal.DataHash).To(BeEmpty())
		Expect(bundleProposal.BundleSize).To(BeZero())
		Expect(bundleProposal.FromKey).To(BeEmpty())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.BundleSummary).To(BeEmpty())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_1)
		Expect(found).To(BeTrue())
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_1, 0)).To(Equal(100 * i.KYVE))
	})

	It("Staker who just left the pool is next uploader of dropped bundle proposal and upload timeout passes", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgress)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposal)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		// remove valaccount directly from pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_1, 0)

		// ACT
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		// check if next uploader got removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_2)
		Expect(found).To(BeTrue())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader not got slashed
		expectedBalance := 100 * i.KYVE

		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_1, i.STAKER_1)))
	})

	It("Staker who just left the pool is next uploader of valid bundle proposal and upload timeout passes", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgress)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposal)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
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

		// remove valaccount directly from pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_1, 0)

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal(""))

		// check if previous bundle got finalized
		finalizedBundle, _ := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		// check if next uploader got removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_2)
		Expect(found).To(BeTrue())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check if next uploader not got slashed
		expectedBalance := 100 * i.KYVE

		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
	})

	It("Staker who just left the pool is next uploader of invalid bundle proposal and upload timeout passes", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgress)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposal)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.VALADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
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
			FromKey:       "0",
			ToKey:         "99",
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

		// remove valaccount directly from pool
		s.App().StakersKeeper.RemovePoolAccountFromPool(s.Ctx(), i.STAKER_0, 0)

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()) + 60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_2))
		Expect(bundleProposal.StorageId).To(Equal(""))

		// check that bundle didn't get finalized
		_, found := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check that next uploader got removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(200 * i.KYVE))

		// check that uploader received upload slash
		expectedBalance := 100*i.KYVE - uint64(s.App().StakersKeeper.GetUploadSlash(s.Ctx()).Mul(math.LegacyNewDec(int64(100*i.KYVE))).TruncateInt64())

		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
	})

	It("Staker with already max points is next uploader of bundle proposal in a second pool and upload timeout passes", func() {
		// ARRANGE
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
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		originalRoundRobinProgressPool2, _ := s.App().BundlesKeeper.GetRoundRobinProgress(s.Ctx(), 1)
		originalBundleProposalPool2, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 1)

		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           1,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(1 * i.T_KYVE),
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        1,
			PoolAddress:   i.VALADDRESS_0_B,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        1,
			PoolAddress:   i.VALADDRESS_1_B,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))
		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgress)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposal)
		s.App().BundlesKeeper.SetRoundRobinProgress(s.Ctx(), originalRoundRobinProgressPool2)
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), originalBundleProposalPool2)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        1,
			PoolAddress:   i.VALADDRESS_2_B,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_1_B,
			Staker:  i.STAKER_1,
			PoolId:  1,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_B,
			Staker:        i.STAKER_1,
			PoolId:        1,
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
			Creator:   i.VALADDRESS_0_B,
			Staker:    i.STAKER_0,
			PoolId:    1,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)

		maxPoints := int(s.App().BundlesKeeper.GetMaxPoints(s.Ctx())) - 1

		for r := 1; r <= maxPoints; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 1)
			bundleProposal.NextUploader = i.STAKER_1
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.VALADDRESS_1_B,
				Staker:        i.STAKER_1,
				PoolId:        1,
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
				Creator:   i.VALADDRESS_0_B,
				Staker:    i.STAKER_0,
				PoolId:    1,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// overwrite next uploader with staker_1
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 1)
		bundleProposal.NextUploader = i.STAKER_2
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ = s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 1)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal(""))

		// check that bundle didn't get finalized
		_, found := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check if next uploader got not removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 1)
		Expect(poolStakers).To(HaveLen(2))

		// check if next uploader received a point
		_, valaccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 1)
		Expect(valaccountActive).To(BeFalse())

		_, found = s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_2)
		Expect(found).To(BeTrue())

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 1)).To(Equal(200 * i.KYVE))

		// check if next uploader not got slashed
		slashAmountRatio := s.App().StakersKeeper.GetTimeoutSlash(s.Ctx())
		expectedBalance := 100*i.KYVE - uint64(math.LegacyNewDec(int64(100*i.KYVE)).Mul(slashAmountRatio).TruncateInt64())

		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.STAKER_2)))
	})
})
