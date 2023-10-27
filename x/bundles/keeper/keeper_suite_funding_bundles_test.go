package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - funding bundles

* Produce a valid bundle with only one funder
* Produce a valid bundle with multiple funders and same funding amounts
* Produce a valid bundle with multiple funders and different funding amounts
* Produce a valid bundle with multiple funders and different funding amounts where not everyone can afford the funds
* Produce a valid bundle although the only funder can not pay for the bundle reward
* Produce a valid bundle although multiple funders with same amount can not pay for the bundle reward
* Produce a valid bundle although multiple funders with different amount can not pay for the bundle reward
* Produce a valid bundle although there are no funders at all

*/

var _ = Describe("funding bundles", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceAlice := s.GetBalanceFromAddress(i.ALICE)
	initialBalanceBob := s.GetBalanceFromAddress(i.BOB)

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "PoolTest",
			MaxBundleSize:  100,
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  10_000,
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
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce a valid bundle with only one funder", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
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
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		totalReward := pool.OperatingCost

		// assert total pool funds
		Expect(pool.TotalFunds).To(Equal(100*i.KYVE - totalReward))
		Expect(pool.Funders).To(HaveLen(1))

		// assert individual funds
		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(100*i.KYVE - totalReward))

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 100*i.KYVE))
	})

	It("Produce a valid bundle with multiple funders and same funding amounts", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  100 * i.KYVE,
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
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		totalReward := pool.OperatingCost

		// assert total pool funds
		Expect(pool.TotalFunds).To(Equal(200*i.KYVE - totalReward))
		Expect(pool.Funders).To(HaveLen(2))

		// assert individual funds
		fundersCharge := uint64(sdk.NewDec(int64(totalReward)).Quo(sdk.NewDec(2)).TruncateInt64())

		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(100*i.KYVE - fundersCharge))

		Expect(pool.GetFunderAmount(i.BOB)).To(Equal(100*i.KYVE - fundersCharge))

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 100*i.KYVE))

		balanceBob := s.GetBalanceFromAddress(i.BOB)
		Expect(balanceBob).To(Equal(initialBalanceBob - 100*i.KYVE))
	})

	It("Produce a valid bundle with multiple funders and different funding amounts", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  150 * i.KYVE,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  50 * i.KYVE,
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
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		totalReward := pool.OperatingCost

		// assert total pool funds
		Expect(pool.TotalFunds).To(Equal(200*i.KYVE - totalReward))
		Expect(pool.Funders).To(HaveLen(2))

		// assert individual funds
		fundersCharge := uint64(sdk.NewDec(int64(totalReward)).Quo(sdk.NewDec(2)).TruncateInt64())

		Expect(pool.GetFunderAmount(i.ALICE)).To(Equal(150*i.KYVE - fundersCharge))

		Expect(pool.GetFunderAmount(i.BOB)).To(Equal(50*i.KYVE - fundersCharge))

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 150*i.KYVE))

		balanceBob := s.GetBalanceFromAddress(i.BOB)
		Expect(balanceBob).To(Equal(initialBalanceBob - 50*i.KYVE))
	})

	It("Produce a valid bundle with multiple funders and different funding amounts where not everyone can afford the funds", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  10,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  100 * i.KYVE,
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
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		totalReward := pool.OperatingCost

		// assert total pool funds
		Expect(pool.TotalFunds).To(Equal(100*i.KYVE - (totalReward / 2)))
		Expect(pool.Funders).To(HaveLen(1))

		Expect(pool.GetFunderAmount(i.BOB)).To(Equal(100*i.KYVE - (totalReward / 2)))

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 10))

		balanceBob := s.GetBalanceFromAddress(i.BOB)
		Expect(balanceBob).To(Equal(initialBalanceBob - 100*i.KYVE))
	})

	It("Produce a valid bundle although the only funder can not pay for the bundle reward", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  10,
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
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 10))
	})

	It("Produce a valid bundle although multiple funders with same amount can not pay for the bundle reward", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  10,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  10,
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
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 10))

		balanceBob := s.GetBalanceFromAddress(i.BOB)
		Expect(balanceBob).To(Equal(initialBalanceBob - 10))
	})

	It("Produce a dropped bundle because multiple funders with different amount can not pay for the bundle reward", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  10,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  20,
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
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 10))

		balanceBob := s.GetBalanceFromAddress(i.BOB)
		Expect(balanceBob).To(Equal(initialBalanceBob - 20))
	})

	It("Produce a valid bundle although there are no funders at all", func() {
		// ARRANGE
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
			Vote:      1,
		})

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert if bundle go finalized
		Expect(pool.TotalBundles).To(Equal(uint64(1)))
		Expect(pool.CurrentKey).To(Equal("99"))

		// assert total pool funds
		Expect(pool.TotalFunds).To(BeZero())
		Expect(pool.Funders).To(BeEmpty())
	})
})
