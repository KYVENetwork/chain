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

TEST CASES - grpc_query_can_vote.go

* Call can vote if pool does not exist
* Call can vote if pool is currently upgrading
* Call can vote if pool is disabled
* Call can vote if pool has not reached the minimum stake
* Call can vote with a valaccount which does not exist
* Call can vote if current bundle was dropped
* Call can vote with a different storage id than the current one
* Call can vote if voter has already voted valid
* Call can vote if voter has already voted invalid
* Call can vote if voter has already voted abstain
* Call can vote on an active pool with a data bundle with valid args
* Call can vote on an active pool with no funds and a data bundle with valid args

*/

var _ = Describe("grpc_query_can_vote.go", Ordered, func() {
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

		s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
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
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Call can vote if pool does not exist", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    1,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(errors.Wrapf(errorsTypes.ErrNotFound, pooltypes.ErrPoolNotFound.Error(), 1).Error()))

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    1,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if pool is currently upgrading", func() {
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
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrPoolCurrentlyUpgrading.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if pool is disabled", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Disabled = true

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrPoolDisabled.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if pool has not reached the minimum stake", func() {
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
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrMinDelegationNotReached.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote with a valaccount which does not exist", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_0,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(stakertypes.ErrValaccountUnauthorized.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_0,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if previous bundle was dropped", func() {
		// ARRANGE
		// wait for timeout so bundle gets dropped
		s.CommitAfterSeconds(s.App().BundlesKeeper.GetUploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(1)

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrBundleDropped.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote with a different storage id than the current one", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "another_test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrInvalidStorageId.Error()))

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "another_test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if voter has already voted valid", func() {
		// ARRANGE
		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).To(BeNil())

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrAlreadyVotedValid.Error()))

		_, txErr = s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if voter has already voted invalid", func() {
		// ARRANGE
		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		Expect(txErr).To(BeNil())

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrAlreadyVotedInvalid.Error()))

		_, txErr = s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if voter has already voted abstain", func() {
		// ARRANGE
		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		Expect(txErr).To(BeNil())

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeTrue())
		Expect(canVote.Reason).To(Equal("KYVE_VOTE_NO_ABSTAIN_ALLOWED"))

		_, txErr = s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).To(BeNil())
	})

	It("Call can vote on an active pool with a data bundle with valid args", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeTrue())
		Expect(canVote.Reason).To(BeEmpty())

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).To(BeNil())
	})

	It("Call can vote on an active pool with no funds and a data bundle with valid args", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&funderstypes.MsgDefundPool{
			Creator: i.ALICE,
			PoolId:  0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1_A,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeTrue())
		Expect(canVote.Reason).To(BeEmpty())

		_, txErr := s.RunTx(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		Expect(txErr).To(BeNil())
	})
})
