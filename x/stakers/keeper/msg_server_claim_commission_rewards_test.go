package keeper_test

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_claim_commission_rewards.go

* Produce a valid bundle and check commission rewards
* Claim with non-staker account
* Claim more rewards than available
* Claim zero rewards
* Claim partial rewards
* Claim partial rewards twice
* Claim all rewards
* Claim multiple coins
* Claim on coin of multiple coins

*/

var _ = Describe("msg_server_claim_commission_rewards.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceStaker0 := s.GetCoinsFromAddress(i.STAKER_0)
	amountPerBundle := int64(10_000)

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// set storage cost to 0.5
		bundleParams := s.App().BundlesKeeper.GetParams(s.Ctx())
		bundleParams.StorageCosts = append(bundleParams.StorageCosts, bundletypes.StorageCost{StorageProviderId: 1, Cost: math.LegacyMustNewDecFromStr("0.5")})
		s.App().BundlesKeeper.SetParams(s.Ctx(), bundleParams)

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
			InflationShareWeight: 10_000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    1,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

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

		initialBalanceStaker0 = s.GetCoinsFromAddress(i.STAKER_0)

		s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		params := funderstypes.DefaultParams()

		params.CoinWhitelist[0].MinFundingAmountPerBundle = uint64(amountPerBundle)
		s.App().FundersKeeper.SetParams(s.Ctx(), params)

		// create a valid bundle so that uploader earns commission rewards
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:          i.ALICE,
			PoolId:           0,
			Amounts:          i.KYVECoins(100 * i.T_KYVE),
			AmountsPerBundle: i.KYVECoins(amountPerBundle),
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
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce a valid bundle and check commission rewards", func() {
		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check uploader rewards
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		balanceUploader := s.GetCoinsFromAddress(i.STAKER_0)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert uploader self delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1)
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.KYVECoins(8865).String()))

		// assert commission rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * 0.1 + (100 * 0.5)
		Expect(uploader.CommissionRewards.String()).To(Equal(i.KYVECoins(1035).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.KYVECoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Claim with non-staker account", func() {
		// ARRANGE
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		commissionRewardsBefore := uploader.CommissionRewards

		// ACT
		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_2,
			Amount:  i.KYVECoins(1),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(errors.Wrapf(errorsTypes.ErrNotFound, stakertypes.ErrNoStaker.Error(), i.STAKER_2).Error()))

		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(uploader.CommissionRewards.String()).To(Equal(commissionRewardsBefore.String()))
		Expect(s.GetCoinsFromAddress(i.STAKER_0).String()).To(Equal(initialBalanceStaker0.String()))
	})

	It("Claim more rewards than available", func() {
		// ARRANGE
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		commissionRewardsBefore := uploader.CommissionRewards

		// ACT
		_, err := s.RunTx(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_1,
			Amount:  uploader.CommissionRewards.Add(i.KYVECoin(1)),
		})

		// ASSERT
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(stakertypes.ErrNotEnoughRewards.Error()))

		// assert commission rewards
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(uploader.CommissionRewards.String()).To(Equal(commissionRewardsBefore.String()))
		Expect(s.GetCoinsFromAddress(i.STAKER_0).String()).To(Equal(initialBalanceStaker0.String()))
	})

	It("Claim zero rewards", func() {
		// ARRANGE
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		commissionRewardsBefore := uploader.CommissionRewards

		// ACT
		s.RunTxStakersSuccess(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  sdk.NewCoins(),
		})

		// ASSERT
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(uploader.CommissionRewards.String()).To(Equal(commissionRewardsBefore.String()))
		Expect(s.GetCoinsFromAddress(i.STAKER_0).String()).To(Equal(initialBalanceStaker0.String()))
	})

	It("Claim partial rewards", func() {
		// ARRANGE
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		commissionRewardsBefore := uploader.CommissionRewards

		// ACT
		s.RunTxStakersSuccess(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  i.KYVECoins(100),
		})

		// ASSERT
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(uploader.CommissionRewards.String()).To(Equal(commissionRewardsBefore.Sub(i.KYVECoin(100)).String()))
		Expect(s.GetCoinsFromAddress(i.STAKER_0).String()).To(Equal(initialBalanceStaker0.Add(i.KYVECoin(100)).String()))
	})

	It("Claim partial rewards twice", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  i.KYVECoins(100),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_0_A,
			Staker:    i.STAKER_0,
			PoolId:    0,
			StorageId: "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "SsdTPx9adtpwAGIjiHilqVPEfoTiq7eRw6khbVxKetQ",
			DataSize:      100,
			DataHash:      "test_hash3",
			FromIndex:     200,
			BundleSize:    100,
			FromKey:       "200",
			ToKey:         "299",
			BundleSummary: "test_value3",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "SsdTPx9adtpwAGIjiHilqVPEfoTiq7eRw6khbVxKetQ",
			Vote:      1,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.VALADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "iW1jN99yH_gdQtRhf5J_lVwOIu8p_i7FyxEgoQAkWxU",
			DataSize:      100,
			DataHash:      "test_hash4",
			FromIndex:     300,
			BundleSize:    100,
			FromKey:       "300",
			ToKey:         "399",
			BundleSummary: "test_value4",
		})

		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		commissionRewardsBefore := uploader.CommissionRewards

		// ACT
		s.RunTxSuccess(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  i.KYVECoins(200),
		})

		// ASSERT
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(uploader.CommissionRewards.String()).To(Equal(commissionRewardsBefore.Sub(i.KYVECoin(200)).String()))
		Expect(s.GetCoinsFromAddress(i.STAKER_0).String()).To(Equal(initialBalanceStaker0.Add(i.KYVECoin(300)).String()))
	})

	It("Claim all rewards", func() {
		// ARRANGE
		uploader, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		commissionRewardsBefore := uploader.CommissionRewards

		// ACT
		s.RunTxSuccess(&stakertypes.MsgClaimCommissionRewards{
			Creator: i.STAKER_0,
			Amount:  uploader.CommissionRewards,
		})

		// ASSERT
		uploader, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(uploader.CommissionRewards).To(BeEmpty())
		Expect(s.GetCoinsFromAddress(i.STAKER_0).String()).To(Equal(initialBalanceStaker0.Add(commissionRewardsBefore...).String()))
	})
})
