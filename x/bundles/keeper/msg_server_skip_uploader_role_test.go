package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_skip_uploader_role.go

* Skip uploader role on data bundle if staker is next uploader
* Skip uploader on data bundle after uploader role has already been skipped
* Skip uploader role on dropped bundle
* Skip uploader role on data bundle with current round containing a valid bundle
* Skip uploader role on data bundle with current round containing an invalid bundle

*/

var _ = Describe("msg_server_skip_uploader_role.go", Ordered, func() {
	var s *i.KeeperTestSuite

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

		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(1 * i.T_KYVE),
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_2,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_2,
			PoolId:     0,
			Valaddress: i.VALADDRESS_2_A,
		})

		s.Commit()
		s.WaitSeconds(60)

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
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Skip uploader role on data bundle if staker is next uploader", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		s.Commit()
		s.WaitSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSkipUploaderRole{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			FromIndex: 100,
		})

		// ASSERT
		bundleProposal, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("0"))
		Expect(bundleProposal.ToKey).To(Equal("99"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_0))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).To(ContainElement(i.STAKER_2))

		// here the next uploader should be always be different after skipping
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))

		// check that the bundle is not finalized
		_, found = s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check that no validator got a point for the second round
		valaccounts := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)
		for _, valaccount := range valaccounts {
			Expect(valaccount.Points).To(BeZero())
		}
	})

	It("Skip uploader on data bundle after uploader role has already been skipped", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		s.Commit()
		s.WaitSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSkipUploaderRole{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			FromIndex: 100,
		})

		s.Commit()
		s.WaitSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSkipUploaderRole{
			Creator:   i.VALADDRESS_0_A,
			Staker:    i.STAKER_0,
			PoolId:    0,
			FromIndex: 100,
		})

		// ASSERT
		bundleProposal, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("0"))
		Expect(bundleProposal.ToKey).To(Equal("99"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_0))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).To(ContainElement(i.STAKER_2))

		// here the next uploader should always be different after skipping
		Expect(bundleProposal.NextUploader).NotTo(Equal(i.STAKER_0))

		// check that the bundle is not finalized
		_, found = s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check that no validator got a point for the second round
		valaccounts := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)
		for _, valaccount := range valaccounts {
			Expect(valaccount.Points).To(BeZero())
		}
	})

	It("Skip uploader role on dropped bundle", func() {
		// ARRANGE
		// create dropped bundle
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// wait for upload interval
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSkipUploaderRole{
			Creator:   i.VALADDRESS_0_A,
			Staker:    i.STAKER_0,
			PoolId:    0,
			FromIndex: 0,
		})

		// ASSERT
		bundleProposal, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.DataSize).To(BeZero())
		Expect(bundleProposal.DataHash).To(BeEmpty())
		Expect(bundleProposal.BundleSize).To(BeZero())
		Expect(bundleProposal.FromKey).To(BeEmpty())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.BundleSummary).To(BeEmpty())
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// here the next uploader should be always be different after skipping
		Expect(bundleProposal.NextUploader).NotTo(Equal(i.STAKER_0))

		// check that the bundle is not finalized
		_, found = s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check that no validator got a point for the second round
		valaccounts := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)
		for _, valaccount := range valaccounts {
			if valaccount.Staker == i.STAKER_0 {
				Expect(valaccount.Points).To(BeZero())
			} else {
				// All others have one point because they did not vote at all
				Expect(valaccount.Points).To(Equal(uint64(1)))
			}
		}
	})

	It("Skip uploader role on data bundle with current round containing a valid bundle", func() {
		// ARRANGE
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

		s.Commit()
		s.WaitSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSkipUploaderRole{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			FromIndex: 100,
		})

		// ASSERT
		bundleProposal, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.DataSize).To(BeZero())
		Expect(bundleProposal.DataHash).To(BeEmpty())
		Expect(bundleProposal.BundleSize).To(BeZero())
		Expect(bundleProposal.FromKey).To(BeEmpty())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.BundleSummary).To(BeEmpty())
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		finalizedBundle, found := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))

		// check if no validator got a point for the second round
		valaccounts := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)
		for _, valaccount := range valaccounts {
			Expect(valaccount.Points).To(BeZero())
		}
	})

	It("Skip uploader role on data bundle with current round containing an invalid bundle", func() {
		// ARRANGE
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

		s.Commit()
		s.WaitSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSkipUploaderRole{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			FromIndex: 100,
		})

		// ASSERT
		bundleProposal, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_2))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.DataSize).To(BeZero())
		Expect(bundleProposal.DataHash).To(BeEmpty())
		Expect(bundleProposal.BundleSize).To(BeZero())
		Expect(bundleProposal.FromKey).To(BeEmpty())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.BundleSummary).To(BeEmpty())
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check that the bundle is not finalized
		_, found = s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(found).To(BeFalse())

		// check that no validator got a point
		valaccounts := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)
		for _, valaccount := range valaccounts {
			Expect(valaccount.Points).To(BeZero())
		}
	})
})
