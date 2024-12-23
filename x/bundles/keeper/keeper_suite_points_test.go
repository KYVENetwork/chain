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

TEST CASES - points

* One validator does not vote for one proposal
* One validator votes after having not voted previously
* One validator does not vote for multiple proposals in a row
* One validator votes after having not voted previously multiple times
* One validator does not vote for multiple proposals and reaches max points
* One validator does not vote for multiple proposals and submits a bundle proposal
* One validator does not vote for multiple proposals and skip the uploader role
* One validator submits a bundle proposal where he reaches max points because he did not vote before

*/

var _ = Describe("points", Ordered, func() {
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

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(50*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CommitAfterSeconds(60)

		// Claim Uploader role for Staker 0
		pool, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		pool.NextUploader = i.STAKER_0
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), pool)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
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
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("One validator does not vote for one proposal", func() {
		// ACT
		// do not vote

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(poolAccountVoter.Points).To(Equal(uint64(1)))
	})

	It("One validator votes after having not voted previously", func() {
		// ARRANGE
		// do not vote

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     200,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(poolAccountVoter.Points).To(BeZero())
	})

	It("One validator does not vote for multiple proposals in a row", func() {
		// ACT
		for r := 1; r <= 3; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.POOL_ADDRESS_0_A,
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
				Creator:   i.POOL_ADDRESS_1_A,
				Staker:    i.STAKER_1,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			// do not vote with voter 3

			s.CommitAfterSeconds(60)
		}

		// ASSERT
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(poolAccountVoter.Points).To(Equal(uint64(3)))
	})

	It("One validator votes after having not voted previously multiple times", func() {
		// ARRANGE
		for r := 1; r <= 3; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.POOL_ADDRESS_0_A,
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
				Creator:   i.POOL_ADDRESS_1_A,
				Staker:    i.STAKER_1,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     400,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(poolAccountVoter.Points).To(BeZero())
	})

	It("One validator does not vote for multiple proposals and reaches max points", func() {
		// ARRANGE
		maxPoints := int(s.App().BundlesKeeper.GetMaxPoints(s.Ctx()))

		// ACT
		for r := 1; r <= maxPoints; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.POOL_ADDRESS_0_A,
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
				Creator:   i.POOL_ADDRESS_1_A,
				Staker:    i.STAKER_1,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// ASSERT
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, stakerFound := s.App().StakersKeeper.GetValidator(s.Ctx(), i.STAKER_2)
		Expect(stakerFound).To(BeTrue())

		_, poolAccountActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(poolAccountActive).To(BeFalse())

		// check if voter got slashed
		slashAmountRatio := s.App().StakersKeeper.GetTimeoutSlash(s.Ctx())
		expectedBalance := 50*i.KYVE - uint64(math.LegacyNewDec(int64(50*i.KYVE)).Mul(slashAmountRatio).TruncateInt64())

		Expect(expectedBalance).To(Equal(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.STAKER_2)))
	})

	It("One validator does not vote for multiple proposals and submits a bundle proposal", func() {
		// ARRANGE
		for r := 1; r <= 3; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.POOL_ADDRESS_0_A,
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
				Creator:   i.POOL_ADDRESS_1_A,
				Staker:    i.STAKER_1,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// ACT
		// overwrite next uploader for test purposes
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_2
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_2_A,
			Staker:        i.STAKER_2,
			PoolId:        0,
			StorageId:     "18SRvVuCrB8vy_OCLBaNbXONMVGeflGcw4gGTZ1oUt4",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     400,
			BundleSize:    100,
			FromKey:       "test_key",
			ToKey:         "test_key",
			BundleSummary: "test_value",
		})

		// ASSERT
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		// points are instantly 1 because node did not vote on this bundle, too
		Expect(poolAccountVoter.Points).To(Equal(uint64(1)))
	})

	It("One validator does not vote for multiple proposals and skip the uploader role", func() {
		// ARRANGE
		for r := 1; r <= 3; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.POOL_ADDRESS_0_A,
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
				Creator:   i.POOL_ADDRESS_1_A,
				Staker:    i.STAKER_1,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// ACT
		// overwrite next uploader for test purposes
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_2
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		s.RunTxBundlesSuccess(&bundletypes.MsgSkipUploaderRole{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			FromIndex: 400,
		})

		// ASSERT
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())
	})

	It("One validator submits a bundle proposal where he reaches max points because he did not vote before", func() {
		// ARRANGE
		maxPoints := int(s.App().BundlesKeeper.GetMaxPoints(s.Ctx())) - 1

		for r := 1; r <= maxPoints; r++ {
			// overwrite next uploader for test purposes
			bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
			bundleProposal.NextUploader = i.STAKER_0
			s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

			s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
				Creator:       i.POOL_ADDRESS_0_A,
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
				Creator:   i.POOL_ADDRESS_1_A,
				Staker:    i.STAKER_1,
				PoolId:    0,
				StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
				Vote:      bundletypes.VOTE_TYPE_VALID,
			})

			s.CommitAfterSeconds(60)

			// do not vote
		}

		// ACT
		// overwrite next uploader for test purposes
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		bundleProposal.NextUploader = i.STAKER_2
		s.App().BundlesKeeper.SetBundleProposal(s.Ctx(), bundleProposal)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_2_A,
			Staker:        i.STAKER_2,
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
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		// points are instantly 1 because node did not vote on this bundle, too
		Expect(poolAccountVoter.Points).To(Equal(uint64(1)))
	})
})
