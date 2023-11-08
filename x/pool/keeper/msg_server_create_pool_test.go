package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_create_pool.go

* Invalid authority (transaction)
* Invalid authority (proposal)
* Create first pool
* Create another pool
* Create pool with invalid binaries

*/

var _ = Describe("msg_server_create_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	votingPeriod := s.App().GovKeeper.GetParams(s.Ctx()).VotingPeriod

	BeforeEach(func() {
		s = i.NewCleanChain()
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Invalid authority (transaction)", func() {
		// ARRANGE
		msg := &types.MsgCreatePool{
			Authority:            i.DUMMY[0],
			Name:                 "TestPool",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Invalid authority (proposal)", func() {
		// ARRANGE
		msg := &types.MsgCreatePool{
			Authority:            i.DUMMY[0],
			Name:                 "TestPool",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}

		proposal, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, err := s.RunTx(&proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Create first pool", func() {
		// ARRANGE
		msg := &types.MsgCreatePool{
			Authority:            gov,
			Name:                 "TestPool",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool).To(Equal(types.Pool{
			Id:                   0,
			Name:                 "TestPool",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			CurrentKey:           "",
			CurrentSummary:       "",
			CurrentIndex:         0,
			TotalBundles:         0,
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Disabled:             false,
			Protocol: &types.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &types.UpgradePlan{
				Version:     "",
				Binaries:    "",
				ScheduledAt: 0,
				Duration:    0,
			},
			CurrentStorageProviderId: 2,
			CurrentCompressionId:     1,
		}))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(fundingState).To(Equal(funderstypes.FundingState{
			PoolId:                0,
			ActiveFunderAddresses: nil,
		}))
	})

	It("Create another pool", func() {
		// ARRANGE
		msg := &types.MsgCreatePool{
			Authority:            gov,
			Name:                 "TestPool",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		// ACT
		msg = &types.MsgCreatePool{
			Authority:            gov,
			Name:                 "TestPool2",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}

		p, v = BuildGovernanceTxs(s, []sdk.Msg{msg})

		_, submitErr = s.RunTx(&p)
		_, voteErr = s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ARRANGE
		proposal, _ = s.App().GovKeeper.GetProposal(s.Ctx(), 2)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		Expect(pool).To(Equal(types.Pool{
			Id:                   1,
			Name:                 "TestPool2",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			CurrentKey:           "",
			CurrentSummary:       "",
			CurrentIndex:         0,
			TotalBundles:         0,
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Disabled:             false,
			Protocol: &types.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &types.UpgradePlan{
				Version:     "",
				Binaries:    "",
				ScheduledAt: 0,
				Duration:    0,
			},
			CurrentStorageProviderId: 2,
			CurrentCompressionId:     1,
		}))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 1)
		Expect(fundingState).To(Equal(funderstypes.FundingState{
			PoolId:                1,
			ActiveFunderAddresses: nil,
		}))
	})

	It("Create pool with invalid binaries", func() {
		// ARRANGE
		msg := &types.MsgCreatePool{
			Authority:            gov,
			Name:                 "TestPool",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{",
			StorageProviderId:    2,
			CompressionId:        1,
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusFailed))

		_, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(found).To(BeFalse())

		_, found = s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
		Expect(found).To(BeFalse())
	})
})
