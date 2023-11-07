package keeper_test

import (
	"cosmossdk.io/errors"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	querytypes "github.com/KYVENetwork/chain/x/query/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - grpc_query_can_propose.go

* Call can propose if pool does not exist
* Call can propose if pool is currently upgrading
* Call can propose if pool is disabled
* Call can propose if pool is out of funds
* Call can propose if pool has not reached the minimum stake
* Call can propose with a valaccount which does not exist
* Call can propose as a staker who is not the next uploader
* Call can propose before the upload interval passed
* Call can propose with an invalid from height
* Call can propose on an active pool as the next uploader with valid args
* Call can propose on an active pool without funds as the next uploader with valid args

*/

var _ = Describe("grpc_query_can_propose.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		s = i.NewCleanChain()

		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:      gov,
			MinDelegation:  200 * i.KYVE,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Binaries:       "{}",
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     0,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     0,
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
			StorageId:     "test_storage_id",
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
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Call can propose if pool does not exist", func() {
		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    1,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeFalse())
		Expect(canPropose.Reason).To(Equal(errors.Wrapf(errorsTypes.ErrNotFound, pooltypes.ErrPoolNotFound.Error(), 1).Error()))

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        1,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canPropose.Reason))
	})

	It("Call can propose if pool is currently upgrading", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.UpgradePlan = &pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: 100,
			Duration:    3600,
		}

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeFalse())
		Expect(canPropose.Reason).To(Equal(bundletypes.ErrPoolCurrentlyUpgrading.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canPropose.Reason))
	})

	It("Call can propose if pool is disabled", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Disabled = true

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeFalse())
		Expect(canPropose.Reason).To(Equal(bundletypes.ErrPoolDisabled.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canPropose.Reason))
	})

	It("Call can propose if pool has not reached the minimum stake", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&delegationtypes.MsgUndelegate{
			Creator: i.STAKER_0,
			Staker:  i.STAKER_0,
			Amount:  50 * i.KYVE,
		})

		// wait for unbonding
		s.CommitAfterSeconds(s.App().DelegationKeeper.GetUnbondingDelegationTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeFalse())
		Expect(canPropose.Reason).To(Equal(bundletypes.ErrMinDelegationNotReached.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canPropose.Reason))
	})

	It("Call can propose with a valaccount which does not exist", func() {
		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_0,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeFalse())
		Expect(canPropose.Reason).To(Equal(stakertypes.ErrValaccountUnauthorized.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canPropose.Reason))
	})

	It("Call can propose as a staker who is not the next uploader", func() {
		// ARRANGE
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		var canPropose *querytypes.QueryCanProposeResponse
		var err error

		// ACT
		if bundleProposal.NextUploader == i.STAKER_0 {
			canPropose, err = s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
				PoolId:    0,
				Staker:    i.STAKER_1,
				Proposer:  i.VALADDRESS_1_A,
				FromIndex: 100,
			})
		} else {
			canPropose, err = s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
				PoolId:    0,
				Staker:    i.STAKER_0,
				Proposer:  i.VALADDRESS_0_A,
				FromIndex: 100,
			})
		}

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeFalse())

		if bundleProposal.NextUploader == i.STAKER_0 {
			Expect(canPropose.Reason).To(Equal(errors.Wrapf(bundletypes.ErrNotDesignatedUploader, "expected %v received %v", i.STAKER_0, i.STAKER_1).Error()))
		} else {
			Expect(canPropose.Reason).To(Equal(errors.Wrapf(bundletypes.ErrNotDesignatedUploader, "expected %v received %v", i.STAKER_1, i.STAKER_0).Error()))
		}

		var txErr error

		if bundleProposal.NextUploader == i.STAKER_0 {
			_, txErr = s.RunTx(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.VALADDRESS_1_A,
				Staker:        i.STAKER_1,
				PoolId:        0,
				StorageId:     "test_storage_id",
				DataSize:      100,
				DataHash:      "test_hash",
				FromIndex:     100,
				BundleSize:    100,
				FromKey:       "100",
				ToKey:         "199",
				BundleSummary: "test_value",
			})
		} else {
			_, txErr = s.RunTx(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.VALADDRESS_0_A,
				Staker:        i.STAKER_0,
				PoolId:        0,
				StorageId:     "test_storage_id",
				DataSize:      100,
				DataHash:      "test_hash",
				FromIndex:     100,
				BundleSize:    100,
				FromKey:       "100",
				ToKey:         "199",
				BundleSummary: "test_value",
			})
		}

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canPropose.Reason))
	})

	It("Call can propose before the upload interval passed", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		// increase upload interval for upload timeout
		pool.UploadInterval = 120

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(canPropose.Possible).To(BeFalse())
		Expect(canPropose.Reason).To(Equal(errors.Wrapf(bundletypes.ErrUploadInterval, "expected %v < %v", s.Ctx().BlockTime().Unix(), bundleProposal.UpdatedAt+pool.UploadInterval).Error()))

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canPropose.Reason))
	})

	It("Call can propose with an invalid from index", func() {
		// ACT
		canPropose_1, err_1 := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 99,
		})

		canPropose_2, err_2 := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 101,
		})

		// ASSERT
		Expect(err_1).To(BeNil())
		Expect(err_2).To(BeNil())

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(canPropose_1.Possible).To(BeFalse())
		Expect(canPropose_1.Reason).To(Equal(errors.Wrapf(bundletypes.ErrFromIndex, "expected %v received %v", pool.CurrentIndex+bundleProposal.BundleSize, 99).Error()))

		Expect(canPropose_2.Possible).To(BeFalse())
		Expect(canPropose_2.Reason).To(Equal(errors.Wrapf(bundletypes.ErrFromIndex, "expected %v received %v", pool.CurrentIndex+bundleProposal.BundleSize, 101).Error()))

		_, txErr_1 := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     99,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr_1).NotTo(BeNil())
		Expect(txErr_1.Error()).To(Equal(canPropose_1.Reason))

		_, txErr_2 := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     101,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr_2).NotTo(BeNil())
		Expect(txErr_2.Error()).To(Equal(canPropose_2.Reason))
	})

	It("Call can propose on an active pool as the next uploader with valid args", func() {
		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeTrue())
		Expect(canPropose.Reason).To(BeEmpty())

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).To(BeNil())
	})

	It("Call can propose on an active pool without funds as the next uploader with valid args", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&funderstypes.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		canPropose, err := s.App().QueryKeeper.CanPropose(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanProposeRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Proposer:  i.VALADDRESS_1_A,
			FromIndex: 100,
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canPropose.Possible).To(BeTrue())
		Expect(canPropose.Reason).To(BeEmpty())

		_, txErr := s.RunTx(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "test_storage_id",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value",
		})

		Expect(txErr).To(BeNil())
	})
})
