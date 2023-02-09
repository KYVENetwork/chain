package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_bundles.go

* Assert pool can run while pool is upgrading
* Assert pool can run while pool is disabled
* Assert pool can run while pool has no funds
* Assert pool can run while min delegation is not reached
* Assert pool can run

* Assert can vote if sender is no staker
* Assert can vote if bundle is dropped
* Assert can vote if storage id does not match
* Assert can vote if sender has already voted valid
* Assert can vote if sender has already voted invalid
* Assert can vote

* Assert can propose if sender is no staker
* Assert can propose if sender is not next uploader
* Assert can propose if upload interval has not passed
* Assert can propose if index does not match
* Assert can propose

*/

var _ = Describe("logic_bundles.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	// ASSERT POOL CAN RUN

	It("Assert pool can run while pool is upgrading", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{
				Version:     "1.0.0",
				Binaries:    "{}",
				ScheduledAt: uint64(s.Ctx().BlockTime().Unix()),
				Duration:    60,
			},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertPoolCanRun(s.Ctx(), 0)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert pool can run while pool is disabled", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Disabled:       true,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertPoolCanRun(s.Ctx(), 0)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert pool can run while pool has no funds", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertPoolCanRun(s.Ctx(), 0)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert pool can run while min delegation is not reached", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  99 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertPoolCanRun(s.Ctx(), 0)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert pool can run", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertPoolCanRun(s.Ctx(), 0)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
	})

	// ASSERT CAN VOTE

	It("Assert can vote if sender is no staker", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundlesTypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.CommitAfterSeconds(60)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertCanVote(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI")

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can vote if bundle is dropped", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundlesTypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     0,
		})

		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(1)

		// ACT
		err := s.App().BundlesKeeper.AssertCanVote(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI")

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can vote if storage id does not match", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundlesTypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "another_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.CommitAfterSeconds(60)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     0,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertCanVote(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI")

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can vote if sender has already voted valid", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundlesTypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.CommitAfterSeconds(60)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      1,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertCanVote(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI")

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can vote if sender has already voted invalid", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundlesTypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.CommitAfterSeconds(60)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      2,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertCanVote(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI")

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can vote", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundlesTypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0,
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

		s.CommitAfterSeconds(60)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     0,
		})

		// ACT
		err := s.App().BundlesKeeper.AssertCanVote(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI")

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
	})

	// ASSERT CAN PROPOSE

	It("Assert can propose if sender is no staker", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		// ACT
		err := s.App().BundlesKeeper.AssertCanPropose(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, 0)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can propose if sender is not next uploader", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     0,
		})

		s.CommitAfterSeconds(60)

		// ACT
		err := s.App().BundlesKeeper.AssertCanPropose(s.Ctx(), 0, i.STAKER_1, i.VALADDRESS_1, 0)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can propose if upload interval has not passed", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(30)

		// ACT
		err := s.App().BundlesKeeper.AssertCanPropose(s.Ctx(), 0, i.STAKER_0, i.VALADDRESS_0, 0)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can propose if index does not match", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		// ACT
		err := s.App().BundlesKeeper.AssertCanPropose(s.Ctx(), 0, i.STAKER_0, i.VALADDRESS_0, 1000)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Assert can propose", func() {
		// ASSERT
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			UploadInterval: 60,
			OperatingCost:  2 * i.KYVE,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     0,
		})

		s.RunTxBundlesSuccess(&bundlesTypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		// ACT
		err := s.App().BundlesKeeper.AssertCanPropose(s.Ctx(), 0, i.STAKER_0, i.VALADDRESS_0, 0)

		// ASSERT
		Expect(err).NotTo(HaveOccurred())
	})
})
